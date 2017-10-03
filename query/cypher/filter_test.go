package cypher_test

import (
	"testing"

	"github.com/RossMerr/Caudex.Graph/expressions"
	"github.com/RossMerr/Caudex.Graph/vertices"

	"github.com/RossMerr/Caudex.Graph/query"
	"github.com/RossMerr/Caudex.Graph/query/cypher"
	"github.com/RossMerr/Caudex.Graph/query/cypher/ast"
)

func Test_Filter(t *testing.T) {
	state := false
	filter := cypher.NewFilter()
	var tests = []struct {
		iterator  query.IteratorFrontier
		predicate ast.Expr
		count     int
	}{
		{
			iterator: func() (*query.Frontier, bool) {
				return nil, false
			},
			predicate: nil,
			count:     0,
		},
		{
			iterator: func() (*query.Frontier, bool) {
				f := query.Frontier{}
				return &f, true
			},
			predicate: ast.NewComparisonExpr(expressions.EQ, &ast.PropertyStmt{Variable: "n", Value: "name"}, &ast.Ident{Data: "foo"}),
			count:     0,
		},
		{
			iterator: func() (*query.Frontier, bool) {
				state = expressions.XORSwap(state)
				v, _ := vertices.NewVertex()
				v.SetProperty("name", "foo")
				f := query.NewFrontier(v, "n")
				return &f, state
			},
			predicate: ast.NewComparisonExpr(expressions.EQ, &ast.PropertyStmt{Variable: "n", Value: "name"}, &ast.Ident{Data: "foo"}),
			count:     1,
		},
		{
			iterator: func() (*query.Frontier, bool) {
				state = expressions.XORSwap(state)
				v, _ := vertices.NewVertex()
				v.SetProperty("name", "foo")
				f := query.NewFrontier(v, "n")
				return &f, state
			},
			predicate: nil,
			count:     1,
		},
	}

	for i, tt := range tests {
		result := filter.Filter(tt.iterator, tt.predicate)
		count := 0
		for v, ok := result(); ok; v, ok = result() {
			count++
			if v == nil {
				t.Errorf("%d %+v", i, v)
			}
		}
		if count != tt.count {
			t.Errorf("%d. expected %d got %d", i, tt.count, count)
		}
	}
}

func Test_ExpressionEvaluator(t *testing.T) {

	filter := cypher.NewFilter()
	var tests = []struct {
		expr     ast.Expr
		variable string
		v        *vertices.Vertex
		result   bool
	}{
		{
			expr:     &ast.PropertyStmt{Variable: "n", Value: "name"},
			variable: "n",
			v: func() *vertices.Vertex {
				x, _ := vertices.NewVertex()
				x.SetProperty("name", "foo")
				return x
			}(),
			result: false,
		},
		{
			expr:     ast.NewComparisonExpr(expressions.EQ, &ast.PropertyStmt{Variable: "n", Value: "name"}, &ast.Ident{Data: "foo"}),
			variable: "n",
			v: func() *vertices.Vertex {
				x, _ := vertices.NewVertex()
				x.SetProperty("name", "foo")
				return x
			}(),
			result: true,
		},
		{
			expr:     ast.NewComparisonExpr(expressions.EQ, &ast.PropertyStmt{Variable: "n", Value: "name"}, &ast.Ident{Data: "foo"}),
			variable: "x",
			v: func() *vertices.Vertex {
				x, _ := vertices.NewVertex()
				x.SetProperty("name", "foo")
				return x
			}(),
			result: false,
		},
		{
			expr:     ast.NewBooleanExpr(expressions.OR, ast.NewComparisonExpr(expressions.EQ, &ast.PropertyStmt{Variable: "n", Value: "name"}, &ast.Ident{Data: "foo"}), nil),
			variable: "n",
			v: func() *vertices.Vertex {
				x, _ := vertices.NewVertex()
				x.SetProperty("name", "foo")
				return x
			}(),
			result: true,
		},
	}
	for i, tt := range tests {
		result := filter.ExpressionEvaluator(tt.expr, tt.variable, tt.v)
		if result != tt.result {
			t.Errorf("%d. exp:\n %+v got:\n %+v", i, tt.result, result)
		}
	}
}
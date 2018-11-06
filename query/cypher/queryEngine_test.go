package cypher_test

import (
	"reflect"
	"testing"

	"github.com/RossMerr/Caudex.Graph/query"
	"github.com/RossMerr/Caudex.Graph/query/cypher"
	"github.com/RossMerr/Caudex.Graph/query/cypher/ast"
	"github.com/RossMerr/Caudex.Graph/query/cypher/parser"
)

type FakeParser struct {
	err error
}

func (p *FakeParser) Parse(s string) (ast.Clauses, error) {
	return nil, p.err
}

func NewFakePaser(err error) parser.Parser {

	return &FakeParser{err: err}
}

type FakeTraversal struct {
}

type FakeParts struct {
}

func (s FakeParts) ToQueryPart(stmt ast.Clauses) ([]*cypher.QueryPart, error) {
	return nil, nil
}

func NewFakeParts() cypher.Parts {
	return &FakeParts{}
}

func NewGraph() *query.Graph {
	return &query.Graph{}
}
func Test_Parser(t *testing.T) {

	tests := []struct {
		e     *cypher.QueryEngine
		p     parser.Parser
		parts cypher.Parts
		path  func(stmt ast.Stmt) ([]cypher.QueryPart, error)
		s     string
		err   string
	}{
		{
			e:     cypher.NewQueryEngine(NewGraph()),
			p:     NewFakePaser(nil),
			parts: NewFakeParts(),
			s:     "str",
		},
		// {
		// 	e:     cypher.NewQueryEngine(NewGraph()),
		// 	p:     NewFakePaser(fmt.Errorf("paser error")),
		// 	parts: NewFakeParts(),
		// 	s:     "str",
		// 	err:   "paser error",
		// },
	}

	for i, tt := range tests {
		tt.e.Parser = tt.p
		tt.e.Parts = tt.parts

		_, err := tt.e.Parse(tt.s)
		if !reflect.DeepEqual(tt.err, errstring(err)) {
			t.Errorf("%d. %q: error mismatch:\n  exp=%s\n  got=%s\n\n", i, tt.s, tt.err, err)
		}
	}
}
package cypher_test

import (
	"reflect"
	"testing"

	"github.com/RossMerr/Caudex.Graph/query"
	"github.com/RossMerr/Caudex.Graph/query/cypher"
	"github.com/RossMerr/Caudex.Graph/query/cypher/ast"
	"github.com/RossMerr/Caudex.Graph/widecolumnstore"
	"github.com/RossMerr/Caudex.Graph/widecolumnstore/storage/memorydb"
)

type unaryMock struct {
}

func (s *unaryMock) Next(i widecolumnstore.Iterator) widecolumnstore.Iterator {
	return i
}

func (s *unaryMock) Op() {

}

type filterMock struct {
	bytes []byte
}

func (s *filterMock) Op() {}

func (s *filterMock) Next(i widecolumnstore.Iterator) widecolumnstore.Iterator {
	return i
}

func TestQueryBuilder_ToPredicateVertexPath(t *testing.T) {
	tests := []struct {
		name    string
		storage widecolumnstore.Storage
		filter  func(storage widecolumnstore.HasPrefix, operator widecolumnstore.Operator, prefix widecolumnstore.Prefix) widecolumnstore.Unary
		patn    *ast.VertexPatn
		last    widecolumnstore.Operator
		want    widecolumnstore.Operator
		err     error
	}{
		{
			name: "Properties filter pattern",
			filter: func(h widecolumnstore.HasPrefix, o widecolumnstore.Operator, p widecolumnstore.Prefix) widecolumnstore.Unary {
				key := widecolumnstore.Key{
					ID: []byte("id"),
				}
				bytes := p(key)
				return &filterMock{
					bytes: bytes,
				}
			},
			storage: func() widecolumnstore.Storage {
				storage, _ := memorydb.NewStorageEngine()
				return storage
			}(),
			patn: &ast.VertexPatn{
				Properties: func() map[string]interface{} {
					prop := make(map[string]interface{}, 0)
					prop["key"] = "value"
					return prop
				}(),
			},
			last: &unaryMock{},
			want: &filterMock{widecolumnstore.NewKey(query.TProperties, &widecolumnstore.Column{[]byte("key"), nil, []byte("id")}).Marshal()},
		},
		{
			name: "No Pattern",
			err:  cypher.ErrNoPattern,
		},
		{
			name: "No Operator",
			patn: &ast.VertexPatn{
				Properties: func() map[string]interface{} {
					prop := make(map[string]interface{}, 0)
					prop["key"] = "value"
					return prop
				}(),
			},
			err: cypher.ErrNoLastOperator,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := cypher.NewQueryBuilder(tt.storage, tt.filter)
			got, err := s.ToPredicateVertexPath(tt.patn, tt.last)

			if err != tt.err {
				t.Errorf("QueryBuilder.ToPredicateVertexPath() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QueryBuilder.ToPredicateVertexPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

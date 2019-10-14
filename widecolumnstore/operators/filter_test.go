package operators

import (
	"testing"

	"github.com/voltable/graph/widecolumnstore"
	"github.com/voltable/graph/widecolumnstore/storage/memorydb"
)

func TestFilter_Next(t *testing.T) {
	type fields struct {
		storage   widecolumnstore.Storage
		operator  widecolumnstore.Operator
		prefix    []byte
		predicate widecolumnstore.Predicate
	}
	tests := []struct {
		name   string
		fields fields
		args   func([]widecolumnstore.KeyValue) widecolumnstore.Iterator
		want   []widecolumnstore.KeyValue
	}{
		{
			name: "Filter",
			fields: func() fields {
				storage, _ := memorydb.NewStorageEngine()
				fields := fields{
					storage:   storage,
					predicate: widecolumnstore.EmptyPredicate,
					prefix:    []byte{},
				}
				return fields
			}(),
			args: func(array []widecolumnstore.KeyValue) widecolumnstore.Iterator {
				l := len(array)
				i := 0
				return func() (widecolumnstore.KeyValue, bool) {
					if i < l {
						old := i
						i++
						return array[old], true
					}

					return widecolumnstore.KeyValue{}, false
				}
			},
			want: func() []widecolumnstore.KeyValue {
				want := []widecolumnstore.KeyValue{}
				want = append(want, widecolumnstore.KeyValue{})
				return want
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := NewFilter(
				tt.fields.predicate,
			)
			got := s.Next(tt.args(tt.want))
			for _, ok := got(); ok; _, ok = got() {
				//	t.Errorf("Filter.Next() = %v", value)
			}
			// if !reflect.DeepEqual(got(), tt.want()) {
			// 	t.Errorf("Filter.Next() = %v, want %v", got, tt.want)
			// }
		})
	}
}

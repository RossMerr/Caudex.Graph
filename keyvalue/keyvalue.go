package keyvalue

import (
	"bytes"

	"github.com/RossMerr/Caudex.Graph/uuid"
)

// NewKeyValue returns a new KeyValue
func newKeyValue(i interface{}, bytes ...[]byte) *KeyValue {
	kv := &KeyValue{
		Value: NewAny(i),
	}
	for _, b := range bytes {
		kv.Key = append(kv.Key, b...)
	}
	return kv
}

// NewKeyValueVertex creates a vertex KeyValue
func NewKeyValueVertex(id *uuid.UUID, label string) *KeyValue {
	return newKeyValue(label, id[:], US, Vertex)
}

// NewKeyValueProperty creates a property KeyValue
func NewKeyValueProperty(id *uuid.UUID, key string, value interface{}) *KeyValue {
	return newKeyValue(value, id[:], US, Properties, US, []byte(key))
}

// NewKeyValueRelationship creates a relationship KeyValue
func NewKeyValueRelationship(from, to *uuid.UUID, relationshipType string) *KeyValue {
	return newKeyValue(to[:], from[:], US, Relationship, US, []byte(relationshipType), US, to[:])
}

// NewKeyValueRelationshipProperty creates a relationship property KeyValue
func NewKeyValueRelationshipProperty(from, to *uuid.UUID, key string, value interface{}) *KeyValue {
	return newKeyValue(value, from[:], US, Relationshipproperties, US, []byte(key), US, to[:])
}

// Transpose

// NewKeyValueVertexTranspose creates a vertex KeyValue
func NewKeyValueVertexTranspose(id *uuid.UUID, label string) *KeyValue {
	return newKeyValue(id[:], Vertex, US, []byte(label), US, id[:])
}

// NewKeyValuePropertyTranspose creates a property KeyValue
func NewKeyValuePropertyTranspose(id *uuid.UUID, key string, value interface{}) *KeyValue {
	return newKeyValue(value, Properties, US, []byte(key), US, id[:])
}

// NewKeyValueRelationshipTranspose creates a relationship KeyValue
func NewKeyValueRelationshipTranspose(from, to *uuid.UUID, relationshipType string) *KeyValue {
	return newKeyValue(to[:], Relationship, US, []byte(relationshipType), US, to[:], US, from[:])
}

// NewKeyValueRelationshipPropertyTranspose creates a relationship property KeyValue
func NewKeyValueRelationshipPropertyTranspose(from, to *uuid.UUID, key string, value interface{}) *KeyValue {
	return newKeyValue(value, Relationshipproperties, US, []byte(key), US, to[:], US, from[:])

}

//NewKeyValue(p, from[:], US, Relationshipproperties, US, []byte(k), US, to[:]))
// UUID looks for the UUID in the KeyValue
func (s *KeyValue) UUID() *uuid.UUID {
	split := bytes.Split(s.Key, US)
	if bytes.Equal(split[1], Vertex) {
		return uuid.SliceToUUID(split[0])
	}

	if bytes.Equal(split[1], Properties) {
		return uuid.SliceToUUID(split[0])
	}

	if bytes.Equal(split[1], Relationship) {
		return uuid.SliceToUUID(split[0])
	}

	if bytes.Equal(split[1], Relationshipproperties) {
		return uuid.SliceToUUID(split[0])
	}

	// Transpose

	if bytes.Equal(split[0], Vertex) {
		return uuid.SliceToUUID(split[2])
	}

	if bytes.Equal(split[0], Properties) {
		return uuid.SliceToUUID(split[2])
	}

	if bytes.Equal(split[0], Relationship) {
		return uuid.SliceToUUID(split[3])
	}

	if bytes.Equal(split[0], Relationshipproperties) {
		return uuid.SliceToUUID(split[3])
	}

	return nil
}

func (s *KeyValue) Interpret(value string) interface{} {
	split := bytes.Split(s.Key, US)
	if len(split) > 1 {
		if bytes.Equal(split[1], Properties) {
			if value == string(split[2]) {
				return s.Value.Unmarshal()
			}
		}
	}

	return nil
}

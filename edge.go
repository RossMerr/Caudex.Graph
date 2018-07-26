package graph

import (
	"strings"

	"github.com/RossMerr/Caudex.Graph/keyvalue"
)

type (
	// An Edge connects two Vertex in a graph.
	Edge struct {
		from UUID
		to   UUID

		isDirected Digraph
		// Weight of a path in a weighted graph
		Weight           float64
		relationshipType string
		properties       map[string]interface{}
	}

	// Edges a array of edges
	Edges []*Edge
)

// NewEdge build a new edge with a id
func NewEdge(from, to *Vertex) *Edge {
	return NewEdgeFromID(from.id, to.id)
}

// NewEdgeFromID creates a edge form the id
func NewEdgeFromID(from, to UUID) *Edge {
	return &Edge{from: from, to: to, properties: make(map[string]interface{})}
}

func (e *Edge) From() string {
	return formatUUID(e.from)
}

func (e *Edge) To() string {
	return formatUUID(e.to)
}

// RelationshipType the type of relationship
func (e *Edge) RelationshipType() string {
	return e.relationshipType
}

// SetRelationshipType the type of relationship
func (e *Edge) SetRelationshipType(label string) {
	e.relationshipType = strings.ToLower(label)
}

// SetProperty set a property to store against this Edge
func (e *Edge) SetProperty(name string, property interface{}) {
	e.properties[name] = property
}

// Property gets a property to store on the Edge
func (e *Edge) Property(name string) interface{} {
	return e.properties[name]
}

func (e *Edge) DeleteProperty(name string) {
	delete(e.properties, name)
}

func (e *Edge) PropertiesCount() int {
	return len(e.properties)
}

func (e *Edge) Properties() map[string]interface{} {
	return e.properties
}

func (a Edges) Len() int           { return len(a) }
func (a Edges) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Edges) Less(i, j int) bool { return a[i].Weight > a[j].Weight }

// MarshalKeyValue marshal a edge into KeyValue
func (e *Edge) MarshalKeyValue() []*keyvalue.KeyValue {
	tt := []*keyvalue.KeyValue{}

	tt = append(tt, keyvalue.NewKeyValue(e.To(), e.from[:], keyvalue.US, keyvalue.Relationship, keyvalue.US, []byte(e.RelationshipType())))

	for k, p := range e.Properties() {
		tt = append(tt, keyvalue.NewKeyValue(p, e.from[:], keyvalue.US, keyvalue.Relationshipproperties, keyvalue.US, []byte(k), keyvalue.US, e.to[:]))
	}

	return tt
}

// MarshalTranspose mashal a Edge into a transposed KeyValue
func (e *Edge) MarshalKeyValueTranspose() []*keyvalue.KeyValue {
	tt := []*keyvalue.KeyValue{}

	tt = append(tt, keyvalue.NewKeyValue(e.To(), keyvalue.Relationship, keyvalue.US, []byte(e.RelationshipType()), keyvalue.US, e.from[:]))

	for k, p := range e.Properties() {
		tt = append(tt, keyvalue.NewKeyValue(p, keyvalue.Relationshipproperties, keyvalue.US, []byte(k), keyvalue.US, e.to[:], keyvalue.US, e.from[:]))
	}

	return tt
}

package ast

import (
	"strings"

	"github.com/RossMerr/Caudex.Graph/query"
	"github.com/RossMerr/Caudex.Graph/vertices"
)

// Patn all pattern nodes implement the Patn interface.
type Patn interface {
	patnNode()
}

type EdgePatn struct {
	Variable string

	Relationship Digraph
	Body         *EdgeBodyStmt

	Vertex *VertexPatn
}

type EdgeBodyStmt struct {
	Variable      string
	Properties    map[string]interface{}
	Type          string
	LengthMinimum uint
	LengthMaximum uint
}

func (*EdgePatn) patnNode() {}

// ToPredicateEdgePath creates a PredicateEdgePath out of the EdgePatn
func (patn *EdgePatn) ToPredicateEdgePath() query.PredicateEdgePath {
	relationshipType := strings.ToLower(patn.Body.Type)
	pvp := query.PredicateEdgePath{PredicateEdge: func(v *vertices.Edge) (string, bool) {
		if relationshipType != v.RelationshipType() {
			return patn.Body.Variable, false
		}

		for key, value := range patn.Body.Properties {
			if v.Property(key) != value {
				return patn.Body.Variable, false
			}
		}

		return patn.Body.Variable, true
	}, Variable: patn.Variable}

	if patn.Body != nil {
		pvp.SetLength(patn.Body.LengthMinimum, patn.Body.LengthMaximum)
	}

	return pvp
}

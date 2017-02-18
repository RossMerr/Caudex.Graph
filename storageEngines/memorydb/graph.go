package memorydb

import "github.com/RossMerr/Caudex.Graph"

type Graph struct {
	vertices map[string]graphs.Vertex
	Options  *graphs.Options
}

func (g *Graph) Close() {

}
func (g *Graph) Query(cypher string) string {
	return ""
}

func BuildGraph(o *graphs.Options) graphs.Graph {
	g := Graph{Options: o}
	return &g
}

func (g *Graph) Command(fn func(*graphs.GraphOperation) error) error {
	op := graphs.CreateGraphOperation(g)
	return fn(op)
}

// Create adds a array of vertices to the persistence
func (g *Graph) Create(c []graphs.Vertex) error {
	for _, v := range c {
		g.vertices[v.ID] = v
	}

	return nil
}

// Delete the array of vertices from the persistence
func (g *Graph) Delete(c []graphs.Vertex) error {
	for _, v := range c {
		delete(g.vertices, v.ID)
	}

	return nil
}

// Find a vertex from the persistence
func (g *Graph) Find(ID string) (*graphs.Vertex, error) {
	v := g.vertices[ID]
	return &v, nil
}

// Update the array of vertices from the persistence
func (g *Graph) Update(c []graphs.Vertex) error {
	g.Delete(c)
	g.Create(c)
	return nil
}
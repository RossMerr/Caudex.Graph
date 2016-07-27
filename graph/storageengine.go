package caudex

import (
	"log"
	"time"

	"github.com/boltdb/bolt"
	"github.com/satori/go.uuid"
)

// Options for the graph
type Options struct {
	Name string
}

// Graph the underlying graph
type Graph struct {
	db      *bolt.DB
	Options *Options
	opend   bool
	ready   bool
	vertexs []Vertex
}

// Open graph
func Open(o *Options) *Graph {
	var st = &Graph{opend: true, Options: o}

	log.Println("Opening " + st.Options.Name)
	// It will be created if it doesn't exist.
	db, err := bolt.Open(st.Options.Name+".db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}

	st.ready = true
	st.db = db

	return st
}

// Close graph
func Close(g *Graph) {
	defer g.db.Close()
}

// Query over the graph using the cypher query language returns JSON
func Query(g *Graph, cypher string) string {
	parse(cypher)
	return "test"
}

// CreateVertex creates a vetex and returns the new vertex.
func CreateVertex(g *Graph) Vertex {
	u1 := uuid.NewV4()
	vertex := Vertex{id: u1.String()}
	g.vertexs = append(g.vertexs, vertex)
	return vertex
}

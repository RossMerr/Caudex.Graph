package memorydb

import (
	"errors"

	"github.com/RossMerr/Caudex.Graph"
	"github.com/RossMerr/Caudex.Graph/storageEngines"
)

func init() {
	storageEngines.RegisterStorageEngine(StorageEngineType, storageEngines.StorageEngineRegistration{
		NewFunc: newStorageEngine,
	})
}

const StorageEngineType = "memory"

var (
	errRecordNotFound = errors.New("Record Not found")
)

type StorageEngine struct {
	vertices map[string]graphs.Vertex
	Options  *graphs.Options
}

func (se *StorageEngine) Close() {

}

func newStorageEngine(o *graphs.Options) (graphs.StorageEngine, error) {
	se := StorageEngine{Options: o, vertices: make(map[string]graphs.Vertex)}
	return &se, nil
}

// Create adds a array of vertices to the persistence
func (se *StorageEngine) Create(c ...*graphs.Vertex) error {
	for _, v := range c {
		se.vertices[v.ID()] = *v
	}

	return nil
}

// Delete the array of vertices from the persistence
func (se *StorageEngine) Delete(c ...*graphs.Vertex) error {
	for _, v := range c {
		delete(se.vertices, v.ID())
	}

	return nil
}

// Find a vertex from the persistence
func (se *StorageEngine) Find(ID string) (*graphs.Vertex, error) {
	if v, ok := se.vertices[ID]; ok {
		return &v, nil
	} else {
		return nil, errRecordNotFound
	}
}

// Update the array of vertices from the persistence
func (se *StorageEngine) Update(c ...*graphs.Vertex) error {
	se.Create(c...)
	return nil
}
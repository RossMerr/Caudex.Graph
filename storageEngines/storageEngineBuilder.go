package storageEngines

import (
	"github.com/RossMerr/Caudex.Graph"
	"github.com/RossMerr/Caudex.Graph/storageEngines/boltdb"
	"github.com/RossMerr/Caudex.Graph/storageEngines/memorydb"
)

// StorageEngineType the backend persistence storage engine
type StorageEngineType int

const (
	// Bolt use Bolt as the storage engine (Default)
	Bolt StorageEngineType = iota
	// Memory use in memory for the storage engine
	Memory StorageEngineType = iota
)

func BuildGraphDefault(o *graphs.Options) graphs.Graph {
	return BuildGraph(Bolt, o)
}

func BuildGraph(e StorageEngineType, o *graphs.Options) graphs.Graph {
	if e == Memory {
		return memorydb.BuildGraph(o)
	}
	return boltdb.BuildGraph(o)
}
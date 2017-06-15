package database

import (
	"time"

	"github.com/mcasviper/yanic/runtime"
)

// Connection interface to use for implementation in e.g. influxdb
type Connection interface {
	// InsertNode stores statistics per node
	InsertNode(node *runtime.Node)

	// InsertGlobals stores global statistics
	InsertGlobals(stats *runtime.GlobalStats, time time.Time)

	// PruneNodes prunes historical per-node data
	PruneNodes(deleteAfter time.Duration)

	// Close closes the database connection
	Close()
}

// Connect function with config to get DB connection interface
type Connect func(config interface{}) (Connection, error)

// Adapters is the list of registered database adapters
var Adapters = map[string]Connect{}

func RegisterAdapter(name string, n Connect) {
	Adapters[name] = n
}

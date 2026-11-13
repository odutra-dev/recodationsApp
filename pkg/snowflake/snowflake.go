package snowflake

import (
	"errors"
	"sync"

	"github.com/bwmarrin/snowflake"
)

var (
	node *snowflake.Node
	once sync.Once
)

// Init initializes the default snowflake generator with the specified node ID.
func Init(nodeID int64) error {
	var err error
	once.Do(func() {
		node, err = snowflake.NewNode(nodeID)
	})
	return err
}

// Generate returns a new unique 64-bit Snowflake ID.
func Generate() int64 {
	if node == nil {
		_ = Init(1)
	}
	return node.Generate().Int64()
}

// IDString returns the Snowflake ID as a string.
func IDString() string {
	if node == nil {
		_ = Init(1)
	}
	return node.Generate().String()
}

// Parse parses a snowflake ID from an int64.
func Parse(id int64) (snowflake.ID, error) {
	if id <= 0 {
		return 0, errors.New("invalid snowflake id")
	}
	return snowflake.ParseInt64(id), nil
}

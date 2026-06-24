package snowflake

import (
	"fmt"
	"sync"

	bwmarrin "github.com/bwmarrin/snowflake"
)

var (
	mu   sync.RWMutex
	node *bwmarrin.Node
)

func Configure(workerID int) error {
	if workerID < 0 || workerID > 1023 {
		return fmt.Errorf("snowflake worker id must be between 0 and 1023")
	}
	next, err := bwmarrin.NewNode(int64(workerID))
	if err != nil {
		return err
	}
	mu.Lock()
	node = next
	mu.Unlock()
	return nil
}

func Next() (uint64, error) {
	mu.RLock()
	current := node
	mu.RUnlock()
	if current == nil {
		if err := Configure(0); err != nil {
			return 0, err
		}
		mu.RLock()
		current = node
		mu.RUnlock()
	}
	return uint64(current.Generate().Int64()), nil
}

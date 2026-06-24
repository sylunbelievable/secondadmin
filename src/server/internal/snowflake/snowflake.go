package snowflake

import (
	"errors"
	"sync"
	"time"
)

const (
	epochSecond  int64  = 1767225600 // 2026-01-01T00:00:00Z
	workerBits   uint   = 10
	sequenceBits uint   = 12
	maxWorkerID  uint64 = 1<<workerBits - 1
	maxSequence  uint64 = 1<<sequenceBits - 1
)

var defaultGenerator = New(0)

type Generator struct {
	mu         sync.Mutex
	workerID   uint64
	lastSecond int64
	sequence   uint64
}

func New(workerID uint64) *Generator {
	return &Generator{workerID: workerID}
}

func Configure(workerID int) error {
	if workerID < 0 || uint64(workerID) > maxWorkerID {
		return errors.New("snowflake worker id must be between 0 and 1023")
	}
	defaultGenerator.mu.Lock()
	defer defaultGenerator.mu.Unlock()
	defaultGenerator.workerID = uint64(workerID)
	defaultGenerator.lastSecond = 0
	defaultGenerator.sequence = 0
	return nil
}

func Next() (uint64, error) {
	return defaultGenerator.Next()
}

func (g *Generator) Next() (uint64, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	now := time.Now().Unix()
	if now < epochSecond {
		return 0, errors.New("clock is before snowflake epoch")
	}
	if now < g.lastSecond {
		return 0, errors.New("clock moved backwards")
	}
	if now == g.lastSecond {
		g.sequence = (g.sequence + 1) & maxSequence
		if g.sequence == 0 {
			for now <= g.lastSecond {
				now = time.Now().Unix()
			}
		}
	} else {
		g.sequence = 0
	}
	g.lastSecond = now

	return uint64(now-epochSecond)<<(workerBits+sequenceBits) |
		g.workerID<<sequenceBits |
		g.sequence, nil
}

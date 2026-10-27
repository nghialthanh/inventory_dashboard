// Package snowflake mints and decodes vehicle addresses.
package snowflake

import (
	"fmt"
	"sync"
	"time"
)

const (
	// Epoch is 2024-01-01T00:00:00Z in Unix milliseconds.
	EpochMs int64 = 1704067200000

	workerBits = 10
	seqBits    = 12

	maxWorker = (1 << workerBits) - 1
	maxSeq    = (1 << seqBits) - 1

	workerShift = seqBits
	timeShift   = workerBits + seqBits
)

// Generator mints Twitter-style snowflake ids for one worker.
type Generator struct {
	workerID int64
	now      func() time.Time

	mu     sync.Mutex
	lastMs int64
	seq    int64
}

// New returns a generator for workerID (0..1023). A nil now uses the wall clock.
func New(workerID int64, now func() time.Time) (*Generator, error) {
	if workerID < 0 || workerID > maxWorker {
		return nil, fmt.Errorf("snowflake worker id %d out of range 0..%d", workerID, maxWorker)
	}
	if now == nil {
		now = time.Now
	}
	return &Generator{workerID: workerID, now: now}, nil
}

// Next mints the next id and the UTC instant encoded in it.
func (g *Generator) Next() (int64, time.Time, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	ms := g.now().UTC().UnixMilli()
	if ms < EpochMs {
		return 0, time.Time{}, fmt.Errorf("snowflake clock before epoch")
	}
	if ms < g.lastMs {
		return 0, time.Time{}, fmt.Errorf("snowflake clock moved backward")
	}
	if ms == g.lastMs {
		if g.seq >= maxSeq {
			for ms <= g.lastMs {
				ms = g.now().UTC().UnixMilli()
			}
			g.seq = 0
		} else {
			g.seq++
		}
	} else {
		g.seq = 0
	}
	g.lastMs = ms
	id := ((ms - EpochMs) << timeShift) | (g.workerID << workerShift) | g.seq
	at := time.UnixMilli(ms).UTC()
	return id, at, nil
}

// Time returns the UTC instant encoded in id.
func Time(id int64) time.Time {
	ms := EpochMs + (id >> timeShift)
	return time.UnixMilli(ms).UTC()
}

// MonthBounds returns [start, end) UTC for the calendar month encoded in id.
func MonthBounds(id int64) (time.Time, time.Time) {
	at := Time(id)
	start := time.Date(at.Year(), at.Month(), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)
	return start, end
}

// Compose builds an id from a UTC instant, worker, and sequence (for seeds and tests).
func Compose(at time.Time, workerID, seq int64) (int64, error) {
	if workerID < 0 || workerID > maxWorker {
		return 0, fmt.Errorf("snowflake worker id %d out of range 0..%d", workerID, maxWorker)
	}
	if seq < 0 || seq > maxSeq {
		return 0, fmt.Errorf("snowflake sequence %d out of range 0..%d", seq, maxSeq)
	}
	ms := at.UTC().UnixMilli()
	if ms < EpochMs {
		return 0, fmt.Errorf("snowflake time before epoch")
	}
	return ((ms - EpochMs) << timeShift) | (workerID << workerShift) | seq, nil
}

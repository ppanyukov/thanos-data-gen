package tsdb

import (
	"github.com/prometheus/prometheus/tsdb/labels"
	"time"
)

// Val is the named value without to write to time series db.
type Val interface {
	Val() float64
	Labels() labels.Labels
}

// ValGenerator is the generator of synthetic values.
type ValGenerator interface {
	Next() Val
}

// Writer is interface to write time series into Prometheus blocks.
type Writer interface {
	// Writes one value, into memory.
	Write(t time.Time, v Val) error

	// Flush writes current block to disk.
	Flush() error

	// Close closes everything
	Close() error
}

// Generator generates synthetic time series using specified profile.
type Generator interface {
	Generate() error
}

// Sample pseudo-code implementation of Generator.
type noddyTsdbGenerator struct {
}

func (g *noddyTsdbGenerator) Generate() error {
	// assume we have these somehow
	var writer Writer
	var generators []ValGenerator

	// these are also given to us at construct time:
	//  - retention: the amount of data to generate time-wise
	//  - sampleInterval: the gap between samples
	//  - flushInterval: the amount of time to go to blocks
	// we want to generate TS data for 30 days
	retention := 30 * 24 * time.Hour

	// all metrics will come in 15s interval, the default?
	sampleInterval := 15 * time.Second

	// flushInterval is the size of the block to write to TSDB
	// NOTE: the blocks need to be precisely 2h, 8h etc.
	//   some extra logic around this needs to be present.
	flushInterval := 2 * time.Hour


	// flush and close on exit
	defer writer.Flush()
	defer writer.Close()


	// write stuff to TSDB from oldest to newest, yes the order matters
	maxt := time.Now()
	mint := maxt.Add(-1 * retention)

	// keep hold of last flush time so we flush at regular intervals
	elapsed := time.Duration(0)

	for t := mint; !t.After(maxt); t = t.Add(sampleInterval) {
		// this is the place where we generate the timestamp
		// since it's the only place that knows anything about time,
		// hence the split between Val and TimedVal
		now := t

		// grab values form generators, timestamp them and shove to
		for _, g := range generators {
			val := g.Next()
			if err := writer.Write(now, val); err != nil {
				return err
			}
		}

		// Flush to disk when written enough data.
		if elapsed >= flushInterval {
			if err := writer.Flush(); err != nil {
				return err
			}
		}
	}

	return nil
}

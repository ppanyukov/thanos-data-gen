package tsdb

import "time"

// generatorT is implementation of Generator.
type generatorT struct {
	retention      time.Duration
	sampleInterval time.Duration
	flushInterval  time.Duration
	valGenerators  []ValGenerator
}

func (g *generatorT) Generate(writer Writer) error {
	// assume we have these somehow
	var generators []ValGenerator

	// these are also given to us at construct time:
	//  - retention: the amount of data to generate time-wise
	//  - sampleInterval: the gap between samples
	//  - flushInterval: the amount of time to go to blocks
	// we want to generate TS data for 30 days
	retention := g.retention

	// all metrics will come in 15s interval, the default?
	sampleInterval := g.sampleInterval

	// flushInterval is the size of the block to write to TSDB
	// NOTE: the blocks need to be precisely 2h, 8h etc.
	//   some extra logic around this needs to be present.
	flushInterval := g.flushInterval

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

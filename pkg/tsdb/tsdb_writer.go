package tsdb

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/prometheus/pkg/timestamp"
	"github.com/prometheus/prometheus/tsdb"
	"github.com/prometheus/prometheus/tsdb/chunkenc"
	"github.com/prometheus/prometheus/tsdb/wal"

	"time"
)

func NewWriter(logger log.Logger, dir string) (Writer, error) {
	res := &writerT{
		logger: logger,
		dir:    dir,
	}

	if err := res.initHeadAndAppender(); err != nil {
		return nil, err
	}

	return res, nil
}

// writerT is implementation of Writer interface
type writerT struct {
	// logger is given to us as arg
	logger log.Logger

	// dir is output directory, given to us as arg
	dir string

	// prometheus specific things, created by us
	head     *tsdb.Head
	appender tsdb.Appender

	// metricCount is incremented internally every time we call Write
	metricCount int64
}

func (w *writerT) Write(t time.Time, v Val) error {
	// Simply write to appender until Flush() is called.
	w.metricCount++

	if _, err := w.appender.Add(v.Labels(), timestamp.FromTime(t), v.Val()); err != nil {
		return errors.Wrap(err, "appender.Add")
	}

	return nil
}

func (w *writerT) Flush() error {
	// Flush should:
	//  - write head to disk
	//  - close head
	//  - open new head and appender
	if err := w.writeHeadToDisk(); err != nil {
		return errors.Wrap(err, "writeHeadToDisk")
	}

	if err := w.head.Close(); err != nil {
		return errors.Wrap(err, "close head")
	}

	if err := w.initHeadAndAppender(); err != nil {
		return errors.Wrap(err, "initHeadAndAppender")
	}

	return nil
}

func (w *writerT) Close() error {
	return w.head.Close()
}

// initHeadAndAppender creates and initialises new head and appender.
func (w *writerT) initHeadAndAppender() error {
	logger := w.logger

	var head *tsdb.Head
	{
		// r and w can be nil as we don't use them
		var r prometheus.Registerer = nil
		var w *wal.WAL = nil

		// chunkRange determines which events are compactable.
		// setting to 1 seems to be the right thing.
		var chunkRange int64 = 1

		h, err := tsdb.NewHead(r, logger, w, chunkRange)
		if err != nil {
			return errors.Wrap(err, "tsdb.NewHead")
		}

		head = h
	}

	w.head = head
	w.appender = head.Appender()
	return nil
}

// writeHeadToDisk commits the appender and writes the head to disk.
func (w *writerT) writeHeadToDisk() error {
	if err := w.appender.Commit(); err != nil {
		return errors.Wrap(err, "appender.Commit")
	}
	seriesCount := w.head.NumSeries()
	mint := timestamp.Time(w.head.MinTime())
	maxt := timestamp.Time(w.head.MaxTime())
	level.Info(w.logger).Log("series_count", seriesCount, "metric_count", w.metricCount, "mint", mint, "maxt", maxt)
	// Step 2. Flush head to disk.
	//
	// copypasta from: github.com/prometheus/prometheus/tsdb/db.go:322
	//
	// Add +1 millisecond to block maxt because block intervals are half-open: [b.MinTime, b.MaxTime).
	// Because of this block intervals are always +1 than the total samples it includes.
	{
		int_mint := timestamp.FromTime(mint)
		int_maxt := timestamp.FromTime(maxt)

		// TODO(ppanyukov): what exactly is "ranges" arg here?
		compactor, err := tsdb.NewLeveledCompactor(context.Background(), nil, w.logger, tsdb.DefaultOptions.BlockRanges, chunkenc.NewPool())
		if err != nil {
			return errors.Wrap(err, "create leveled compactor")
		}

		_, err = compactor.Write(w.dir, w.head, int_mint, int_maxt+1, nil)
		return errors.Wrap(err, "writing WAL")
	}
}

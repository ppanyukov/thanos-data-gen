package tsdb

import (
	"github.com/go-kit/kit/log"
	"github.com/prometheus/prometheus/tsdb/labels"
	"os"
	"testing"
	"time"
)

func Test_Writer_Write(t *testing.T) {
	logger := log.NewLogfmtLogger(os.Stdout)
	w, err := NewWriter(logger, "/Users/philip/zzz-prom-data/zzz")
	if err != nil {
		t.Error(err)
		return
	}

	// generate 4h worth of metrics.
	maxt := time.Now()
	mint := maxt.Add(-8 * time.Hour)
	step := 15 * time.Second

	flushInterval := 2 * time.Hour
	elapsed := time.Duration(0)

	val := &testVal{
		labels:labels.FromStrings("__name__", "foo_metric_total"),
	}

	count := 0
	for now := mint; !now.After(maxt); now = now.Add(step) {
		val.val = float64(count)
		if err := w.Write(now, val); err != nil {
			t.Error(err)
			return
		}

		elapsed += step
		count += 1

		if elapsed >= flushInterval {
			if err := w.Flush(); err != nil {
				t.Error(err)
				return
			}

			elapsed = time.Duration(0)
		}
	}

	if err := w.Flush(); err != nil {
		t.Error(err)
	}

	if err := w.Close(); err != nil {
		t.Error(err)
	}
}

type testVal struct {
	val    float64
	labels labels.Labels
}

func (t *testVal) Val() float64 {
	return t.val
}

func (t *testVal) Labels() labels.Labels {
	return t.labels
}


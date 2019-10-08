package tsdb

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func Test_Everything(t *testing.T) {
	if err := runGenerator(); err != nil {
		t.Errorf("Failure: %v", err)
	}
}

// A really noddy quick test. Not sure it's even a good
// idea to have it.
func runGenerator() error {
	// set this to false to retain the output dir
	// for any manual examination etc.
	removeDir := true

	dir, err := ioutil.TempDir("", "thanos-data-test")
	if err != nil {
		return errors.Wrap(err, "create temp dir")
	}

	// delete temp dir if required
	defer func() {
		if removeDir {
			// ignore errors
			os.RemoveAll(dir)
		} else {
			fmt.Fprintf(os.Stderr, "\n")
			fmt.Fprintf(os.Stderr, "Output directory: %s\n", dir)
			fmt.Fprintf(os.Stderr, "       directory retained, delete manually\n")
		}
	}()

	logger := log.NewLogfmtLogger(os.Stderr)

	// Generate 2 metrics from 3 targets.
	valProviderConfig := ValProviderConfig{
		MetricCount: 2,
		TargetCount: 3,
	}

	valProvider := NewValProvider(valProviderConfig)

	// Custom generator config to make it faster :)
	generatorConfig := DefaultGeneratorConfig(10 * time.Minute)
	generatorConfig.SampleInterval = 1 * time.Second
	generatorConfig.FlushInterval = 2 * time.Minute
	generator := NewGenerator(2 * time.Hour)

	// Create block writer.
	blockWriter, err := NewWriter(logger, dir)
	if err != nil {
		return err
	}

	// Go and hope for the best :)
	if err := generator.Generate(blockWriter, valProvider); err != nil {
		return err
	}

	return nil
}

package randval

import (
	"fmt"
	"os"
	"testing"
)

func Test_randCounterValT_Next(t *testing.T) {
	config := Config{
		MinValue:       10,
		MaxValue:       100,
		MaxChangeValue: 18,
		ChangeRandSeed: 156,
	}

	counter := NewRandCounterVal(config)

	i := 0
	for i < 10 {
		i++
		val := counter.Next()
		fmt.Fprintf(os.Stdout, "Counter %d: %d\n", i, int(val))
	}
}

func Test_randGaugeValT_Next(t *testing.T) {
	config := Config{
		MinValue:       10,
		MaxValue:       100,
		MaxChangeValue: 18,
		ChangeRandSeed: 86755,
	}

	counter := NewRandGaugeVal(config)

	i := 0
	for i < 10 {
		i++
		val := counter.Next()
		fmt.Fprintf(os.Stdout, "Gauge %d: %d\n", i, int(val))
	}
}

// Package randval generates pseudo-random counter and gauge sequences.
package randval

import (
	"math"
	"math/rand"
)

// ValSeq is the interface for getting an infinite sequence of values.
type ValSeq interface {
	Next() float64
}

// Config is the configuration for the value generators.
type Config struct {
	MinValue float64 `yaml:"minValue"`
	MaxValue float64 `yaml:"maxValue"`

	// MaxChangeValue is the maximum change of value as the
	// time goes by. The actual change will be randomised
	// to be in range of [0, +MaxChangeValue] for counters
	// and [-MaxChangeValue, +MaxChangeValue] for gauges.
	MaxChangeValue float64 `yaml:"changeBaseValue"`

	// ChangeRandSeed is the random number generator seed
	// for generating the sequence of changes. Use `0` for
	// the seed based on current time and completely random
	// sequences.
	ChangeRandSeed int64 `yaml:"changeRandSeed"`
}

// DefaultConfig returns a copy of default config.
// The random seed is based on current time.
func DefaultConfig() Config {
	return Config{
		MinValue:       0,
		MaxValue:       10000,
		MaxChangeValue: 10,
		ChangeRandSeed: 0,
	}
}

// NewRandCounterVal creates new random counter sequence.
func NewRandCounterVal(config Config) ValSeq {
	changeRandSource := rand.NewSource(config.ChangeRandSeed)
	changeRand := rand.New(changeRandSource)

	return &randCounterValT{
		config:       config,
		currentValue: config.MinValue,
		changeRand:   changeRand,
	}
}

// NewRandGaugeVal creates new random gauge sequence.
func NewRandGaugeVal(config Config) ValSeq {
	changeRandSource := rand.NewSource(config.ChangeRandSeed)
	changeRand := rand.New(changeRandSource)

	return &randGaugeValT{
		config:       config,
		currentValue: config.MinValue,
		changeRand:   changeRand,
	}
}

// randCounterValT implements counter `ValSeq`: monotonic increase in value.
type randCounterValT struct {
	config       Config
	currentValue float64
	changeRand   *rand.Rand
}

func (c *randCounterValT) Next() float64 {
	// monotonic increase like so:
	//  nextValue = currentValue + (rand baseChange)
	actualChange := c.config.MaxChangeValue * c.changeRand.Float64()
	nextValue := c.currentValue + actualChange

	// reset to min if out of bounds
	if nextValue > c.config.MaxValue || nextValue < c.config.MinValue {
		nextValue = c.config.MinValue
	}

	c.currentValue = nextValue
	return c.currentValue
}

// randCounterValT implements gauge `ValSeq`: value which goes between min and max.
type randGaugeValT struct {
	config       Config
	currentValue float64
	changeRand   *rand.Rand
}

func (c *randGaugeValT) Next() float64 {
	// fluctuate like so:
	//  nextValue = currentValue +/- (baseChange +/- jitter)
	actualChange := c.config.MaxChangeValue * 2 * (c.changeRand.Float64() - 0.5)
	nextValue := c.currentValue + actualChange
	nextValue = math.Min(nextValue, c.config.MaxValue)
	nextValue = math.Max(nextValue, c.config.MinValue)

	c.currentValue = nextValue
	return c.currentValue
}

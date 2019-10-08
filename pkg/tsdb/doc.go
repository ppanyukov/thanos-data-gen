//
// Package tsdb generates synthetic time series in raw Prometheus block format.
//
// It is mainly designed for performance benchmarking of Thanos components.
//
// Quick start:
//
//
//	// Create random value provider
//	valProvider := NewValProvider(valProviderConfig)
//
//	// Create block writer to write to dir
//	blockWriter, _ := NewWriter(logger, dir)
//
//	// Specify how much data to generate and go
// 	generator := NewGenerator(4 * time.Hour)
//	generator.Generate(blockWriter, valProvider)
package tsdb

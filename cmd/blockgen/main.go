package main

import (
	"flag"
	"github.com/pkg/errors"
	"github.com/ppanyukov/thanos-data-gen/pkg/blockgen"
	"log"
	"os"
	"time"
)

// Hacky hacky script to generate TSDB

// The only thing we need to change
var profiles = map[string]runProfile{
	"zzz": {
		name:      "zzz",
		outDir:    os.ExpandEnv("${HOME}/zzz-prom-data/zzz"),
		deleteDir: true,
		genConfig: blockgen.GeneratorConfig{
			StartTime:      time.Date(2019, time.September, 30, 0, 0, 0, 0, time.Local),
			SampleInterval: 15 * time.Second,
			FlushInterval:  2 * time.Hour,
			Retention:      10 * time.Hour,
		},
		valConfig: blockgen.ValProviderConfig{
			MetricCount: 200,
			TargetCount: 100,
		},
	},
}

type runProfile struct {
	name      string
	outDir    string
	deleteDir bool
	genConfig blockgen.GeneratorConfig
	valConfig blockgen.ValProviderConfig
}

var (
	d1d = 24 * time.Hour
	d1w = 7 * d1d
	d1m = 4 * d1w
)

func main() {
	// args is one place for command-line args
	args := struct {
		profileName string
	}{}

	{
		flag.StringVar(&args.profileName, "profile", "zzz", "Profile name to use")
		flag.Parse()

		if args.profileName == "" {
			log.Fatal("profile arg is not specified")
		}
	}

	profile, found := profiles[args.profileName]
	if !found {
		log.Fatalf("Profile with name '%s' not found", args.profileName)
	}

	if err := run(profile); err != nil {
		log.Fatalf("ERROR: %v", err)
	}

	log.Printf("GREAT SUCCESS!")
	log.Printf("Data generated into: %s", profile.outDir)
}

func run(p runProfile) error {
	// remove dir if asked to do so
	if p.deleteDir {
		log.Printf("Deleting outDir %s", p.outDir)
		if err := os.RemoveAll(p.outDir); err != nil {
			return errors.Wrapf(err, "delete dir %s", p.outDir)
		}
	}

	writer, err := blockgen.NewBlockWriter(p.outDir)
	if err != nil {
		return errors.Wrap(err, "blockgen.NewBlockWriter")
	}

	valProvider := blockgen.NewValProvider(p.valConfig)
	generator := blockgen.NewGeneratorWithConfig(p.genConfig)

	log.Printf("Writing to dir: %s", p.outDir)
	return generator.Generate(writer, valProvider)
}

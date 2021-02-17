package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/hermo/registry-check/pkg/analyzer"
	"github.com/hermo/registry-check/pkg/presentation"
	"github.com/hermo/registry-check/pkg/presentation/json"
	"github.com/hermo/registry-check/pkg/presentation/text"
)

func usage() {
	fmt.Printf("registry-check provides a list of registry URLs used in a given NPM or Composer lockfile in text or JSON format.\n\n")
	fmt.Printf("USAGE:\n  registry-check [OPTIONS] LOCKFILE\n\n")
	fmt.Printf("ARGS:\n  <LOCKFILE>   package-lock.json or composer.lock file.\n\n")
	fmt.Println("OPTIONS:")
	flag.PrintDefaults()
	fmt.Println("\nEXAMPLES:")
	fmt.Printf("  registry-check composer.lock                    List registries in composer.lock. Output in text format (default).\n\n")
	fmt.Printf("  registry-check -json package-lock.json          List registries in package-lock.json. Outputs JSON.\n\n")
	fmt.Printf("  registry-check -type npm -json mylock.json      List registries in package-lock.json and force NPM lockfile format\n\n")
}

func main() {
	filename := ""
	lockfileType := ""
	jsonOutput := false

	flag.BoolVar(&jsonOutput, "json", false, "Enable JSON output. (defaults to false)")
	flag.StringVar(&lockfileType, "type", "", "Force lockfile type. Possible values: \"npm\", \"composer\". (defaults to guessing from filename)")
	flag.Usage = usage
	flag.Parse()

	filename = flag.Arg(0)

	if filename == "" {
		flag.Usage()
		os.Exit(2)
	}

	out := presentation.Output{
		Filename:   filename,
		Registries: make([]string, 0),
	}

	var err error

	if lockfileType == "" {
		lockfileType, err = analyzer.DetermineLockfileType(filename)
		if err != nil {
			out.ParsedSuccessfully = false
			out.Error = err
		}
	}

	if lockfileType != "" {
		if result, err := analyzer.AnalyzeLockfile(filename, lockfileType); err != nil {
			out.ParsedSuccessfully = false
			out.Error = err
		} else {
			out.ParsedSuccessfully = true
			out.LockfileType = result.LockfileType
			out.NumPackages = result.NumPackages
			out.Registries = result.Registries
		}
	}

	var presenter presentation.Presenter

	if jsonOutput {
		presenter = json.NewJSONPresenter()
	} else {
		presenter = text.NewTextPresenter()
	}
	presenter.Present(&out)
	if !out.ParsedSuccessfully {
		os.Exit(1)
	}
}

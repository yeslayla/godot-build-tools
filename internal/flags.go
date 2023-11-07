package internal

import (
	"flag"
	"strings"

	"github.com/yeslayla/godot-build-tools/logging"
)

type BuildFlags struct {
	stepsRaw string
	DebugLog bool
}

// Steps returns the steps to run as a slice of strings
func (f *BuildFlags) Steps() []string {
	return strings.Split(f.stepsRaw, ",")
}

// HasStep returns true if the given step is in the list of steps to run
func (f *BuildFlags) HasStep(step string) bool {
	steps := f.Steps()
	for _, s := range steps {
		if s == step {
			return true
		}
	}
	return false
}

// Parse parses the flags
func (f *BuildFlags) Parse() {
	flag.Parse()
}

// NewBuildFlags creates a new BuildFlags instance
func NewBuildFlags(logger logging.Logger) *BuildFlags {
	flags := &BuildFlags{}

	flag.StringVar(&flags.stepsRaw, "steps", "godot-setup", "Comma-separated list of build steps to run")
	flag.BoolVar(&flags.DebugLog, "verbose", false, "Enable debug logging")

	return flags
}

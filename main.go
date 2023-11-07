package main

import (
	"github.com/yeslayla/godot-build-tools/internal"
	"github.com/yeslayla/godot-build-tools/logging"
	"github.com/yeslayla/godot-build-tools/steps"
)

func main() {
	logger := logging.NewLogger(&logging.LoggerOptions{})

	flags := internal.NewBuildFlags(logger)
	flags.Parse()

	if flags.DebugLog {
		logger = logging.NewLogger(&logging.LoggerOptions{
			Debug: true,
		})
	}

	buildConfig := internal.LoadBuildConfig(logger)

	var targetOS internal.TargetOS = internal.CurrentTargetOS()

	if flags.HasStep("godot-setup") {
		steps.GodotSetup(logger, targetOS, buildConfig.Godot.Version, buildConfig.Godot.Release)
	} else {
		logger.Debugf("Skipping godot-setup step")
	}

}

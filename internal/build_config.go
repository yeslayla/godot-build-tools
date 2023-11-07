package internal

import (
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/yeslayla/godot-build-tools/logging"
)

const defaultGodotVersion = "4.1.3"
const defaultGodotRelease = "stable"

type BuildConfig struct {
	Godot BuildConfigGodot `toml:"godot"`
}

type BuildConfigGodot struct {
	Version string `toml:"version"`
	Release string `toml:"release"`
}

func LoadBuildConfig(logger logging.Logger) BuildConfig {
	config := BuildConfig{}

	content, err := ioutil.ReadFile(".godot-build.toml")
	if err != nil {
		if os.IsNotExist(err) {
			logger.Errorf("Build config not found, please run `gbt init`")
			os.Exit(1)
		} else {
			logger.Errorf("Failed to read build config: %s", err)
		}
		return config
	}

	_, err = toml.Decode(string(content), &config)
	if err != nil {
		logger.Errorf("Failed to parse build config: %s", err)
		return config
	}

	if config.Godot.Release == "" {
		logger.Warnf("Godot release not specified, defaulting to %s", defaultGodotRelease)
		config.Godot.Release = defaultGodotRelease
	}

	if config.Godot.Version == "" {
		logger.Warnf("Godot version not specified, defaulting to %s", defaultGodotVersion)
		config.Godot.Version = defaultGodotVersion
	}

	return config
}

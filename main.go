package main

import (
	"os"
	"runtime"

	"github.com/yeslayla/godot-build-tools/internal"
	"github.com/yeslayla/godot-build-tools/logging"
)

func main() {
	logger := logging.NewLogger(&logging.LoggerOptions{})

	var targetOS internal.TargetOS
	switch runtime.GOOS {
	case "linux":
		targetOS = internal.TargetOSLinux
	case "windows":
		targetOS = internal.TargetOSWindows
	case "darwin":
		targetOS = internal.TargetOSMacOS
	}

	GodotSetup(logger, targetOS, "3.3.2", "stable")

}

func GodotSetup(logger logging.Logger, targetOS internal.TargetOS, version string, release string) (string, bool) {
	logger.StartGroup("Godot Setup")
	defer logger.EndGroup()
	downloader := internal.NewDownloader(internal.TargetOSLinux, logger, &internal.DownloaderOptions{})

	logger.Infof("Downloading Godot")
	godotPackage, err := downloader.DownloadGodot(internal.TargetOSLinux, version, release)
	if err != nil {
		logger.Errorf("Failed to download Godot: %s", err)
		return "", false
	}
	defer os.Remove(godotPackage)
	logger.Infof("Godot package: %s", godotPackage)

	logger.Infof("Installing Godot")
	godotBin, err := downloader.InstallGodot(godotPackage, internal.TargetOSLinux, version, release)
	if err != nil {
		logger.Errorf("Failed to install Godot: %s", err)
		return "", false
	}
	logger.Infof("Godot binary: %s", godotBin)

	return godotBin, true
}

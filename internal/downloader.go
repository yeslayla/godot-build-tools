package internal

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/yeslayla/godot-build-tools/logging"
	"github.com/yeslayla/godot-build-tools/utils"
)

type DownloaderOptions struct {
	DownloadRepositoryURL string
	BinDir                string
}

type Downloader struct {
	downloadRepositoryURL string
	bin                   string

	logger logging.Logger
}

func NewDownloader(targetOS TargetOS, logger logging.Logger, options *DownloaderOptions) *Downloader {
	var url string = options.DownloadRepositoryURL
	if url == "" {
		url = "https://downloads.tuxfamily.org/godotengine/"
	}
	var binDir string = options.BinDir
	if binDir == "" {
		switch targetOS {
		case TargetOSLinux:
			home, _ := os.UserHomeDir()
			binDir = filepath.Join(home, "/.local/bin")
		case TargetOSWindows:
			binDir = "C:\\Program Files (x86)\\Godot"
		case TargetOSMacOS:
			binDir = "/Applications/Godot"
		}
	}

	return &Downloader{
		downloadRepositoryURL: url,
		bin:                   binDir,
		logger:                logger,
	}
}

func getRemoteFileFormat(targetOS TargetOS, version string) string {
	switch targetOS {
	case TargetOSLinux:
		if version[0] == '3' {
			return "Godot_v%s-%s_x11.64.zip"
		}
		return "Godot_v%s-%s_linux.x86_64.zip"
	case TargetOSWindows:
		return "Godot_v%s-%s_win64.exe.zip"
	case TargetOSMacOS:
		return "Godot_v%s-%s_macos.universal.zip"
	}

	return ""
}

func (d *Downloader) DownloadGodot(targetOS TargetOS, version string, release string) (string, error) {

	var fileName string = fmt.Sprintf(getRemoteFileFormat(targetOS, version), version, release)

	tempDir, _ := os.MkdirTemp("", "godot-build-tools")
	outFile := filepath.Join(tempDir, fileName)
	out, err := os.Create(outFile)
	if err != nil {
		return "", fmt.Errorf("failed to create output file: %s", err)
	}
	defer out.Close()

	downloadURL, err := url.Parse(d.downloadRepositoryURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse download repository URL: %s", err)
	}

	downloadURL.Path = path.Join(downloadURL.Path, version)
	if release != "stable" {
		downloadURL.Path = path.Join(downloadURL.Path, release)
	}

	downloadURL.Path = path.Join(downloadURL.Path, fileName)
	d.logger.Debugf("Download URL: %s", downloadURL.String())

	resp, err := http.Get(downloadURL.String())
	if err != nil {
		return "", fmt.Errorf("failed to download Godot: %s", err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to write Godot package to output file: %s", err)
	}

	return outFile, nil
}

func (d *Downloader) UnzipGodot(targetOS TargetOS, godotPackage string) (string, error) {
	files, err := utils.Unzip(godotPackage)
	if err != nil {
		return "", fmt.Errorf("failed to unzip Godot package: %s", err)
	}

	// Look for godot binary
	for _, file := range files {
		switch targetOS {
		case TargetOSLinux:
			if path.Ext(file) == ".x86_64" || path.Ext(file) == ".64" {
				return file, nil
			}
		case TargetOSWindows:
			if path.Ext(file) == ".exe" {
				return file, nil
			}
		case TargetOSMacOS:
			if path.Ext(file) == ".universal" {
				return file, nil
			}
		}
	}

	return "", fmt.Errorf("failed to find godot binary in Godot package")
}

func (d *Downloader) InstallGodot(godotPackage string, targetOS TargetOS, version string, release string) (string, error) {

	// Unzip package
	godotUnzipBinPath, err := d.UnzipGodot(targetOS, godotPackage)
	if err != nil {
		return "", fmt.Errorf("failed to unzip Godot package: %s", err)
	}

	godotBin := path.Base(godotUnzipBinPath)
	godotBinPath := filepath.Join(d.bin, godotBin)

	// Copy Godot binary to bin directory
	data, err := ioutil.ReadFile(godotUnzipBinPath)
	if err != nil {
		return "", fmt.Errorf("failed to read Godot binary: %s", err)
	}
	err = ioutil.WriteFile(godotBinPath, data, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to write Godot binary: %s", err)
	}

	_ = os.Remove(godotUnzipBinPath)

	return godotBinPath, nil
}

package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Unzip unzips a zip archive and returns the paths of the unzipped files.
func Unzip(archivePath string) ([]string, error) {
	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open Godot package: %s", err)
	}
	defer reader.Close()

	destDir, err := filepath.Abs(path.Dir(archivePath))
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path of Godot package: %s", err)
	}

	var unzippedFiles []string = make([]string, 0)
	for _, file := range reader.File {
		filePath := filepath.Join(destDir, file.Name)
		if !strings.HasPrefix(filePath, filepath.Clean(destDir)+string(os.PathSeparator)) {
			return unzippedFiles, fmt.Errorf("%s: illegal file path", filePath)
		}

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
				return unzippedFiles, fmt.Errorf("failed to create directory: %s", err)
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return unzippedFiles, fmt.Errorf("failed to create directory: %s", err)
		}

		destFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return unzippedFiles, fmt.Errorf("failed to open file: %s", err)
		}
		defer destFile.Close()

		zippedFile, err := file.Open()
		if err != nil {
			return unzippedFiles, fmt.Errorf("failed to open zipped file: %s", err)
		}
		defer zippedFile.Close()

		if _, err := io.Copy(destFile, zippedFile); err != nil {
			return unzippedFiles, fmt.Errorf("failed to copy file: %s", err)
		}

		unzippedFiles = append(unzippedFiles, filePath)
	}

	return unzippedFiles, nil
}

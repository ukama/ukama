package utils

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// CopyDir recursively copies a directory and its contents to a destination directory.
func CopyDir(srcDir, destDir string) error {
	err := filepath.Walk(srcDir, func(srcPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Create the corresponding directory in the destination.
		destPath := filepath.Join(destDir, srcPath[len(srcDir):])
		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		// Copy the file to the destination.
		srcFile, err := os.Open(srcPath)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		destFile, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer destFile.Close()

		_, err = io.Copy(destFile, srcFile)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func RemoveDir(path string) error {
	return os.RemoveAll(path)
}

func CreateDir(path string, perm fs.FileMode) error {
	return os.MkdirAll(path, perm)
}

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package utils

import (
	"io"
	"io/fs"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func RandomDirName() string {
	return StringWithCharset(10, charset)
}

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

		defer func() {
			err := srcFile.Close()
			if err != nil {
				log.Warnf("failed to close source file: %v", err)
			}
		}()

		destFile, err := os.Create(destPath)
		if err != nil {
			return err
		}

		defer func() {
			err := destFile.Close()
			if err != nil {
				log.Warnf("failed to close destination file: %v", err)
			}
		}()

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

func FindDir(name string, root string) (*string, error) {
	var foundDirectory string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() && info.Name() == name {
			// If a directory with the target name is found, store its path and stop walking.
			foundDirectory = path
			return filepath.SkipDir
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &foundDirectory, nil
}

func GetFiles(path string) ([]string, error) {
	var files []string

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			// If it's not a directory, it's a file, so add its path to the slice.
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		log.Println(err)
		return nil, err
	}
	return files, nil
}

// Copyright 2013 cae authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package zip

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

func extractFile(f *zip.File, destPath string) error {
	// Create diretory before create file
	os.MkdirAll(path.Join(destPath, path.Dir(f.Name)), os.ModePerm)

	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	fw, _ := os.Create(path.Join(destPath, f.Name))
	if err != nil {
		return err
	}
	_, err = io.Copy(fw, rc)
	if err != nil {
		return err
	}
	return nil
}

func isEntry(name string, entries []string) bool {
	for _, e := range entries {
		if e == name {
			return true
		}
	}
	return false
}

// ExtractTo extracts the complete archive or the given files to the specified destination.
// Call Flush() to apply changes before this.
func (z *ZipArchive) ExtractTo(destPath string, entries ...string) (err error) {
	isHasEntry := len(entries) > 0
	fmt.Println("Unzipping " + z.FileName + "...")
	os.MkdirAll(destPath, os.ModePerm)

	for _, f := range z.File {
		fmt.Println("Unzipping file..." + f.Name)
		// Directory.
		if strings.HasSuffix(f.Name, "/") {
			os.MkdirAll(path.Join(destPath, f.Name), os.ModePerm)
			continue
		}

		// File.
		if isHasEntry {
			if isEntry(f.Name, entries) {
				err = extractFile(f, destPath)
			}
			continue
		} else {
			err = extractFile(f, destPath)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// Flush saves changes to original zip file if any.
func (z *ZipArchive) Flush() error {
	if !z.isHasChanged {
		return nil
	}

	// Extract to tmp path and pack back.

	return z.Open(z.FileName, z.Flag, z.Permission)
}

// Close opened or created archive and save changes.
func (z *ZipArchive) Close() (err error) {
	if z.ReadCloser != nil {
		err = z.ReadCloser.Close()
		if err != nil {
			return err
		}
		z.ReadCloser = nil
	}

	return z.Flush()
}

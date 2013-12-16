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
	"path/filepath"
	"strings"
)

var Verbose = true

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
		// Directory.
		if strings.HasSuffix(f.Name, "/") {
			if isHasEntry {
				if isEntry(f.Name, entries) {
					fmt.Println("Unzipping file..." + f.Name)
					os.MkdirAll(path.Join(destPath, f.Name), os.ModePerm)
				}
				continue
			}
			fmt.Println("Unzipping file..." + f.Name)
			os.MkdirAll(path.Join(destPath, f.Name), os.ModePerm)
			continue
		}

		// File.
		if isHasEntry {
			if isEntry(f.Name, entries) {
				fmt.Println("Unzipping file..." + f.Name)
				err = extractFile(f, destPath)
			}
			continue
		} else {
			fmt.Println("Unzipping file..." + f.Name)
			err = extractFile(f, destPath)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (z *ZipArchive) extractFile(f *File) error {
	for _, zf := range z.ReadCloser.File {
		if f.Name == zf.Name {
			extractFile(zf, f.absPath)
			return nil
		}
	}

	return copy(f.Name, f.absPath)
}

// Flush saves changes to original zip file if any.
func (z *ZipArchive) Flush() error {
	if !z.isHasChanged || z.ReadCloser == nil {
		return nil
	}

	// Extract to tmp path and pack back.
	tmpPath := path.Join(os.TempDir(), "cae", path.Base(z.FileName))
	os.RemoveAll(tmpPath)
	defer os.RemoveAll(tmpPath)

	for _, f := range z.files {
		if strings.HasSuffix(f.Name, "/") {
			os.MkdirAll(path.Join(tmpPath, f.Name), os.ModePerm)
			continue
		}

		f.absPath = path.Join(tmpPath, f.Name)
		err := z.extractFile(f)
		if err != nil {
			return err
		}
	}

	err := PackTo(tmpPath, z.FileName)
	if err != nil {
		return err
	}
	return z.Open(z.FileName, os.O_RDWR|os.O_TRUNC, z.Permission)
}

func packDir(srcPath string, recPath string, zw *zip.Writer, fn func(fullName string, fi os.FileInfo) error) error {
	dir, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer dir.Close()

	// Get file info slice
	fis, err := dir.Readdir(0)
	if err != nil {
		return err
	}

	for _, fi := range fis {
		if globalFilter(fi.Name()) {
			continue
		}
		// Append path
		curPath := srcPath + "/" + fi.Name()
		tmpRecPath := filepath.Join(recPath, fi.Name())
		err = fn(curPath, fi)
		if err != nil {
			return err
		}

		// Check it is directory or file
		if fi.IsDir() {
			err = packFile(srcPath, tmpRecPath, zw, fi)
			if err != nil {
				return err
			}

			err = packDir(curPath, tmpRecPath, zw, fn)
			if err != nil {
				return err
			}
		} else {
			err = packFile(curPath, tmpRecPath, zw, fi)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func packFile(srcFile string, recPath string, zw *zip.Writer, fi os.FileInfo) error {
	if fi.IsDir() {
		// Create zip header
		fh := new(zip.FileHeader)
		fh.Name = recPath + "/"
		fh.UncompressedSize = 0

		_, err := zw.CreateHeader(fh)
		if err != nil {
			return err
		}
	} else {
		// Create zip header
		fh := new(zip.FileHeader)
		fh.Name = recPath
		fh.UncompressedSize = uint32(fi.Size())
		fw, err := zw.CreateHeader(fh)
		if err != nil {
			return err
		}

		f, err := os.Open(srcFile)
		if err != nil {
			return err
		}
		_, err = io.Copy(fw, f)
		if err != nil {
			return err
		}
	}
	return nil
}

func packTo(srcPath, destPath string, fn func(fullName string, fi os.FileInfo) error, includeDir bool) error {
	fw, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer fw.Close()

	zw := zip.NewWriter(fw)
	defer zw.Close()

	f, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	fi, err := f.Stat()
	if err != nil {
		return err
	}

	basePath := path.Base(srcPath)

	if fi.IsDir() {
		if includeDir {
			packFile(srcPath, basePath, zw, fi)
		} else {
			basePath = ""
		}
		return packDir(srcPath, basePath, zw, fn)
	}

	return packFile(srcPath, basePath, zw, fi)
}

var defaultPackFunc = func(fullName string, fi os.FileInfo) error {
	if !Verbose {
		return nil
	}

	if fi.IsDir() {
		fmt.Printf("Adding dir...%s\n", fullName)
	} else {
		fmt.Printf("Adding file...%s\n", fullName)
	}

	return nil
}

// PackTo packs the complete archive to the specified destination.
// It accepts a function as a middleware for custom-operations.
func PackToFunc(srcPath, destPath string, fn func(fullName string, fi os.FileInfo) error, includeDir ...bool) error {
	isIncludeDir := false
	if len(includeDir) > 0 && includeDir[0] {
		isIncludeDir = true
	}

	return packTo(srcPath, destPath, fn, isIncludeDir)
}

// PackTo packs the complete archive to the specified destination.
// Call Flush() will automatically call this in the end.
func PackTo(srcPath, destPath string, includeDir ...bool) error {
	return PackToFunc(srcPath, destPath, defaultPackFunc, includeDir...)
}

// Close opened or created archive and save changes.
func (z *ZipArchive) Close() (err error) {
	err = z.Flush()
	if err != nil {
		return err
	}

	if z.ReadCloser != nil {
		err = z.ReadCloser.Close()
		if err != nil {
			return err
		}
		z.ReadCloser = nil
	}
	return nil
}

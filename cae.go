// Copyright 2013-2014 Unknown
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

// Package cae implements PHP-like Compression and Archive Extensions.
package cae

import (
	"io"
	"os"
	"strings"
)

func HasPrefix(name string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(name, prefix) {
			return true
		}
	}
	return false
}

func IsEntry(name string, entries []string) bool {
	for _, e := range entries {
		if e == name {
			return true
		}
	}
	return false
}

func GlobalFilter(name string) bool {
	if strings.Contains(name, ".DS_Store") {
		return true
	}
	return false
}

// Copy copies file from source to target path.
// It returns false and error when error occurs in underlying functions.
func Copy(destPath, srcPath string) error {

	si, err := os.Lstat(srcPath)
	if err != nil {
		return err
	}

	// Symbolic link.
	if si.Mode()&os.ModeSymlink != 0 {
		target, err := os.Readlink(srcPath)
		if err != nil {
			return err
		}
		return os.Symlink(target, destPath)
	}

	sf, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer sf.Close()

	df, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer df.Close()

	// buffer reader, do chunk transfer
	buf := make([]byte, 1024)
	for {
		// read a chunk
		n, err := sf.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		// write a chunk
		if _, err := df.Write(buf[:n]); err != nil {
			return err
		}
	}

	return os.Chmod(destPath, si.Mode())
}

func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

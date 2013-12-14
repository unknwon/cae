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

// Package zip enables you to transparently read or write ZIP compressed archives and the files inside them.
package zip

import (
	"archive/zip"
	"os"
)

// File represents a file in archive.
type File struct {
	*zip.FileHeader
	oldName    string
	oldComment string
	absPath    string
}

// ZipArchive represents a file archive, compressed with Zip.
type ZipArchive struct {
	*zip.ReadCloser
	FileName   string
	Comment    string
	NumFiles   int
	Flag       int
	Permission os.FileMode

	files        []*File
	isHasChanged bool
}

// Create creates the named zip file, truncating
// it if it already exists. If successful, methods on the returned
// ZipArchive can be used for I/O; the associated file descriptor has mode
// O_RDWR.
// If there is an error, it will be of type *PathError.
func Create(fileName string) (zip *ZipArchive, err error) {
	return OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
}

// Open opens the named zip file for reading.  If successful, methods on
// the returned ZipArchive can be used for reading; the associated file
// descriptor has mode O_RDONLY.
// If there is an error, it will be of type *PathError.
func Open(fileName string) (zip *ZipArchive, err error) {
	return OpenFile(fileName, os.O_RDONLY, 0)
}

// OpenFile is the generalized open call; most users will use Open
// instead. It opens the named zip file with specified flag
// (O_RDONLY etc.) if applicable. If successful,
// methods on the returned ZipArchive can be used for I/O.
// If there is an error, it will be of type *PathError.
func OpenFile(fileName string, flag int, perm os.FileMode) (zip *ZipArchive, err error) {
	zip = &ZipArchive{}
	err = zip.Open(fileName, flag, perm)
	return zip, err
}

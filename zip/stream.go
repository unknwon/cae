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

package zip

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

// StreamArchive represents a streaming archive.
type StreamArchive struct {
	*zip.Writer
}

// NewStreamArachive returns a new straming archive with given io.Writer.
// It's caller's responsibility to close io.Writer after operation.
func NewStreamArachive(writer io.Writer) *StreamArchive {
	return &StreamArchive{zip.NewWriter(writer)}
}

// StreamFile strams a file entry to StreamArchive.
func (sr *StreamArchive) StreamFile(relPath string, fi os.FileInfo, data []byte) (err error) {
	if fi.IsDir() {
		fh, err := zip.FileInfoHeader(fi)
		if err != nil {
			return err
		}
		fh.Name = relPath + "/"

		_, err = sr.Writer.CreateHeader(fh)
	} else {
		fh, err := zip.FileInfoHeader(fi)
		if err != nil {
			return err
		}
		fh.Name = filepath.Join(relPath, fi.Name())
		fh.Method = zip.Deflate

		fw, err := sr.Writer.CreateHeader(fh)
		if err != nil {
			return err
		}
		_, err = fw.Write(data)
	}
	return err
}

// StreamReader streams data from io.Reader to StreamArchive.
func (sr *StreamArchive) StreamReader(relPath string, fi os.FileInfo, reader io.Reader) (err error) {
	fh, err := zip.FileInfoHeader(fi)
	if err != nil {
		return err
	}
	fh.Name = filepath.Join(relPath, fi.Name())

	fw, err := sr.Writer.CreateHeader(fh)
	if err != nil {
		return err
	}
	_, err = io.Copy(fw, reader)
	return err
}

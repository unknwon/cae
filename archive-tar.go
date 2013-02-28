// Copyright 2013 The Author - Unknown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
)

func handleError(e error) {
	if e != nil {
		log.Fatal(e)
	}

}

// Deal with files
func HandleFile(srcFile string, tw *tar.Writer, fi os.FileInfo) {
	if fi.IsDir() {
		// Create tar header
		hdr := new(tar.Header)
		hdr.Name = srcFile + "/"

		// Write hander
		err := tw.WriteHeader(hdr)
		handleError(err)
	} else {
		// File reader
		fr, err := os.Open(srcFile)
		handleError(err)
		defer fr.Close()

		// Create tar header
		hdr := new(tar.Header)
		hdr.Name = srcFile
		hdr.Size = fi.Size()
		hdr.Mode = int64(fi.Mode())
		hdr.ModTime = fi.ModTime()

		// Write hander
		err = tw.WriteHeader(hdr)
		handleError(err)

		// Write file data
		_, err = io.Copy(tw, fr)
		handleError(err)
	}
}

// Deal with directories
// if find files, handle them with HandleFile
func HandleDir(srcDirPath string, tw *tar.Writer) {
	// Open source diretory
	dir, err := os.Open(srcDirPath)
	handleError(err)
	defer dir.Close()

	// Get file info slice
	fis, err := dir.Readdir(0)
	handleError(err)
	for _, fi := range fis {
		// Append path
		curPath := srcDirPath + "/" + fi.Name()
		// Check it is directory or file
		if fi.IsDir() {
			// Directory
			HandleDir(curPath, tw)
		} else {
			// File
			fmt.Printf("Adding file...%s\n", curPath)
		}

		HandleFile(curPath, tw, fi)
	}
}

// Gzip and tar from source directory to destination file
func TarGz(srcDirPath string, destFilePath string) {
	// File writer
	fw, err := os.Create(destFilePath)
	handleError(err)
	defer fw.Close()

	// Gzip writer
	gw := gzip.NewWriter(fw)
	defer gw.Close()

	// Tar writer
	tw := tar.NewWriter(gw)
	defer tw.Close()

	// handle source directory
	HandleDir(srcDirPath, tw)
}

func main() {
	targetFilePath := "/home/unknown/Applications/Go/src/test.tar.gz"
	inputDirPath := "/home/unknown/Applications/Go/src/test"
	TarGz(inputDirPath, targetFilePath)
	fmt.Println("Finish!")
}

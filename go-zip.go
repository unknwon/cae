// Copyright 2013 The Author - Unknown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"archive/zip"
	"fmt"
	"log"
	"os"
)

func main() {
	Zip("/Users/Joe/Applications/Go/src/go-zip.go", "ok.zip")
}

// Zip from source directory or file to destination file
// you need check file exist before you call this function
func Zip(srcDirPath string, destFilePath string) {
	fw, err := os.Create(destFilePath)
	defer fw.Close()

	// Zip writer
	zw := zip.NewWriter(fw)
	defer zw.Close()

	// Check if it's a file or a directory
	f, err := os.Open(srcDirPath)
	handleError(err)
	fi, err := f.Stat()
	handleError(err)
	switch {
	case fi.IsDir():
		// handle source directory
		fmt.Println("Cerating zip from directory...")
		//tarGzDir(srcDirPath, path.Base(srcDirPath), tw)
		fallthrough
	case 0 == (os.ModeType & fi.Mode()):
		// handle file directly
		fmt.Println("Cerating zip from " + fi.Name() + "...")
		zipFile(srcDirPath, fi.Name(), zw, f, fi)
	}
	fmt.Println("Well done!")
}

// Deal with file
func zipFile(srcFile string, recPath string, zw *zip.Writer, f *os.File, fi os.FileInfo) {
	if fi.IsDir() {

	} else {
		// File header
		fh := new(zip.FileHeader)
		fh.Name = fi.Name()
		fh.UncompressedSize = uint32(fi.Size())
		fw, err := zw.CreateHeader(fh)
		handleError(err)
		buf := make([]byte, fi.Size())
		_, err = f.Read(buf)
		handleError(err)
		_, err = fw.Write(buf)
		handleError(err)
	}
}

func handleError(e error) {
	if e != nil {
		log.Fatal(e)
	}

}

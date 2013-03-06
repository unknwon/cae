// Copyright 2013 The Author - Unknown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"archive/zip"
	"fmt"
	"log"
	"os"
	"path"
)

func main() {
	Zip("/Users/Joe/Applications/Go/src/demos", "ok.zip")
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
		zipDir(srcDirPath, path.Base(srcDirPath), zw)
		fallthrough
	case 0 == (os.ModeType & fi.Mode()):
		// handle file directly
		fmt.Println("Cerating zip from " + fi.Name() + "...")
		zipFile(srcDirPath, fi.Name(), zw, fi)
	}
	fmt.Println("Well done!")
}

// Deal with directories
// if find files, handle them with zipFile
// Every recurrence append the base path to the recPath
// recPath is the path inside of zip
func zipDir(srcDirPath string, recPath string, zw *zip.Writer) {
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
			// (Directory won't add unitl all subfiles are added)
			fmt.Printf("Adding path...%s\n", curPath)
			zipDir(curPath, recPath+"/"+fi.Name(), zw)
		} else {
			// File
			fmt.Printf("Adding file...%s\n", curPath)
		}

		zipFile(curPath, recPath+"/"+fi.Name(), zw, fi)
	}
}

// Deal with file
func zipFile(srcFile string, recPath string, zw *zip.Writer, fi os.FileInfo) {
	if fi.IsDir() {

	} else {
		// File header
		fh := new(zip.FileHeader)
		fh.Name = fi.Name()
		fh.UncompressedSize = uint32(fi.Size())
		fw, err := zw.CreateHeader(fh)
		handleError(err)
		buf := make([]byte, fi.Size())
		f, err := os.Open(srcFile)
		handleError(err)
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

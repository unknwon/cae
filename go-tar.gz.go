// Copyright 2013 The Author - Unknown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Help linls:
// https://codereview.appspot.com/7305072
package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path"
)

// Copy from http://golang.org/misc/dist/bindist.go
// Mode constants from the tar spec.
const (
	c_ISDIR  = 040000
	c_ISFIFO = 010000
	c_ISREG  = 0100000
	c_ISLNK  = 0120000
	c_ISBLK  = 060000
	c_ISCHR  = 020000
	c_ISSOCK = 0140000
)

func main() {
	targetFilePath := "test.tar.gz"
	srcDirPath := "test"
	TarGz(srcDirPath, targetFilePath)
	UnTarGz(targetFilePath, srcDirPath+"_temp")
	fmt.Println("Finish!")
}

// Gzip and tar from source directory to destination file
// you need check file exist before you call this function
func TarGz(srcDirPath string, destFilePath string) {
	fmt.Println("Cerating tar.gz...")
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
	TarGzDir(srcDirPath, path.Base(srcDirPath), tw)
}

// Deal with directories
// if find files, handle them with HandleFile
// Every recurrence append the base path to the recPath
// recPath is the path inside of tar.gz
func TarGzDir(srcDirPath string, recPath string, tw *tar.Writer) {
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
			TarGzDir(curPath, recPath+"/"+fi.Name(), tw)
		} else {
			// File
			fmt.Printf("Adding file...%s\n", curPath)
		}

		TarGzFile(curPath, recPath+"/"+fi.Name(), tw, fi)
	}
}

// Deal with files
func TarGzFile(srcFile string, recPath string, tw *tar.Writer, fi os.FileInfo) {
	if fi.IsDir() {
		// Create tar header
		hdr := new(tar.Header)
		hdr.Name = recPath + "/"
		hdr.Typeflag = tar.TypeDir
		hdr.Size = 0
		hdr.Mode = 0755 | c_ISDIR
		hdr.ModTime = fi.ModTime()
		// if last character of header name is '/' it also can be directory

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
		hdr.Name = recPath
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

// Ungzip and untar from source file to destination directory
// you need check file exist before you call this function
func UnTarGz(srcFilePath string, destDirPath string) {
	fmt.Println("UnTarGzing " + srcFilePath + "...")
	// Create destination directory
	os.Mkdir(destDirPath, os.ModePerm)

	fr, err := os.Open(srcFilePath)
	handleError(err)
	defer fr.Close()

	// Gzip reader
	gr, err := gzip.NewReader(fr)

	// Tar reader
	tr := tar.NewReader(gr)

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			// End of tar archive
			break
		}
		//handleError(err)
		fmt.Println("UnTarGzing file..." + hdr.Name)

		// If header name ends with '/' means it is directory ,
		// and doesn't need extract
		c := []byte(hdr.Name)
		if c[len(hdr.Name)-1] != '/' {
			// Get files from archive
			os.MkdirAll(destDirPath+"/"+path.Dir(hdr.Name), os.ModePerm)
			// Write data to file
			fw, _ := os.Create(destDirPath + "/" + hdr.Name)
			handleError(err)
			_, err = io.Copy(fw, tr)
			handleError(err)
		}
	}
}

func handleError(e error) {
	if e != nil {
		log.Fatal(e)
	}

}

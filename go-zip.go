// Copyright 2013 The Author - Unknown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gocompresser

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path"
)

// main functions shows how to Zip a directory/file and
// UnZip a file
//func main() {
//	targetFilePath := "demos.zip"
//	srcDirPath := "../demos"
//	Zip(srcDirPath, targetFilePath)
//	UnZip(targetFilePath, "./")
//}

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
	if fi.IsDir() {
		// handle source directory
		fmt.Println("Cerating zip from directory...")
		zipDir(srcDirPath, path.Base(srcDirPath), zw)
	} else {
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
		// Create zip header
		fh := new(zip.FileHeader)
		fh.Name = recPath + "/"
		fh.UncompressedSize = 0

		_, err := zw.CreateHeader(fh)
		handleError(err)
	} else {
		// Create zip header
		fh := new(zip.FileHeader)
		fh.Name = recPath
		fh.UncompressedSize = uint32(fi.Size())
		fw, err := zw.CreateHeader(fh)
		handleError(err)

		// Read file data
		buf := make([]byte, fi.Size())
		f, err := os.Open(srcFile)
		handleError(err)
		_, err = f.Read(buf)
		handleError(err)

		// Write file data to zip
		_, err = fw.Write(buf)
		handleError(err)
	}
}

// Unzip from source file to destination directory
// you need check file exist before you call this function
func UnZip(srcFilePath string, destDirPath string) {
	fmt.Println("Unzipping " + srcFilePath + "...")
	// Create destination directory
	os.Mkdir(destDirPath, os.ModePerm)

	// Open a zip archive for reading
	r, err := zip.OpenReader(srcFilePath)
	handleError(err)
	defer r.Close()

	// Iterate through the files in the archive
	for _, f := range r.File {
		fmt.Println("Unzipping file..." + f.FileInfo().Name())
		// Get files from archive
		rc, err := f.Open()
		handleError(err)
		// Create diretory before create file
		os.MkdirAll(destDirPath+"/"+path.Dir(f.FileInfo().Name()), os.ModePerm)
		// Write data to file
		fw, _ := os.Create(destDirPath + "/" + f.FileInfo().Name())
		handleError(err)
		_, err = io.Copy(fw, rc)
		handleError(err)
	}
	fmt.Println("Well done!")
}

func handleError(e error) {
	if e != nil {
		log.Fatal(e)
	}

}

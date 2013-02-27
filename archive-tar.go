// Copyright 2013 The Author - Unknown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"archive/tar"
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

func main() {
	// Open file
	fr, err := os.Open("small.txt")
	handleError(err)
	defer fr.Close()

	// Get file info
	fi, err := fr.Stat()
	handleError(err)

	// Create tar header
	hdr := new(tar.Header)
	hdr.Name = "small.txt"
	hdr.Size = fi.Size()
	hdr.Mode = int64(fi.Mode())
	hdr.ModTime = fi.ModTime()

	// Create tar package
	fw, err := os.Create("gnu.tar")
	if err != nil {
		fmt.Println("SHIT! There are some errors!")
	}
	defer fw.Close()

	tw := tar.NewWriter(fw)
	err = tw.WriteHeader(hdr)
	handleError(err)
	defer tw.Close()

	// Write file data
	io.Copy(tw, fr)
	fmt.Println("finished!")
}

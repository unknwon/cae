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

// Package tz enables you to transparently read or write tar.gz archives and the files inside them.
package tz

import (
	"archive/tar"
	// "compress/gzip"
	// "fmt"
	"errors"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/Unknwon/cae"
)

// File represents a file in archive.
type File struct {
	*tar.Header
	oldName string
	absPath string
}

// TzArchive represents a file archive, compressed with Zip.
type TzArchive struct {
	*ReadCloser
	FileName   string
	NumFiles   int
	Flag       int
	Permission os.FileMode

	files        []*File
	isHasChanged bool

	// For supporting to flush to io.Writer.
	writer      io.Writer
	isHasWriter bool
}

// Create creates the named tar.gz file, truncating
// it if it already exists. If successful, methods on the returned
// TzArchive can be used for I/O; the associated file descriptor has mode
// O_RDWR.
// If there is an error, it will be of type *PathError.
func Create(fileName string) (tz *TzArchive, err error) {
	os.MkdirAll(path.Dir(fileName), os.ModePerm)
	return OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
}

// Open opens the named tar.gz file for reading.  If successful, methods on
// the returned TzArchive can be used for reading; the associated file
// descriptor has mode O_RDONLY.
// If there is an error, it will be of type *PathError.
func Open(fileName string) (tz *TzArchive, err error) {
	return OpenFile(fileName, os.O_RDONLY, 0)
}

// OpenFile is the generalized open call; most users will use Open
// instead. It opens the named tar.gz file with specified flag
// (O_RDONLY etc.) if applicable. If successful,
// methods on the returned TzArchive can be used for I/O.
// If there is an error, it will be of type *PathError.
func OpenFile(fileName string, flag int, perm os.FileMode) (*TzArchive, error) {
	tz := &TzArchive{}
	err := tz.Open(fileName, flag, perm)
	return tz, err
}

// New accepts a variable that implemented interface io.Writer
// for write-only purpose operations.
func New(w io.Writer) (tz *TzArchive) {
	return &TzArchive{
		writer:      w,
		isHasWriter: true,
	}
}

// ListName returns a string slice of files' name in TzArchive.
func (tz *TzArchive) ListName(prefixes ...string) []string {
	isHasPrefix := len(prefixes) > 0
	names := make([]string, 0, tz.NumFiles)
	for _, f := range tz.files {
		if isHasPrefix {
			if cae.HasPrefix(f.Name, prefixes) {
				names = append(names, f.Name)
			}
			continue
		}
		names = append(names, f.Name)
	}
	return names
}

// AddEmptyDir adds a directory entry to TzArchive,
// it returns false when directory already existed.
func (tz *TzArchive) AddEmptyDir(dirPath string) bool {
	if !strings.HasSuffix(dirPath, "/") {
		dirPath += "/"
	}

	for _, f := range tz.files {
		if dirPath == f.Name {
			return false
		}
	}

	dirPath = strings.TrimSuffix(dirPath, "/")
	if strings.Contains(dirPath, "/") {
		// Auto add all upper level directory.
		tmpPath := path.Dir(dirPath)
		tz.AddEmptyDir(tmpPath)
	}
	tz.files = append(tz.files, &File{
		Header: &tar.Header{
			Name: dirPath + "/",
		},
	})
	tz.updateStat()
	return true
}

// AddFile adds a directory and subdirectories entries to TzArchive,
func (tz *TzArchive) AddDir(dirPath, absPath string) error {
	dir, err := os.Open(absPath)
	if err != nil {
		return err
	}
	defer dir.Close()

	tz.AddEmptyDir(dirPath)

	// Get file info slice
	fis, err := dir.Readdir(0)
	if err != nil {
		return err
	}

	for _, fi := range fis {
		curPath := strings.Replace(absPath+"/"+fi.Name(), "\\", "/", -1)
		tmpRecPath := strings.Replace(filepath.Join(dirPath, fi.Name()), "\\", "/", -1)
		if fi.IsDir() {
			err = tz.AddDir(tmpRecPath, curPath)
			if err != nil {
				return err
			}
		} else {
			err = tz.AddFile(tmpRecPath, curPath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// AddFile adds a file entry to TzArchive,
func (tz *TzArchive) AddFile(fileName, absPath string) error {
	if cae.GlobalFilter(absPath) {
		return nil
	}

	si, err := os.Lstat(absPath)
	if err != nil {
		return err
	}

	target := ""
	if si.Mode()&os.ModeSymlink != 0 {
		target, err = os.Readlink(absPath)
		if err != nil {
			return err
		}
	}

	file := new(File)
	file.Header, err = tar.FileInfoHeader(si, target)
	if err != nil {
		return err
	}
	file.Name = fileName
	file.absPath = absPath

	tz.AddEmptyDir(path.Dir(fileName))

	isExist := false
	for _, f := range tz.files {
		if fileName == f.Name {
			f = file
			isExist = true
			break
		}
	}

	if !isExist {
		tz.files = append(tz.files, file)
	}
	tz.updateStat()
	return nil
}

// DeleteIndex deletes an entry in the archive using its index.
func (tz *TzArchive) DeleteIndex(index int) error {
	if index >= tz.NumFiles {
		return errors.New("index out of range of number of files")
	}

	tz.files = append(tz.files[:index], tz.files[index+1:]...)
	return nil
}

// DeleteName deletes an entry in the archive using its name.
func (tz *TzArchive) DeleteName(name string) error {
	for i, f := range tz.files {
		if f.Name == name {
			return tz.DeleteIndex(i)
		}
	}
	return errors.New("entry with given name not found")
}

func (tz *TzArchive) updateStat() {
	tz.NumFiles = len(tz.files)
	tz.isHasChanged = true
}

// main functions shows how to TarGz a directory/file and
// UnTarGz a file
//func main() {
//	targetFilePath := "testdata.tar.gz"
//	srcDirPath := "testdata"
//	TarGz(srcDirPath, targetFilePath)
//	UnTarGz(targetFilePath, srcDirPath+"_temp")
//}

// Gzip and tar from source directory or file to destination file
// you need check file exist before you call this function
// func TarGz(srcDirPath string, destFilePath string) {
// 	fw, err := os.Create(destFilePath)
// 	handleError(err)
// 	defer fw.Close()

// 	// Gzip writer
// 	gw := gzip.NewWriter(fw)
// 	defer gw.Close()

// 	// Tar writer
// 	tw := tar.NewWriter(gw)
// 	defer tw.Close()

// 	// Check if it's a file or a directory
// 	f, err := os.Open(srcDirPath)
// 	handleError(err)
// 	fi, err := f.Stat()
// 	handleError(err)
// 	if fi.IsDir() {
// 		// handle source directory
// 		fmt.Println("Cerating tar.gz from directory...")
// 		tarGzDir(srcDirPath, path.Base(srcDirPath), tw)
// 	} else {
// 		// handle file directly
// 		fmt.Println("Cerating tar.gz from " + fi.Name() + "...")
// 		tarGzFile(srcDirPath, fi.Name(), tw, fi)
// 	}
// 	fmt.Println("Well done!")
// }

// Deal with directories
// if find files, handle them with tarGzFile
// Every recurrence append the base path to the recPath
// recPath is the path inside of tar.gz
// func tarGzDir(srcDirPath string, recPath string, tw *tar.Writer) {
// 	// Open source diretory
// 	dir, err := os.Open(srcDirPath)
// 	handleError(err)
// 	defer dir.Close()

// 	// Get file info slice
// 	fis, err := dir.Readdir(0)
// 	handleError(err)
// 	for _, fi := range fis {
// 		// Append path
// 		curPath := srcDirPath + "/" + fi.Name()
// 		// Check it is directory or file
// 		if fi.IsDir() {
// 			// Directory
// 			// (Directory won't add unitl all subfiles are added)
// 			fmt.Printf("Adding path...%s\n", curPath)
// 			tarGzDir(curPath, recPath+"/"+fi.Name(), tw)
// 		} else {
// 			// File
// 			fmt.Printf("Adding file...%s\n", curPath)
// 		}

// 		tarGzFile(curPath, recPath+"/"+fi.Name(), tw, fi)
// 	}
// }

// Deal with files
// func tarGzFile(srcFile string, recPath string, tw *tar.Writer, fi os.FileInfo) {
// 	if fi.IsDir() {
// 		// Create tar header
// 		hdr := new(tar.Header)
// 		// if last character of header name is '/' it also can be directory
// 		// but if you don't set Typeflag, error will occur when you untargz
// 		hdr.Name = recPath + "/"
// 		hdr.Typeflag = tar.TypeDir
// 		hdr.Size = 0
// 		//hdr.Mode = 0755 | c_ISDIR
// 		hdr.Mode = int64(fi.Mode())
// 		hdr.ModTime = fi.ModTime()

// 		// Write hander
// 		err := tw.WriteHeader(hdr)
// 		handleError(err)
// 	} else {
// 		// File reader
// 		fr, err := os.Open(srcFile)
// 		handleError(err)
// 		defer fr.Close()

// 		// Create tar header
// 		hdr := new(tar.Header)
// 		hdr.Name = recPath
// 		hdr.Size = fi.Size()
// 		hdr.Mode = int64(fi.Mode())
// 		hdr.ModTime = fi.ModTime()

// 		// Write hander
// 		err = tw.WriteHeader(hdr)
// 		handleError(err)

// 		// Write file data
// 		_, err = io.Copy(tw, fr)
// 		handleError(err)
// 	}
// }

// Ungzip and untar from source file to destination directory
// you need check file exist before you call this function
// func UnTarGz(srcFilePath string, destDirPath string) {
// 	fmt.Println("UnTarGzing " + srcFilePath + "...")
// 	// Create destination directory
// 	os.Mkdir(destDirPath, os.ModePerm)

// 	fr, err := os.Open(srcFilePath)
// 	handleError(err)
// 	defer fr.Close()

// 	// Gzip reader
// 	gr, err := gzip.NewReader(fr)
// 	handleError(err)
// 	defer gr.Close()

// 	// Tar reader
// 	tr := tar.NewReader(gr)

// 	for {
// 		hdr, err := tr.Next()
// 		if err == io.EOF {
// 			// End of tar archive
// 			break
// 		}
// 		//handleError(err)
// 		fmt.Println("UnTarGzing file..." + hdr.Name)
// 		// Check if it is diretory or file
// 		if hdr.Typeflag != tar.TypeDir {
// 			// Get files from archive
// 			// Create diretory before create file
// 			os.MkdirAll(destDirPath+"/"+path.Dir(hdr.Name), os.ModePerm)
// 			// Write data to file
// 			fw, _ := os.Create(destDirPath + "/" + hdr.Name)
// 			handleError(err)
// 			_, err = io.Copy(fw, tr)
// 			handleError(err)
// 		}
// 	}
// 	fmt.Println("Well done!")
// }

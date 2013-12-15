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

package zip_test

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/Unknwon/cae/zip"
	"github.com/Unknwon/com"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCreate(t *testing.T) {
	Convey("Create a zip file", t, func() {
		_, err := zip.Create(path.Join(os.TempDir(), "testdata/TestCreate.zip"))
		So(err, ShouldBeNil)
	})
}

func TestOpen(t *testing.T) {
	Convey("Open a zip file normally with read-only flag", t, func() {
		z, err := zip.Open("testdata/test.zip")
		So(err, ShouldBeNil)
		So(z.FileName, ShouldEqual, "testdata/test.zip")
		So(z.Comment, ShouldEqual, "This is the comment for test.zip")
		So(z.NumFiles, ShouldEqual, 5)
		So(z.Flag, ShouldEqual, os.O_RDONLY)
		So(z.Permission, ShouldEqual, 0)
	})

	Convey("Open a zip file that does not exist", t, func() {
		_, err := zip.Open("testdata/404.zip")
		So(err, ShouldNotBeNil)
	})

	Convey("Open a file that is not a zip file", t, func() {
		_, err := zip.Open("testdata/readme.notzip")
		So(err, ShouldNotBeNil)
	})
}

func TestListName(t *testing.T) {
	Convey("Open a zip file and get list of file/dir name", t, func() {
		z, err := zip.Open("testdata/test.zip")
		So(err, ShouldBeNil)

		Convey("List name without prefix", func() {
			So(strings.Join(z.ListName(), " "), ShouldEqual, "dir/ dir/bar dir/empty/ hello readonly")
		})

		Convey("List name with prefix", func() {
			So(strings.Join(z.ListName("h"), " "), ShouldEqual, "hello")
		})
	})
}

func TestAddEmptyDir(t *testing.T) {
	Convey("Open a zip file and add empty dirs", t, func() {
		z, err := zip.Create(path.Join(os.TempDir(), "testdata/TestAddEmptyDir.zip"))
		So(err, ShouldBeNil)

		Convey("Add dir that does not exist in list", func() {
			isExist := z.AddEmptyDir("level1")
			So(isExist, ShouldBeTrue)
		})

		Convey("Add dir that does exist in list", func() {
			isExist := z.AddEmptyDir("level1")
			So(!isExist, ShouldBeTrue)
		})

		Convey("Add multiple-level dir", func() {
			z.AddEmptyDir("level1/level2/level3/level4")
			So(strings.Join(z.ListName(), " "), ShouldEqual,
				"level1/ level1/level2/ level1/level2/level3/ level1/level2/level3/level4/")
		})
	})
}

func TestAddDir(t *testing.T) {
	Convey("Open a zip file and add dir with files", t, func() {
		z, err := zip.Create(path.Join(os.TempDir(), "testdata/TestAddDir.zip"))
		So(err, ShouldBeNil)

		Convey("Add a dir that does exist", func() {
			err := z.AddDir("testdata/testdir", "testdata/testdir")
			So(err, ShouldBeNil)
			So(strings.Join(z.ListName(), " "), ShouldEqual,
				"testdata/ testdata/testdir/ testdata/testdir/gophercolor16x16.png testdata/testdir/level1/ testdata/testdir/level1/bar/ testdata/testdir/level1/README.txt")
		})
	})
}

func TestAddFile(t *testing.T) {
	Convey("Open a zip file and add files", t, func() {
		z, err := zip.Create(path.Join(os.TempDir(), "testdata/TestAddFile.zip"))
		So(err, ShouldBeNil)

		Convey("Add a file that does exist", func() {
			err := z.AddFile("testdata/README.txt", "testdata/README.txt")
			So(err, ShouldBeNil)
			So(strings.Join(z.ListName(), " "), ShouldEqual,
				"testdata/ testdata/README.txt")
		})

		Convey("Add a file that does not exist", func() {
			err := z.AddFile("testdata/README.txt", "testdata/README_404.txt")
			So(err, ShouldNotBeNil)
		})
	})
}

func TestExtractTo(t *testing.T) {
	Convey("Extract a zip file to given path", t, func() {
		z, err := zip.Open("testdata/test.zip")
		So(err, ShouldBeNil)

		Convey("Extract the zip file without entries", func() {
			err := z.ExtractTo(path.Join(os.TempDir(), "testdata/test1"))
			So(err, ShouldBeNil)
			list, err := com.StatDir(path.Join(os.TempDir(), "testdata/test1"), true)
			So(strings.Join(list, " "), ShouldEqual, "dir/ dir/bar dir/empty/ hello readonly")
		})

		Convey("Extract the zip file with entries", func() {
			err := z.ExtractTo(path.Join(os.TempDir(), "testdata/test2"),
				"dir/", "dir/bar", "readonly")
			So(err, ShouldBeNil)
			list, err := com.StatDir(path.Join(os.TempDir(), "testdata/test2"), true)
			So(strings.Join(list, " "), ShouldEqual, "dir/ dir/bar readonly")
		})
	})
}

func TestFlush(t *testing.T) {
	Convey("Do some operations and flush to file system", t, func() {
		z, err := zip.Create(path.Join(os.TempDir(), "testdata/TestFlush.zip"))
		So(err, ShouldBeNil)

		z.AddEmptyDir("level1/level2/level3/level4")
		err = z.AddFile("testdata/README.txt", "testdata/README.txt")
		So(err, ShouldBeNil)

		err = z.Flush()
		So(err, ShouldBeNil)
	})
}

func TestClose(t *testing.T) {
	Convey("Close the zip file currently operating", t, func() {
		z, err := zip.Open("testdata/test.zip")
		So(err, ShouldBeNil)
		err = z.Close()
		So(err, ShouldBeNil)
	})
}

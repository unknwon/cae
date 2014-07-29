// Copyright 2013 Unknwon
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
	"fmt"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/Unknwon/cae"
	"github.com/Unknwon/com"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCreate(t *testing.T) {
	Convey("Create a zip file", t, func() {
		_, err := Create(path.Join(os.TempDir(), "testdata/TestCreate.zip"))
		So(err, ShouldBeNil)
	})
}

func TestOpen(t *testing.T) {
	Convey("Open a zip file normally with read-only flag", t, func() {
		z, err := Open("testdata/test.zip")
		So(err, ShouldBeNil)
		So(z.FileName, ShouldEqual, "testdata/test.zip")
		So(z.Comment, ShouldEqual, "This is the comment for test.zip")
		So(z.NumFiles, ShouldEqual, 5)
		So(z.Flag, ShouldEqual, os.O_RDONLY)
		So(z.Permission, ShouldEqual, 0)
	})

	Convey("Open a zip file that does not exist", t, func() {
		_, err := Open("testdata/404.zip")
		So(err, ShouldNotBeNil)
	})

	Convey("Open a file that is not a zip file", t, func() {
		_, err := Open("testdata/readme.notzip")
		So(err, ShouldNotBeNil)
	})
}

func TestList(t *testing.T) {
	Convey("Open a zip file and get list of file/dir name", t, func() {
		z, err := Open("testdata/test.zip")
		So(err, ShouldBeNil)

		Convey("List without prefix", func() {
			So(strings.Join(z.List(), " "), ShouldEqual,
				"dir/ dir/bar dir/empty/ hello readonly")
		})

		Convey("List with prefix", func() {
			So(strings.Join(z.List("h"), " "), ShouldEqual, "hello")
		})
	})
}

func TestAddEmptyDir(t *testing.T) {
	Convey("Open a zip file and add empty dirs", t, func() {
		z, err := Create(path.Join(os.TempDir(), "testdata/TestAddEmptyDir.zip"))
		So(err, ShouldBeNil)

		Convey("Add dir that does not exist and then add again", func() {
			So(z.AddEmptyDir("level1"), ShouldBeTrue)
			So(!z.AddEmptyDir("level1"), ShouldBeTrue)
		})

		Convey("Add multiple-level dir", func() {
			z.AddEmptyDir("level1/level2/level3/level4")
			So(strings.Join(z.List(), " "), ShouldEqual,
				"level1/ level1/level2/ level1/level2/level3/ level1/level2/level3/level4/")
		})
	})
}

func TestAddDir(t *testing.T) {
	Convey("Open a zip file and add dir with files", t, func() {
		z, err := Create(path.Join(os.TempDir(), "testdata/TestAddDir.zip"))
		So(err, ShouldBeNil)

		Convey("Add a dir that does exist", func() {
			So(z.AddDir("testdata/testdir", "testdata/testdir"), ShouldBeNil)
			So(strings.Join(z.List(), " "), ShouldEqual,
				"testdata/ testdata/testdir/ testdata/testdir/gophercolor16x16.png "+
					"testdata/testdir/level1/ testdata/testdir/level1/README.txt")
		})

		Convey("Add a dir that does not exist", func() {
			So(z.AddDir("testdata/testdir", "testdata/testdir404"), ShouldNotBeNil)
		})

		Convey("Add a dir that is a file", func() {
			So(z.AddDir("testdata/testdir", "testdata/README.txt"), ShouldNotBeNil)
		})
	})
}

func TestAddFile(t *testing.T) {
	Convey("Open a zip file and add files", t, func() {
		z, err := Create(path.Join(os.TempDir(), "testdata/TestAddFile.zip"))
		So(err, ShouldBeNil)

		Convey("Add a file that does exist", func() {
			So(z.AddFile("testdata/README.txt", "testdata/gophercolor16x16.png"), ShouldBeNil)
			So(strings.Join(z.List(), " "), ShouldEqual, "testdata/ testdata/README.txt")
		})

		Convey("Add a file that does not exist", func() {
			So(z.AddFile("testdata/README.txt", "testdata/README_404.txt"), ShouldNotBeNil)
		})

		Convey("Add a file that is exist in list", func() {
			So(z.AddFile("testdata/README.txt", "testdata/README.txt"), ShouldBeNil)
			So(strings.Join(z.List(), " "), ShouldEqual,
				"testdata/ testdata/README.txt")
		})
	})
}

func TestExtractTo(t *testing.T) {
	Convey("Extract a zip file to given path", t, func() {
		z, err := Open("testdata/test.zip")
		So(err, ShouldBeNil)

		Convey("Extract the zip file without entries", func() {
			os.RemoveAll(path.Join(os.TempDir(), "testdata/test1"))
			So(z.ExtractTo(path.Join(os.TempDir(), "testdata/test1")), ShouldBeNil)
			list, err := com.StatDir(path.Join(os.TempDir(), "testdata/test1"), true)
			So(err, ShouldBeNil)
			So(com.CompareSliceStrU(list,
				strings.Split("dir/ dir/bar dir/empty/ hello readonly", " ")), ShouldBeTrue)
		})

		Convey("Extract the zip file with entries", func() {
			os.RemoveAll(path.Join(os.TempDir(), "testdata/test2"))
			So(z.ExtractTo(
				path.Join(os.TempDir(), "testdata/test2"),
				"dir/", "dir/bar", "readonly"), ShouldBeNil)
			list, err := com.StatDir(path.Join(os.TempDir(), "testdata/test2"), true)
			So(err, ShouldBeNil)
			So(com.CompareSliceStrU(
				list, strings.Split("dir/ dir/bar readonly", " ")), ShouldBeTrue)
		})
	})
}

func TestFlush(t *testing.T) {
	Convey("Do some operations and flush to file system", t, func() {
		z, err := Create(path.Join(os.TempDir(), "testdata/TestFlush.zip"))
		So(err, ShouldBeNil)

		z.AddEmptyDir("level1/level2/level3/level4")
		So(z.AddFile("testdata/README.txt", "testdata/README.txt"), ShouldBeNil)

		// Add symbolic links.
		So(z.AddFile("testdata/test.lnk", "testdata/test.lnk"), ShouldBeNil)
		So(z.AddFile("testdata/testdir.lnk", "testdata/testdir.lnk"), ShouldBeNil)

		fmt.Println("Flushing to local file system...")
		// So(z.Flush(), ShouldBeNil)
	})

	Convey("Do some operation and flush to io.Writer", t, func() {
		f, err := os.Create(path.Join(os.TempDir(), "testdata/TestFlush2.zip"))
		So(err, ShouldBeNil)
		So(f, ShouldNotBeNil)

		z := New(f)
		z.AddEmptyDir("level1/level2/level3/level4")
		So(z.AddFile("testdata/README.txt", "testdata/README.txt"), ShouldBeNil)

		fmt.Println("Flushing to local io.Writer...")
		So(z.Flush(), ShouldBeNil)
		So(z.Flush(), ShouldBeNil) // In case first time is lucky.
	})
}

func TestPackTo(t *testing.T) {
	Convey("Pack a dir or file to zip file", t, func() {
		Convey("Pack a dir that does exist and includir root dir", func() {
			So(PackTo("testdata/testdir",
				path.Join(os.TempDir(), "testdata/testdir1.zip"), true), ShouldBeNil)
		})

		Convey("Pack a dir that does exist and does not includir root dir", func() {
			So(PackTo("testdata/testdir",
				path.Join(os.TempDir(), "testdata/testdir2.zip")), ShouldBeNil)
		})

		Convey("Pack a dir that does not exist and does not includir root dir", func() {
			So(PackTo("testdata/testdir404",
				path.Join(os.TempDir(), "testdata/testdir3.zip")), ShouldNotBeNil)
		})

		Convey("Pack a file that does exist", func() {
			So(PackTo("testdata/README.txt",
				path.Join(os.TempDir(), "testdata/testdir4.zip")), ShouldBeNil)
		})

		Convey("Pack a file that does not exist", func() {
			So(PackTo("testdata/README404.txt",
				path.Join(os.TempDir(), "testdata/testdir5.zip")), ShouldNotBeNil)
		})
	})
}

func TestClose(t *testing.T) {
	Convey("Close the zip file currently operating", t, func() {
		z, err := Open("testdata/test.zip")
		So(err, ShouldBeNil)
		So(z.Close(), ShouldBeNil)
	})
}

func TestDeleteIndex(t *testing.T) {
	Convey("Delete an entry with given index", t, func() {
		z, err := Create(path.Join(os.TempDir(), "testdata/TestDeleteIndex.zip"))
		So(err, ShouldBeNil)

		z.AddEmptyDir("level1/level2/level3/level4")

		Convey("Delete an entry with valid index", func() {
			So(z.DeleteIndex(3), ShouldBeNil)
			So(strings.Join(z.List(), " "), ShouldEqual,
				"level1/ level1/level2/ level1/level2/level3/")
		})

		Convey("Delete an entry with invalid index", func() {
			So(z.DeleteIndex(5), ShouldNotBeNil)
		})

		Convey("Test after flush", func() {
			So(z.Flush(), ShouldBeNil)
			So(strings.Join(z.List(), " "), ShouldEqual,
				"level1/ level1/level2/ level1/level2/level3/ level1/level2/level3/level4/")
		})
	})
}

func TestDeleteName(t *testing.T) {
	Convey("Delete an entry with given name", t, func() {
		z, err := Create(path.Join(os.TempDir(), "testdata/TestDeleteName.zip"))
		So(err, ShouldBeNil)

		z.AddEmptyDir("level1/level2/level3/level4")

		Convey("Delete an entry with valid name", func() {
			So(z.DeleteName("level1/level2/level3/level4/"), ShouldBeNil)
			So(strings.Join(z.List(), " "), ShouldEqual,
				"level1/ level1/level2/ level1/level2/level3/")
		})

		Convey("Delete an entry with invalid name", func() {
			So(z.DeleteName("level1/level2/level3/level"), ShouldNotBeNil)
		})

		Convey("Test after flush", func() {
			So(z.Flush(), ShouldBeNil)
			So(strings.Join(z.List(), " "), ShouldEqual,
				"level1/ level1/level2/ level1/level2/level3/ level1/level2/level3/level4/")
		})
	})
}

func TestCopy(t *testing.T) {
	Convey("Copy file from A to B", t, func() {
		Convey("Copy a file that does exist", func() {
			tmpPath := path.Join(os.TempDir(), "testdata/README.txt")
			So(cae.Copy(tmpPath, "testdata/README.txt"), ShouldBeNil)
			So(cae.IsExist(tmpPath), ShouldBeTrue)
		})

		Convey("Copy a file that does not exist", func() {
			So(cae.Copy(
				path.Join(os.TempDir(), "testdata/README.txt"),
				"testdata/404.txt"), ShouldNotBeNil)
		})
	})
}

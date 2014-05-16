// Copyright 2014 Unknown
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

package tz

import (
	"os"
	"path"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestStream(t *testing.T) {
	Convey("Create a stream archive", t, func() {
		fw, err := os.Create(path.Join(os.TempDir(), "testdata/TestStream.tar.gz"))
		So(err, ShouldBeNil)
		s := NewStreamArachive(fw)

		Convey("Stream a file", func() {
			f, err := os.Open("testdata/gophercolor16x16.png")
			So(err, ShouldBeNil)

			fi, err := f.Stat()
			So(err, ShouldBeNil)

			data := make([]byte, fi.Size())
			_, err = f.Read(data)
			So(err, ShouldBeNil)

			So(s.StreamFile("", fi, data), ShouldBeNil)
		})

		Convey("Stream a file with type directory", func() {
			f, err := os.Open("testdata")
			So(err, ShouldBeNil)

			fi, err := f.Stat()
			So(err, ShouldBeNil)

			So(s.StreamFile("", fi, nil), ShouldBeNil)
		})

		Convey("Stream a reader", func() {
			f, err := os.Open("testdata/gophercolor16x16.png")
			So(err, ShouldBeNil)

			fi, err := f.Stat()
			So(err, ShouldBeNil)

			So(s.StreamReader("", fi, f), ShouldBeNil)
		})
	})
}

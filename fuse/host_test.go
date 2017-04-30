/*
 * host_test.go
 *
 * Copyright 2017 Bill Zissimopoulos
 */
/*
 * This file is part of Cgofuse.
 *
 * It is licensed under the MIT license. The full license text can be found
 * in the License.txt file at the root of this project.
 */

package fuse

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

type testfs struct {
	FileSystemBase
}

func (self *testfs) Getattr(path string, stat *Stat_t, fh uint64) (errc int) {
	switch path {
	case "/":
		stat.Mode = S_IFDIR | 0555
		return 0
	default:
		return -ENOENT
	}
}

func (self *testfs) Readdir(path string,
	fill func(name string, stat *Stat_t, ofst int64) bool,
	ofst int64,
	fh uint64) (errc int) {
	fill(".", nil, 0)
	fill("..", nil, 0)
	return 0
}

func TestHost(t *testing.T) {
	path, err := ioutil.TempDir("", "test")
	if nil != err {
		panic(err)
	}
	defer os.Remove(path)
	mntp := filepath.Join(path, "m")
	if "windows" != runtime.GOOS {
		err = os.Mkdir(mntp, os.FileMode(0755))
		if nil != err {
			panic(err)
		}
		defer os.Remove(mntp)
	}
	done := make(chan bool)
	tmch := time.After(3 * time.Second)
	host := NewFileSystemHost(&testfs{})
	mres := false
	go func() {
		mres = host.Mount([]string{"test", mntp})
		done <- true
	}()
	<-tmch
	ures := host.Unmount()
	<-done
	if !mres {
		t.Error("Mount failed")
	}
	if !ures {
		t.Error("Unmount failed")
	}
}

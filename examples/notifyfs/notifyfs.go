// +build windows

/*
 * notifyfs.go
 *
 * Copyright 2017-2022 Bill Zissimopoulos
 */
/*
 * This file is part of Cgofuse.
 *
 * It is licensed under the MIT license. The full license text can be found
 * in the License.txt file at the root of this project.
 */

package main

import (
	"fmt"
	"os"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/winfsp/cgofuse/fuse"
)

type Notifyfs struct {
	fuse.FileSystemBase
	ticks uint64
}

func countFromTicks(ticks uint64) uint64 {
	/*
	 * The formula below produces the periodic sequence:
	 *     0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1,
	 *     0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1,
	 *     ...
	 */
	div10 := (ticks % 20) / 10
	mod10 := ticks % 10
	mdv10 := 1 - div10
	mmd10 := 10 - mod10
	return mdv10*mod10 + div10*mmd10
}

func (self *Notifyfs) count() uint64 {
	return countFromTicks(atomic.LoadUint64(&self.ticks))
}

func (self *Notifyfs) lookup(path string) uint64 {
	count := self.count()
	index, _ := strconv.ParseUint(path[1:], 10, 0)
	if 0 < index && index <= count {
		return index
	}
	return ^uint64(0)
}

func (self *Notifyfs) Open(path string, flags int) (errc int, fh uint64) {
	index := self.lookup(path)
	if ^uint64(0) == index {
		return -fuse.ENOENT, ^uint64(0)
	}
	return 0, index
}

func (self *Notifyfs) Getattr(path string, stat *fuse.Stat_t, fh uint64) (errc int) {
	if "/" == path {
		stat.Mode = fuse.S_IFDIR | 0555
		return 0
	}
	index := self.lookup(path)
	if ^uint64(0) == index {
		return -fuse.ENOENT
	}
	contents := strconv.FormatUint(index, 10) + "\n"
	stat.Mode = fuse.S_IFREG | 0444
	stat.Size = int64(len(contents))
	return 0
}

func (self *Notifyfs) Read(path string, buff []byte, ofst int64, fh uint64) (n int) {
	index := self.lookup(path)
	if ^uint64(0) == index {
		return -fuse.ENOENT
	}
	contents := strconv.FormatUint(index, 10) + "\n"
	endofst := ofst + int64(len(buff))
	if endofst > int64(len(contents)) {
		endofst = int64(len(contents))
	}
	if endofst < ofst {
		return 0
	}
	n = copy(buff, contents[ofst:endofst])
	return
}

func (self *Notifyfs) Readdir(path string,
	fill func(name string, stat *fuse.Stat_t, ofst int64) bool,
	ofst int64,
	fh uint64) (errc int) {
	fill(".", nil, 0)
	fill("..", nil, 0)
	count := self.count()
	for u := uint64(1); count >= u; u++ {
		fill(strconv.FormatUint(u, 10), nil, 0)
	}
	return 0
}

func (self *Notifyfs) tick(host *fuse.FileSystemHost) {
	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {
		ticks := atomic.AddUint64(&self.ticks, 1)
		oldcount := countFromTicks(ticks - 1)
		newcount := countFromTicks(ticks)
		if oldcount < newcount {
			fmt.Println("CREATE", "/"+strconv.FormatUint(newcount, 10))
			host.Notify("/"+strconv.FormatUint(newcount, 10), fuse.NOTIFY_CREATE)
		} else if oldcount > newcount {
			fmt.Println("UNLINK", "/"+strconv.FormatUint(oldcount, 10))
			host.Notify("/"+strconv.FormatUint(oldcount, 10), fuse.NOTIFY_UNLINK)
		}
	}
}

func main() {
	notifyfs := &Notifyfs{}
	host := fuse.NewFileSystemHost(notifyfs)
	go notifyfs.tick(host)
	host.Mount("", os.Args[1:])
}

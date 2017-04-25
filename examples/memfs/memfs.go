/*
 * memfs.go
 *
 * Copyright 2017 Bill Zissimopoulos
 */
/*
 * This file is part of Cgofuse.
 *
 * You can redistribute it and/or modify it under the terms of the GNU
 * General Public License version 3 as published by the Free Software
 * Foundation.
 */

package main

import (
	"fmt"
	"github.com/billziss-gh/cgofuse/examples/shared"
	"github.com/billziss-gh/cgofuse/fuse"
	"os"
	"strings"
	"sync"
	"time"
)

func trace(vals ...interface{}) func(vals ...interface{}) {
	uid, gid, _ := fuse.Getcontext()
	return shared.Trace(1, fmt.Sprintf("[uid=%v,gid=%v]", uid, gid), vals...)
}

func split(path string) []string {
	return strings.Split(path, "/")
}

func resize(slice []byte, size int64, zeroinit bool) []byte {
	const allocunit = 64 * 1024
	allocsize := (size + allocunit - 1) / allocunit * allocunit
	if cap(slice) != int(allocsize) {
		newslice := make([]byte, size, allocsize)
		copy(newslice, slice)
		slice = newslice
	} else if zeroinit {
		i := len(slice)
		slice = slice[:size]
		for ; len(slice) > i; i++ {
			slice[i] = 0
		}
	}
	return slice
}

type node_t struct {
	stat    fuse.Stat_t
	xatr    map[string]string
	chld    map[string]*node_t
	data    []byte
	opencnt int
}

func newNode(dev uint64, ino uint64, mode uint32, uid uint32, gid uint32) *node_t {
	nano := time.Now().UnixNano()
	tmsp := fuse.Timespec{nano / 1e9, nano % 1e9}
	self := node_t{
		fuse.Stat_t{
			Dev:      dev,
			Ino:      ino,
			Mode:     mode,
			Nlink:    1,
			Uid:      uid,
			Gid:      gid,
			Atim:     tmsp,
			Mtim:     tmsp,
			Ctim:     tmsp,
			Birthtim: tmsp,
		},
		nil,
		nil,
		nil,
		0}
	if 0040000 == self.stat.Mode&0170000 {
		self.chld = map[string]*node_t{}
	}
	return &self
}

type Memfs struct {
	fuse.FileSystemBase
	lock    sync.Mutex
	ino     uint64
	root    *node_t
	openmap map[uint64]*node_t
}

func (self *Memfs) Mknod(path string, mode uint32, dev uint64) (errc int) {
	defer trace(path, mode, dev)(&errc)
	defer self.synchronize()()
	return self.makeNode(path, mode, dev, nil)
}

func (self *Memfs) Mkdir(path string, mode uint32) (errc int) {
	defer trace(path, mode)(&errc)
	defer self.synchronize()()
	return self.makeNode(path, 0040000|(mode&07777), 0, nil)
}

func (self *Memfs) Unlink(path string) (errc int) {
	defer trace(path)(&errc)
	defer self.synchronize()()
	return self.removeNode(path, false)
}

func (self *Memfs) Rmdir(path string) (errc int) {
	defer trace(path)(&errc)
	defer self.synchronize()()
	return self.removeNode(path, true)
}

func (self *Memfs) Link(oldpath string, newpath string) (errc int) {
	defer trace(oldpath, newpath)(&errc)
	defer self.synchronize()()
	_, _, oldnode := self.lookupNode(oldpath, nil)
	if nil == oldnode {
		return -fuse.ENOENT
	}
	newprnt, newname, newnode := self.lookupNode(newpath, nil)
	if nil == newprnt {
		return -fuse.ENOENT
	}
	if nil != newnode {
		return -fuse.EEXIST
	}
	oldnode.stat.Nlink++
	newprnt.chld[newname] = oldnode
	return 0
}

func (self *Memfs) Symlink(target string, newpath string) (errc int) {
	defer trace(target, newpath)(&errc)
	defer self.synchronize()()
	return self.makeNode(newpath, 0120777, 0, []byte(target))
}

func (self *Memfs) Readlink(path string) (errc int, target string) {
	defer trace(path)(&errc, &target)
	defer self.synchronize()()
	_, _, node := self.lookupNode(path, nil)
	if nil == node {
		return -fuse.ENOENT, ""
	}
	if 0120000 != node.stat.Mode&0170000 {
		return -fuse.EINVAL, ""
	}
	return 0, string(node.data)
}

func (self *Memfs) Rename(oldpath string, newpath string) (errc int) {
	defer trace(oldpath, newpath)(&errc)
	defer self.synchronize()()
	oldprnt, oldname, oldnode := self.lookupNode(oldpath, nil)
	if nil == oldnode {
		return -fuse.ENOENT
	}
	newprnt, newname, newnode := self.lookupNode(newpath, oldnode)
	if nil == newprnt {
		return -fuse.ENOENT
	}
	if "" == newname {
		// guard against directory loop creation
		return -fuse.EINVAL
	}
	if oldprnt == newprnt && oldname == newname {
		return 0
	}
	if nil != newnode {
		errc = self.removeNode(newpath, 0040000 == oldnode.stat.Mode&0170000)
		if 0 != errc {
			return errc
		}
	}
	delete(oldprnt.chld, oldname)
	newprnt.chld[newname] = oldnode
	return 0
}

func (self *Memfs) Chmod(path string, mode uint32) (errc int) {
	defer trace(path, mode)(&errc)
	defer self.synchronize()()
	_, _, node := self.lookupNode(path, nil)
	if nil == node {
		return -fuse.ENOENT
	}
	node.stat.Mode = mode & 07777
	return 0
}

func (self *Memfs) Chown(path string, uid uint32, gid uint32) (errc int) {
	defer trace(path, uid, gid)(&errc)
	defer self.synchronize()()
	_, _, node := self.lookupNode(path, nil)
	if nil == node {
		return -fuse.ENOENT
	}
	if ^uint32(0) != uid {
		node.stat.Uid = uid
	}
	if ^uint32(0) != gid {
		node.stat.Gid = gid
	}
	return 0
}

func (self *Memfs) Utimens(path string, tmsp []fuse.Timespec) (errc int) {
	defer trace(path, tmsp)(&errc)
	defer self.synchronize()()
	_, _, node := self.lookupNode(path, nil)
	if nil == node {
		return -fuse.ENOENT
	}
	if nil == tmsp {
		nano := time.Now().UnixNano()
		tmsa := [2]fuse.Timespec{
			fuse.Timespec{nano / 1e9, nano % 1e9},
			fuse.Timespec{nano / 1e9, nano % 1e9},
		}
		tmsp = tmsa[:]
	}
	node.stat.Atim = tmsp[0]
	node.stat.Mtim = tmsp[1]
	return 0
}

func (self *Memfs) Open(path string, flags int) (errc int, fh uint64) {
	defer trace(path, flags)(&errc, &fh)
	defer self.synchronize()()
	return self.openNode(path, false)
}

func (self *Memfs) Getattr(path string, stat *fuse.Stat_t, fh uint64) (errc int) {
	defer trace(path, fh)(&errc, stat)
	defer self.synchronize()()
	node := self.getNode(path, fh)
	if nil == node {
		return -fuse.ENOENT
	}
	*stat = node.stat
	return 0
}

func (self *Memfs) Truncate(path string, size int64, fh uint64) (errc int) {
	defer trace(path, size, fh)(&errc)
	defer self.synchronize()()
	node := self.getNode(path, fh)
	if nil == node {
		return -fuse.ENOENT
	}
	resize(node.data, size, true)
	node.stat.Size = size
	return 0
}

func (self *Memfs) Read(path string, buff []byte, ofst int64, fh uint64) (n int) {
	defer trace(path, buff, ofst, fh)(&n)
	defer self.synchronize()()
	node := self.getNode(path, fh)
	if nil == node {
		return -fuse.ENOENT
	}
	endofst := ofst + int64(len(buff))
	if endofst < node.stat.Size {
		endofst = node.stat.Size
	}
	if endofst < ofst {
		return 0
	}
	return copy(buff, node.data[ofst:endofst])
}

func (self *Memfs) Write(path string, buff []byte, ofst int64, fh uint64) (n int) {
	defer trace(path, buff, ofst, fh)(&n)
	defer self.synchronize()()
	node := self.getNode(path, fh)
	if nil == node {
		return -fuse.ENOENT
	}
	endofst := ofst + int64(len(buff))
	if endofst > node.stat.Size {
		resize(node.data, endofst, false)
		node.stat.Size = endofst
	}
	return copy(node.data[ofst:endofst], buff)
}

func (self *Memfs) Release(path string, fh uint64) (errc int) {
	defer trace(path, fh)(&errc)
	defer self.synchronize()()
	return self.closeNode(fh)
}

func (self *Memfs) Opendir(path string) (errc int, fh uint64) {
	defer trace(path)(&errc, &fh)
	defer self.synchronize()()
	return self.openNode(path, true)
}

func (self *Memfs) Readdir(path string,
	fill func(name string, stat *fuse.Stat_t, ofst int64) bool,
	ofst int64,
	fh uint64) (errc int) {
	defer trace(path, fill, ofst, fh)(&errc)
	defer self.synchronize()()
	node := self.openmap[fh]
	fill(".", &node.stat, 0)
	fill("..", nil, 0)
	for name, chld := range node.chld {
		if !fill(name, &chld.stat, 0) {
			break
		}
	}
	return 0
}

func (self *Memfs) Releasedir(path string, fh uint64) (errc int) {
	defer trace(path, fh)(&errc)
	defer self.synchronize()()
	return self.closeNode(fh)
}

func (self *Memfs) lookupNode(path string, ancestor *node_t) (prnt *node_t, name string, node *node_t) {
	prnt = self.root
	name = ""
	node = self.root
	for _, c := range split(path) {
		if "" != c {
			prnt, name = node, c
			node = node.chld[c]
			if nil != ancestor && node == ancestor {
				name = "" // special case loop condition
				return
			}
		}
	}
	return
}

func (self *Memfs) makeNode(path string, mode uint32, dev uint64, data []byte) int {
	prnt, name, node := self.lookupNode(path, nil)
	if nil == prnt {
		return -fuse.ENOENT
	}
	if nil != node {
		return -fuse.EEXIST
	}
	self.ino++
	uid, gid, _ := fuse.Getcontext()
	node = newNode(dev, self.ino, mode, uid, gid)
	if nil != data {
		node.data = make([]byte, len(data))
		node.stat.Size = int64(len(data))
		copy(node.data, data)
	}
	prnt.chld[name] = node
	return 0
}

func (self *Memfs) removeNode(path string, dir bool) int {
	prnt, name, node := self.lookupNode(path, nil)
	if nil == node {
		return -fuse.ENOENT
	}
	if !dir && 0040000 == node.stat.Mode&0170000 {
		return -fuse.EISDIR
	}
	if dir && 0040000 != node.stat.Mode&0170000 {
		return -fuse.ENOTDIR
	}
	if 0 < len(node.chld) {
		return -fuse.ENOTEMPTY
	}
	node.stat.Nlink--
	if 0 == node.stat.Nlink {
		delete(prnt.chld, name)
	}
	return 0
}

func (self *Memfs) openNode(path string, dir bool) (int, uint64) {
	_, _, node := self.lookupNode(path, nil)
	if nil == node {
		return -fuse.ENOENT, ^uint64(0)
	}
	if !dir && 0040000 == node.stat.Mode&0170000 {
		return -fuse.EISDIR, ^uint64(0)
	}
	if dir && 0040000 != node.stat.Mode&0170000 {
		return -fuse.ENOTDIR, ^uint64(0)
	}
	node.opencnt++
	if 1 == node.opencnt {
		self.openmap[node.stat.Ino] = node
	}
	return 0, node.stat.Ino
}

func (self *Memfs) closeNode(fh uint64) int {
	node := self.openmap[fh]
	node.opencnt--
	if 0 == node.opencnt {
		delete(self.openmap, node.stat.Ino)
	}
	return 0
}

func (self *Memfs) getNode(path string, fh uint64) *node_t {
	if ^uint64(0) == fh {
		_, _, node := self.lookupNode(path, nil)
		return node
	} else {
		return self.openmap[fh]
	}
}

func (self *Memfs) synchronize() func() {
	self.lock.Lock()
	return func() {
		self.lock.Unlock()
	}
}

func NewMemfs() *Memfs {
	self := Memfs{}
	defer self.synchronize()()
	self.ino++
	self.root = newNode(0, self.ino, 0040777, 0, 0)
	self.openmap = map[uint64]*node_t{}
	return &self
}

func main() {
	memfs := NewMemfs()
	host := fuse.NewFileSystemHost(memfs)
	host.Mount(os.Args)
}

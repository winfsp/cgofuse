//go:build memfs3
// +build memfs3

/*
 * memfs3.go
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
	"github.com/winfsp/cgofuse/fuse"
)

func (self *Memfs) Rename3(oldpath string, newpath string, flags uint32) (errc int) {
	defer trace(oldpath, newpath)(&errc)
	defer self.synchronize()()
	if 0 != flags&^fuse.RENAME_NOREPLACE {
		// we only support NOREPLACE
		return -fuse.EINVAL
	}
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
		if fuse.RENAME_NOREPLACE == flags&fuse.RENAME_NOREPLACE {
			return -fuse.EEXIST
		}
		errc = self.removeNode(newpath, fuse.S_IFDIR == oldnode.stat.Mode&fuse.S_IFMT)
		if 0 != errc {
			return errc
		}
	}
	delete(oldprnt.chld, oldname)
	newprnt.chld[newname] = oldnode
	return 0
}

func (self *Memfs) Chmod3(path string, mode uint32, fh uint64) (errc int) {
	defer trace(path, mode, fh)(&errc)
	defer self.synchronize()()
	node := self.getNode(path, fh)
	if nil == node {
		return -fuse.ENOENT
	}
	node.stat.Mode = (node.stat.Mode & fuse.S_IFMT) | mode&07777
	node.stat.Ctim = fuse.Now()
	return 0
}

func (self *Memfs) Chown3(path string, uid uint32, gid uint32, fh uint64) (errc int) {
	defer trace(path, uid, gid, fh)(&errc)
	defer self.synchronize()()
	node := self.getNode(path, fh)
	if nil == node {
		return -fuse.ENOENT
	}
	if ^uint32(0) != uid {
		node.stat.Uid = uid
	}
	if ^uint32(0) != gid {
		node.stat.Gid = gid
	}
	node.stat.Ctim = fuse.Now()
	return 0
}

func (self *Memfs) Utimens3(path string, tmsp []fuse.Timespec, fh uint64) (errc int) {
	defer trace(path, tmsp, fh)(&errc)
	defer self.synchronize()()
	node := self.getNode(path, fh)
	if nil == node {
		return -fuse.ENOENT
	}
	node.stat.Ctim = fuse.Now()
	if nil == tmsp {
		tmsp0 := node.stat.Ctim
		tmsa := [2]fuse.Timespec{tmsp0, tmsp0}
		tmsp = tmsa[:]
	}
	node.stat.Atim = tmsp[0]
	node.stat.Mtim = tmsp[1]
	return 0
}

var _ fuse.FileSystemRename3 = (*Memfs)(nil)
var _ fuse.FileSystemChmod3 = (*Memfs)(nil)
var _ fuse.FileSystemChown3 = (*Memfs)(nil)
var _ fuse.FileSystemRename3 = (*Memfs)(nil)

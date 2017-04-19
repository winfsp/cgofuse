/*
 * passthrough.go
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
	"cgofuse/fuse"
	"os"
	"path/filepath"
	"syscall"
)

type Ptfs struct {
	fuse.FileSystemBase
	root string
}

func errno(err error) int {
	if nil != err {
		return int(err.(syscall.Errno))
	} else {
		return 0
	}
}

func (self *Ptfs) Statfs(path string, stat *fuse.Statfs_t) int {
	path = filepath.Join(self.root, path)
	stgo := syscall.Statfs_t{}
	errc := errno(syscall.Statfs(path, &stgo))
	copyFusestatfsFromGostatfs(stat, &stgo)
	return errc
}

func (self *Ptfs) Getattr(path string, stat *fuse.Stat_t, fh uint64) (errc int) {
	stgo := syscall.Stat_t{}
	if ^uint64(0) == fh {
		path = filepath.Join(self.root, path)
		errc = errno(syscall.Stat(path, &stgo))
	} else {
		errc = errno(syscall.Fstat(int(fh), &stgo))
	}
	copyFusestatFromGostat(stat, &stgo)
	return
}

func main() {
	ptfs := Ptfs{}
	args := os.Args
	if 3 <= len(args) && '-' != args[len(args)-2][0] && '-' != args[len(args)-1][0] {
		ptfs.root, _ = filepath.Abs(args[len(args)-2])
		args = append(args[:len(args)-2], args[len(args)-1])
	}
	host := fuse.NewFileSystemHost(&ptfs)
	host.Mount(args)
}

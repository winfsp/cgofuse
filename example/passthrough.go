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

func errno(err error) syscall.Errno {
	if (nil != err) {
		return err.(syscall.Errno)
	} else {
		return 0
	}
}

func (self *Ptfs) Statfs(path string, stbuf *syscall.Statfs_t) syscall.Errno {
	path = filepath.Join(self.root, path)
	return errno(syscall.Statfs(path, stbuf))
}

func (self *Ptfs) Getattr(path string, stbuf *syscall.Stat_t, fh uint64) syscall.Errno {
	if ^uint64(0) == fh {
		path = filepath.Join(self.root, path)
		return errno(syscall.Stat(path, stbuf))
	} else {
		return errno(syscall.Fstat(int(fh), stbuf))
	}
}

func main() {
	ptfs := Ptfs{}
	args := os.Args
	if 3 <= len(args) && '-' != args[len(args) - 2][0] && '-' != args[len(args) - 1][0] {
		ptfs.root, _ = filepath.Abs(args[len(args) - 2])
		args = append(args[:len(args) - 2], args[len(args) - 1])
	}
	host := fuse.NewFileSystemHost(&ptfs)
	host.Mount(args)
}

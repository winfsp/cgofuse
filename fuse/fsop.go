/*
 * fsop.go
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

// Package fuse allows the creation of user mode file systems in Go.
//
// A user mode file system must implement the methods in FileSystemInterface
// and be hosted (mounted) by a FileSystemHost.
// Alternatively a user mode file system can use the FileSystemBase struct which
// provides default implementations of the methods in FileSystemInterface.
package fuse

import "syscall"

// FileSystemInterface is the interface that Cgofuse file systems must implement.
type FileSystemInterface interface {
	// Init is called when the file system is mounted.
	Init()

	// Destroy is called when the file system is unmounted.
	Destroy()

	// Statfs gets file system statistics.
	Statfs(path string) (syscall.Errno, *syscall.Statfs_t)

	// Mknod creates a file node.
	Mknod(path string, mode uint32, dev uint64) syscall.Errno

	// Mkdir creates a directory.
	Mkdir(path string, mode uint32) syscall.Errno

	// Unlink removes a file.
	Unlink(path string) syscall.Errno

	// Rmdir removes a directory.
	Rmdir(path string) syscall.Errno

	// Link creates a hard link to a file.
	Link(srcpath string, dstpath string) syscall.Errno

	// Symlink creates a symbolic link.
	Symlink(dstpath string, srcpath string) syscall.Errno

	// Readlink reads the target of a symbolic link.
	Readlink(path string) (syscall.Errno, string)

	// Rename renames a file.
	Rename(oldpath string, newpath string) syscall.Errno

	// Chmod changes the permission bits of a file.
	Chmod(path string, mode uint32) syscall.Errno

	// Chown changes the owner and group of a file.
	Chown(path string, uid uint32, gid uint32) syscall.Errno

	// Utimens changes the access and modification times of a file.
	Utimens(path string, tv []syscall.Timespec) syscall.Errno

	// Access checks file access permissions.
	Access(path string, mask int) syscall.Errno

	// Create creates and opens a file.
	Create(path string, mode uint32) (syscall.Errno, uint64)

	// Open opens a file.
	Open(path string) (syscall.Errno, uint64)

	// Getattr gets file attributes.
	Getattr(path string, fh uint64) (syscall.Errno, *syscall.Stat_t)

	// Truncate changes the size of a file.
	Truncate(path string, size uint64, fh uint64) syscall.Errno

	// Read reads data from a file.
	Read(path string, buf []byte, off uint64, fh uint64) syscall.Errno

	// Write writes data to a file.
	Write(path string, buf []byte, off uint64, fh uint64) syscall.Errno

	// Flush flushes cached file data.
	Flush(path string, fh uint64) syscall.Errno

	// Release closes an open file.
	Release(path string, fh uint64) syscall.Errno

	// Fsync synchronizes file contents.
	Fsync(path string, datasync bool, fh uint64) syscall.Errno

	//Lock(path string, fh uint64, cmd int, lock syscall.Flock_t) syscall.Errno

	// Opendir opens a directory.
	Opendir(path string) (syscall.Errno, uint64)

	// Readdir reads a directory.
	Readdir(path string, fh uint64) (syscall.Errno, []string)

	// Releasedir closes an open directory.
	Releasedir(path string, fh uint64) syscall.Errno

	// Fsyncdir synchronizes directory contents.
	Fsyncdir(path string, datasync bool, fh uint64) syscall.Errno

	// Setxattr sets extended attributes.
	Setxattr(path string, name string, value []byte, flags int) syscall.Errno

	// Getxattr gets extended attributes.
	Getxattr(path string, name string, value []byte) syscall.Errno

	// Removexattr removes extended attributes.
	Removexattr(path string, name string) syscall.Errno

	// Listxattr lists extended attributes.
	Listxattr(path string) (syscall.Errno, []string)
}

// FileSystemBase provides default implementations of the methods in FileSystemInterface.
type FileSystemBase struct {
}

func (*FileSystemBase) Init() {
}

func (*FileSystemBase) Destroy() {
}

func (*FileSystemBase) Statfs(path string) (syscall.Errno, *syscall.Statfs_t) {
	return syscall.ENOSYS, nil
}

func (*FileSystemBase) Mknod(path string, mode uint32, dev uint64) syscall.Errno {
	return syscall.ENOSYS
}

func (*FileSystemBase) Mkdir(path string, mode uint32) syscall.Errno {
	return syscall.ENOSYS
}

func (*FileSystemBase) Unlink(path string) syscall.Errno {
	return syscall.ENOSYS
}

func (*FileSystemBase) Rmdir(path string) syscall.Errno {
	return syscall.ENOSYS
}

func (*FileSystemBase) Link(srcpath string, dstpath string) syscall.Errno {
	return syscall.ENOSYS
}

func (*FileSystemBase) Symlink(dstpath string, srcpath string) syscall.Errno {
	return syscall.ENOSYS
}

func (*FileSystemBase) Readlink(path string) (syscall.Errno, string) {
	return syscall.ENOSYS, ""
}

func (*FileSystemBase) Rename(oldpath string, newpath string) syscall.Errno {
	return syscall.ENOSYS
}

func (*FileSystemBase) Chmod(path string, mode uint32) syscall.Errno {
	return syscall.ENOSYS
}

func (*FileSystemBase) Chown(path string, uid uint32, gid uint32) syscall.Errno {
	return syscall.ENOSYS
}

func (*FileSystemBase) Utimens(path string, tv []syscall.Timespec) syscall.Errno {
	return syscall.ENOSYS
}

func (*FileSystemBase) Access(path string, mask int) syscall.Errno {
	return syscall.ENOSYS
}

func (*FileSystemBase) Create(path string, mode uint32) (syscall.Errno, uint64) {
	return syscall.ENOSYS, ^uint64(0)
}

func (*FileSystemBase) Open(path string) (syscall.Errno, uint64) {
	return syscall.ENOSYS, ^uint64(0)
}

func (*FileSystemBase) Getattr(path string, fh uint64) (syscall.Errno, *syscall.Stat_t) {
	return syscall.ENOSYS, nil
}

func (*FileSystemBase) Truncate(path string, size uint64, fh uint64) syscall.Errno {
	return syscall.ENOSYS
}

func (*FileSystemBase) Read(path string, buf []byte, off uint64, fh uint64) syscall.Errno {
	return syscall.ENOSYS
}

func (*FileSystemBase) Write(path string, buf []byte, off uint64, fh uint64) syscall.Errno {
	return syscall.ENOSYS
}

func (*FileSystemBase) Flush(path string, fh uint64) syscall.Errno {
	return syscall.ENOSYS
}

func (*FileSystemBase) Release(path string, fh uint64) syscall.Errno {
	return syscall.ENOSYS
}

func (*FileSystemBase) Fsync(path string, datasync bool, fh uint64) syscall.Errno {
	return syscall.ENOSYS
}

/*
func (*FileSystemBase) Lock(path string, fh uint64, cmd int, lock syscall.Flock_t) syscall.Errno {
	return syscall.ENOSYS
}
*/

func (*FileSystemBase) Opendir(path string) (syscall.Errno, uint64) {
	return syscall.ENOSYS, ^uint64(0)
}

func (*FileSystemBase) Readdir(path string, fh uint64) (syscall.Errno, []string) {
	return syscall.ENOSYS, nil
}

func (*FileSystemBase) Releasedir(path string, fh uint64) syscall.Errno {
	return syscall.ENOSYS
}

func (*FileSystemBase) Fsyncdir(path string, datasync bool, fh uint64) syscall.Errno {
	return syscall.ENOSYS
}

func (*FileSystemBase) Setxattr(path string, name string, value []byte, flags int) syscall.Errno {
	return syscall.ENOSYS
}

func (*FileSystemBase) Getxattr(path string, name string, value []byte) syscall.Errno {
	return syscall.ENOSYS
}

func (*FileSystemBase) Removexattr(path string, name string) syscall.Errno {
	return syscall.ENOSYS
}

func (*FileSystemBase) Listxattr(path string) (syscall.Errno, []string) {
	return syscall.ENOSYS, nil
}

var _ FileSystemInterface = (*FileSystemBase)(nil)

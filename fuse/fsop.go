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

/*
#if !(defined(__APPLE__) || defined(__linux__) || defined(_WIN32))
#error platform not supported
#endif

#include <errno.h>

#if defined(__APPLE__) || defined(__linux__)

#include <sys/xattr.h>

#elif defined(_WIN32)

#define ENOATTR ENODATA

#define XATTR_CREATE  1
#define XATTR_REPLACE 2

#endif
*/
import "C"

const (
	E2BIG           = C.E2BIG
	EACCES          = C.EACCES
	EADDRINUSE      = C.EADDRINUSE
	EADDRNOTAVAIL   = C.EADDRNOTAVAIL
	EAFNOSUPPORT    = C.EAFNOSUPPORT
	EAGAIN          = C.EAGAIN
	EALREADY        = C.EALREADY
	EBADF           = C.EBADF
	EBADMSG         = C.EBADMSG
	EBUSY           = C.EBUSY
	ECANCELED       = C.ECANCELED
	ECHILD          = C.ECHILD
	ECONNABORTED    = C.ECONNABORTED
	ECONNREFUSED    = C.ECONNREFUSED
	ECONNRESET      = C.ECONNRESET
	EDEADLK         = C.EDEADLK
	EDESTADDRREQ    = C.EDESTADDRREQ
	EDOM            = C.EDOM
	EEXIST          = C.EEXIST
	EFAULT          = C.EFAULT
	EFBIG           = C.EFBIG
	EHOSTUNREACH    = C.EHOSTUNREACH
	EIDRM           = C.EIDRM
	EILSEQ          = C.EILSEQ
	EINPROGRESS     = C.EINPROGRESS
	EINTR           = C.EINTR
	EINVAL          = C.EINVAL
	EIO             = C.EIO
	EISCONN         = C.EISCONN
	EISDIR          = C.EISDIR
	ELOOP           = C.ELOOP
	EMFILE          = C.EMFILE
	EMLINK          = C.EMLINK
	EMSGSIZE        = C.EMSGSIZE
	ENAMETOOLONG    = C.ENAMETOOLONG
	ENETDOWN        = C.ENETDOWN
	ENETRESET       = C.ENETRESET
	ENETUNREACH     = C.ENETUNREACH
	ENFILE          = C.ENFILE
	ENOATTR         = C.ENOATTR
	ENOBUFS         = C.ENOBUFS
	ENODATA         = C.ENODATA
	ENODEV          = C.ENODEV
	ENOENT          = C.ENOENT
	ENOEXEC         = C.ENOEXEC
	ENOLCK          = C.ENOLCK
	ENOLINK         = C.ENOLINK
	ENOMEM          = C.ENOMEM
	ENOMSG          = C.ENOMSG
	ENOPROTOOPT     = C.ENOPROTOOPT
	ENOSPC          = C.ENOSPC
	ENOSR           = C.ENOSR
	ENOSTR          = C.ENOSTR
	ENOSYS          = C.ENOSYS
	ENOTCONN        = C.ENOTCONN
	ENOTDIR         = C.ENOTDIR
	ENOTEMPTY       = C.ENOTEMPTY
	ENOTRECOVERABLE = C.ENOTRECOVERABLE
	ENOTSOCK        = C.ENOTSOCK
	ENOTSUP         = C.ENOTSUP
	ENOTTY          = C.ENOTTY
	ENXIO           = C.ENXIO
	EOPNOTSUPP      = C.EOPNOTSUPP
	EOVERFLOW       = C.EOVERFLOW
	EOWNERDEAD      = C.EOWNERDEAD
	EPERM           = C.EPERM
	EPIPE           = C.EPIPE
	EPROTO          = C.EPROTO
	EPROTONOSUPPORT = C.EPROTONOSUPPORT
	EPROTOTYPE      = C.EPROTOTYPE
	ERANGE          = C.ERANGE
	EROFS           = C.EROFS
	ESPIPE          = C.ESPIPE
	ESRCH           = C.ESRCH
	ETIME           = C.ETIME
	ETIMEDOUT       = C.ETIMEDOUT
	ETXTBSY         = C.ETXTBSY
	EWOULDBLOCK     = C.EWOULDBLOCK
	EXDEV           = C.EXDEV
)

const (
	XATTR_CREATE  = int(C.XATTR_CREATE)
	XATTR_REPLACE = int(C.XATTR_REPLACE)
)

type Timespec struct {
	Sec  int64
	Nsec int64
}

type Statfs_t struct {
	Bsize   uint64
	Frsize  uint64
	Blocks  uint64
	Bfree   uint64
	Bavail  uint64
	Files   uint64
	Ffree   uint64
	Favail  uint64
	Fsid    uint64
	Flag    uint64
	Namemax uint64
}

type Stat_t struct {
	Dev      uint64
	Ino      uint64
	Mode     uint32
	Nlink    uint32
	Uid      uint32
	Gid      uint32
	Rdev     uint64
	Size     int64
	Atim     Timespec
	Mtim     Timespec
	Ctim     Timespec
	Blksize  int64
	Blocks   int64
	Birthtim Timespec
}

// FileSystemInterface is the interface that Cgofuse file systems must implement.
type FileSystemInterface interface {
	// Init is called when the file system is mounted.
	Init()

	// Destroy is called when the file system is unmounted.
	Destroy()

	// Statfs gets file system statistics.
	Statfs(path string, stat *Statfs_t) int

	// Mknod creates a file node.
	Mknod(path string, mode uint32, dev uint64) int

	// Mkdir creates a directory.
	Mkdir(path string, mode uint32) int

	// Unlink removes a file.
	Unlink(path string) int

	// Rmdir removes a directory.
	Rmdir(path string) int

	// Link creates a hard link to a file.
	Link(oldpath string, newpath string) int

	// Symlink creates a symbolic link.
	Symlink(target string, newpath string) int

	// Readlink reads the target of a symbolic link.
	Readlink(path string) (int, string)

	// Rename renames a file.
	Rename(oldpath string, newpath string) int

	// Chmod changes the permission bits of a file.
	Chmod(path string, mode uint32) int

	// Chown changes the owner and group of a file.
	Chown(path string, uid uint32, gid uint32) int

	// Utimens changes the access and modification times of a file.
	Utimens(path string, tmsp []Timespec) int

	// Access checks file access permissions.
	Access(path string, mask uint32) int

	// Create creates and opens a file.
	Create(path string, mode uint32) (int, uint64)

	// Open opens a file.
	Open(path string, flags int) (int, uint64)

	// Getattr gets file attributes.
	Getattr(path string, stat *Stat_t, fh uint64) int

	// Truncate changes the size of a file.
	Truncate(path string, size int64, fh uint64) int

	// Read reads data from a file.
	Read(path string, buff []byte, ofst int64, fh uint64) int

	// Write writes data to a file.
	Write(path string, buff []byte, ofst int64, fh uint64) int

	// Flush flushes cached file data.
	Flush(path string, fh uint64) int

	// Release closes an open file.
	Release(path string, fh uint64) int

	// Fsync synchronizes file contents.
	Fsync(path string, datasync bool, fh uint64) int

	//Lock(path string, fh uint64, cmd int, lock Flock_t) int

	// Opendir opens a directory.
	Opendir(path string) (int, uint64)

	// Readdir reads a directory.
	Readdir(path string,
		fill func(name string, stat *Stat_t, ofst int64) bool,
		ofst int64,
		fh uint64) int

	// Releasedir closes an open directory.
	Releasedir(path string, fh uint64) int

	// Fsyncdir synchronizes directory contents.
	Fsyncdir(path string, datasync bool, fh uint64) int

	// Setxattr sets extended attributes.
	Setxattr(path string, name string, value []byte, flags int) int

	// Getxattr gets extended attributes.
	Getxattr(path string, name string, fill func(value []byte) bool) int

	// Removexattr removes extended attributes.
	Removexattr(path string, name string) int

	// Listxattr lists extended attributes.
	Listxattr(path string, fill func(name string) bool) int
}

// FileSystemBase provides default implementations of the methods in FileSystemInterface.
type FileSystemBase struct {
}

func (*FileSystemBase) Init() {
}

func (*FileSystemBase) Destroy() {
}

func (*FileSystemBase) Statfs(path string, stat *Statfs_t) int {
	return -ENOSYS
}

func (*FileSystemBase) Mknod(path string, mode uint32, dev uint64) int {
	return -ENOSYS
}

func (*FileSystemBase) Mkdir(path string, mode uint32) int {
	return -ENOSYS
}

func (*FileSystemBase) Unlink(path string) int {
	return -ENOSYS
}

func (*FileSystemBase) Rmdir(path string) int {
	return -ENOSYS
}

func (*FileSystemBase) Link(oldpath string, newpath string) int {
	return -ENOSYS
}

func (*FileSystemBase) Symlink(target string, newpath string) int {
	return -ENOSYS
}

func (*FileSystemBase) Readlink(path string) (int, string) {
	return -ENOSYS, ""
}

func (*FileSystemBase) Rename(oldpath string, newpath string) int {
	return -ENOSYS
}

func (*FileSystemBase) Chmod(path string, mode uint32) int {
	return -ENOSYS
}

func (*FileSystemBase) Chown(path string, uid uint32, gid uint32) int {
	return -ENOSYS
}

func (*FileSystemBase) Utimens(path string, tmsp []Timespec) int {
	return -ENOSYS
}

func (*FileSystemBase) Access(path string, mask uint32) int {
	return -ENOSYS
}

func (*FileSystemBase) Create(path string, mode uint32) (int, uint64) {
	return -ENOSYS, ^uint64(0)
}

func (*FileSystemBase) Open(path string, flags int) (int, uint64) {
	return -ENOSYS, ^uint64(0)
}

func (*FileSystemBase) Getattr(path string, stat *Stat_t, fh uint64) int {
	return -ENOSYS
}

func (*FileSystemBase) Truncate(path string, size int64, fh uint64) int {
	return -ENOSYS
}

func (*FileSystemBase) Read(path string, buff []byte, ofst int64, fh uint64) int {
	return -ENOSYS
}

func (*FileSystemBase) Write(path string, buff []byte, ofst int64, fh uint64) int {
	return -ENOSYS
}

func (*FileSystemBase) Flush(path string, fh uint64) int {
	return -ENOSYS
}

func (*FileSystemBase) Release(path string, fh uint64) int {
	return -ENOSYS
}

func (*FileSystemBase) Fsync(path string, datasync bool, fh uint64) int {
	return -ENOSYS
}

/*
func (*FileSystemBase) Lock(path string, fh uint64, cmd int, lock Flock_t) int {
	return -ENOSYS
}
*/

func (*FileSystemBase) Opendir(path string) (int, uint64) {
	return -ENOSYS, ^uint64(0)
}

func (*FileSystemBase) Readdir(path string,
	fill func(name string, stat *Stat_t, ofst int64) bool,
	ofst int64,
	fh uint64) int {
	return -ENOSYS
}

func (*FileSystemBase) Releasedir(path string, fh uint64) int {
	return -ENOSYS
}

func (*FileSystemBase) Fsyncdir(path string, datasync bool, fh uint64) int {
	return -ENOSYS
}

func (*FileSystemBase) Setxattr(path string, name string, value []byte, flags int) int {
	return -ENOSYS
}

func (*FileSystemBase) Getxattr(path string, name string, fill func(value []byte) bool) int {
	return -ENOSYS
}

func (*FileSystemBase) Removexattr(path string, name string) int {
	return -ENOSYS
}

func (*FileSystemBase) Listxattr(path string, fill func(name string) bool) int {
	return -ENOSYS
}

var _ FileSystemInterface = (*FileSystemBase)(nil)

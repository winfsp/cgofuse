/*
 * fsop.go
 *
 * Copyright 2017 Bill Zissimopoulos
 */
/*
 * This file is part of Cgofuse.
 *
 * It is licensed under the MIT license. The full license text can be found
 * in the License.txt file at the root of this project.
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

#if defined(__APPLE__) || defined(__linux__)

#include <errno.h>
#include <fcntl.h>

#elif defined(_WIN32)

#define EPERM           1
#define ENOENT          2
#define ESRCH           3
#define EINTR           4
#define EIO             5
#define ENXIO           6
#define E2BIG           7
#define ENOEXEC         8
#define EBADF           9
#define ECHILD          10
#define EAGAIN          11
#define ENOMEM          12
#define EACCES          13
#define EFAULT          14
#define EBUSY           16
#define EEXIST          17
#define EXDEV           18
#define ENODEV          19
#define ENOTDIR         20
#define EISDIR          21
#define ENFILE          23
#define EMFILE          24
#define ENOTTY          25
#define EFBIG           27
#define ENOSPC          28
#define ESPIPE          29
#define EROFS           30
#define EMLINK          31
#define EPIPE           32
#define EDOM            33
#define EDEADLK         36
#define ENAMETOOLONG    38
#define ENOLCK          39
#define ENOSYS          40
#define ENOTEMPTY       41
#define EINVAL          22
#define ERANGE          34
#define EILSEQ          42
#define EADDRINUSE      100
#define EADDRNOTAVAIL   101
#define EAFNOSUPPORT    102
#define EALREADY        103
#define EBADMSG         104
#define ECANCELED       105
#define ECONNABORTED    106
#define ECONNREFUSED    107
#define ECONNRESET      108
#define EDESTADDRREQ    109
#define EHOSTUNREACH    110
#define EIDRM           111
#define EINPROGRESS     112
#define EISCONN         113
#define ELOOP           114
#define EMSGSIZE        115
#define ENETDOWN        116
#define ENETRESET       117
#define ENETUNREACH     118
#define ENOBUFS         119
#define ENODATA         120
#define ENOLINK         121
#define ENOMSG          122
#define ENOPROTOOPT     123
#define ENOSR           124
#define ENOSTR          125
#define ENOTCONN        126
#define ENOTRECOVERABLE 127
#define ENOTSOCK        128
#define ENOTSUP         129
#define EOPNOTSUPP      130
#define EOTHER          131
#define EOVERFLOW       132
#define EOWNERDEAD      133
#define EPROTO          134
#define EPROTONOSUPPORT 135
#define EPROTOTYPE      136
#define ETIME           137
#define ETIMEDOUT       138
#define ETXTBSY         139
#define EWOULDBLOCK     140

#include <fcntl.h>
#define O_RDONLY        _O_RDONLY
#define O_WRONLY        _O_WRONLY
#define O_RDWR          _O_RDWR
#define O_APPEND        _O_APPEND
#define O_CREAT         _O_CREAT
#define O_EXCL          _O_EXCL
#define O_TRUNC         _O_TRUNC

#endif

#if defined(__linux__) || defined(_WIN32)
// incantation needed for cgo to figure out "kind of name" for ENOATTR
#define ENOATTR ((int)ENODATA)
#endif

#if defined(__APPLE__) || defined(__linux__)
#include <sys/xattr.h>
#elif defined(_WIN32)
#define XATTR_CREATE  1
#define XATTR_REPLACE 2
#endif
*/
import "C"
import (
	"strconv"
	"time"
)

const (
	E2BIG           = int(C.E2BIG)
	EACCES          = int(C.EACCES)
	EADDRINUSE      = int(C.EADDRINUSE)
	EADDRNOTAVAIL   = int(C.EADDRNOTAVAIL)
	EAFNOSUPPORT    = int(C.EAFNOSUPPORT)
	EAGAIN          = int(C.EAGAIN)
	EALREADY        = int(C.EALREADY)
	EBADF           = int(C.EBADF)
	EBADMSG         = int(C.EBADMSG)
	EBUSY           = int(C.EBUSY)
	ECANCELED       = int(C.ECANCELED)
	ECHILD          = int(C.ECHILD)
	ECONNABORTED    = int(C.ECONNABORTED)
	ECONNREFUSED    = int(C.ECONNREFUSED)
	ECONNRESET      = int(C.ECONNRESET)
	EDEADLK         = int(C.EDEADLK)
	EDESTADDRREQ    = int(C.EDESTADDRREQ)
	EDOM            = int(C.EDOM)
	EEXIST          = int(C.EEXIST)
	EFAULT          = int(C.EFAULT)
	EFBIG           = int(C.EFBIG)
	EHOSTUNREACH    = int(C.EHOSTUNREACH)
	EIDRM           = int(C.EIDRM)
	EILSEQ          = int(C.EILSEQ)
	EINPROGRESS     = int(C.EINPROGRESS)
	EINTR           = int(C.EINTR)
	EINVAL          = int(C.EINVAL)
	EIO             = int(C.EIO)
	EISCONN         = int(C.EISCONN)
	EISDIR          = int(C.EISDIR)
	ELOOP           = int(C.ELOOP)
	EMFILE          = int(C.EMFILE)
	EMLINK          = int(C.EMLINK)
	EMSGSIZE        = int(C.EMSGSIZE)
	ENAMETOOLONG    = int(C.ENAMETOOLONG)
	ENETDOWN        = int(C.ENETDOWN)
	ENETRESET       = int(C.ENETRESET)
	ENETUNREACH     = int(C.ENETUNREACH)
	ENFILE          = int(C.ENFILE)
	ENOATTR         = int(C.ENOATTR)
	ENOBUFS         = int(C.ENOBUFS)
	ENODATA         = int(C.ENODATA)
	ENODEV          = int(C.ENODEV)
	ENOENT          = int(C.ENOENT)
	ENOEXEC         = int(C.ENOEXEC)
	ENOLCK          = int(C.ENOLCK)
	ENOLINK         = int(C.ENOLINK)
	ENOMEM          = int(C.ENOMEM)
	ENOMSG          = int(C.ENOMSG)
	ENOPROTOOPT     = int(C.ENOPROTOOPT)
	ENOSPC          = int(C.ENOSPC)
	ENOSR           = int(C.ENOSR)
	ENOSTR          = int(C.ENOSTR)
	ENOSYS          = int(C.ENOSYS)
	ENOTCONN        = int(C.ENOTCONN)
	ENOTDIR         = int(C.ENOTDIR)
	ENOTEMPTY       = int(C.ENOTEMPTY)
	ENOTRECOVERABLE = int(C.ENOTRECOVERABLE)
	ENOTSOCK        = int(C.ENOTSOCK)
	ENOTSUP         = int(C.ENOTSUP)
	ENOTTY          = int(C.ENOTTY)
	ENXIO           = int(C.ENXIO)
	EOPNOTSUPP      = int(C.EOPNOTSUPP)
	EOVERFLOW       = int(C.EOVERFLOW)
	EOWNERDEAD      = int(C.EOWNERDEAD)
	EPERM           = int(C.EPERM)
	EPIPE           = int(C.EPIPE)
	EPROTO          = int(C.EPROTO)
	EPROTONOSUPPORT = int(C.EPROTONOSUPPORT)
	EPROTOTYPE      = int(C.EPROTOTYPE)
	ERANGE          = int(C.ERANGE)
	EROFS           = int(C.EROFS)
	ESPIPE          = int(C.ESPIPE)
	ESRCH           = int(C.ESRCH)
	ETIME           = int(C.ETIME)
	ETIMEDOUT       = int(C.ETIMEDOUT)
	ETXTBSY         = int(C.ETXTBSY)
	EWOULDBLOCK     = int(C.EWOULDBLOCK)
	EXDEV           = int(C.EXDEV)
)

const (
	O_RDONLY = int(C.O_RDONLY)
	O_WRONLY = int(C.O_WRONLY)
	O_RDWR   = int(C.O_RDWR)
	O_APPEND = int(C.O_APPEND)
	O_CREAT  = int(C.O_CREAT)
	O_EXCL   = int(C.O_EXCL)
	O_TRUNC  = int(C.O_TRUNC)
)

const (
	S_IFMT   = 0170000
	S_IFBLK  = 0060000
	S_IFCHR  = 0020000
	S_IFIFO  = 0010000
	S_IFREG  = 0100000
	S_IFDIR  = 0040000
	S_IFLNK  = 0120000
	S_IFSOCK = 0140000

	S_IRWXU = 00700
	S_IRUSR = 00400
	S_IWUSR = 00200
	S_IXUSR = 00100
	S_IRWXG = 00070
	S_IRGRP = 00040
	S_IWGRP = 00020
	S_IXGRP = 00010
	S_IRWXO = 00007
	S_IROTH = 00004
	S_IWOTH = 00002
	S_IXOTH = 00001
	S_ISUID = 04000
	S_ISGID = 02000
	S_ISVTX = 01000
)

const (
	XATTR_CREATE  = int(C.XATTR_CREATE)
	XATTR_REPLACE = int(C.XATTR_REPLACE)
)

// Timespec contains a time as the UNIX time in seconds and nanoseconds.
type Timespec struct {
	Sec  int64
	Nsec int64
}

// NewTimespec creates a Timespec from a time.Time.
func NewTimespec(t time.Time) Timespec {
	return Timespec{t.Unix(), int64(t.Nanosecond())}
}

// Now creates a Timespec that contains the current time.
func Now() Timespec {
	return NewTimespec(time.Now())
}

// Time returns the Timespec as a time.Time.
func (ts *Timespec) Time() time.Time {
	return time.Unix(ts.Sec, ts.Nsec)
}

// Statfs_t contains file system information.
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

// Stat contains file metadata information.
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
	Create(path string, flags int, mode uint32) (int, uint64)

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

// Error encapsulates a FUSE error code. In some rare circumstances it is useful
// to signal an error to the FUSE layer by boxing the error code using Error and
// calling panic(). The FUSE layer will recover and report the boxed error code
// to the OS.
type Error int

func (self Error) Error() string {
	return "fuse.Error(" + strconv.Itoa(int(self)) + ")"
}

var _ error = (*Error)(nil)

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

func (*FileSystemBase) Create(path string, flags int, mode uint32) (int, uint64) {
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

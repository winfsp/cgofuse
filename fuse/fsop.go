package fuse

import "syscall"

/*
 * FileSystemInterface
 */

type FileSystemInterface interface {
	Init()
	Destroy()
	Statfs(path string) (syscall.Errno, *syscall.Statfs_t)
	Mknod(path string, mode uint32, dev uint64) syscall.Errno
	Mkdir(path string, mode uint32) syscall.Errno
	Unlink(path string) syscall.Errno
	Rmdir(path string) syscall.Errno
	Link(srcpath string, dstpath string) syscall.Errno
	Symlink(dstpath string, srcpath string) syscall.Errno
	Readlink(path string) (syscall.Errno, string)
	Rename(oldpath string, newpath string) syscall.Errno
	Chmod(path string, mode uint32) syscall.Errno
	Chown(path string, uid uint32, gid uint32) syscall.Errno
	Utimens(path string, tv []syscall.Timespec) syscall.Errno
	Access(path string, mask int) syscall.Errno
	Create(path string, mode uint32) (syscall.Errno, uint64)
	Open(path string) (syscall.Errno, uint64)
	Getattr(path string, fh uint64) (syscall.Errno, *syscall.Stat_t)
	Truncate(path string, size uint64, fh uint64) syscall.Errno
	Read(path string, buf []byte, off uint64, fh uint64) syscall.Errno
	Write(path string, buf []byte, off uint64, fh uint64) syscall.Errno
	Flush(path string, fh uint64) syscall.Errno
	Release(path string, fh uint64) syscall.Errno
	Fsync(path string, datasync bool, fh uint64) syscall.Errno
	//Lock(path string, fh uint64, cmd int, lock syscall.Flock_t) syscall.Errno
	Opendir(path string) (syscall.Errno, uint64)
	Readdir(path string, fh uint64) (syscall.Errno, []string)
	Releasedir(path string, fh uint64) syscall.Errno
	Fsyncdir(path string, datasync bool, fh uint64) syscall.Errno
	Setxattr(path string, name string, value []byte, flags int) syscall.Errno
	Getxattr(path string, name string, value []byte) syscall.Errno
	Removexattr(path string, name string) syscall.Errno
	Listxattr(path string) (syscall.Errno, []string)
}

/*
 * FileSystemBase
 */

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

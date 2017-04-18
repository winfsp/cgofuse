package fuse

/*
#cgo CFLAGS: -DFUSE_USE_VERSION=28 -D_FILE_OFFSET_BITS=64 -I/usr/local/include/osxfuse
#cgo LDFLAGS: -L/usr/local/lib -losxfuse

#include <stdlib.h>
#include <fuse.h>

#define fuse_stat stat
#define fuse_statvfs statvfs
#define fuse_timespec timespec
#define fuse_mode_t mode_t
#define fuse_dev_t dev_t
#define fuse_uid_t uid_t
#define fuse_gid_t gid_t
#define fuse_off_t off_t

extern int hostGetattr(char *path, struct fuse_stat *stbuf);
extern int hostReadlink(char *path, char *buf, size_t size);
extern int hostMknod(char *path, fuse_mode_t mode, fuse_dev_t dev);
extern int hostMkdir(char *path, fuse_mode_t mode);
extern int hostUnlink(char *path);
extern int hostRmdir(char *path);
extern int hostSymlink(char *dstpath, char *srcpath);
extern int hostRename(char *oldpath, char *newpath);
extern int hostLink(char *srcpath, char *dstpath);
extern int hostChmod(char *path, fuse_mode_t mode);
extern int hostChown(char *path, fuse_uid_t uid, fuse_gid_t gid);
extern int hostTruncate(char *path, fuse_off_t size);
extern int hostOpen(char *path, struct fuse_file_info *fi);
extern int hostRead(char *path, char *buf, size_t size, fuse_off_t off,
    struct fuse_file_info *fi);
extern int hostWrite(char *path, char *buf, size_t size, fuse_off_t off,
    struct fuse_file_info *fi);
extern int hostStatfs(char *path, struct fuse_statvfs *stbuf);
extern int hostFlush(char *path, struct fuse_file_info *fi);
extern int hostRelease(char *path, struct fuse_file_info *fi);
extern int hostFsync(char *path, int datasync, struct fuse_file_info *fi);
extern int hostSetxattr(char *path, char *name, char *value, size_t size,
    int flags);
extern int hostGetxattr(char *path, char *name, char *value, size_t size);
extern int hostListxattr(char *path, char *namebuf, size_t size);
extern int hostRemovexattr(char *path, char *name);
extern int hostOpendir(char *path, struct fuse_file_info *fi);
extern int hostReaddir(char *path, char *buf, fuse_fill_dir_t filler, fuse_off_t off,
    struct fuse_file_info *fi);
extern int hostReleasedir(char *path, struct fuse_file_info *fi);
extern int hostFsyncdir(char *path, int datasync, struct fuse_file_info *fi);
extern void *hostInit(struct fuse_conn_info *conn);
extern void hostDestroy(void *data);
extern int hostAccess(char *path, int mask);
extern int hostCreate(char *path, fuse_mode_t mode, struct fuse_file_info *fi);
extern int hostFtruncate(char *path, fuse_off_t off, struct fuse_file_info *fi);
extern int hostFgetattr(char *path, struct fuse_stat *stbuf, struct fuse_file_info *fi);
//extern int hostLock(char *path, struct fuse_file_info *fi, int cmd, struct fuse_flock *lock);
extern int hostUtimens(char *path, struct fuse_timespec tv[2]);

static int hostFilldir(char *buf, fuse_fill_dir_t filler, char *name)
{
    return filler(buf, name, 0, 0);
}

static struct fuse_operations *hostFsop(void)
{
#if defined(__GNUC__)
#pragma GCC diagnostic push
#pragma GCC diagnostic ignored "-Wincompatible-pointer-types"
#endif
    static struct fuse_operations fsop =
    {
        .getattr = hostGetattr,
        .readlink = hostReadlink,
        .mknod = hostMknod,
        .mkdir = hostMkdir,
        .unlink = hostUnlink,
        .rmdir = hostRmdir,
        .symlink = hostSymlink,
        .rename = hostRename,
        .link = hostLink,
        .chmod = hostChmod,
        .chown = hostChown,
        .truncate = hostTruncate,
        .open = hostOpen,
        .read = hostRead,
        .write = hostWrite,
        .statfs = hostStatfs,
        .flush = hostFlush,
        .release = hostRelease,
        .fsync = hostFsync,
        .setxattr = hostSetxattr,
        .getxattr = hostGetxattr,
        .listxattr = hostListxattr,
        .removexattr = hostRemovexattr,
        .opendir = hostOpendir,
        .readdir = hostReaddir,
        .releasedir = hostReleasedir,
        .fsyncdir = hostFsyncdir,
        .init = hostInit,
        .destroy = hostDestroy,
        .access = hostAccess,
        .create = hostCreate,
        .ftruncate = hostFtruncate,
        .fgetattr = hostFgetattr,
        //.lock = hostFlock,
        .utimens = hostUtimens,
    };
    return &fsop;
#if defined(__GNUC__)
#pragma GCC diagnostic pop
#endif
}

static size_t hostFsopSize(void)
{
    return sizeof(struct fuse_operations);
}
*/
import "C"

import (
	"syscall"
	"unsafe"
)

/*
 * FileSystemHost
 */

type FileSystemHost struct {
	fsop FileSystemInterface
}

func copyCstatvfsFromGostatfs(dst *C.struct_statvfs, src *syscall.Statfs_t) {
	if nil == src {
		return
	}
	*dst = C.struct_statvfs{}
	// !!!
}

func copyCstatFromGostat(dst *C.struct_stat, src *syscall.Stat_t) {
	if nil == src {
		return
	}
	*dst = C.struct_stat{}
	// !!!
}

func copyGotimespecFromCtimespec(dst *syscall.Timespec, src *C.struct_timespec) {
	if nil == src {
		return
	}
	*dst = syscall.Timespec{}
	// !!!
}

//export hostGetattr
func hostGetattr(path0 *C.char, stbuf0 *C.struct_stat) (errno C.int) {
	defer func() {
		if r := recover(); r != nil {
			errno = -C.int(syscall.EIO)
		}
	}()
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	serr, rslt := fsop.Getattr(path, ^uint64(0))
	copyCstatFromGostat(stbuf0, rslt)
	return -C.int(serr)
}

//export hostReadlink
func hostReadlink(path0 *C.char, buf0 *C.char, size0 C.size_t) (errno C.int) {
	defer func() {
		if r := recover(); r != nil {
			errno = -C.int(syscall.EIO)
		}
	}()
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	serr, rslt := fsop.Readlink(path)
	buf := (*[1 << 30]byte)(unsafe.Pointer(buf0))
	copy(buf[:size0], rslt)
	return -C.int(serr)
}

//export hostMknod
func hostMknod(path0 *C.char, mode0 C.mode_t, dev0 C.dev_t) (errno C.int) {
	defer func() {
		if r := recover(); r != nil {
			errno = -C.int(syscall.EIO)
		}
	}()
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	serr := fsop.Mknod(path, uint32(mode0), uint64(dev0))
	return -C.int(serr)
}

//export hostMkdir
func hostMkdir(path0 *C.char, mode0 C.mode_t) (errno C.int) {
	defer func() {
		if r := recover(); r != nil {
			errno = -C.int(syscall.EIO)
		}
	}()
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	serr := fsop.Mkdir(path, uint32(mode0))
	return -C.int(serr)
}

//export hostUnlink
func hostUnlink(path0 *C.char) (errno C.int) {
	defer func() {
		if r := recover(); r != nil {
			errno = -C.int(syscall.EIO)
		}
	}()
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	serr := fsop.Unlink(path)
	return -C.int(serr)
}

//export hostRmdir
func hostRmdir(path0 *C.char) (errno C.int) {
	defer func() {
		if r := recover(); r != nil {
			errno = -C.int(syscall.EIO)
		}
	}()
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	serr := fsop.Rmdir(path)
	return -C.int(serr)
}

//export hostSymlink
func hostSymlink(dstpath0 *C.char, srcpath0 *C.char) (errno C.int) {
	defer func() {
		if r := recover(); r != nil {
			errno = -C.int(syscall.EIO)
		}
	}()
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	srcpath, dstpath := C.GoString(srcpath0), C.GoString(dstpath0)
	serr := fsop.Symlink(srcpath, dstpath)
	return -C.int(serr)
}

//export hostRename
func hostRename(oldpath0 *C.char, newpath0 *C.char) (errno C.int) {
	defer func() {
		if r := recover(); r != nil {
			errno = -C.int(syscall.EIO)
		}
	}()
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	oldpath, newpath := C.GoString(oldpath0), C.GoString(newpath0)
	serr := fsop.Rename(oldpath, newpath)
	return -C.int(serr)
}

//export hostLink
func hostLink(dstpath0 *C.char, srcpath0 *C.char) (errno C.int) {
	defer func() {
		if r := recover(); r != nil {
			errno = -C.int(syscall.EIO)
		}
	}()
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	srcpath, dstpath := C.GoString(srcpath0), C.GoString(dstpath0)
	serr := fsop.Link(srcpath, dstpath)
	return -C.int(serr)
}

//export hostChmod
func hostChmod(path0 *C.char, mode0 C.mode_t) (errno C.int) {
	defer func() {
		if r := recover(); r != nil {
			errno = -C.int(syscall.EIO)
		}
	}()
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	serr := fsop.Chmod(path, uint32(mode0))
	return -C.int(serr)
}

//export hostChown
func hostChown(path0 *C.char, uid0 C.uid_t, gid0 C.gid_t) (errno C.int) {
	defer func() {
		if r := recover(); r != nil {
			errno = -C.int(syscall.EIO)
		}
	}()
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	serr := fsop.Chown(path, uint32(uid0), uint32(gid0))
	return -C.int(serr)
}

//export hostTruncate
func hostTruncate(path0 *C.char, size0 C.off_t) (errno C.int) {
	defer func() {
		if r := recover(); r != nil {
			errno = -C.int(syscall.EIO)
		}
	}()
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	serr := fsop.Truncate(path, uint64(size0), ^uint64(0))
	return -C.int(serr)
}

//export hostOpen
func hostOpen(path0 *C.char, fi0 *C.struct_fuse_file_info) (errno C.int) {
	defer func() {
		if r := recover(); r != nil {
			errno = -C.int(syscall.EIO)
		}
	}()
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	serr, rslt := fsop.Open(path)
	fi0.fh = C.uint64_t(rslt)
	return -C.int(serr)
}

//export hostRead
func hostRead(path0 *C.char, buf0 *C.char, size0 C.size_t, off0 C.off_t,
	fi0 *C.struct_fuse_file_info) (errno C.int) {
	defer func() {
		if r := recover(); r != nil {
			errno = -C.int(syscall.EIO)
		}
	}()
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	buf := (*[1 << 30]byte)(unsafe.Pointer(buf0))
	serr := fsop.Read(path, buf[:size0], uint64(off0), uint64(fi0.fh))
	return -C.int(serr)
}

//export hostWrite
func hostWrite(path0 *C.char, buf0 *C.char, size0 C.size_t, off0 C.off_t,
	fi0 *C.struct_fuse_file_info) (errno C.int) {
	defer func() {
		if r := recover(); r != nil {
			errno = -C.int(syscall.EIO)
		}
	}()
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	buf := (*[1 << 30]byte)(unsafe.Pointer(buf0))
	serr := fsop.Write(path, buf[:size0], uint64(off0), uint64(fi0.fh))
	return -C.int(serr)
}

//export hostStatfs
func hostStatfs(path0 *C.char, stbuf0 *C.struct_statvfs) (errno C.int) {
	defer func() {
		if r := recover(); r != nil {
			errno = -C.int(syscall.EIO)
		}
	}()
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	serr, rslt := fsop.Statfs(path)
	copyCstatvfsFromGostatfs(stbuf0, rslt)
	return -C.int(serr)
}

//export hostFlush
func hostFlush(path0 *C.char, fi0 *C.struct_fuse_file_info) (errno C.int) {
	defer func() {
		if r := recover(); r != nil {
			errno = -C.int(syscall.EIO)
		}
	}()
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	serr := fsop.Flush(path, uint64(fi0.fh))
	return -C.int(serr)
}

//export hostRelease
func hostRelease(path0 *C.char, fi0 *C.struct_fuse_file_info) (errno C.int) {
	defer func() {
		if r := recover(); r != nil {
			errno = -C.int(syscall.EIO)
		}
	}()
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	serr := fsop.Release(path, uint64(fi0.fh))
	return -C.int(serr)
}

//export hostFsync
func hostFsync(path0 *C.char, datasync C.int, fi0 *C.struct_fuse_file_info) (errno C.int) {
	defer func() {
		if r := recover(); r != nil {
			errno = -C.int(syscall.EIO)
		}
	}()
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	serr := fsop.Fsync(path, 0 != datasync, uint64(fi0.fh))
	return -C.int(serr)
}

//export hostSetxattr
func hostSetxattr(path0 *C.char, name0 *C.char, buf0 *C.char, size0 C.size_t,
	flags C.int) (errno C.int) {
	defer func() {
		if r := recover(); r != nil {
			errno = -C.int(syscall.EIO)
		}
	}()
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	name := C.GoString(name0)
	buf := (*[1 << 30]byte)(unsafe.Pointer(buf0))
	serr := fsop.Setxattr(path, name, buf[:size0], int(flags))
	return -C.int(serr)
}

//export hostGetxattr
func hostGetxattr(path0 *C.char, name0 *C.char, buf0 *C.char, size0 C.size_t) (errno C.int) {
	defer func() {
		if r := recover(); r != nil {
			errno = -C.int(syscall.EIO)
		}
	}()
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	name := C.GoString(name0)
	buf := (*[1 << 30]byte)(unsafe.Pointer(buf0))
	serr := fsop.Getxattr(path, name, buf[:size0])
	return -C.int(serr)
}

//export hostListxattr
func hostListxattr(path0 *C.char, buf0 *C.char, size0 C.size_t) (errno C.int) {
	// !!!
	return -C.int(syscall.ENOSYS)
}

//export hostRemovexattr
func hostRemovexattr(path0 *C.char, name0 *C.char) (errno C.int) {
	defer func() {
		if r := recover(); r != nil {
			errno = -C.int(syscall.EIO)
		}
	}()
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	name := C.GoString(name0)
	serr := fsop.Removexattr(path, name)
	return -C.int(serr)
}

//export hostOpendir
func hostOpendir(path0 *C.char, fi0 *C.struct_fuse_file_info) (errno C.int) {
	defer func() {
		if r := recover(); r != nil {
			errno = -C.int(syscall.EIO)
		}
	}()
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	serr, rslt := fsop.Opendir(path)
	fi0.fh = C.uint64_t(rslt)
	return -C.int(serr)
}

//export hostReaddir
func hostReaddir(path0 *C.char, buf0 *C.char, filler0 C.fuse_fill_dir_t, off C.off_t,
	fi0 *C.struct_fuse_file_info) (errno C.int) {
	defer func() {
		if r := recover(); r != nil {
			errno = -C.int(syscall.EIO)
		}
	}()
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	serr, rslt := fsop.Readdir(path, uint64(fi0.fh))
	if serr != 0 {
		return -C.int(serr)
	}
	for _, i := range rslt {
		if C.hostFilldir(buf0, filler0, C.CString(i)) != 0 {
			break
		}
	}
	return 0
}

//export hostReleasedir
func hostReleasedir(path0 *C.char, fi0 *C.struct_fuse_file_info) (errno C.int) {
	defer func() {
		if r := recover(); r != nil {
			errno = -C.int(syscall.EIO)
		}
	}()
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	serr := fsop.Releasedir(path, uint64(fi0.fh))
	return -C.int(serr)
}

//export hostFsyncdir
func hostFsyncdir(path0 *C.char, datasync C.int, fi0 *C.struct_fuse_file_info) (errno C.int) {
	defer func() {
		if r := recover(); r != nil {
			errno = -C.int(syscall.EIO)
		}
	}()
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	serr := fsop.Fsyncdir(path, 0 != datasync, uint64(fi0.fh))
	return -C.int(serr)
}

//export hostInit
func hostInit(conn0 *C.struct_fuse_conn_info) (user_data unsafe.Pointer) {
	defer func() {
		recover()
	}()
	user_data = C.fuse_get_context().private_data
	fsop := getInterfaceForPointer(user_data).(FileSystemInterface)
	fsop.Init()
	return
}

//export hostDestroy
func hostDestroy(data0 unsafe.Pointer) {
	defer func() {
		recover()
	}()
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	fsop.Destroy()
}

//export hostAccess
func hostAccess(path0 *C.char, mask0 C.int) (errno C.int) {
	defer func() {
		if r := recover(); r != nil {
			errno = -C.int(syscall.EIO)
		}
	}()
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	serr := fsop.Access(path, int(mask0))
	return -C.int(serr)
}

//export hostCreate
func hostCreate(path0 *C.char, mode0 C.mode_t, fi0 *C.struct_fuse_file_info) (errno C.int) {
	defer func() {
		if r := recover(); r != nil {
			errno = -C.int(syscall.EIO)
		}
	}()
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	serr, rslt := fsop.Create(path, uint32(mode0))
	fi0.fh = C.uint64_t(rslt)
	return -C.int(serr)
}

//export hostFtruncate
func hostFtruncate(path0 *C.char, size0 C.off_t, fi0 *C.struct_fuse_file_info) (errno C.int) {
	defer func() {
		if r := recover(); r != nil {
			errno = -C.int(syscall.EIO)
		}
	}()
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	serr := fsop.Truncate(path, uint64(size0), uint64(fi0.fh))
	return -C.int(serr)
}

//export hostFgetattr
func hostFgetattr(path0 *C.char, stbuf0 *C.struct_stat,
	fi0 *C.struct_fuse_file_info) (errno C.int) {
	defer func() {
		if r := recover(); r != nil {
			errno = -C.int(syscall.EIO)
		}
	}()
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	serr, rslt := fsop.Getattr(path, uint64(fi0.fh))
	copyCstatFromGostat(stbuf0, rslt)
	return -C.int(serr)
}

//export hostUtimens
func hostUtimens(path0 *C.char, tv0 *C.struct_fuse_timespec) (errno C.int) {
	defer func() {
		if r := recover(); r != nil {
			errno = -C.int(syscall.EIO)
		}
	}()
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	ts := ([]syscall.Timespec)(nil)
	if tv0 != nil {
		tsar := [2]syscall.Timespec{}
		tvar := (*[2]C.struct_fuse_timespec)(unsafe.Pointer(tv0))
		copyGotimespecFromCtimespec(&tsar[0], &tvar[0])
		copyGotimespecFromCtimespec(&tsar[1], &tvar[1])
		ts = tsar[:]
	}
	serr := fsop.Utimens(path, ts)
	return -C.int(serr)
}

func NewFileSystemHost(fsop FileSystemInterface) *FileSystemHost {
	return &FileSystemHost{fsop}
}

func (host *FileSystemHost) Mount(args []string) bool {
	argv := make([]*C.char, len(args))
	for i, s := range args {
		argv[i] = C.CString(s)
	}
	p := getPointerForInterface(host.fsop)
	defer func() {
		delInterfaceFromPointer(p)
		for _, v := range argv {
			C.free(unsafe.Pointer(v))
		}
	}()
	return 0 == C.fuse_main_real(C.int(len(args)), &argv[0], C.hostFsop(), C.hostFsopSize(), p)
}

func (host *FileSystemHost) Unmount() {
	// !!!: NOTIMPL
}

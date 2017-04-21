/*
 * host.go
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

package fuse

/*
#cgo CFLAGS: -DFUSE_USE_VERSION=28 -D_FILE_OFFSET_BITS=64 -I/usr/local/include/osxfuse
#cgo LDFLAGS: -L/usr/local/lib -losxfuse

#include <stdlib.h>
#include <string.h>
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
extern int hostSymlink(char *target, char *newpath);
extern int hostRename(char *oldpath, char *newpath);
extern int hostLink(char *oldpath, char *newpath);
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
extern int hostSetxattr(char *path, char *name, char *value, size_t size, int flags);
extern int hostGetxattr(char *path, char *name, char *value, size_t size);
extern int hostListxattr(char *path, char *namebuf, size_t size);
extern int hostRemovexattr(char *path, char *name);
extern int hostOpendir(char *path, struct fuse_file_info *fi);
extern int hostReaddir(char *path, void *buf, fuse_fill_dir_t filler, fuse_off_t off,
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

static inline void hostCstatvfsFromFusestatfs(struct fuse_statvfs *stbuf,
    uint64_t bsize,
    uint64_t frsize,
    uint64_t blocks,
    uint64_t bfree,
    uint64_t bavail,
    uint64_t files,
    uint64_t ffree,
    uint64_t favail,
    uint64_t fsid,
    uint64_t flag,
    uint64_t namemax)
{
    memset(stbuf, 0, sizeof *stbuf);
    stbuf->f_bsize = bsize;
    stbuf->f_frsize = frsize;
    stbuf->f_blocks = blocks;
    stbuf->f_bfree = bfree;
    stbuf->f_bavail = bavail;
    stbuf->f_files = files;
    stbuf->f_ffree = ffree;
    stbuf->f_favail = favail;
    stbuf->f_fsid = fsid;
    stbuf->f_flag = flag;
    stbuf->f_namemax = namemax;
}

static inline void hostCstatFromFusestat(struct fuse_stat *stbuf,
    uint64_t dev,
    uint64_t ino,
    uint32_t mode,
    uint32_t nlink,
    uint32_t uid,
    uint32_t gid,
    uint64_t rdev,
    int64_t size,
    int64_t atimSec, int64_t atimNsec,
    int64_t mtimSec, int64_t mtimNsec,
    int64_t ctimSec, int64_t ctimNsec,
    int64_t blksize,
    int64_t blocks,
    int64_t birthtimSec, int64_t birthtimNsec)
{
    memset(stbuf, 0, sizeof *stbuf);
    stbuf->st_dev = dev;
    stbuf->st_ino = ino;
    stbuf->st_mode = mode;
    stbuf->st_nlink = nlink;
    stbuf->st_uid = uid;
    stbuf->st_gid = gid;
    stbuf->st_rdev = rdev;
    stbuf->st_size = size;
    stbuf->st_atime = atimSec;
    stbuf->st_mtime = mtimSec;
    stbuf->st_ctime = ctimSec;
    stbuf->st_blksize = blksize;
    stbuf->st_blocks = blocks;
}

static inline int hostFilldir(fuse_fill_dir_t filler, void *buf,
    char *name, struct fuse_stat *stbuf, fuse_off_t off)
{
    return filler(buf, name, stbuf, off);
}

#if !defined(__APPLE__)
#define _hostSetxattr hostSetxattr
#define _hostGetxattr hostGetxattr
#else
static int _hostSetxattr(char *path, char *name, char *value, size_t size, int flags,
    uint32_t position)
{
    // OSX uses position only for the resource fork; we do not support it!
    return hostSetxattr(path, name, value, size, flags);
}
static int _hostGetxattr(char *path, char *name, char *value, size_t size,
    uint32_t position)
{
    // OSX uses position only for the resource fork; we do not support it!
    return hostGetxattr(path, name, value, size);
}
#endif

static inline struct fuse_operations *hostFsop(void)
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
        .setxattr = _hostSetxattr,
        .getxattr = _hostGetxattr,
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

static inline size_t hostFsopSize(void)
{
    return sizeof(struct fuse_operations);
}
*/
import "C"
import "unsafe"

// FileSystemHost is used to host a Cgofuse file system.
type FileSystemHost struct {
	fsop FileSystemInterface
}

func copyCstatvfsFromFusestatfs(dst *C.struct_statvfs, src *Statfs_t) {
	C.hostCstatvfsFromFusestatfs(dst,
		C.uint64_t(src.Bsize),
		C.uint64_t(src.Frsize),
		C.uint64_t(src.Blocks),
		C.uint64_t(src.Bfree),
		C.uint64_t(src.Bavail),
		C.uint64_t(src.Files),
		C.uint64_t(src.Ffree),
		C.uint64_t(src.Favail),
		C.uint64_t(src.Fsid),
		C.uint64_t(src.Flag),
		C.uint64_t(src.Namemax))
}

func copyCstatFromFusestat(dst *C.struct_stat, src *Stat_t) {
	C.hostCstatFromFusestat(dst,
		C.uint64_t(src.Dev),
		C.uint64_t(src.Ino),
		C.uint32_t(src.Mode),
		C.uint32_t(src.Nlink),
		C.uint32_t(src.Uid),
		C.uint32_t(src.Gid),
		C.uint64_t(src.Rdev),
		C.int64_t(src.Size),
		C.int64_t(src.Atim.Sec), C.int64_t(src.Atim.Nsec),
		C.int64_t(src.Mtim.Sec), C.int64_t(src.Mtim.Nsec),
		C.int64_t(src.Ctim.Sec), C.int64_t(src.Ctim.Nsec),
		C.int64_t(src.Blksize),
		C.int64_t(src.Blocks),
		C.int64_t(src.Birthtim.Sec), C.int64_t(src.Birthtim.Nsec))
}

func copyFusetimespecFromCtimespec(dst *Timespec, src *C.struct_timespec) {
	dst.Sec = int64(src.tv_sec)
	dst.Nsec = int64(src.tv_nsec)
}

func recoverAsErrno(errc0 *C.int) {
	if r := recover(); nil != r {
		*errc0 = -C.int(EIO)
	}
}

//export hostGetattr
func hostGetattr(path0 *C.char, stat0 *C.struct_stat) (errc0 C.int) {
	defer recoverAsErrno(&errc0)
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	stat := &Stat_t{}
	errc := fsop.Getattr(path, stat, ^uint64(0))
	copyCstatFromFusestat(stat0, stat)
	return C.int(errc)
}

//export hostReadlink
func hostReadlink(path0 *C.char, buff0 *C.char, size0 C.size_t) (errc0 C.int) {
	defer recoverAsErrno(&errc0)
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	errc, rslt := fsop.Readlink(path)
	buff := (*[1 << 30]byte)(unsafe.Pointer(buff0))
	copy(buff[:size0-1], rslt)
	rlen := len(rslt)
	if C.size_t(rlen) < size0 {
		buff[rlen] = 0
	}
	return C.int(errc)
}

//export hostMknod
func hostMknod(path0 *C.char, mode0 C.mode_t, dev0 C.dev_t) (errc0 C.int) {
	defer recoverAsErrno(&errc0)
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	errc := fsop.Mknod(path, uint32(mode0), uint64(dev0))
	return C.int(errc)
}

//export hostMkdir
func hostMkdir(path0 *C.char, mode0 C.mode_t) (errc0 C.int) {
	defer recoverAsErrno(&errc0)
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	errc := fsop.Mkdir(path, uint32(mode0))
	return C.int(errc)
}

//export hostUnlink
func hostUnlink(path0 *C.char) (errc0 C.int) {
	defer recoverAsErrno(&errc0)
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	errc := fsop.Unlink(path)
	return C.int(errc)
}

//export hostRmdir
func hostRmdir(path0 *C.char) (errc0 C.int) {
	defer recoverAsErrno(&errc0)
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	errc := fsop.Rmdir(path)
	return C.int(errc)
}

//export hostSymlink
func hostSymlink(target0 *C.char, newpath0 *C.char) (errc0 C.int) {
	defer recoverAsErrno(&errc0)
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	target, newpath := C.GoString(target0), C.GoString(newpath0)
	errc := fsop.Symlink(target, newpath)
	return C.int(errc)
}

//export hostRename
func hostRename(oldpath0 *C.char, newpath0 *C.char) (errc0 C.int) {
	defer recoverAsErrno(&errc0)
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	oldpath, newpath := C.GoString(oldpath0), C.GoString(newpath0)
	errc := fsop.Rename(oldpath, newpath)
	return C.int(errc)
}

//export hostLink
func hostLink(oldpath0 *C.char, newpath0 *C.char) (errc0 C.int) {
	defer recoverAsErrno(&errc0)
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	oldpath, newpath := C.GoString(oldpath0), C.GoString(newpath0)
	errc := fsop.Link(oldpath, newpath)
	return C.int(errc)
}

//export hostChmod
func hostChmod(path0 *C.char, mode0 C.mode_t) (errc0 C.int) {
	defer recoverAsErrno(&errc0)
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	errc := fsop.Chmod(path, uint32(mode0))
	return C.int(errc)
}

//export hostChown
func hostChown(path0 *C.char, uid0 C.uid_t, gid0 C.gid_t) (errc0 C.int) {
	defer recoverAsErrno(&errc0)
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	errc := fsop.Chown(path, uint32(uid0), uint32(gid0))
	return C.int(errc)
}

//export hostTruncate
func hostTruncate(path0 *C.char, size0 C.off_t) (errc0 C.int) {
	defer recoverAsErrno(&errc0)
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	errc := fsop.Truncate(path, int64(size0), ^uint64(0))
	return C.int(errc)
}

//export hostOpen
func hostOpen(path0 *C.char, fi0 *C.struct_fuse_file_info) (errc0 C.int) {
	defer recoverAsErrno(&errc0)
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	errc, rslt := fsop.Open(path, int(fi0.flags))
	fi0.fh = C.uint64_t(rslt)
	return C.int(errc)
}

//export hostRead
func hostRead(path0 *C.char, buff0 *C.char, size0 C.size_t, ofst0 C.off_t,
	fi0 *C.struct_fuse_file_info) (nbyt0 C.int) {
	defer recoverAsErrno(&nbyt0)
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	buff := (*[1 << 30]byte)(unsafe.Pointer(buff0))
	nbyt := fsop.Read(path, buff[:size0], int64(ofst0), uint64(fi0.fh))
	return C.int(nbyt)
}

//export hostWrite
func hostWrite(path0 *C.char, buff0 *C.char, size0 C.size_t, ofst0 C.off_t,
	fi0 *C.struct_fuse_file_info) (nbyt0 C.int) {
	defer recoverAsErrno(&nbyt0)
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	buff := (*[1 << 30]byte)(unsafe.Pointer(buff0))
	nbyt := fsop.Write(path, buff[:size0], int64(ofst0), uint64(fi0.fh))
	return C.int(nbyt)
}

//export hostStatfs
func hostStatfs(path0 *C.char, stat0 *C.struct_statvfs) (errc0 C.int) {
	defer recoverAsErrno(&errc0)
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	stat := &Statfs_t{}
	errc := fsop.Statfs(path, stat)
	copyCstatvfsFromFusestatfs(stat0, stat)
	return C.int(errc)
}

//export hostFlush
func hostFlush(path0 *C.char, fi0 *C.struct_fuse_file_info) (errc0 C.int) {
	defer recoverAsErrno(&errc0)
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	errc := fsop.Flush(path, uint64(fi0.fh))
	return C.int(errc)
}

//export hostRelease
func hostRelease(path0 *C.char, fi0 *C.struct_fuse_file_info) (errc0 C.int) {
	defer recoverAsErrno(&errc0)
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	errc := fsop.Release(path, uint64(fi0.fh))
	return C.int(errc)
}

//export hostFsync
func hostFsync(path0 *C.char, datasync C.int, fi0 *C.struct_fuse_file_info) (errc0 C.int) {
	defer recoverAsErrno(&errc0)
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	errc := fsop.Fsync(path, 0 != datasync, uint64(fi0.fh))
	return C.int(errc)
}

//export hostSetxattr
func hostSetxattr(path0 *C.char, name0 *C.char, buff0 *C.char, size0 C.size_t,
	flags C.int) (errc0 C.int) {
	defer recoverAsErrno(&errc0)
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	name := C.GoString(name0)
	buff := (*[1 << 30]byte)(unsafe.Pointer(buff0))
	errc := fsop.Setxattr(path, name, buff[:size0], int(flags))
	return C.int(errc)
}

//export hostGetxattr
func hostGetxattr(path0 *C.char, name0 *C.char, buff0 *C.char, size0 C.size_t) (nbyt0 C.int) {
	defer recoverAsErrno(&nbyt0)
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	name := C.GoString(name0)
	buff := (*[1 << 30]byte)(unsafe.Pointer(buff0))
	nbyt := fsop.Getxattr(path, name, buff[:size0])
	return C.int(nbyt)
}

//export hostListxattr
func hostListxattr(path0 *C.char, buff0 *C.char, size0 C.size_t) (nbyt0 C.int) {
	defer recoverAsErrno(&nbyt0)
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	buff := (*[1 << 30]byte)(unsafe.Pointer(buff0))
	size := int(size0)
	nbyt := 0
	fill := func(name1 string) bool {
		nlen := len(name1)
		if 0 != size {
			if nbyt+nlen+1 > size {
				return false
			}
			copy(buff[nbyt:nbyt+nlen], name1)
			buff[nbyt+nlen] = 0
		}
		nbyt += nlen + 1
		return true
	}
	errc := fsop.Listxattr(path, fill)
	if 0 != errc {
		return C.int(errc)
	}
	return C.int(nbyt)
}

//export hostRemovexattr
func hostRemovexattr(path0 *C.char, name0 *C.char) (errc0 C.int) {
	defer recoverAsErrno(&errc0)
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	name := C.GoString(name0)
	errc := fsop.Removexattr(path, name)
	return C.int(errc)
}

//export hostOpendir
func hostOpendir(path0 *C.char, fi0 *C.struct_fuse_file_info) (errc0 C.int) {
	defer recoverAsErrno(&errc0)
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	errc, rslt := fsop.Opendir(path)
	fi0.fh = C.uint64_t(rslt)
	return C.int(errc)
}

//export hostReaddir
func hostReaddir(path0 *C.char, buff0 unsafe.Pointer, fill0 C.fuse_fill_dir_t, ofst0 C.off_t,
	fi0 *C.struct_fuse_file_info) (errc0 C.int) {
	defer recoverAsErrno(&errc0)
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	fill := func(name1 string, stat1 *Stat_t, off1 int64) bool {
		name := C.CString(name1)
		defer C.free(unsafe.Pointer(name))
		if nil == stat1 {
			return 0 == C.hostFilldir(fill0, buff0, name, nil, C.off_t(off1))
		} else {
			stat := C.struct_stat{}
			copyCstatFromFusestat(&stat, stat1)
			return 0 == C.hostFilldir(fill0, buff0, name, &stat, C.off_t(off1))
		}
	}
	errc := fsop.Readdir(path, fill, int64(ofst0), uint64(fi0.fh))
	return C.int(errc)
}

//export hostReleasedir
func hostReleasedir(path0 *C.char, fi0 *C.struct_fuse_file_info) (errc0 C.int) {
	defer recoverAsErrno(&errc0)
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	errc := fsop.Releasedir(path, uint64(fi0.fh))
	return C.int(errc)
}

//export hostFsyncdir
func hostFsyncdir(path0 *C.char, datasync C.int, fi0 *C.struct_fuse_file_info) (errc0 C.int) {
	defer recoverAsErrno(&errc0)
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	errc := fsop.Fsyncdir(path, 0 != datasync, uint64(fi0.fh))
	return C.int(errc)
}

//export hostInit
func hostInit(conn0 *C.struct_fuse_conn_info) (user_data unsafe.Pointer) {
	defer recover()
	user_data = C.fuse_get_context().private_data
	fsop := getInterfaceForPointer(user_data).(FileSystemInterface)
	fsop.Init()
	return
}

//export hostDestroy
func hostDestroy(data0 unsafe.Pointer) {
	defer recover()
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	fsop.Destroy()
}

//export hostAccess
func hostAccess(path0 *C.char, mask0 C.int) (errc0 C.int) {
	defer recoverAsErrno(&errc0)
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	errc := fsop.Access(path, uint32(mask0))
	return C.int(errc)
}

//export hostCreate
func hostCreate(path0 *C.char, mode0 C.mode_t, fi0 *C.struct_fuse_file_info) (errc0 C.int) {
	defer recoverAsErrno(&errc0)
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	errc, rslt := fsop.Create(path, uint32(mode0))
	fi0.fh = C.uint64_t(rslt)
	return C.int(errc)
}

//export hostFtruncate
func hostFtruncate(path0 *C.char, size0 C.off_t, fi0 *C.struct_fuse_file_info) (errc0 C.int) {
	defer recoverAsErrno(&errc0)
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	errc := fsop.Truncate(path, int64(size0), uint64(fi0.fh))
	return C.int(errc)
}

//export hostFgetattr
func hostFgetattr(path0 *C.char, stat0 *C.struct_stat,
	fi0 *C.struct_fuse_file_info) (errc0 C.int) {
	defer recoverAsErrno(&errc0)
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	stat := &Stat_t{}
	errc := fsop.Getattr(path, stat, uint64(fi0.fh))
	copyCstatFromFusestat(stat0, stat)
	return C.int(errc)
}

//export hostUtimens
func hostUtimens(path0 *C.char, tmsp0 *C.struct_fuse_timespec) (errc0 C.int) {
	defer recoverAsErrno(&errc0)
	fsop := getInterfaceForPointer(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	if nil == tmsp0 {
		errc := fsop.Utimens(path, nil)
		return C.int(errc)
	} else {
		tmsp := [2]Timespec{}
		tmsa := (*[2]C.struct_fuse_timespec)(unsafe.Pointer(tmsp0))
		copyFusetimespecFromCtimespec(&tmsp[0], &tmsa[0])
		copyFusetimespecFromCtimespec(&tmsp[1], &tmsa[1])
		errc := fsop.Utimens(path, tmsp[:])
		return C.int(errc)
	}
}

// NewFileSystemHost creates a file system host.
func NewFileSystemHost(fsop FileSystemInterface) *FileSystemHost {
	return &FileSystemHost{fsop}
}

// Mount mounts a file system.
func (host *FileSystemHost) Mount(args []string) bool {
	argc := len(args) + 1
	argv := make([]*C.char, argc+1)
	argv[0] = C.CString(args[0])
	defer C.free(unsafe.Pointer(argv[0]))
	argv[1] = C.CString("-f") // do not daemonize; Go cannot handle it (at least on OSX)
	defer C.free(unsafe.Pointer(argv[1]))
	for i := 1; len(args) > i; i++ {
		argv[i+1] = C.CString(args[i])
		defer C.free(unsafe.Pointer(argv[i+1]))
	}
	p := newPointerForInterface(host.fsop)
	defer delPointerForInterface(p)
	return 0 == C.fuse_main_real(C.int(argc), &argv[0], C.hostFsop(), C.hostFsopSize(), p)
}

// Mount unmounts a file system.
func (host *FileSystemHost) Unmount() {
	// !!!: NOTIMPL
}

// Getcontext gets information related to a file system operation.
func Getcontext() (uid uint32, gid uint32, pid int) {
	uid = uint32(C.fuse_get_context().uid)
	gid = uint32(C.fuse_get_context().gid)
	pid = int(C.fuse_get_context().pid)
	return
}

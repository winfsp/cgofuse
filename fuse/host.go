/*
 * host.go
 *
 * Copyright 2017-2018 Bill Zissimopoulos
 */
/*
 * This file is part of Cgofuse.
 *
 * It is licensed under the MIT license. The full license text can be found
 * in the License.txt file at the root of this project.
 */

package fuse

import (
	"errors"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"unsafe"
)

// FileSystemHost is used to host a file system.
type FileSystemHost struct {
	fsop FileSystemInterface
	fuse *C_struct_fuse
	mntp *C_char
	sigc chan os.Signal

	capCaseInsensitive, capReaddirPlus bool
}

var (
	hostGuard = sync.Mutex{}
	hostTable = map[unsafe.Pointer]*FileSystemHost{}
)

func hostHandleNew(host *FileSystemHost) unsafe.Pointer {
	p := C_malloc(1)
	hostGuard.Lock()
	defer hostGuard.Unlock()
	hostTable[p] = host
	return p
}

func hostHandleDel(p unsafe.Pointer) *FileSystemHost {
	hostGuard.Lock()
	defer hostGuard.Unlock()
	if host, ok := hostTable[p]; ok {
		delete(hostTable, p)
		C_free(p)
		return host
	}
	return nil
}

func hostHandleGet(p unsafe.Pointer) *FileSystemHost {
	hostGuard.Lock()
	defer hostGuard.Unlock()
	if host, ok := hostTable[p]; ok {
		return host
	}
	return nil
}

func copyCstatvfsFromFusestatfs(dst *C_fuse_statvfs_t, src *Statfs_t) {
	C_hostCstatvfsFromFusestatfs(dst,
		C_uint64_t(src.Bsize),
		C_uint64_t(src.Frsize),
		C_uint64_t(src.Blocks),
		C_uint64_t(src.Bfree),
		C_uint64_t(src.Bavail),
		C_uint64_t(src.Files),
		C_uint64_t(src.Ffree),
		C_uint64_t(src.Favail),
		C_uint64_t(src.Fsid),
		C_uint64_t(src.Flag),
		C_uint64_t(src.Namemax))
}

func copyCstatFromFusestat(dst *C_fuse_stat_t, src *Stat_t) {
	C_hostCstatFromFusestat(dst,
		C_uint64_t(src.Dev),
		C_uint64_t(src.Ino),
		C_uint32_t(src.Mode),
		C_uint32_t(src.Nlink),
		C_uint32_t(src.Uid),
		C_uint32_t(src.Gid),
		C_uint64_t(src.Rdev),
		C_int64_t(src.Size),
		C_int64_t(src.Atim.Sec), C_int64_t(src.Atim.Nsec),
		C_int64_t(src.Mtim.Sec), C_int64_t(src.Mtim.Nsec),
		C_int64_t(src.Ctim.Sec), C_int64_t(src.Ctim.Nsec),
		C_int64_t(src.Blksize),
		C_int64_t(src.Blocks),
		C_int64_t(src.Birthtim.Sec), C_int64_t(src.Birthtim.Nsec),
		C_uint32_t(src.Flags))
}

func copyFusetimespecFromCtimespec(dst *Timespec, src *C_fuse_timespec_t) {
	dst.Sec = int64(src.tv_sec)
	dst.Nsec = int64(src.tv_nsec)
}

func recoverAsErrno(errc0 *C_int) {
	if r := recover(); nil != r {
		switch e := r.(type) {
		case Error:
			*errc0 = C_int(e)
		default:
			*errc0 = -C_int(EIO)
		}
	}
}

func hostGetattr(path0 *C_char, stat0 *C_fuse_stat_t) (errc0 C_int) {
	defer recoverAsErrno(&errc0)
	fsop := hostHandleGet(C_fuse_get_context().private_data).fsop
	path := C_GoString(path0)
	stat := &Stat_t{}
	errc := fsop.Getattr(path, stat, ^uint64(0))
	copyCstatFromFusestat(stat0, stat)
	return C_int(errc)
}

func hostReadlink(path0 *C_char, buff0 *C_char, size0 C_size_t) (errc0 C_int) {
	defer recoverAsErrno(&errc0)
	fsop := hostHandleGet(C_fuse_get_context().private_data).fsop
	path := C_GoString(path0)
	errc, rslt := fsop.Readlink(path)
	buff := (*[1 << 30]byte)(unsafe.Pointer(buff0))
	copy(buff[:size0-1], rslt)
	rlen := len(rslt)
	if C_size_t(rlen) < size0 {
		buff[rlen] = 0
	}
	return C_int(errc)
}

func hostMknod(path0 *C_char, mode0 C_fuse_mode_t, dev0 C_fuse_dev_t) (errc0 C_int) {
	defer recoverAsErrno(&errc0)
	fsop := hostHandleGet(C_fuse_get_context().private_data).fsop
	path := C_GoString(path0)
	errc := fsop.Mknod(path, uint32(mode0), uint64(dev0))
	return C_int(errc)
}

func hostMkdir(path0 *C_char, mode0 C_fuse_mode_t) (errc0 C_int) {
	defer recoverAsErrno(&errc0)
	fsop := hostHandleGet(C_fuse_get_context().private_data).fsop
	path := C_GoString(path0)
	errc := fsop.Mkdir(path, uint32(mode0))
	return C_int(errc)
}

func hostUnlink(path0 *C_char) (errc0 C_int) {
	defer recoverAsErrno(&errc0)
	fsop := hostHandleGet(C_fuse_get_context().private_data).fsop
	path := C_GoString(path0)
	errc := fsop.Unlink(path)
	return C_int(errc)
}

func hostRmdir(path0 *C_char) (errc0 C_int) {
	defer recoverAsErrno(&errc0)
	fsop := hostHandleGet(C_fuse_get_context().private_data).fsop
	path := C_GoString(path0)
	errc := fsop.Rmdir(path)
	return C_int(errc)
}

func hostSymlink(target0 *C_char, newpath0 *C_char) (errc0 C_int) {
	defer recoverAsErrno(&errc0)
	fsop := hostHandleGet(C_fuse_get_context().private_data).fsop
	target, newpath := C_GoString(target0), C_GoString(newpath0)
	errc := fsop.Symlink(target, newpath)
	return C_int(errc)
}

func hostRename(oldpath0 *C_char, newpath0 *C_char) (errc0 C_int) {
	defer recoverAsErrno(&errc0)
	fsop := hostHandleGet(C_fuse_get_context().private_data).fsop
	oldpath, newpath := C_GoString(oldpath0), C_GoString(newpath0)
	errc := fsop.Rename(oldpath, newpath)
	return C_int(errc)
}

func hostLink(oldpath0 *C_char, newpath0 *C_char) (errc0 C_int) {
	defer recoverAsErrno(&errc0)
	fsop := hostHandleGet(C_fuse_get_context().private_data).fsop
	oldpath, newpath := C_GoString(oldpath0), C_GoString(newpath0)
	errc := fsop.Link(oldpath, newpath)
	return C_int(errc)
}

func hostChmod(path0 *C_char, mode0 C_fuse_mode_t) (errc0 C_int) {
	defer recoverAsErrno(&errc0)
	fsop := hostHandleGet(C_fuse_get_context().private_data).fsop
	path := C_GoString(path0)
	errc := fsop.Chmod(path, uint32(mode0))
	return C_int(errc)
}

func hostChown(path0 *C_char, uid0 C_fuse_uid_t, gid0 C_fuse_gid_t) (errc0 C_int) {
	defer recoverAsErrno(&errc0)
	fsop := hostHandleGet(C_fuse_get_context().private_data).fsop
	path := C_GoString(path0)
	errc := fsop.Chown(path, uint32(uid0), uint32(gid0))
	return C_int(errc)
}

func hostTruncate(path0 *C_char, size0 C_fuse_off_t) (errc0 C_int) {
	defer recoverAsErrno(&errc0)
	fsop := hostHandleGet(C_fuse_get_context().private_data).fsop
	path := C_GoString(path0)
	errc := fsop.Truncate(path, int64(size0), ^uint64(0))
	return C_int(errc)
}

func hostOpen(path0 *C_char, fi0 *C_struct_fuse_file_info) (errc0 C_int) {
	defer recoverAsErrno(&errc0)
	fsop := hostHandleGet(C_fuse_get_context().private_data).fsop
	path := C_GoString(path0)
	errc, rslt := fsop.Open(path, int(fi0.flags))
	fi0.fh = C_uint64_t(rslt)
	return C_int(errc)
}

func hostRead(path0 *C_char, buff0 *C_char, size0 C_size_t, ofst0 C_fuse_off_t,
	fi0 *C_struct_fuse_file_info) (nbyt0 C_int) {
	defer recoverAsErrno(&nbyt0)
	fsop := hostHandleGet(C_fuse_get_context().private_data).fsop
	path := C_GoString(path0)
	buff := (*[1 << 30]byte)(unsafe.Pointer(buff0))
	nbyt := fsop.Read(path, buff[:size0], int64(ofst0), uint64(fi0.fh))
	return C_int(nbyt)
}

func hostWrite(path0 *C_char, buff0 *C_char, size0 C_size_t, ofst0 C_fuse_off_t,
	fi0 *C_struct_fuse_file_info) (nbyt0 C_int) {
	defer recoverAsErrno(&nbyt0)
	fsop := hostHandleGet(C_fuse_get_context().private_data).fsop
	path := C_GoString(path0)
	buff := (*[1 << 30]byte)(unsafe.Pointer(buff0))
	nbyt := fsop.Write(path, buff[:size0], int64(ofst0), uint64(fi0.fh))
	return C_int(nbyt)
}

func hostStatfs(path0 *C_char, stat0 *C_fuse_statvfs_t) (errc0 C_int) {
	defer recoverAsErrno(&errc0)
	fsop := hostHandleGet(C_fuse_get_context().private_data).fsop
	path := C_GoString(path0)
	stat := &Statfs_t{}
	errc := fsop.Statfs(path, stat)
	if -ENOSYS == errc {
		stat = &Statfs_t{}
		errc = 0
	}
	copyCstatvfsFromFusestatfs(stat0, stat)
	return C_int(errc)
}

func hostFlush(path0 *C_char, fi0 *C_struct_fuse_file_info) (errc0 C_int) {
	defer recoverAsErrno(&errc0)
	fsop := hostHandleGet(C_fuse_get_context().private_data).fsop
	path := C_GoString(path0)
	errc := fsop.Flush(path, uint64(fi0.fh))
	return C_int(errc)
}

func hostRelease(path0 *C_char, fi0 *C_struct_fuse_file_info) (errc0 C_int) {
	defer recoverAsErrno(&errc0)
	fsop := hostHandleGet(C_fuse_get_context().private_data).fsop
	path := C_GoString(path0)
	errc := fsop.Release(path, uint64(fi0.fh))
	return C_int(errc)
}

func hostFsync(path0 *C_char, datasync C_int, fi0 *C_struct_fuse_file_info) (errc0 C_int) {
	defer recoverAsErrno(&errc0)
	fsop := hostHandleGet(C_fuse_get_context().private_data).fsop
	path := C_GoString(path0)
	errc := fsop.Fsync(path, 0 != datasync, uint64(fi0.fh))
	if -ENOSYS == errc {
		errc = 0
	}
	return C_int(errc)
}

func hostSetxattr(path0 *C_char, name0 *C_char, buff0 *C_char, size0 C_size_t,
	flags C_int) (errc0 C_int) {
	defer recoverAsErrno(&errc0)
	fsop := hostHandleGet(C_fuse_get_context().private_data).fsop
	path := C_GoString(path0)
	name := C_GoString(name0)
	buff := (*[1 << 30]byte)(unsafe.Pointer(buff0))
	errc := fsop.Setxattr(path, name, buff[:size0], int(flags))
	return C_int(errc)
}

func hostGetxattr(path0 *C_char, name0 *C_char, buff0 *C_char, size0 C_size_t) (nbyt0 C_int) {
	defer recoverAsErrno(&nbyt0)
	fsop := hostHandleGet(C_fuse_get_context().private_data).fsop
	path := C_GoString(path0)
	name := C_GoString(name0)
	errc, rslt := fsop.Getxattr(path, name)
	if 0 != errc {
		return C_int(errc)
	}
	if 0 != size0 {
		if len(rslt) > int(size0) {
			return -C_int(ERANGE)
		}
		buff := (*[1 << 30]byte)(unsafe.Pointer(buff0))
		copy(buff[:size0], rslt)
	}
	return C_int(len(rslt))
}

func hostListxattr(path0 *C_char, buff0 *C_char, size0 C_size_t) (nbyt0 C_int) {
	defer recoverAsErrno(&nbyt0)
	fsop := hostHandleGet(C_fuse_get_context().private_data).fsop
	path := C_GoString(path0)
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
		return C_int(errc)
	}
	return C_int(nbyt)
}

func hostRemovexattr(path0 *C_char, name0 *C_char) (errc0 C_int) {
	defer recoverAsErrno(&errc0)
	fsop := hostHandleGet(C_fuse_get_context().private_data).fsop
	path := C_GoString(path0)
	name := C_GoString(name0)
	errc := fsop.Removexattr(path, name)
	return C_int(errc)
}

func hostOpendir(path0 *C_char, fi0 *C_struct_fuse_file_info) (errc0 C_int) {
	defer recoverAsErrno(&errc0)
	fsop := hostHandleGet(C_fuse_get_context().private_data).fsop
	path := C_GoString(path0)
	errc, rslt := fsop.Opendir(path)
	if -ENOSYS == errc {
		errc = 0
	}
	fi0.fh = C_uint64_t(rslt)
	return C_int(errc)
}

func hostReaddir(path0 *C_char, buff0 unsafe.Pointer, fill0 C_fuse_fill_dir_t, ofst0 C_fuse_off_t,
	fi0 *C_struct_fuse_file_info) (errc0 C_int) {
	defer recoverAsErrno(&errc0)
	fsop := hostHandleGet(C_fuse_get_context().private_data).fsop
	path := C_GoString(path0)
	fill := func(name1 string, stat1 *Stat_t, off1 int64) bool {
		name := C_CString(name1)
		defer C_free(unsafe.Pointer(name))
		if nil == stat1 {
			return 0 == C_hostFilldir(fill0, buff0, name, nil, C_fuse_off_t(off1))
		} else {
			stat := C_fuse_stat_t{}
			copyCstatFromFusestat(&stat, stat1)
			return 0 == C_hostFilldir(fill0, buff0, name, &stat, C_fuse_off_t(off1))
		}
	}
	errc := fsop.Readdir(path, fill, int64(ofst0), uint64(fi0.fh))
	return C_int(errc)
}

func hostReleasedir(path0 *C_char, fi0 *C_struct_fuse_file_info) (errc0 C_int) {
	defer recoverAsErrno(&errc0)
	fsop := hostHandleGet(C_fuse_get_context().private_data).fsop
	path := C_GoString(path0)
	errc := fsop.Releasedir(path, uint64(fi0.fh))
	return C_int(errc)
}

func hostFsyncdir(path0 *C_char, datasync C_int, fi0 *C_struct_fuse_file_info) (errc0 C_int) {
	defer recoverAsErrno(&errc0)
	fsop := hostHandleGet(C_fuse_get_context().private_data).fsop
	path := C_GoString(path0)
	errc := fsop.Fsyncdir(path, 0 != datasync, uint64(fi0.fh))
	if -ENOSYS == errc {
		errc = 0
	}
	return C_int(errc)
}

func hostInit(conn0 *C_struct_fuse_conn_info) (user_data unsafe.Pointer) {
	defer recover()
	fctx := C_fuse_get_context()
	user_data = fctx.private_data
	host := hostHandleGet(user_data)
	host.fuse = fctx.fuse
	C_hostAsgnCconninfo(conn0,
		C_bool(host.capCaseInsensitive),
		C_bool(host.capReaddirPlus))
	if nil != host.sigc {
		signal.Notify(host.sigc, syscall.SIGINT, syscall.SIGTERM)
	}
	host.fsop.Init()
	return
}

func hostDestroy(user_data unsafe.Pointer) {
	defer recover()
	host := hostHandleGet(user_data)
	host.fsop.Destroy()
	if nil != host.sigc {
		signal.Stop(host.sigc)
	}
	host.fuse = nil
}

func hostAccess(path0 *C_char, mask0 C_int) (errc0 C_int) {
	defer recoverAsErrno(&errc0)
	fsop := hostHandleGet(C_fuse_get_context().private_data).fsop
	path := C_GoString(path0)
	errc := fsop.Access(path, uint32(mask0))
	return C_int(errc)
}

func hostCreate(path0 *C_char, mode0 C_fuse_mode_t, fi0 *C_struct_fuse_file_info) (errc0 C_int) {
	defer recoverAsErrno(&errc0)
	fsop := hostHandleGet(C_fuse_get_context().private_data).fsop
	path := C_GoString(path0)
	errc, rslt := fsop.Create(path, int(fi0.flags), uint32(mode0))
	if -ENOSYS == errc {
		errc = fsop.Mknod(path, S_IFREG|uint32(mode0), 0)
		if 0 == errc {
			errc, rslt = fsop.Open(path, int(fi0.flags))
		}
	}
	fi0.fh = C_uint64_t(rslt)
	return C_int(errc)
}

func hostFtruncate(path0 *C_char, size0 C_fuse_off_t, fi0 *C_struct_fuse_file_info) (errc0 C_int) {
	defer recoverAsErrno(&errc0)
	fsop := hostHandleGet(C_fuse_get_context().private_data).fsop
	path := C_GoString(path0)
	errc := fsop.Truncate(path, int64(size0), uint64(fi0.fh))
	return C_int(errc)
}

func hostFgetattr(path0 *C_char, stat0 *C_fuse_stat_t,
	fi0 *C_struct_fuse_file_info) (errc0 C_int) {
	defer recoverAsErrno(&errc0)
	fsop := hostHandleGet(C_fuse_get_context().private_data).fsop
	path := C_GoString(path0)
	stat := &Stat_t{}
	errc := fsop.Getattr(path, stat, uint64(fi0.fh))
	copyCstatFromFusestat(stat0, stat)
	return C_int(errc)
}

func hostUtimens(path0 *C_char, tmsp0 *C_fuse_timespec_t) (errc0 C_int) {
	defer recoverAsErrno(&errc0)
	fsop := hostHandleGet(C_fuse_get_context().private_data).fsop
	path := C_GoString(path0)
	if nil == tmsp0 {
		errc := fsop.Utimens(path, nil)
		return C_int(errc)
	} else {
		tmsp := [2]Timespec{}
		tmsa := (*[2]C_fuse_timespec_t)(unsafe.Pointer(tmsp0))
		copyFusetimespecFromCtimespec(&tmsp[0], &tmsa[0])
		copyFusetimespecFromCtimespec(&tmsp[1], &tmsa[1])
		errc := fsop.Utimens(path, tmsp[:])
		return C_int(errc)
	}
}

func hostSetchgtime(path0 *C_char, tmsp0 *C_fuse_timespec_t) (errc0 C_int) {
	defer recoverAsErrno(&errc0)
	fsop := hostHandleGet(C_fuse_get_context().private_data).fsop
	intf, ok := fsop.(FileSystemSetchgtime)
	if !ok {
		// say we did it!
		return 0
	}
	path := C_GoString(path0)
	tmsp := Timespec{}
	copyFusetimespecFromCtimespec(&tmsp, tmsp0)
	errc := intf.Setchgtime(path, tmsp)
	return C_int(errc)
}

func hostSetcrtime(path0 *C_char, tmsp0 *C_fuse_timespec_t) (errc0 C_int) {
	defer recoverAsErrno(&errc0)
	fsop := hostHandleGet(C_fuse_get_context().private_data).fsop
	intf, ok := fsop.(FileSystemSetcrtime)
	if !ok {
		// say we did it!
		return 0
	}
	path := C_GoString(path0)
	tmsp := Timespec{}
	copyFusetimespecFromCtimespec(&tmsp, tmsp0)
	errc := intf.Setcrtime(path, tmsp)
	return C_int(errc)
}

func hostChflags(path0 *C_char, flags C_uint32_t) (errc0 C_int) {
	defer recoverAsErrno(&errc0)
	fsop := hostHandleGet(C_fuse_get_context().private_data).fsop
	intf, ok := fsop.(FileSystemChflags)
	if !ok {
		// say we did it!
		return 0
	}
	path := C_GoString(path0)
	errc := intf.Chflags(path, uint32(flags))
	return C_int(errc)
}

// NewFileSystemHost creates a file system host.
func NewFileSystemHost(fsop FileSystemInterface) *FileSystemHost {
	host := &FileSystemHost{}
	host.fsop = fsop
	return host
}

// SetCapCaseInsensitive informs the host that the hosted file system is case insensitive
// [OSX and Windows only].
func (host *FileSystemHost) SetCapCaseInsensitive(value bool) {
	host.capCaseInsensitive = value
}

// SetCapReaddirPlus informs the host that the hosted file system has the readdir-plus
// capability [Windows only]. A file system that has the readdir-plus capability can send
// full stat information during Readdir, thus avoiding extraneous Getattr calls.
func (host *FileSystemHost) SetCapReaddirPlus(value bool) {
	host.capReaddirPlus = value
}

// Mount mounts a file system on the given mountpoint with the mount options in opts.
//
// Many of the mount options in opts are specific to the underlying FUSE implementation.
// Some of the common options include:
//
//     -h   --help            print help
//     -V   --version         print FUSE version
//     -d   -o debug          enable FUSE debug output
//     -s                     disable multi-threaded operation
//
// Please refer to the individual FUSE implementation documentation for additional options.
//
// It is allowed for the mountpoint to be the empty string ("") in which case opts is assumed
// to contain the mountpoint. It is also allowed for opts to be nil, although in this case the
// mountpoint must be non-empty.
func (host *FileSystemHost) Mount(mountpoint string, opts []string) bool {
	if 0 == C_hostFuseInit() {
		panic("cgofuse: cannot find winfsp")
	}

	/*
	 * Command line handling
	 *
	 * We must prepare a command line to send to FUSE. This command line will look like this:
	 *
	 *     execname [mountpoint] "-f" [opts...] NULL
	 *
	 * We add the "-f" option because Go cannot handle daemonization (at least on OSX).
	 */
	exec := "<UNKNOWN>"
	if 0 < len(os.Args) {
		exec = os.Args[0]
	}
	argc := len(opts) + 2
	if "" != mountpoint {
		argc++
	}
	argv := make([]*C_char, argc+1)
	argv[0] = C_CString(exec)
	defer C_free(unsafe.Pointer(argv[0]))
	opti := 1
	if "" != mountpoint {
		argv[1] = C_CString(mountpoint)
		defer C_free(unsafe.Pointer(argv[1]))
		opti++
	}
	argv[opti] = C_CString("-f")
	defer C_free(unsafe.Pointer(argv[opti]))
	opti++
	for i := 0; len(opts) > i; i++ {
		argv[i+opti] = C_CString(opts[i])
		defer C_free(unsafe.Pointer(argv[i+opti]))
	}

	/*
	 * Mountpoint extraction
	 *
	 * We need to determine the mountpoint that FUSE is going (to try) to use, so that we
	 * can unmount later.
	 */
	host.mntp = C_hostMountpoint(C_int(argc), &argv[0])
	defer func() {
		C_free(unsafe.Pointer(host.mntp))
		host.mntp = nil
	}()

	/*
	 * Handle zombie mounts
	 *
	 * FUSE on UNIX does not automatically unmount the file system, leaving behind "zombie"
	 * mounts. So set things up to always unmount the file system (unless forcibly terminated).
	 * This has the added benefit that the file system Destroy() always gets called.
	 *
	 * On Windows (WinFsp) this is handled by the FUSE layer and we do not have to do anything.
	 */
	if "windows" != runtime.GOOS {
		done := make(chan bool)
		defer func() {
			<-done
		}()
		host.sigc = make(chan os.Signal, 1)
		defer close(host.sigc)
		go func() {
			_, ok := <-host.sigc
			if ok {
				host.Unmount()
			}
			close(done)
		}()
	}

	/*
	 * Tell FUSE to do its job!
	 */
	hndl := hostHandleNew(host)
	defer hostHandleDel(hndl)
	return 0 != C_hostMount(C_int(argc), &argv[0], hndl)
}

// Unmount unmounts a mounted file system.
// Unmount may be called at any time after the Init() method has been called
// and before the Destroy() method has been called.
func (host *FileSystemHost) Unmount() bool {
	if nil == host.fuse {
		return false
	}
	return 0 != C_hostUnmount(host.fuse, host.mntp)
}

// Getcontext gets information related to a file system operation.
func Getcontext() (uid uint32, gid uint32, pid int) {
	uid = uint32(C_fuse_get_context().uid)
	gid = uint32(C_fuse_get_context().gid)
	pid = int(C_fuse_get_context().pid)
	return
}

func optNormBool(opt string) string {
	if i := strings.Index(opt, "=%"); -1 != i {
		switch opt[i+2:] {
		case "d", "o", "x", "X":
			return opt
		case "v":
			return opt[:i+1]
		default:
			panic("unknown format " + opt[i+1:])
		}
	} else {
		return opt
	}
}

func optNormInt(opt string, modf string) string {
	if i := strings.Index(opt, "=%"); -1 != i {
		switch opt[i+2:] {
		case "d", "o", "x", "X":
			return opt[:i+2] + modf + opt[i+2:]
		case "v":
			return opt[:i+2] + modf + "i"
		default:
			panic("unknown format " + opt[i+1:])
		}
	} else if strings.HasSuffix(opt, "=") {
		return opt + "%" + modf + "i"
	} else {
		return opt + "=%" + modf + "i"
	}
}

func optNormStr(opt string) string {
	if i := strings.Index(opt, "=%"); -1 != i {
		switch opt[i+2:] {
		case "s", "v":
			return opt[:i+2] + "s"
		default:
			panic("unknown format " + opt[i+1:])
		}
	} else if strings.HasSuffix(opt, "=") {
		return opt + "%s"
	} else {
		return opt + "=%s"
	}
}

// OptParse parses the FUSE command line arguments in args as determined by format
// and stores the resulting values in vals, which must be pointers. It returns a
// list of unparsed arguments or nil if an error happens.
//
// The format may be empty or non-empty. An empty format is taken as a special
// instruction to OptParse to only return all non-option arguments in outargs.
//
// A non-empty format is a space separated list of acceptable FUSE options. Each
// option is matched with a corresponding pointer value in vals. The combination
// of the option and the type of the corresponding pointer value, determines how
// the option is used. The allowed pointer types are pointer to bool, pointer to
// an integer type and pointer to string.
//
// For pointer to bool types:
//
//     -x                       Match -x without parameter.
//     -foo --foo               As above for -foo or --foo.
//     foo                      Match "-o foo".
//     -x= -foo= --foo= foo=    Match option with parameter.
//     -x=%VERB ... foo=%VERB   Match option with parameter of syntax.
//                              Allowed verbs: d,o,x,X,v
//                              - d,o,x,X: set to true if parameter non-0.
//                              - v: set to true if parameter present.
//
//     The formats -x=, and -x=%v are equivalent.
//
// For pointer to other types:
//
//     -x                       Match -x with parameter (-x=PARAM).
//     -foo --foo               As above for -foo or --foo.
//     foo                      Match "-o foo=PARAM".
//     -x= -foo= --foo= foo=    Match option with parameter.
//     -x=%VERB ... foo=%VERB   Match option with parameter of syntax.
//                              Allowed verbs for pointer to int types: d,o,x,X,v
//                              Allowed verbs for pointer to string types: s,v
//
//     The formats -x, -x=, and -x=%v are equivalent.
//
// For example:
//
//     var f bool
//     var set_attr_timeout bool
//     var attr_timeout int
//     var umask uint32
//     outargs, err := OptParse(args, "-f attr_timeout= attr_timeout umask=%o",
//         &f, &set_attr_timeout, &attr_timeout, &umask)
//
// Will accept a command line of:
//
//     $ program -f -o attr_timeout=42,umask=077
//
// And will set variables as follows:
//
//     f == true
//     set_attr_timeout == true
//     attr_timeout == 42
//     umask == 077
//
func OptParse(args []string, format string, vals ...interface{}) (outargs []string, err error) {
	if 0 == C_hostFuseInit() {
		panic("cgofuse: cannot find winfsp")
	}

	defer func() {
		if r := recover(); nil != r {
			if s, ok := r.(string); ok {
				outargs = nil
				err = errors.New("OptParse: " + s)
			} else {
				panic(r)
			}
		}
	}()

	var opts []string
	var nonopts bool
	if "" == format {
		opts = make([]string, 0)
		nonopts = true
	} else {
		opts = strings.Split(format, " ")
	}

	align := int(2 * unsafe.Sizeof(C_size_t(0))) // match malloc alignment (usually 8 or 16)

	fuse_opts := make([]C_struct_fuse_opt, len(opts)+1)
	for i := 0; len(opts) > i; i++ {
		switch vals[i].(type) {
		case *bool:
			fuse_opts[i].templ = C_CString(optNormBool(opts[i]))
		case *int:
			fuse_opts[i].templ = C_CString(optNormInt(opts[i], ""))
		case *int8:
			fuse_opts[i].templ = C_CString(optNormInt(opts[i], "hh"))
		case *int16:
			fuse_opts[i].templ = C_CString(optNormInt(opts[i], "h"))
		case *int32:
			fuse_opts[i].templ = C_CString(optNormInt(opts[i], ""))
		case *int64:
			fuse_opts[i].templ = C_CString(optNormInt(opts[i], "ll"))
		case *uint:
			fuse_opts[i].templ = C_CString(optNormInt(opts[i], ""))
		case *uint8:
			fuse_opts[i].templ = C_CString(optNormInt(opts[i], "hh"))
		case *uint16:
			fuse_opts[i].templ = C_CString(optNormInt(opts[i], "h"))
		case *uint32:
			fuse_opts[i].templ = C_CString(optNormInt(opts[i], ""))
		case *uint64:
			fuse_opts[i].templ = C_CString(optNormInt(opts[i], "ll"))
		case *uintptr:
			fuse_opts[i].templ = C_CString(optNormInt(opts[i], "ll"))
		case *string:
			fuse_opts[i].templ = C_CString(optNormStr(opts[i]))
		}
		defer C_free(unsafe.Pointer(fuse_opts[i].templ))

		// Work around Go pre-1.10 limitation. See golang issue:
		// https://github.com/golang/go/issues/21809
		*(*C_fuse_opt_offset_t)(unsafe.Pointer(&fuse_opts[i].offset)) =
			C_fuse_opt_offset_t(i * align)

		fuse_opts[i].value = 1
	}

	fuse_args := C_struct_fuse_args{}
	defer C_fuse_opt_free_args(&fuse_args)
	argc := 1 + len(args)
	argp := C_calloc(C_size_t(argc+1), C_size_t(unsafe.Sizeof((*C_char)(nil))))
	defer C_free(argp)
	argv := (*[1 << 16]*C_char)(argp)
	argv[0] = C_CString("<UNKNOWN>")
	defer C_free(unsafe.Pointer(argv[0]))
	for i := 0; len(args) > i; i++ {
		argv[1+i] = C_CString(args[i])
		defer C_free(unsafe.Pointer(argv[1+i]))
	}
	fuse_args.argc = C_int(argc)
	fuse_args.argv = (**C_char)(&argv[0])

	data := C_calloc(C_size_t(len(opts)), C_size_t(align))
	defer C_free(data)

	if -1 == C_hostOptParse(&fuse_args, data, &fuse_opts[0], C_bool(nonopts)) {
		panic("failed")
	}

	for i := 0; len(opts) > i; i++ {
		switch v := vals[i].(type) {
		case *bool:
			*v = 0 != int(*(*C_int)(unsafe.Pointer(uintptr(data) + uintptr(i*align))))
		case *int:
			*v = int(*(*C_int)(unsafe.Pointer(uintptr(data) + uintptr(i*align))))
		case *int8:
			*v = int8(*(*C_int8_t)(unsafe.Pointer(uintptr(data) + uintptr(i*align))))
		case *int16:
			*v = int16(*(*C_int16_t)(unsafe.Pointer(uintptr(data) + uintptr(i*align))))
		case *int32:
			*v = int32(*(*C_int32_t)(unsafe.Pointer(uintptr(data) + uintptr(i*align))))
		case *int64:
			*v = int64(*(*C_int64_t)(unsafe.Pointer(uintptr(data) + uintptr(i*align))))
		case *uint:
			*v = uint(*(*C_unsigned)(unsafe.Pointer(uintptr(data) + uintptr(i*align))))
		case *uint8:
			*v = uint8(*(*C_uint8_t)(unsafe.Pointer(uintptr(data) + uintptr(i*align))))
		case *uint16:
			*v = uint16(*(*C_uint16_t)(unsafe.Pointer(uintptr(data) + uintptr(i*align))))
		case *uint32:
			*v = uint32(*(*C_uint32_t)(unsafe.Pointer(uintptr(data) + uintptr(i*align))))
		case *uint64:
			*v = uint64(*(*C_uint64_t)(unsafe.Pointer(uintptr(data) + uintptr(i*align))))
		case *uintptr:
			*v = uintptr(*(*C_uintptr_t)(unsafe.Pointer(uintptr(data) + uintptr(i*align))))
		case *string:
			s := *(**C_char)(unsafe.Pointer(uintptr(data) + uintptr(i*align)))
			*v = C_GoString(s)
			C_free(unsafe.Pointer(s))
		}
	}

	if 1 >= fuse_args.argc {
		outargs = make([]string, 0)
	} else {
		outargs = make([]string, fuse_args.argc-1)
		for i := 1; int(fuse_args.argc) > i; i++ {
			outargs[i-1] = C_GoString((*[1 << 16]*C_char)(unsafe.Pointer(fuse_args.argv))[i])
		}
	}

	if nonopts && 1 <= len(outargs) && "--" == outargs[0] {
		outargs = outargs[1:]
	}

	return
}

func init() {
	C_hostStaticInit()
}

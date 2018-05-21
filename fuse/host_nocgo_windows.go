// +build !cgo,windows

/*
 * host_nocgo_windows.go
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
	"path/filepath"
	"sync"
	"syscall"
	"unsafe"
)

type fuse_operations struct {
	getattr     uintptr
	getdir      uintptr
	readlink    uintptr
	mknod       uintptr
	mkdir       uintptr
	unlink      uintptr
	rmdir       uintptr
	symlink     uintptr
	rename      uintptr
	link        uintptr
	chmod       uintptr
	chown       uintptr
	truncate    uintptr
	utime       uintptr
	open        uintptr
	read        uintptr
	write       uintptr
	statfs      uintptr
	flush       uintptr
	release     uintptr
	fsync       uintptr
	setxattr    uintptr
	getxattr    uintptr
	listxattr   uintptr
	removexattr uintptr
	opendir     uintptr
	readdir     uintptr
	releasedir  uintptr
	fsyncdir    uintptr
	init        uintptr
	destroy     uintptr
	access      uintptr
	create      uintptr
	ftruncate   uintptr
	fgetattr    uintptr
	lock        uintptr
	utimens     uintptr
	bmap        uintptr
	ioctl       uintptr
	poll        uintptr
	write_buf   uintptr
	read_buf    uintptr
	flock       uintptr
	fallocate   uintptr
	reserved00  uintptr
	reserved01  uintptr
	reserved02  uintptr
	statfs_x    uintptr
	setvolname  uintptr
	exchange    uintptr
	getxtimes   uintptr
	setbkuptime uintptr
	setchgtime  uintptr
	setcrtime   uintptr
	chflags     uintptr
	setattr_x   uintptr
	fsetattr_x  uintptr
}

type fuse_fill_dir_t struct {
}

type fuse_stat_t struct {
	st_dev      c_fuse_dev_t
	st_ino      c_fuse_ino_t
	st_mode     c_fuse_mode_t
	st_nlink    c_fuse_nlink_t
	st_uid      c_fuse_uid_t
	st_gid      c_fuse_gid_t
	st_rdev     c_fuse_dev_t
	st_size     c_fuse_off_t
	st_atim     c_fuse_timespec_t
	st_mtim     c_fuse_timespec_t
	st_ctim     c_fuse_timespec_t
	st_blksize  c_fuse_blksize_t
	st_blocks   c_fuse_blkcnt_t
	st_birthtim c_fuse_timespec_t
}

type fuse_stat_ex_t struct {
	fuse_stat_t
	st_flags      c_uint32_t
	st_reserved32 [3]c_uint32_t
	st_reserved64 [2]c_uint64_t
}

type fuse_statvfs_t struct {
	f_bsize   c_uint64_t
	f_frsize  c_uint64_t
	f_blocks  c_fuse_fsblkcnt_t
	f_bfree   c_fuse_fsblkcnt_t
	f_bavail  c_fuse_fsblkcnt_t
	f_files   c_fuse_fsfilcnt_t
	f_ffree   c_fuse_fsfilcnt_t
	f_favail  c_fuse_fsfilcnt_t
	f_fsid    c_uint64_t
	f_flag    c_uint64_t
	f_namemax c_uint64_t
}

type fuse_timespec_t struct {
	tv_sec  uintptr
	tv_nsec uintptr
}

type struct_fuse struct {
}

type struct_fuse_args struct {
	argc      c_int
	argv      **c_char
	allocated c_int
}

type struct_fuse_conn_info struct {
	proto_major   c_unsigned
	proto_minor   c_unsigned
	async_read    c_unsigned
	max_write     c_unsigned
	max_readahead c_unsigned
	capable       c_unsigned
	want          c_unsigned
	reserved      [25]c_unsigned
}

type struct_fuse_context struct {
	fuse         *c_struct_fuse
	uid          c_fuse_uid_t
	gid          c_fuse_gid_t
	pid          c_fuse_pid_t
	private_data unsafe.Pointer
	umask        c_fuse_mode_t
}

type struct_fuse_file_info struct {
	flags      c_int
	fh_old     c_unsigned
	writepage  c_int
	bits       c_uint32_t
	fh         c_uint64_t
	lock_owner c_uint64_t
}

type struct_fuse_opt struct {
	templ  *c_char
	offset c_fuse_opt_offset_t
	value  c_int
}

type (
	c_bool                  = bool
	c_char                  = byte
	c_fuse_blkcnt_t         = int64
	c_fuse_blksize_t        = int32
	c_fuse_dev_t            = uint32
	c_fuse_fill_dir_t       = fuse_fill_dir_t
	c_fuse_fsblkcnt_t       = uintptr
	c_fuse_fsfilcnt_t       = uintptr
	c_fuse_gid_t            = uint32
	c_fuse_ino_t            = uint64
	c_fuse_mode_t           = uint32
	c_fuse_nlink_t          = uint16
	c_fuse_off_t            = int64
	c_fuse_opt_offset_t     = uint32
	c_fuse_pid_t            = int32
	c_fuse_stat_t           = fuse_stat_t
	c_fuse_statvfs_t        = fuse_statvfs_t
	c_fuse_timespec_t       = fuse_timespec_t
	c_fuse_uid_t            = uint32
	c_int                   = int32
	c_int16_t               = int16
	c_int32_t               = int32
	c_int64_t               = int64
	c_int8_t                = int8
	c_size_t                = uintptr
	c_struct_fuse           = struct_fuse
	c_struct_fuse_args      = struct_fuse_args
	c_struct_fuse_conn_info = struct_fuse_conn_info
	c_struct_fuse_context   = struct_fuse_context
	c_struct_fuse_file_info = struct_fuse_file_info
	c_struct_fuse_opt       = struct_fuse_opt
	c_uint16_t              = uint16
	c_uint32_t              = uint32
	c_uint64_t              = uint64
	c_uint8_t               = uint8
	c_uintptr_t             = uintptr
	c_unsigned              = uint32
)

var (
	kernel32       = syscall.MustLoadDLL("kernel32.dll")
	getProcessHeap = kernel32.MustFindProc("GetProcessHeap")
	heapAlloc      = kernel32.MustFindProc("HeapAlloc")
	heapFree       = kernel32.MustFindProc("HeapFree")
	processHeap    uintptr

	/*
	 * It appears safe to call cdecl functions from Go. Is it really?
	 * https://codereview.appspot.com/4961045/
	 */
	fuseOnce                 sync.Once
	fuseDll                  *syscall.DLL
	fuse_version             *syscall.Proc
	fuse_mount               *syscall.Proc
	fuse_unmount             *syscall.Proc
	fuse_parse_cmdline       *syscall.Proc
	fuse_ntstatus_from_errno *syscall.Proc
	fuse_main_real           *syscall.Proc
	fuse_is_lib_option       *syscall.Proc
	fuse_new                 *syscall.Proc
	fuse_destroy             *syscall.Proc
	fuse_loop                *syscall.Proc
	fuse_loop_mt             *syscall.Proc
	fuse_exit                *syscall.Proc
	fuse_get_context         *syscall.Proc
	fuse_opt_parse           *syscall.Proc
	fuse_opt_add_arg         *syscall.Proc
	fuse_opt_insert_arg      *syscall.Proc
	fuse_opt_free_args       *syscall.Proc
	fuse_opt_add_opt         *syscall.Proc
	fuse_opt_add_opt_escaped *syscall.Proc
	fuse_opt_match           *syscall.Proc

	fsop = fuse_operations{
		getattr:     syscall.NewCallbackCDecl(hostGetattr),
		readlink:    syscall.NewCallbackCDecl(hostReadlink),
		mknod:       syscall.NewCallbackCDecl(hostMknod),
		mkdir:       syscall.NewCallbackCDecl(hostMkdir),
		unlink:      syscall.NewCallbackCDecl(hostUnlink),
		rmdir:       syscall.NewCallbackCDecl(hostRmdir),
		symlink:     syscall.NewCallbackCDecl(hostSymlink),
		rename:      syscall.NewCallbackCDecl(hostRename),
		link:        syscall.NewCallbackCDecl(hostLink),
		chmod:       syscall.NewCallbackCDecl(hostChmod),
		chown:       syscall.NewCallbackCDecl(hostChown),
		truncate:    syscall.NewCallbackCDecl(hostTruncate),
		open:        syscall.NewCallbackCDecl(hostOpen),
		read:        syscall.NewCallbackCDecl(hostRead),
		write:       syscall.NewCallbackCDecl(hostWrite),
		statfs:      syscall.NewCallbackCDecl(hostStatfs),
		flush:       syscall.NewCallbackCDecl(hostFlush),
		release:     syscall.NewCallbackCDecl(hostRelease),
		fsync:       syscall.NewCallbackCDecl(hostFsync),
		setxattr:    syscall.NewCallbackCDecl(hostSetxattr),
		getxattr:    syscall.NewCallbackCDecl(hostGetxattr),
		listxattr:   syscall.NewCallbackCDecl(hostListxattr),
		removexattr: syscall.NewCallbackCDecl(hostRemovexattr),
		opendir:     syscall.NewCallbackCDecl(hostOpendir),
		readdir:     syscall.NewCallbackCDecl(hostReaddir),
		releasedir:  syscall.NewCallbackCDecl(hostReleasedir),
		fsyncdir:    syscall.NewCallbackCDecl(hostFsyncdir),
		init:        syscall.NewCallbackCDecl(hostInit),
		destroy:     syscall.NewCallbackCDecl(hostDestroy),
		access:      syscall.NewCallbackCDecl(hostAccess),
		create:      syscall.NewCallbackCDecl(hostCreate),
		ftruncate:   syscall.NewCallbackCDecl(hostFtruncate),
		fgetattr:    syscall.NewCallbackCDecl(hostFgetattr),
		utimens:     syscall.NewCallbackCDecl(hostUtimens),
		setchgtime:  syscall.NewCallbackCDecl(hostSetchgtime),
		setcrtime:   syscall.NewCallbackCDecl(hostSetcrtime),
		chflags:     syscall.NewCallbackCDecl(hostChflags),
	}

	cgofuse_stat_ex bool
)

const (
	FSP_FUSE_CAP_CASE_INSENSITIVE = 1 << 29
	FSP_FUSE_CAP_READDIR_PLUS     = 1 << 21
	FSP_FUSE_CAP_STAT_EX          = 1 << 23
)

func init() {
	processHeap, _, _ = getProcessHeap.Call()
}

func c_GoString(s *c_char) string {
	if nil == s {
		return ""
	}
	q := (*[1 << 30]c_char)(unsafe.Pointer(s))
	l := 0
	for 0 != q[l] {
		l++
	}
	return string(q[:l])
}
func c_CString(s string) *c_char {
	p := c_malloc(c_size_t(len(s) + 1))
	q := (*[1 << 30]c_char)(p)
	copy(q[:], s)
	q[len(s)] = 0
	return (*c_char)(p)
}

func c_malloc(size c_size_t) unsafe.Pointer {
	p, _, _ := heapAlloc.Call(processHeap, 0, size)
	if 0 == p {
		panic("runtime: C malloc failed")
	}
	return unsafe.Pointer(p)
}
func c_calloc(count c_size_t, size c_size_t) unsafe.Pointer {
	p, _, _ := heapAlloc.Call(processHeap, 8 /*HEAP_ZERO_MEMORY*/, count*size)
	return unsafe.Pointer(p)
}
func c_free(p unsafe.Pointer) {
	if nil != p {
		heapFree.Call(processHeap, 0, uintptr(p))
	}
}

func c_fuse_get_context() *c_struct_fuse_context {
	p, _, _ := fuse_get_context.Call()
	return (*c_struct_fuse_context)(unsafe.Pointer(p))
}
func c_fuse_opt_free_args(args *c_struct_fuse_args) {
	fuse_opt_free_args.Call(uintptr(unsafe.Pointer(args)))
}

func c_hostAsgnCconninfo(conn *c_struct_fuse_conn_info,
	capCaseInsensitive c_bool,
	capReaddirPlus c_bool) {
	conn.want |= conn.capable & FSP_FUSE_CAP_STAT_EX
	cgofuse_stat_ex = 0 != conn.want&FSP_FUSE_CAP_STAT_EX // hack!
	if capCaseInsensitive {
		conn.want |= conn.capable & FSP_FUSE_CAP_CASE_INSENSITIVE
	}
	if capReaddirPlus {
		conn.want |= conn.capable & FSP_FUSE_CAP_READDIR_PLUS
	}
}
func c_hostCstatvfsFromFusestatfs(stbuf *c_fuse_statvfs_t,
	bsize c_uint64_t,
	frsize c_uint64_t,
	blocks c_uint64_t,
	bfree c_uint64_t,
	bavail c_uint64_t,
	files c_uint64_t,
	ffree c_uint64_t,
	favail c_uint64_t,
	fsid c_uint64_t,
	flag c_uint64_t,
	namemax c_uint64_t) {
	*stbuf = c_fuse_statvfs_t{
		f_bsize:   bsize,
		f_frsize:  frsize,
		f_blocks:  c_fuse_fsblkcnt_t(blocks),
		f_bfree:   c_fuse_fsblkcnt_t(bfree),
		f_bavail:  c_fuse_fsblkcnt_t(bavail),
		f_files:   c_fuse_fsfilcnt_t(files),
		f_ffree:   c_fuse_fsfilcnt_t(ffree),
		f_favail:  c_fuse_fsfilcnt_t(favail),
		f_fsid:    fsid,
		f_flag:    flag,
		f_namemax: namemax,
	}
}
func c_hostCstatFromFusestat(stbuf *c_fuse_stat_t,
	dev c_uint64_t,
	ino c_uint64_t,
	mode c_uint32_t,
	nlink c_uint32_t,
	uid c_uint32_t,
	gid c_uint32_t,
	rdev c_uint64_t,
	size c_int64_t,
	atimSec c_int64_t, atimNsec c_int64_t,
	mtimSec c_int64_t, mtimNsec c_int64_t,
	ctimSec c_int64_t, ctimNsec c_int64_t,
	blksize c_int64_t,
	blocks c_int64_t,
	birthtimSec c_int64_t, birthtimNsec c_int64_t,
	flags c_uint32_t) {
	if !cgofuse_stat_ex {
		*stbuf = c_fuse_stat_t{
			st_dev:     c_fuse_dev_t(dev),
			st_ino:     c_fuse_ino_t(ino),
			st_mode:    c_fuse_mode_t(mode),
			st_nlink:   c_fuse_nlink_t(nlink),
			st_uid:     c_fuse_uid_t(uid),
			st_gid:     c_fuse_gid_t(gid),
			st_rdev:    c_fuse_dev_t(rdev),
			st_size:    c_fuse_off_t(size),
			st_blksize: c_fuse_blksize_t(blksize),
			st_blocks:  c_fuse_blkcnt_t(blocks),
			st_atim: c_fuse_timespec_t{
				tv_sec:  uintptr(atimSec),
				tv_nsec: uintptr(atimNsec),
			},
			st_mtim: c_fuse_timespec_t{
				tv_sec:  uintptr(mtimSec),
				tv_nsec: uintptr(mtimNsec),
			},
			st_ctim: c_fuse_timespec_t{
				tv_sec:  uintptr(ctimSec),
				tv_nsec: uintptr(ctimNsec),
			},
		}
	} else {
		*(*fuse_stat_ex_t)(unsafe.Pointer(stbuf)) = fuse_stat_ex_t{
			fuse_stat_t: c_fuse_stat_t{
				st_dev:     c_fuse_dev_t(dev),
				st_ino:     c_fuse_ino_t(ino),
				st_mode:    c_fuse_mode_t(mode),
				st_nlink:   c_fuse_nlink_t(nlink),
				st_uid:     c_fuse_uid_t(uid),
				st_gid:     c_fuse_gid_t(gid),
				st_rdev:    c_fuse_dev_t(rdev),
				st_size:    c_fuse_off_t(size),
				st_blksize: c_fuse_blksize_t(blksize),
				st_blocks:  c_fuse_blkcnt_t(blocks),
				st_atim: c_fuse_timespec_t{
					tv_sec:  uintptr(atimSec),
					tv_nsec: uintptr(atimNsec),
				},
				st_mtim: c_fuse_timespec_t{
					tv_sec:  uintptr(mtimSec),
					tv_nsec: uintptr(mtimNsec),
				},
				st_ctim: c_fuse_timespec_t{
					tv_sec:  uintptr(ctimSec),
					tv_nsec: uintptr(ctimNsec),
				},
			},
			st_flags: flags,
		}
	}
	if 0 != birthtimSec {
		stbuf.st_birthtim.tv_sec = uintptr(birthtimSec)
		stbuf.st_birthtim.tv_nsec = uintptr(birthtimNsec)
	} else {
		stbuf.st_birthtim.tv_sec = uintptr(ctimSec)
		stbuf.st_birthtim.tv_nsec = uintptr(ctimNsec)
	}
}
func c_hostFilldir(filler c_fuse_fill_dir_t,
	buf unsafe.Pointer, name *c_char, stbuf *c_fuse_stat_t, off c_fuse_off_t) c_int {
	return 0
}
func c_hostStaticInit() {
}
func c_hostFuseInit() c_int {
	fuseOnce.Do(func() {
		fuseDll, _ = fspload()
		if nil != fuseDll {
			fuse_version = fuseDll.MustFindProc("fuse_version")
			fuse_mount = fuseDll.MustFindProc("fuse_mount")
			fuse_unmount = fuseDll.MustFindProc("fuse_unmount")
			fuse_parse_cmdline = fuseDll.MustFindProc("fuse_parse_cmdline")
			fuse_ntstatus_from_errno = fuseDll.MustFindProc("fuse_ntstatus_from_errno")
			fuse_main_real = fuseDll.MustFindProc("fuse_main_real")
			fuse_is_lib_option = fuseDll.MustFindProc("fuse_is_lib_option")
			fuse_new = fuseDll.MustFindProc("fuse_new")
			fuse_destroy = fuseDll.MustFindProc("fuse_destroy")
			fuse_loop = fuseDll.MustFindProc("fuse_loop")
			fuse_loop_mt = fuseDll.MustFindProc("fuse_loop_mt")
			fuse_exit = fuseDll.MustFindProc("fuse_exit")
			fuse_get_context = fuseDll.MustFindProc("fuse_get_context")
			fuse_opt_parse = fuseDll.MustFindProc("fuse_opt_parse")
			fuse_opt_add_arg = fuseDll.MustFindProc("fuse_opt_add_arg")
			fuse_opt_insert_arg = fuseDll.MustFindProc("fuse_opt_insert_arg")
			fuse_opt_free_args = fuseDll.MustFindProc("fuse_opt_free_args")
			fuse_opt_add_opt = fuseDll.MustFindProc("fuse_opt_add_opt")
			fuse_opt_add_opt_escaped = fuseDll.MustFindProc("fuse_opt_add_opt_escaped")
			fuse_opt_match = fuseDll.MustFindProc("fuse_opt_match")
		}
	})
	if nil == fuseDll {
		return 0
	}
	return 1
}
func c_hostMountpoint(argc c_int, argv **c_char) *c_char {
	return nil
}
func c_hostMount(argc c_int, argv **c_char, data unsafe.Pointer) c_int {
	return 0
}
func c_hostUnmount(fuse *c_struct_fuse, mountpoint *c_char) c_int {
	return 0
}
func c_hostOptParse(args *c_struct_fuse_args, data unsafe.Pointer, opts *c_struct_fuse_opt,
	nonopts c_bool) c_int {
	return 0
}

func fspload() (dll *syscall.DLL, err error) {
	dllname := ""
	if uint64(0xffffffff) < uint64(^uintptr(0)) {
		dllname = "winfsp-x64.dll"
	} else {
		dllname = "winfsp-x86.dll"
	}

	dll, err = syscall.LoadDLL(dllname)
	if nil == dll {
		var pathbuf [syscall.MAX_PATH]uint16
		var regkey syscall.Handle
		var regtype, size uint32

		kname, _ := syscall.UTF16PtrFromString("Software\\WinFsp")
		err = syscall.RegOpenKeyEx(syscall.HKEY_LOCAL_MACHINE, kname,
			0, syscall.KEY_READ|syscall.KEY_WOW64_32KEY, &regkey)
		if nil != err {
			err = syscall.ERROR_MOD_NOT_FOUND
			return
		}

		vname, _ := syscall.UTF16PtrFromString("InstallDir")
		size = uint32(len(pathbuf) * 2)
		err = syscall.RegQueryValueEx(regkey, vname,
			nil, &regtype, (*byte)(unsafe.Pointer(&pathbuf)), &size)
		syscall.RegCloseKey(regkey)
		if nil != err || syscall.REG_SZ != regtype {
			err = syscall.ERROR_MOD_NOT_FOUND
			return
		}

		if 0 < size && 0 == pathbuf[size/2-1] {
			size -= 2
		}

		path := syscall.UTF16ToString(pathbuf[:size/2])
		dllpath := filepath.Join(path, "bin", dllname)

		dll, err = syscall.LoadDLL(dllpath)
		if nil != err {
			err = syscall.ERROR_MOD_NOT_FOUND
			return
		}
	}

	return
}

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
	"syscall"
	"unsafe"
)

type fuse_fill_dir_t struct {
}

type fuse_stat_t struct {
	st_dev        c_fuse_dev_t
	st_ino        c_fuse_ino_t
	st_mode       c_fuse_mode_t
	st_nlink      c_fuse_nlink_t
	st_uid        c_fuse_uid_t
	st_gid        c_fuse_gid_t
	st_rdev       c_fuse_dev_t
	st_size       c_fuse_off_t
	st_atim       c_fuse_timespec_t
	st_mtim       c_fuse_timespec_t
	st_ctim       c_fuse_timespec_t
	st_blksize    c_fuse_blksize_t
	st_blocks     c_fuse_blkcnt_t
	st_birthtim   c_fuse_timespec_t
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
)

func init() {
	processHeap, _, _ = getProcessHeap.Call()
}

func c_GoString(s *c_char) string {
	return ""
}
func c_CString(s string) *c_char {
	return nil
}

func c_malloc(size c_size_t) unsafe.Pointer {
	p, _, _ := heapAlloc.Call(processHeap, 0, size)
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
	return nil
}
func c_fuse_opt_free_args(args *c_struct_fuse_args) {
}

func c_hostAsgnCconninfo(conn *c_struct_fuse_conn_info,
	capCaseInsensitive c_bool,
	capReaddirPlus c_bool) {
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
}
func c_hostFilldir(filler c_fuse_fill_dir_t,
	buf unsafe.Pointer, name *c_char, stbuf *c_fuse_stat_t, off c_fuse_off_t) c_int {
	return 0
}
func c_hostStaticInit() {
}
func c_hostFuseInit() c_int {
	return 0
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

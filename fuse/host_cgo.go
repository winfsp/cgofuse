// +build cgo

/*
 * host_cgo.go
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

/*
#cgo darwin CFLAGS: -DFUSE_USE_VERSION=28 -D_FILE_OFFSET_BITS=64 -I/usr/local/include/osxfuse/fuse
#cgo darwin LDFLAGS: -L/usr/local/lib -losxfuse

#cgo freebsd CFLAGS: -DFUSE_USE_VERSION=28 -D_FILE_OFFSET_BITS=64 -I/usr/local/include/fuse
#cgo freebsd LDFLAGS: -L/usr/local/lib -lfuse

#cgo linux CFLAGS: -DFUSE_USE_VERSION=28 -D_FILE_OFFSET_BITS=64 -I/usr/include/fuse
#cgo linux LDFLAGS: -lfuse

// Use `set CPATH=C:\Program Files (x86)\WinFsp\inc\fuse` on Windows.
// The flag `I/usr/local/include/winfsp` only works on xgo and docker.
#cgo windows CFLAGS: -DFUSE_USE_VERSION=28 -I/usr/local/include/winfsp

#if !(defined(__APPLE__) || defined(__FreeBSD__) || defined(__linux__) || defined(_WIN32))
#error platform not supported
#endif

#include <stdbool.h>
#include <stdlib.h>
#include <string.h>

#if defined(__APPLE__) || defined(__FreeBSD__) || defined(__linux__)

#include <spawn.h>
#include <sys/mount.h>
#include <sys/wait.h>
#include <fuse.h>

#elif defined(_WIN32)

#include <windows.h>

static PVOID cgofuse_init_slow(int hardfail);
static VOID  cgofuse_init_fail(VOID);
static PVOID cgofuse_init_winfsp(VOID);

static CRITICAL_SECTION cgofuse_lock;
static PVOID cgofuse_module = 0;
static BOOLEAN cgofuse_stat_ex = FALSE;

static inline PVOID cgofuse_init_fast(int hardfail)
{
	PVOID Module = cgofuse_module;
	MemoryBarrier();
	if (0 == Module)
		Module = cgofuse_init_slow(hardfail);
	return Module;
}

static PVOID cgofuse_init_slow(int hardfail)
{
	PVOID Module;
	EnterCriticalSection(&cgofuse_lock);
	Module = cgofuse_module;
	if (0 == Module)
	{
		Module = cgofuse_init_winfsp();
		MemoryBarrier();
		cgofuse_module = Module;
	}
	LeaveCriticalSection(&cgofuse_lock);
	if (0 == Module && hardfail)
		cgofuse_init_fail();
	return Module;
}

static VOID cgofuse_init_fail(VOID)
{
	static const char *message = "cgofuse: cannot find winfsp\n";
	DWORD BytesTransferred;
	WriteFile(GetStdHandle(STD_ERROR_HANDLE), message, lstrlenA(message), &BytesTransferred, 0);
	ExitProcess(ERROR_DLL_NOT_FOUND);
}

#define FSP_FUSE_API                    static
#define FSP_FUSE_API_NAME(api)          (* pfn_ ## api)
#define FSP_FUSE_API_CALL(api)          (cgofuse_init_fast(1), pfn_ ## api)
#define FSP_FUSE_SYM(proto, ...)        static inline proto { __VA_ARGS__ }
#include <fuse_common.h>
#include <fuse.h>
#include <fuse_opt.h>

static NTSTATUS FspLoad(PVOID *PModule)
{
#if defined(_WIN64)
#define FSP_DLLNAME                     "winfsp-x64.dll"
#else
#define FSP_DLLNAME                     "winfsp-x86.dll"
#endif
#define FSP_DLLPATH                     "bin\\" FSP_DLLNAME

	WCHAR PathBuf[MAX_PATH];
	DWORD Size;
	DWORD RegType;
	HKEY RegKey;
	LONG Result;
	HMODULE Module;

	if (0 != PModule)
		*PModule = 0;

	Module = LoadLibraryW(L"" FSP_DLLNAME);
	if (0 == Module)
	{
		Result = RegOpenKeyExW(HKEY_LOCAL_MACHINE, L"Software\\WinFsp",
			0, KEY_READ | KEY_WOW64_32KEY, &RegKey);
		if (ERROR_SUCCESS == Result)
		{
			Size = sizeof PathBuf - sizeof L"" FSP_DLLPATH + sizeof(WCHAR);
			Result = RegQueryValueExW(RegKey, L"InstallDir", 0,
				&RegType, (LPBYTE)PathBuf, &Size);
			RegCloseKey(RegKey);
			if (ERROR_SUCCESS == Result && REG_SZ != RegType)
				Result = ERROR_FILE_NOT_FOUND;
		}
		if (ERROR_SUCCESS != Result)
			return 0xC0000034;//STATUS_OBJECT_NAME_NOT_FOUND

		if (0 < Size && L'\0' == PathBuf[Size / sizeof(WCHAR) - 1])
			Size -= sizeof(WCHAR);

		RtlCopyMemory(PathBuf + Size / sizeof(WCHAR),
			L"" FSP_DLLPATH, sizeof L"" FSP_DLLPATH);
		Module = LoadLibraryW(PathBuf);
		if (0 == Module)
			return 0xC0000135;//STATUS_DLL_NOT_FOUND
	}

	if (0 != PModule)
		*PModule = Module;

	return 0;//STATUS_SUCCESS

#undef FSP_DLLNAME
#undef FSP_DLLPATH
}

#define CGOFUSE_GET_API(h, n)           \
	if (0 == (*(void **)&(pfn_ ## n) = GetProcAddress(Module, #n)))\
		return 0;

static PVOID cgofuse_init_winfsp(VOID)
{
	PVOID Module;
	NTSTATUS Result;

	Result = FspLoad(&Module);
	if (0 > Result)
		return 0;

	// fuse_common.h
	CGOFUSE_GET_API(h, fsp_fuse_version);
	CGOFUSE_GET_API(h, fsp_fuse_mount);
	CGOFUSE_GET_API(h, fsp_fuse_unmount);
	CGOFUSE_GET_API(h, fsp_fuse_parse_cmdline);
	CGOFUSE_GET_API(h, fsp_fuse_ntstatus_from_errno);

	// fuse.h
	CGOFUSE_GET_API(h, fsp_fuse_main_real);
	CGOFUSE_GET_API(h, fsp_fuse_is_lib_option);
	CGOFUSE_GET_API(h, fsp_fuse_new);
	CGOFUSE_GET_API(h, fsp_fuse_destroy);
	CGOFUSE_GET_API(h, fsp_fuse_loop);
	CGOFUSE_GET_API(h, fsp_fuse_loop_mt);
	CGOFUSE_GET_API(h, fsp_fuse_exit);
	CGOFUSE_GET_API(h, fsp_fuse_get_context);

	// fuse_opt.h
	CGOFUSE_GET_API(h, fsp_fuse_opt_parse);
	CGOFUSE_GET_API(h, fsp_fuse_opt_add_arg);
	CGOFUSE_GET_API(h, fsp_fuse_opt_insert_arg);
	CGOFUSE_GET_API(h, fsp_fuse_opt_free_args);
	CGOFUSE_GET_API(h, fsp_fuse_opt_add_opt);
	CGOFUSE_GET_API(h, fsp_fuse_opt_add_opt_escaped);
	CGOFUSE_GET_API(h, fsp_fuse_opt_match);

	return Module;
}

#endif

#if defined(__APPLE__) || defined(__FreeBSD__) || defined(__linux__)
typedef struct stat fuse_stat_t;
typedef struct statvfs fuse_statvfs_t;
typedef struct timespec fuse_timespec_t;
typedef mode_t fuse_mode_t;
typedef dev_t fuse_dev_t;
typedef uid_t fuse_uid_t;
typedef gid_t fuse_gid_t;
typedef off_t fuse_off_t;
typedef unsigned long fuse_opt_offset_t;
#elif defined(_WIN32)
typedef struct fuse_stat fuse_stat_t;
typedef struct fuse_statvfs fuse_statvfs_t;
typedef struct fuse_timespec fuse_timespec_t;
typedef unsigned int fuse_opt_offset_t;
#endif

extern int go_hostGetattr(char *path, fuse_stat_t *stbuf);
extern int go_hostReadlink(char *path, char *buf, size_t size);
extern int go_hostMknod(char *path, fuse_mode_t mode, fuse_dev_t dev);
extern int go_hostMkdir(char *path, fuse_mode_t mode);
extern int go_hostUnlink(char *path);
extern int go_hostRmdir(char *path);
extern int go_hostSymlink(char *target, char *newpath);
extern int go_hostRename(char *oldpath, char *newpath);
extern int go_hostLink(char *oldpath, char *newpath);
extern int go_hostChmod(char *path, fuse_mode_t mode);
extern int go_hostChown(char *path, fuse_uid_t uid, fuse_gid_t gid);
extern int go_hostTruncate(char *path, fuse_off_t size);
extern int go_hostOpen(char *path, struct fuse_file_info *fi);
extern int go_hostRead(char *path, char *buf, size_t size, fuse_off_t off,
	struct fuse_file_info *fi);
extern int go_hostWrite(char *path, char *buf, size_t size, fuse_off_t off,
	struct fuse_file_info *fi);
extern int go_hostStatfs(char *path, fuse_statvfs_t *stbuf);
extern int go_hostFlush(char *path, struct fuse_file_info *fi);
extern int go_hostRelease(char *path, struct fuse_file_info *fi);
extern int go_hostFsync(char *path, int datasync, struct fuse_file_info *fi);
extern int go_hostSetxattr(char *path, char *name, char *value, size_t size, int flags);
extern int go_hostGetxattr(char *path, char *name, char *value, size_t size);
extern int go_hostListxattr(char *path, char *namebuf, size_t size);
extern int go_hostRemovexattr(char *path, char *name);
extern int go_hostOpendir(char *path, struct fuse_file_info *fi);
extern int go_hostReaddir(char *path, void *buf, fuse_fill_dir_t filler, fuse_off_t off,
	struct fuse_file_info *fi);
extern int go_hostReleasedir(char *path, struct fuse_file_info *fi);
extern int go_hostFsyncdir(char *path, int datasync, struct fuse_file_info *fi);
extern void *go_hostInit(struct fuse_conn_info *conn);
extern void go_hostDestroy(void *data);
extern int go_hostAccess(char *path, int mask);
extern int go_hostCreate(char *path, fuse_mode_t mode, struct fuse_file_info *fi);
extern int go_hostFtruncate(char *path, fuse_off_t off, struct fuse_file_info *fi);
extern int go_hostFgetattr(char *path, fuse_stat_t *stbuf, struct fuse_file_info *fi);
//extern int go_hostLock(char *path, struct fuse_file_info *fi, int cmd, struct fuse_flock *lock);
extern int go_hostUtimens(char *path, fuse_timespec_t tv[2]);
extern int go_hostSetchgtime(char *path, fuse_timespec_t *tv);
extern int go_hostSetcrtime(char *path, fuse_timespec_t *tv);
extern int go_hostChflags(char *path, uint32_t flags);

static inline void hostAsgnCconninfo(struct fuse_conn_info *conn,
	bool capCaseInsensitive,
	bool capReaddirPlus)
{
#if defined(__APPLE__)
	if (capCaseInsensitive)
		FUSE_ENABLE_CASE_INSENSITIVE(conn);
#elif defined(__FreeBSD__) || defined(__linux__)
#elif defined(_WIN32)
#if defined(FSP_FUSE_CAP_STAT_EX)
	conn->want |= conn->capable & FSP_FUSE_CAP_STAT_EX;
	cgofuse_stat_ex = 0 != (conn->want & FSP_FUSE_CAP_STAT_EX); // hack!
#endif
	if (capCaseInsensitive)
		conn->want |= conn->capable & FSP_FUSE_CAP_CASE_INSENSITIVE;
	if (capReaddirPlus)
		conn->want |= conn->capable & FSP_FUSE_CAP_READDIR_PLUS;
#endif
}

static inline void hostCstatvfsFromFusestatfs(fuse_statvfs_t *stbuf,
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

static inline void hostCstatFromFusestat(fuse_stat_t *stbuf,
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
	int64_t birthtimSec, int64_t birthtimNsec,
	uint32_t flags)
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
	stbuf->st_blksize = blksize;
	stbuf->st_blocks = blocks;
#if defined(__APPLE__)
	stbuf->st_atimespec.tv_sec = atimSec; stbuf->st_atimespec.tv_nsec = atimNsec;
	stbuf->st_mtimespec.tv_sec = mtimSec; stbuf->st_mtimespec.tv_nsec = mtimNsec;
	stbuf->st_ctimespec.tv_sec = ctimSec; stbuf->st_ctimespec.tv_nsec = ctimNsec;
	if (0 != birthtimSec)
	{
		stbuf->st_birthtimespec.tv_sec = birthtimSec;
		stbuf->st_birthtimespec.tv_nsec = birthtimNsec;
	}
	else
	{
		stbuf->st_birthtimespec.tv_sec = ctimSec;
		stbuf->st_birthtimespec.tv_nsec = ctimNsec;
	}
	stbuf->st_flags = flags;
#elif defined(_WIN32)
	stbuf->st_atim.tv_sec = atimSec; stbuf->st_atim.tv_nsec = atimNsec;
	stbuf->st_mtim.tv_sec = mtimSec; stbuf->st_mtim.tv_nsec = mtimNsec;
	stbuf->st_ctim.tv_sec = ctimSec; stbuf->st_ctim.tv_nsec = ctimNsec;
	if (0 != birthtimSec)
	{
		stbuf->st_birthtim.tv_sec = birthtimSec;
		stbuf->st_birthtim.tv_nsec = birthtimNsec;
	}
	else
	{
		stbuf->st_birthtim.tv_sec = ctimSec;
		stbuf->st_birthtim.tv_nsec = ctimNsec;
	}
#if defined(FSP_FUSE_CAP_STAT_EX)
	if (cgofuse_stat_ex)
		((struct fuse_stat_ex *)stbuf)->st_flags = flags;
#endif
#else
	stbuf->st_atim.tv_sec = atimSec; stbuf->st_atim.tv_nsec = atimNsec;
	stbuf->st_mtim.tv_sec = mtimSec; stbuf->st_mtim.tv_nsec = mtimNsec;
	stbuf->st_ctim.tv_sec = ctimSec; stbuf->st_ctim.tv_nsec = ctimNsec;
#endif
}

static inline int hostFilldir(fuse_fill_dir_t filler, void *buf,
	char *name, fuse_stat_t *stbuf, fuse_off_t off)
{
	return filler(buf, name, stbuf, off);
}

#if defined(__APPLE__)
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
#else
#define _hostSetxattr hostSetxattr
#define _hostGetxattr hostGetxattr
#endif

// hostStaticInit, hostFuseInit and hostInit serve different purposes.
//
// hostStaticInit and hostFuseInit are needed to provide static and dynamic initialization
// of the FUSE layer. This is currently useful on Windows only.
//
// hostInit is simply the .init implementation of struct fuse_operations.

static void hostStaticInit(void)
{
#if defined(__APPLE__) || defined(__FreeBSD__) || defined(__linux__)
#elif defined(_WIN32)
	InitializeCriticalSection(&cgofuse_lock);
#endif
}

static int hostFuseInit(void)
{
#if defined(__APPLE__) || defined(__FreeBSD__) || defined(__linux__)
	return 1;
#elif defined(_WIN32)
	return 0 != cgofuse_init_fast(0);
#endif
}

static int hostMountpointOptProc(void *opt_data, const char *arg, int key,
	struct fuse_args *outargs)
{
	char **pmountpoint = opt_data;
	switch (key)
	{
	default:
		return 1;
	case FUSE_OPT_KEY_NONOPT:
		if (0 == *pmountpoint)
		{
			size_t size = strlen(arg) + 1;
			*pmountpoint = malloc(size);
			if (0 == *pmountpoint)
				return -1;
			memcpy(*pmountpoint, arg, size);
		}
		return 1;
	}
}

static const char *hostMountpoint(int argc, char *argv[])
{
	static struct fuse_opt opts[] = { FUSE_OPT_END };
	struct fuse_args args = FUSE_ARGS_INIT(argc, argv);
	char *mountpoint = 0;
	if (-1 == fuse_opt_parse(&args, &mountpoint, opts, hostMountpointOptProc))
		return 0;
	fuse_opt_free_args(&args);
	return mountpoint;
}

static int hostMount(int argc, char *argv[], void *data)
{
#if defined(__GNUC__)
#pragma GCC diagnostic push
#pragma GCC diagnostic ignored "-Wincompatible-pointer-types"
#endif
	static struct fuse_operations fsop =
	{
		.getattr = (int (*)())go_hostGetattr,
		.readlink = (int (*)())go_hostReadlink,
		.mknod = (int (*)())go_hostMknod,
		.mkdir = (int (*)())go_hostMkdir,
		.unlink = (int (*)())go_hostUnlink,
		.rmdir = (int (*)())go_hostRmdir,
		.symlink = (int (*)())go_hostSymlink,
		.rename = (int (*)())go_hostRename,
		.link = (int (*)())go_hostLink,
		.chmod = (int (*)())go_hostChmod,
		.chown = (int (*)())go_hostChown,
		.truncate = (int (*)())go_hostTruncate,
		.open = (int (*)())go_hostOpen,
		.read = (int (*)())go_hostRead,
		.write = (int (*)())go_hostWrite,
		.statfs = (int (*)())go_hostStatfs,
		.flush = (int (*)())go_hostFlush,
		.release = (int (*)())go_hostRelease,
		.fsync = (int (*)())go_hostFsync,
		.setxattr = (int (*)())go_hostSetxattr,
		.getxattr = (int (*)())go_hostGetxattr,
		.listxattr = (int (*)())go_hostListxattr,
		.removexattr = (int (*)())go_hostRemovexattr,
		.opendir = (int (*)())go_hostOpendir,
		.readdir = (int (*)())go_hostReaddir,
		.releasedir = (int (*)())go_hostReleasedir,
		.fsyncdir = (int (*)())go_hostFsyncdir,
		.init = (void *(*)())go_hostInit,
		.destroy = (void (*)())go_hostDestroy,
		.access = (int (*)())go_hostAccess,
		.create = (int (*)())go_hostCreate,
		.ftruncate = (int (*)())go_hostFtruncate,
		.fgetattr = (int (*)())go_hostFgetattr,
		//.lock = (int (*)())go_hostFlock,
		.utimens = (int (*)())go_hostUtimens,
#if defined(__APPLE__) || (defined(_WIN32) && defined(FSP_FUSE_CAP_STAT_EX))
		.setchgtime = (int (*)())go_hostSetchgtime,
		.setcrtime = (int (*)())go_hostSetcrtime,
		.chflags = (int (*)())go_hostChflags,
#endif
	};
#if defined(__GNUC__)
#pragma GCC diagnostic pop
#endif
	return 0 == fuse_main_real(argc, argv, &fsop, sizeof fsop, data);
}

static int hostUnmount(struct fuse *fuse, char *mountpoint)
{
#if defined(__APPLE__) || defined(__FreeBSD__)
	if (0 == mountpoint)
		return 0;
	// darwin,freebsd: unmount is available to non-root
	return 0 == unmount(mountpoint, MNT_FORCE);
#elif defined(__linux__)
	if (0 == mountpoint)
		return 0;
	// linux: try umount2 first in case we are root
	if (0 == umount2(mountpoint, MNT_DETACH))
		return 1;
	// linux: umount2 failed; try fusermount
	char *argv[] =
	{
		"/bin/fusermount",
		"-z",
		"-u",
		mountpoint,
		0,
	};
	pid_t pid = 0;
	int status = 0;
	return
		0 == posix_spawn(&pid, argv[0], 0, 0, argv, 0) &&
		pid == waitpid(pid, &status, 0) &&
		WIFEXITED(status) && 0 == WEXITSTATUS(status);
#elif defined(_WIN32)
	// windows/winfsp: fuse_exit just works from anywhere
	fuse_exit(fuse);
	return 1;
#endif
}

static int hostOptParseOptProc(void *opt_data, const char *arg, int key,
	struct fuse_args *outargs)
{
	switch (key)
	{
	default:
		return 0;
	case FUSE_OPT_KEY_NONOPT:
		return 1;
	}
}

static int hostOptParse(struct fuse_args *args, void *data, const struct fuse_opt opts[],
	bool nonopts)
{
	return fuse_opt_parse(args, data, opts, nonopts ? hostOptParseOptProc : 0);
}
*/
import "C"
import "unsafe"

type (
	C_bool                  = C.bool
	C_char                  = C.char
	C_fuse_dev_t            = C.fuse_dev_t
	C_fuse_fill_dir_t       = C.fuse_fill_dir_t
	C_fuse_gid_t            = C.fuse_gid_t
	C_fuse_mode_t           = C.fuse_mode_t
	C_fuse_off_t            = C.fuse_off_t
	C_fuse_opt_offset_t     = C.fuse_opt_offset_t
	C_fuse_stat_t           = C.fuse_stat_t
	C_fuse_statvfs_t        = C.fuse_statvfs_t
	C_fuse_timespec_t       = C.fuse_timespec_t
	C_fuse_uid_t            = C.fuse_uid_t
	C_int                   = C.int
	C_int16_t               = C.int16_t
	C_int32_t               = C.int32_t
	C_int64_t               = C.int64_t
	C_int8_t                = C.int8_t
	C_size_t                = C.size_t
	C_struct_fuse           = C.struct_fuse
	C_struct_fuse_args      = C.struct_fuse_args
	C_struct_fuse_conn_info = C.struct_fuse_conn_info
	C_struct_fuse_context   = C.struct_fuse_context
	C_struct_fuse_file_info = C.struct_fuse_file_info
	C_struct_fuse_opt       = C.struct_fuse_opt
	C_uint16_t              = C.uint16_t
	C_uint32_t              = C.uint32_t
	C_uint64_t              = C.uint64_t
	C_uint8_t               = C.uint8_t
	C_uintptr_t             = C.uintptr_t
	C_unsigned              = C.unsigned
)

// I would like to do the following, but it is not allowed (Go 1.10):
//     var C_GoString = C.GoString
//
// See https://go-review.googlesource.com/c/gofrontend/+/12543

func C_GoString(s *C_char) string {
	return C.GoString(s)
}
func C_CString(s string) *C_char {
	return C.CString(s)
}

func C_malloc(size C_size_t) unsafe.Pointer {
	return C.malloc(size)
}
func C_calloc(count C_size_t, size C_size_t) unsafe.Pointer {
	return C.calloc(count, size)
}
func C_free(p unsafe.Pointer) {
	C.free(p)
}

func C_fuse_get_context() *C_struct_fuse_context {
	return C.fuse_get_context()
}
func C_fuse_opt_free_args(args *C_struct_fuse_args) {
	C.fuse_opt_free_args(args)
}

func C_hostAsgnCconninfo(conn *C_struct_fuse_conn_info,
	capCaseInsensitive C_bool,
	capReaddirPlus C_bool) {
	C.hostAsgnCconninfo(conn, capCaseInsensitive, capReaddirPlus)
}
func C_hostCstatvfsFromFusestatfs(stbuf *C_fuse_statvfs_t,
	bsize C_uint64_t,
	frsize C_uint64_t,
	blocks C_uint64_t,
	bfree C_uint64_t,
	bavail C_uint64_t,
	files C_uint64_t,
	ffree C_uint64_t,
	favail C_uint64_t,
	fsid C_uint64_t,
	flag C_uint64_t,
	namemax C_uint64_t) {
	C.hostCstatvfsFromFusestatfs(stbuf,
		bsize,
		frsize,
		blocks,
		bfree,
		bavail,
		files,
		ffree,
		favail,
		fsid,
		flag,
		namemax)
}
func C_hostCstatFromFusestat(stbuf *C_fuse_stat_t,
	dev C_uint64_t,
	ino C_uint64_t,
	mode C_uint32_t,
	nlink C_uint32_t,
	uid C_uint32_t,
	gid C_uint32_t,
	rdev C_uint64_t,
	size C_int64_t,
	atimSec C_int64_t, atimNsec C_int64_t,
	mtimSec C_int64_t, mtimNsec C_int64_t,
	ctimSec C_int64_t, ctimNsec C_int64_t,
	blksize C_int64_t,
	blocks C_int64_t,
	birthtimSec C_int64_t, birthtimNsec C_int64_t,
	flags C_uint32_t) {
	C_hostCstatFromFusestat(stbuf,
		dev,
		ino,
		mode,
		nlink,
		uid,
		gid,
		rdev,
		size,
		atimSec,
		atimNsec,
		mtimSec,
		mtimNsec,
		ctimSec,
		ctimNsec,
		blksize,
		blocks,
		birthtimSec,
		birthtimNsec,
		flags)
}
func C_hostFilldir(filler C_fuse_fill_dir_t,
	buf unsafe.Pointer, name *C_char, stbuf *C_fuse_stat_t, off C_fuse_off_t) C_int {
	return C.hostFilldir(filler, buf, name, stbuf, off)
}
func C_hostStaticInit() {
	C.hostStaticInit()
}
func C_hostFuseInit() C_int {
	return C.hostFuseInit()
}
func C_hostMountpoint(argc C_int, argv **C_char) *C_char {
	return C.hostMountpoint(argc, argv)
}
func C_hostMount(argc C_int, argv **C_char, data unsafe.Pointer) C_int {
	return C.hostMount(argc, argv, data)
}
func C_hostUnmount(fuse *C_struct_fuse, mountpoint *C_char) C_int {
	return C.hostUnmount(fuse, mountpoint)
}
func C_hostOptParse(args *C_struct_fuse_args, data unsafe.Pointer, opts *C_struct_fuse_opt,
	nonopts C_bool) C_int {
	return C.hostOptParse(args, data, opts, nonopts)
}

//export go_hostGetattr
func go_hostGetattr(path0 *C_char, stat0 *C_fuse_stat_t) (errc0 C_int) {
	return hostGetattr(path0, stat0)
}

//export go_hostReadlink
func go_hostReadlink(path0 *C_char, buff0 *C_char, size0 C_size_t) (errc0 C_int) {
	return hostReadlink(path0, buff0, size0)
}

//export go_hostMknod
func go_hostMknod(path0 *C_char, mode0 C_fuse_mode_t, dev0 C_fuse_dev_t) (errc0 C_int) {
	return hostMknod(path0, mode0, dev0)
}

//export go_hostMkdir
func go_hostMkdir(path0 *C_char, mode0 C_fuse_mode_t) (errc0 C_int) {
	return hostMkdir(path0, mode0)
}

//export go_hostUnlink
func go_hostUnlink(path0 *C_char) (errc0 C_int) {
	return hostUnlink(path0)
}

//export go_hostRmdir
func go_hostRmdir(path0 *C_char) (errc0 C_int) {
	return hostRmdir(path0)
}

//export go_hostSymlink
func go_hostSymlink(target0 *C_char, newpath0 *C_char) (errc0 C_int) {
	return hostSymlink(target0, newpath0)
}

//export go_hostRename
func go_hostRename(oldpath0 *C_char, newpath0 *C_char) (errc0 C_int) {
	return hostRename(oldpath0, newpath0)
}

//export go_hostLink
func go_hostLink(oldpath0 *C_char, newpath0 *C_char) (errc0 C_int) {
	return hostLink(oldpath0, newpath0)
}

//export go_hostChmod
func go_hostChmod(path0 *C_char, mode0 C_fuse_mode_t) (errc0 C_int) {
	return hostChmod(path0, mode0)
}

//export go_hostChown
func go_hostChown(path0 *C_char, uid0 C_fuse_uid_t, gid0 C_fuse_gid_t) (errc0 C_int) {
	return hostChown(path0, uid0, gid0)
}

//export go_hostTruncate
func go_hostTruncate(path0 *C_char, size0 C_fuse_off_t) (errc0 C_int) {
	return hostTruncate(path0, size0)
}

//export go_hostOpen
func go_hostOpen(path0 *C_char, fi0 *C_struct_fuse_file_info) (errc0 C_int) {
	return hostOpen(path0, fi0)
}

//export go_hostRead
func go_hostRead(path0 *C_char, buff0 *C_char, size0 C_size_t, ofst0 C_fuse_off_t,
	fi0 *C_struct_fuse_file_info) (nbyt0 C_int) {
	return hostRead(path0, buff0, size0, ofst0, fi0)
}

//export go_hostWrite
func go_hostWrite(path0 *C_char, buff0 *C_char, size0 C_size_t, ofst0 C_fuse_off_t,
	fi0 *C_struct_fuse_file_info) (nbyt0 C_int) {
	return hostWrite(path0, buff0, size0, ofst0, fi0)
}

//export go_hostStatfs
func go_hostStatfs(path0 *C_char, stat0 *C_fuse_statvfs_t) (errc0 C_int) {
	return hostStatfs(path0, stat0)
}

//export go_hostFlush
func go_hostFlush(path0 *C_char, fi0 *C_struct_fuse_file_info) (errc0 C_int) {
	return hostFlush(path0, fi0)
}

//export go_hostRelease
func go_hostRelease(path0 *C_char, fi0 *C_struct_fuse_file_info) (errc0 C_int) {
	return hostRelease(path0, fi0)
}

//export go_hostFsync
func go_hostFsync(path0 *C_char, datasync C_int, fi0 *C_struct_fuse_file_info) (errc0 C_int) {
	return hostFsync(path0, datasync, fi0)
}

//export go_hostSetxattr
func go_hostSetxattr(path0 *C_char, name0 *C_char, buff0 *C_char, size0 C_size_t,
	flags C_int) (errc0 C_int) {
	return hostSetxattr(path0, name0, buff0, size0, flags)
}

//export go_hostGetxattr
func go_hostGetxattr(path0 *C_char, name0 *C_char, buff0 *C_char, size0 C_size_t) (nbyt0 C_int) {
	return hostGetxattr(path0, name0, buff0, size0)
}

//export go_hostListxattr
func go_hostListxattr(path0 *C_char, buff0 *C_char, size0 C_size_t) (nbyt0 C_int) {
	return hostListxattr(path0, buff0, size0)
}

//export go_hostRemovexattr
func go_hostRemovexattr(path0 *C_char, name0 *C_char) (errc0 C_int) {
	return hostRemovexattr(path0, name0)
}

//export go_hostOpendir
func go_hostOpendir(path0 *C_char, fi0 *C_struct_fuse_file_info) (errc0 C_int) {
	return hostOpendir(path0, fi0)
}

//export go_hostReaddir
func go_hostReaddir(path0 *C_char,
	buff0 unsafe.Pointer, fill0 C_fuse_fill_dir_t, ofst0 C_fuse_off_t,
	fi0 *C_struct_fuse_file_info) (errc0 C_int) {
	return hostReaddir(path0, buff0, fill0, ofst0, fi0)
}

//export go_hostReleasedir
func go_hostReleasedir(path0 *C_char, fi0 *C_struct_fuse_file_info) (errc0 C_int) {
	return hostReleasedir(path0, fi0)
}

//export go_hostFsyncdir
func go_hostFsyncdir(path0 *C_char, datasync C_int, fi0 *C_struct_fuse_file_info) (errc0 C_int) {
	return hostFsyncdir(path0, datasync, fi0)
}

//export go_hostInit
func go_hostInit(conn0 *C_struct_fuse_conn_info) (user_data unsafe.Pointer) {
	return hostInit(conn0)
}

//export go_hostDestroy
func go_hostDestroy(user_data unsafe.Pointer) {
	hostDestroy(user_data)
}

//export go_hostAccess
func go_hostAccess(path0 *C_char, mask0 C_int) (errc0 C_int) {
	return hostAccess(path0, mask0)
}

//export go_hostCreate
func go_hostCreate(path0 *C_char, mode0 C_fuse_mode_t, fi0 *C_struct_fuse_file_info) (errc0 C_int) {
	return hostCreate(path0, mode0, fi0)
}

//export go_hostFtruncate
func go_hostFtruncate(path0 *C_char, size0 C_fuse_off_t,
	fi0 *C_struct_fuse_file_info) (errc0 C_int) {
	return hostFtruncate(path0, size0, fi0)
}

//export go_hostFgetattr
func go_hostFgetattr(path0 *C_char, stat0 *C_fuse_stat_t,
	fi0 *C_struct_fuse_file_info) (errc0 C_int) {
	return hostFgetattr(path0, stat0, fi0)
}

//export go_hostUtimens
func go_hostUtimens(path0 *C_char, tmsp0 *C_fuse_timespec_t) (errc0 C_int) {
	return hostUtimens(path0, tmsp0)
}

//export go_hostSetchgtime
func go_hostSetchgtime(path0 *C_char, tmsp0 *C_fuse_timespec_t) (errc0 C_int) {
	return hostSetchgtime(path0, tmsp0)
}

//export go_hostSetcrtime
func go_hostSetcrtime(path0 *C_char, tmsp0 *C_fuse_timespec_t) (errc0 C_int) {
	return hostSetcrtime(path0, tmsp0)
}

//export go_hostChflags
func go_hostChflags(path0 *C_char, flags C_uint32_t) (errc0 C_int) {
	return hostChflags(path0, flags)
}

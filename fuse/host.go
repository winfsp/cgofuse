/*
 * host.go
 *
 * Copyright 2017 Bill Zissimopoulos
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
#cgo linux CFLAGS: -DFUSE_USE_VERSION=28 -D_FILE_OFFSET_BITS=64 -I/usr/include/fuse
#cgo linux LDFLAGS: -lfuse
#cgo windows CFLAGS: -D_WIN32_WINNT=0x0600 -DFUSE_USE_VERSION=28

#if !(defined(__APPLE__) || defined(__linux__) || defined(_WIN32))
#error platform not supported
#endif

#include <stdlib.h>
#include <string.h>

#if defined(__APPLE__) || defined(__linux__)

#include <spawn.h>
#include <sys/mount.h>
#include <sys/wait.h>
#include <fuse.h>

#elif defined(_WIN32)

#include <windows.h>

static PVOID cgofuse_init_winfsp(VOID);
static PVOID cgofuse_init_fail();
static inline VOID cgofuse_init(VOID)
{
	static SRWLOCK Lock = SRWLOCK_INIT;
	static PVOID _Module = 0;
	PVOID Module = _Module;
	MemoryBarrier();
	if (0 == Module)
	{
		AcquireSRWLockExclusive(&Lock);
		Module = _Module;
		if (0 == Module)
		{
			Module = cgofuse_init_winfsp();
			MemoryBarrier();
			_Module = Module;
		}
		ReleaseSRWLockExclusive(&Lock);
	}
}

#define FSP_FUSE_API                    static
#define FSP_FUSE_API_NAME(api)          (* pfn_ ## api)
#define FSP_FUSE_API_CALL(api)          (cgofuse_init(), pfn_ ## api)
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

	WINADVAPI
	LSTATUS
	APIENTRY
	RegGetValueW(
		HKEY hkey,
		LPCWSTR lpSubKey,
		LPCWSTR lpValue,
		DWORD dwFlags,
		LPDWORD pdwType,
		PVOID pvData,
		LPDWORD pcbData);

	WCHAR PathBuf[MAX_PATH];
	DWORD Size;
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
			Result = RegGetValueW(RegKey, 0, L"InstallDir",
				RRF_RT_REG_SZ, 0, PathBuf, &Size);
			RegCloseKey(RegKey);
		}
		if (ERROR_SUCCESS != Result)
			return 0xC0000034;//STATUS_OBJECT_NAME_NOT_FOUND

		RtlCopyMemory(PathBuf + (Size / sizeof(WCHAR) - 1), L"" FSP_DLLPATH, sizeof L"" FSP_DLLPATH);
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
		return cgofuse_init_fail();

static PVOID cgofuse_init_winfsp(VOID)
{
	PVOID Module;
	NTSTATUS Result;

	Result = FspLoad(&Module);
	if (0 > Result)
		return cgofuse_init_fail();

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

static PVOID cgofuse_init_fail()
{
	return 0;
}

#endif

#if defined(__APPLE__) || defined(__linux__)
typedef struct stat fuse_stat_t;
typedef struct statvfs fuse_statvfs_t;
typedef struct timespec fuse_timespec_t;
typedef mode_t fuse_mode_t;
typedef dev_t fuse_dev_t;
typedef uid_t fuse_uid_t;
typedef gid_t fuse_gid_t;
typedef off_t fuse_off_t;
#elif defined(_WIN32)
typedef struct fuse_stat fuse_stat_t;
typedef struct fuse_statvfs fuse_statvfs_t;
typedef struct fuse_timespec fuse_timespec_t;
#endif

extern int hostGetattr(char *path, fuse_stat_t *stbuf);
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
extern int hostStatfs(char *path, fuse_statvfs_t *stbuf);
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
extern int hostFgetattr(char *path, fuse_stat_t *stbuf, struct fuse_file_info *fi);
//extern int hostLock(char *path, struct fuse_file_info *fi, int cmd, struct fuse_flock *lock);
extern int hostUtimens(char *path, fuse_timespec_t tv[2]);

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
#if defined(__APPLE__)
	stbuf->st_atimespec.tv_sec = atimSec; stbuf->st_atimespec.tv_nsec = atimNsec;
	stbuf->st_mtimespec.tv_sec = mtimSec; stbuf->st_mtimespec.tv_nsec = mtimNsec;
	stbuf->st_ctimespec.tv_sec = ctimSec; stbuf->st_ctimespec.tv_nsec = ctimNsec;
	stbuf->st_birthtimespec.tv_sec = birthtimSec; stbuf->st_birthtimespec.tv_nsec = birthtimNsec;
#else
	stbuf->st_atim.tv_sec = atimSec; stbuf->st_atim.tv_nsec = atimNsec;
	stbuf->st_mtim.tv_sec = mtimSec; stbuf->st_mtim.tv_nsec = mtimNsec;
	stbuf->st_ctim.tv_sec = ctimSec; stbuf->st_ctim.tv_nsec = ctimNsec;
#endif
	stbuf->st_blksize = blksize;
	stbuf->st_blocks = blocks;
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

static const char *hostMountpoint(int argc, char *argv[])
{
	struct fuse_args args = FUSE_ARGS_INIT(argc, argv);
	char *mountpoint;
	if (-1 == fuse_parse_cmdline(&args, &mountpoint, 0, 0))
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
		.getattr = (int (*)())hostGetattr,
		.readlink = (int (*)())hostReadlink,
		.mknod = (int (*)())hostMknod,
		.mkdir = (int (*)())hostMkdir,
		.unlink = (int (*)())hostUnlink,
		.rmdir = (int (*)())hostRmdir,
		.symlink = (int (*)())hostSymlink,
		.rename = (int (*)())hostRename,
		.link = (int (*)())hostLink,
		.chmod = (int (*)())hostChmod,
		.chown = (int (*)())hostChown,
		.truncate = (int (*)())hostTruncate,
		.open = (int (*)())hostOpen,
		.read = (int (*)())hostRead,
		.write = (int (*)())hostWrite,
		.statfs = (int (*)())hostStatfs,
		.flush = (int (*)())hostFlush,
		.release = (int (*)())hostRelease,
		.fsync = (int (*)())hostFsync,
		.setxattr = (int (*)())_hostSetxattr,
		.getxattr = (int (*)())_hostGetxattr,
		.listxattr = (int (*)())hostListxattr,
		.removexattr = (int (*)())hostRemovexattr,
		.opendir = (int (*)())hostOpendir,
		.readdir = (int (*)())hostReaddir,
		.releasedir = (int (*)())hostReleasedir,
		.fsyncdir = (int (*)())hostFsyncdir,
		.init = (void *(*)())hostInit,
		.destroy = (void (*)())hostDestroy,
		.access = (int (*)())hostAccess,
		.create = (int (*)())hostCreate,
		.ftruncate = (int (*)())hostFtruncate,
		.fgetattr = (int (*)())hostFgetattr,
		//.lock = (int (*)())hostFlock,
		.utimens = (int (*)())hostUtimens,
	};
#if defined(__GNUC__)
#pragma GCC diagnostic pop
#endif
	return 0 == fuse_main_real(argc, argv, &fsop, sizeof fsop, data);
}

static int hostUnmount(struct fuse *fuse, char *mountpoint)
{
#if defined(__APPLE__)
	if (0 == mountpoint)
		return 0;
	// darwin: unmount is available to non-root
	return 0 == unmount(mountpoint, MNT_FORCE);
#elif defined(__linux__)
	if (0 == mountpoint)
		return 0;
	// linux: try umount2 first in case we are root
	if (0 == umount2(mountpoint, MNT_FORCE))
		return 1;
	// linux: umount2 failed; try fusermount
	char *argv[] =
	{
		"/bin/fusermount",
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
*/
import "C"
import "unsafe"

// FileSystemHost is used to host a file system.
type FileSystemHost struct {
	fsop FileSystemInterface
	hndl unsafe.Pointer
	fuse *C.struct_fuse
	mntp *C.char
}

func copyCstatvfsFromFusestatfs(dst *C.fuse_statvfs_t, src *Statfs_t) {
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

func copyCstatFromFusestat(dst *C.fuse_stat_t, src *Stat_t) {
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

func copyFusetimespecFromCtimespec(dst *Timespec, src *C.fuse_timespec_t) {
	dst.Sec = int64(src.tv_sec)
	dst.Nsec = int64(src.tv_nsec)
}

func recoverAsErrno(errc0 *C.int) {
	if r := recover(); nil != r {
		switch e := r.(type) {
		case Error:
			*errc0 = C.int(e)
		default:
			*errc0 = -C.int(EIO)
		}
	}
}

//export hostGetattr
func hostGetattr(path0 *C.char, stat0 *C.fuse_stat_t) (errc0 C.int) {
	defer recoverAsErrno(&errc0)
	fsop := getInterfaceForHandle(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	stat := &Stat_t{}
	errc := fsop.Getattr(path, stat, ^uint64(0))
	copyCstatFromFusestat(stat0, stat)
	return C.int(errc)
}

//export hostReadlink
func hostReadlink(path0 *C.char, buff0 *C.char, size0 C.size_t) (errc0 C.int) {
	defer recoverAsErrno(&errc0)
	fsop := getInterfaceForHandle(C.fuse_get_context().private_data).(FileSystemInterface)
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
func hostMknod(path0 *C.char, mode0 C.fuse_mode_t, dev0 C.fuse_dev_t) (errc0 C.int) {
	defer recoverAsErrno(&errc0)
	fsop := getInterfaceForHandle(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	errc := fsop.Mknod(path, uint32(mode0), uint64(dev0))
	return C.int(errc)
}

//export hostMkdir
func hostMkdir(path0 *C.char, mode0 C.fuse_mode_t) (errc0 C.int) {
	defer recoverAsErrno(&errc0)
	fsop := getInterfaceForHandle(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	errc := fsop.Mkdir(path, uint32(mode0))
	return C.int(errc)
}

//export hostUnlink
func hostUnlink(path0 *C.char) (errc0 C.int) {
	defer recoverAsErrno(&errc0)
	fsop := getInterfaceForHandle(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	errc := fsop.Unlink(path)
	return C.int(errc)
}

//export hostRmdir
func hostRmdir(path0 *C.char) (errc0 C.int) {
	defer recoverAsErrno(&errc0)
	fsop := getInterfaceForHandle(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	errc := fsop.Rmdir(path)
	return C.int(errc)
}

//export hostSymlink
func hostSymlink(target0 *C.char, newpath0 *C.char) (errc0 C.int) {
	defer recoverAsErrno(&errc0)
	fsop := getInterfaceForHandle(C.fuse_get_context().private_data).(FileSystemInterface)
	target, newpath := C.GoString(target0), C.GoString(newpath0)
	errc := fsop.Symlink(target, newpath)
	return C.int(errc)
}

//export hostRename
func hostRename(oldpath0 *C.char, newpath0 *C.char) (errc0 C.int) {
	defer recoverAsErrno(&errc0)
	fsop := getInterfaceForHandle(C.fuse_get_context().private_data).(FileSystemInterface)
	oldpath, newpath := C.GoString(oldpath0), C.GoString(newpath0)
	errc := fsop.Rename(oldpath, newpath)
	return C.int(errc)
}

//export hostLink
func hostLink(oldpath0 *C.char, newpath0 *C.char) (errc0 C.int) {
	defer recoverAsErrno(&errc0)
	fsop := getInterfaceForHandle(C.fuse_get_context().private_data).(FileSystemInterface)
	oldpath, newpath := C.GoString(oldpath0), C.GoString(newpath0)
	errc := fsop.Link(oldpath, newpath)
	return C.int(errc)
}

//export hostChmod
func hostChmod(path0 *C.char, mode0 C.fuse_mode_t) (errc0 C.int) {
	defer recoverAsErrno(&errc0)
	fsop := getInterfaceForHandle(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	errc := fsop.Chmod(path, uint32(mode0))
	return C.int(errc)
}

//export hostChown
func hostChown(path0 *C.char, uid0 C.fuse_uid_t, gid0 C.fuse_gid_t) (errc0 C.int) {
	defer recoverAsErrno(&errc0)
	fsop := getInterfaceForHandle(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	errc := fsop.Chown(path, uint32(uid0), uint32(gid0))
	return C.int(errc)
}

//export hostTruncate
func hostTruncate(path0 *C.char, size0 C.fuse_off_t) (errc0 C.int) {
	defer recoverAsErrno(&errc0)
	fsop := getInterfaceForHandle(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	errc := fsop.Truncate(path, int64(size0), ^uint64(0))
	return C.int(errc)
}

//export hostOpen
func hostOpen(path0 *C.char, fi0 *C.struct_fuse_file_info) (errc0 C.int) {
	defer recoverAsErrno(&errc0)
	fsop := getInterfaceForHandle(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	errc, rslt := fsop.Open(path, int(fi0.flags))
	fi0.fh = C.uint64_t(rslt)
	return C.int(errc)
}

//export hostRead
func hostRead(path0 *C.char, buff0 *C.char, size0 C.size_t, ofst0 C.fuse_off_t,
	fi0 *C.struct_fuse_file_info) (nbyt0 C.int) {
	defer recoverAsErrno(&nbyt0)
	fsop := getInterfaceForHandle(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	buff := (*[1 << 30]byte)(unsafe.Pointer(buff0))
	nbyt := fsop.Read(path, buff[:size0], int64(ofst0), uint64(fi0.fh))
	return C.int(nbyt)
}

//export hostWrite
func hostWrite(path0 *C.char, buff0 *C.char, size0 C.size_t, ofst0 C.fuse_off_t,
	fi0 *C.struct_fuse_file_info) (nbyt0 C.int) {
	defer recoverAsErrno(&nbyt0)
	fsop := getInterfaceForHandle(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	buff := (*[1 << 30]byte)(unsafe.Pointer(buff0))
	nbyt := fsop.Write(path, buff[:size0], int64(ofst0), uint64(fi0.fh))
	return C.int(nbyt)
}

//export hostStatfs
func hostStatfs(path0 *C.char, stat0 *C.fuse_statvfs_t) (errc0 C.int) {
	defer recoverAsErrno(&errc0)
	fsop := getInterfaceForHandle(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	stat := &Statfs_t{}
	errc := fsop.Statfs(path, stat)
	if -ENOSYS == errc {
		stat = &Statfs_t{}
		errc = 0
	}
	copyCstatvfsFromFusestatfs(stat0, stat)
	return C.int(errc)
}

//export hostFlush
func hostFlush(path0 *C.char, fi0 *C.struct_fuse_file_info) (errc0 C.int) {
	defer recoverAsErrno(&errc0)
	fsop := getInterfaceForHandle(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	errc := fsop.Flush(path, uint64(fi0.fh))
	return C.int(errc)
}

//export hostRelease
func hostRelease(path0 *C.char, fi0 *C.struct_fuse_file_info) (errc0 C.int) {
	defer recoverAsErrno(&errc0)
	fsop := getInterfaceForHandle(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	errc := fsop.Release(path, uint64(fi0.fh))
	return C.int(errc)
}

//export hostFsync
func hostFsync(path0 *C.char, datasync C.int, fi0 *C.struct_fuse_file_info) (errc0 C.int) {
	defer recoverAsErrno(&errc0)
	fsop := getInterfaceForHandle(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	errc := fsop.Fsync(path, 0 != datasync, uint64(fi0.fh))
	if -ENOSYS == errc {
		errc = 0
	}
	return C.int(errc)
}

//export hostSetxattr
func hostSetxattr(path0 *C.char, name0 *C.char, buff0 *C.char, size0 C.size_t,
	flags C.int) (errc0 C.int) {
	defer recoverAsErrno(&errc0)
	fsop := getInterfaceForHandle(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	name := C.GoString(name0)
	buff := (*[1 << 30]byte)(unsafe.Pointer(buff0))
	errc := fsop.Setxattr(path, name, buff[:size0], int(flags))
	return C.int(errc)
}

//export hostGetxattr
func hostGetxattr(path0 *C.char, name0 *C.char, buff0 *C.char, size0 C.size_t) (nbyt0 C.int) {
	defer recoverAsErrno(&nbyt0)
	fsop := getInterfaceForHandle(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	name := C.GoString(name0)
	buff := (*[1 << 30]byte)(unsafe.Pointer(buff0))
	size := int(size0)
	nbyt := 0
	fill := func(value []byte) bool {
		nbyt = len(value)
		if 0 != size {
			if nbyt > size {
				return false
			}
			copy(buff[:size], value)
		}
		return true
	}
	errc := fsop.Getxattr(path, name, fill)
	if 0 != errc {
		return C.int(errc)
	}
	return C.int(nbyt)
}

//export hostListxattr
func hostListxattr(path0 *C.char, buff0 *C.char, size0 C.size_t) (nbyt0 C.int) {
	defer recoverAsErrno(&nbyt0)
	fsop := getInterfaceForHandle(C.fuse_get_context().private_data).(FileSystemInterface)
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
	fsop := getInterfaceForHandle(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	name := C.GoString(name0)
	errc := fsop.Removexattr(path, name)
	return C.int(errc)
}

//export hostOpendir
func hostOpendir(path0 *C.char, fi0 *C.struct_fuse_file_info) (errc0 C.int) {
	defer recoverAsErrno(&errc0)
	fsop := getInterfaceForHandle(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	errc, rslt := fsop.Opendir(path)
	if -ENOSYS == errc {
		errc = 0
	}
	fi0.fh = C.uint64_t(rslt)
	return C.int(errc)
}

//export hostReaddir
func hostReaddir(path0 *C.char, buff0 unsafe.Pointer, fill0 C.fuse_fill_dir_t, ofst0 C.fuse_off_t,
	fi0 *C.struct_fuse_file_info) (errc0 C.int) {
	defer recoverAsErrno(&errc0)
	fsop := getInterfaceForHandle(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	fill := func(name1 string, stat1 *Stat_t, off1 int64) bool {
		name := C.CString(name1)
		defer C.free(unsafe.Pointer(name))
		if nil == stat1 {
			return 0 == C.hostFilldir(fill0, buff0, name, nil, C.fuse_off_t(off1))
		} else {
			stat := C.fuse_stat_t{}
			copyCstatFromFusestat(&stat, stat1)
			return 0 == C.hostFilldir(fill0, buff0, name, &stat, C.fuse_off_t(off1))
		}
	}
	errc := fsop.Readdir(path, fill, int64(ofst0), uint64(fi0.fh))
	return C.int(errc)
}

//export hostReleasedir
func hostReleasedir(path0 *C.char, fi0 *C.struct_fuse_file_info) (errc0 C.int) {
	defer recoverAsErrno(&errc0)
	fsop := getInterfaceForHandle(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	errc := fsop.Releasedir(path, uint64(fi0.fh))
	return C.int(errc)
}

//export hostFsyncdir
func hostFsyncdir(path0 *C.char, datasync C.int, fi0 *C.struct_fuse_file_info) (errc0 C.int) {
	defer recoverAsErrno(&errc0)
	fsop := getInterfaceForHandle(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	errc := fsop.Fsyncdir(path, 0 != datasync, uint64(fi0.fh))
	if -ENOSYS == errc {
		errc = 0
	}
	return C.int(errc)
}

//export hostInit
func hostInit(conn0 *C.struct_fuse_conn_info) (user_data unsafe.Pointer) {
	defer recover()
	fctx := C.fuse_get_context()
	host := getInterfaceForHandle(fctx.private_data).(*FileSystemHost)
	host.fuse = fctx.fuse
	user_data = host.hndl
	host.fsop.Init()
	return
}

//export hostDestroy
func hostDestroy(user_data unsafe.Pointer) {
	defer recover()
	fsop := getInterfaceForHandle(user_data).(FileSystemInterface)
	fsop.Destroy()
}

//export hostAccess
func hostAccess(path0 *C.char, mask0 C.int) (errc0 C.int) {
	defer recoverAsErrno(&errc0)
	fsop := getInterfaceForHandle(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	errc := fsop.Access(path, uint32(mask0))
	return C.int(errc)
}

//export hostCreate
func hostCreate(path0 *C.char, mode0 C.fuse_mode_t, fi0 *C.struct_fuse_file_info) (errc0 C.int) {
	defer recoverAsErrno(&errc0)
	fsop := getInterfaceForHandle(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	errc, rslt := fsop.Create(path, int(fi0.flags), uint32(mode0))
	if -ENOSYS == errc {
		errc = fsop.Mknod(path, S_IFREG|uint32(mode0), 0)
		if 0 == errc {
			errc, rslt = fsop.Open(path, int(fi0.flags))
		}
	}
	fi0.fh = C.uint64_t(rslt)
	return C.int(errc)
}

//export hostFtruncate
func hostFtruncate(path0 *C.char, size0 C.fuse_off_t, fi0 *C.struct_fuse_file_info) (errc0 C.int) {
	defer recoverAsErrno(&errc0)
	fsop := getInterfaceForHandle(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	errc := fsop.Truncate(path, int64(size0), uint64(fi0.fh))
	return C.int(errc)
}

//export hostFgetattr
func hostFgetattr(path0 *C.char, stat0 *C.fuse_stat_t,
	fi0 *C.struct_fuse_file_info) (errc0 C.int) {
	defer recoverAsErrno(&errc0)
	fsop := getInterfaceForHandle(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	stat := &Stat_t{}
	errc := fsop.Getattr(path, stat, uint64(fi0.fh))
	copyCstatFromFusestat(stat0, stat)
	return C.int(errc)
}

//export hostUtimens
func hostUtimens(path0 *C.char, tmsp0 *C.fuse_timespec_t) (errc0 C.int) {
	defer recoverAsErrno(&errc0)
	fsop := getInterfaceForHandle(C.fuse_get_context().private_data).(FileSystemInterface)
	path := C.GoString(path0)
	if nil == tmsp0 {
		errc := fsop.Utimens(path, nil)
		return C.int(errc)
	} else {
		tmsp := [2]Timespec{}
		tmsa := (*[2]C.fuse_timespec_t)(unsafe.Pointer(tmsp0))
		copyFusetimespecFromCtimespec(&tmsp[0], &tmsa[0])
		copyFusetimespecFromCtimespec(&tmsp[1], &tmsa[1])
		errc := fsop.Utimens(path, tmsp[:])
		return C.int(errc)
	}
}

// NewFileSystemHost creates a file system host.
func NewFileSystemHost(fsop FileSystemInterface) *FileSystemHost {
	return &FileSystemHost{fsop, nil, nil, nil}
}

// Mount mounts a file system.
// The file system is considered mounted only after its Init() method has been called.
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
	host.hndl = newHandleForInterface(host.fsop)
	defer delHandleForInterface(host.hndl)
	hosthndl := newHandleForInterface(host)
	defer delHandleForInterface(hosthndl)
	host.mntp = C.hostMountpoint(C.int(argc), &argv[0])
	defer func() {
		C.free(unsafe.Pointer(host.mntp))
		host.mntp = nil
		host.fuse = nil
	}()
	return 0 != C.hostMount(C.int(argc), &argv[0], hosthndl)
}

// Unmount unmounts a mounted file system.
// Unmount may be called at any time after the Init() method has been called.
func (host *FileSystemHost) Unmount() bool {
	if nil == host.fuse {
		return false
	}
	return 0 != C.hostUnmount(host.fuse, host.mntp)
}

// Getcontext gets information related to a file system operation.
func Getcontext() (uid uint32, gid uint32, pid int) {
	uid = uint32(C.fuse_get_context().uid)
	gid = uint32(C.fuse_get_context().gid)
	pid = int(C.fuse_get_context().pid)
	return
}

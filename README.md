<div align="center">

# Cross-platform FUSE library for Go

<a href="https://godoc.org/github.com/winfsp/cgofuse/fuse">
    <img src="https://godoc.org/github.com/winfsp/cgofuse/fuse?status.svg"/>
</a>

_Cgofuse is a cross-platform FUSE library for Go. It is supported on multiple platforms and can be ported to any platform that has a FUSE implementation. It has [cgo](https://pkg.go.dev/cmd/cgo), [!cgo](https://go.dev/wiki/WindowsDLLs) (aka "nocgo") and [FUSE](https://github.com/libfuse/libfuse/tree/fuse-2.9.9) (aka "FUSE2"), [FUSE3](https://github.com/libfuse/libfuse) variants depending on the platform._

|            |**cgo**     |**!cgo**    |**FUSE**    |**FUSE3**   |
|:----------:|:----------:|:----------:|:----------:|:----------:|
|**Windows** |&#x2713;    |&#x2713;    |&#x2713;    |            |
|**macOS**   |&#x2713;    |            |&#x2713;    |&#x2713;    |
|**Linux**   |&#x2713;    |            |&#x2713;    |&#x2713;    |
|**FreeBSD** |&#x2713;    |            |&#x2713;    |&#x2713;    |
|**NetBSD**  |&#x2713;    |            |&#x2713;    |            |
|**OpenBSD** |&#x2713;    |            |&#x2713;    |            |

</div>

## How to build

**Windows**
- Prerequisites: [WinFsp](https://github.com/winfsp/winfsp), gcc (e.g. from [Mingw-builds](http://mingw-w64.org/doku.php/download))
- Build **cgo**:
    ```
    > set CPATH=C:\Program Files (x86)\WinFsp\inc\fuse
    > go install -v ./fuse ./examples/memfs
    ```
- Build **!cgo**:
    ```
    > set CGO_ENABLED=0
    > go install -v ./fuse ./examples/memfs
    ```

**macOS**
- Prerequisites: [macFUSE](https://macfuse.github.io/) or [FUSE-T](https://www.fuse-t.org/), [command line tools](https://developer.apple.com/library/content/technotes/tn2339/_index.html). FUSE3 requires a distribution that includes the FUSE3 headers and library.
- Build **FUSE** (default):
    ```
    $ go install -v ./fuse ./examples/memfs ./examples/passthrough
    ```
- Build **FUSE3**:
    ```
    $ go install -tags=fuse3 -v ./fuse ./examples/memfs ./examples/passthrough
    ```
- FUSE3 on macOS uses the vanilla FUSE3 API, so the macOS-only `Setchgtime`, `Setcrtime` and `Chflags` callbacks are never invoked. Use a FUSE2 build if your file system needs them.
- A FUSE3 build can mount through [FSKit](#fskit-backend-macos) instead of NFS, which produces a local rather than a network volume.

**Linux**
- Prerequisites: libfuse-dev, libfuse3-dev, gcc
- Build **FUSE**:
    ```
    $ go install -v ./fuse ./examples/memfs ./examples/passthrough
    ```
- Build **FUSE3**:
    ```
    $ go install -tags=fuse3 -v ./fuse ./examples/memfs ./examples/passthrough
    ```

**FreeBSD**
- Prerequisites: fusefs-libs, fusefs-libs3
- Build **FUSE**:
    ```
    $ go install -v ./fuse ./examples/memfs ./examples/passthrough
    ```
- Build **FUSE3**:
    ```
    $ go install -tags=fuse3 -v ./fuse ./examples/memfs ./examples/passthrough
    ```
- Note: You may also need the following in order to run FUSE file system. Commands must be run as root.
    ```
    $ vi /boot/loader.conf                      # add: fuse_load="YES"
    $ sysctl vfs.usermount=1                    # allow user mounts
    $ pw usermod USERNAME -G operator           # allow user to open /dev/fuse
    ```

**NetBSD**
- Prerequisites: NONE
- Build:
    ```
    $ go install -v ./fuse ./examples/memfs ./examples/passthrough
    ```
- Note: You may also need the following in order to run FUSE file systems. Commands must be run as root.
    ```
    $ chmod go+rw /dev/puffs
    $ sysctl -w vfs.generic.usermount=1
    ```
- Note: Support is experimental. There are known issues that stem from the differences in the NetBSD [librefuse](https://github.com/NetBSD/src/tree/bbc46b99bff565d75f55fb23b51eff511068b183/lib/librefuse) implementation from the reference [libfuse](https://github.com/libfuse/libfuse) implementation.

**OpenBSD**
- Prerequisites: NONE
- Build:
    ```
    $ go install -v ./fuse ./examples/memfs ./examples/passthrough
    ```
- Note: OpenBSD 6 removed the `kern.usermount` option, which allowed non-root users to mount file systems [[link](https://undeadly.org/cgi?action=article&sid=20160715125022&mode=expanded&count=0)]. Therefore you must be root in order to use FUSE and cgofuse.
- Note: Support is experimental. There are known issues that stem from the differences in the OpenBSD [libfuse](https://github.com/openbsd/src/tree/dae5ffec5618b0b660e9064e3b0991bb4ab1b1e8/lib/libfuse) implementation from the reference [libfuse](https://github.com/libfuse/libfuse) implementation.


## How to use

User mode file systems are expected to implement `fuse.FileSystemInterface`. To make implementation simpler a file system can embed ("inherit") a `fuse.FileSystemBase` which provides default implementations for all operations. To mount a file system one must instantiate a `fuse.FileSystemHost` using `fuse.NewFileSystemHost`.

The full documentation is available at GoDoc.org: [package fuse](https://godoc.org/github.com/winfsp/cgofuse/fuse)

There are currently three example file systems:

- [Hellofs](examples/hellofs/hellofs.go) is an extremely simple file system. Runs on all OS'es.
- [Memfs](examples/memfs/memfs.go) is an in memory file system. Runs on all OS'es.
- [Passthrough](examples/passthrough/passthrough.go) is a file system that passes all operations to the underlying file system. Runs on all OS'es except Windows.
- [Notifyfs](examples/notifyfs/notifyfs.go) is a file system that can issue file change notifications. Runs on Windows only.

## FSKit backend (macOS)

FSKit requires macOS 15 or later, a FUSE3 build of cgofuse, and either macFUSE 5.2 or later or
FUSE-T 1.2 or later. It is selected with a mount option:

```
$ ./memfs -o backend=fskit,volname=MyVolume mnt
```

The option is interpreted by the FUSE3 library rather than by cgofuse, so the cgofuse API is
unchanged: pass it in the `opts` argument of `Mount`, or on the command line of a file system
that forwards its arguments.

## How it is tested

The following software is being used to test cgofuse.

**Windows (cgo and !cgo)**
- [winfsp-tests](https://github.com/winfsp/winfsp/tree/master/tst/winfsp-tests)
- [fsx](https://github.com/billziss-gh/secfs.test/tree/master/fstools/src/fsx)

**macOS**
- [fstest](https://github.com/billziss-gh/secfs.test/tree/master/fstest/ntfs-3g-pjd-fstest-8af5670)
- [fsx](https://github.com/billziss-gh/secfs.test/tree/master/fstools/src/fsx)

**Linux**
- [fstest](https://github.com/billziss-gh/secfs.test/tree/master/fstest/ntfs-3g-pjd-fstest-8af5670)
- [fsx](https://github.com/billziss-gh/secfs.test/tree/master/fstools/src/fsx)

**FreeBSD**
- [fsx](https://github.com/billziss-gh/secfs.test/tree/master/fstools/src/fsx)

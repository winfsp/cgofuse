<h1 align="center">
    Cross-platform FUSE library for Go
</h1>

<p align="center">
    <a href="https://godoc.org/github.com/winfsp/cgofuse/fuse">
        <img src="https://godoc.org/github.com/winfsp/cgofuse/fuse?status.svg"/>
    </a>
</p>

Cgofuse is a cross-platform FUSE library for Go. It is supported on multiple platforms and can be ported to any platform that has a FUSE implementation. It has [cgo](https://golang.org/cmd/cgo/) and [!cgo](https://github.com/golang/go/wiki/WindowsDLLs) ("nocgo") variants depending on the platform.

|       |Windows<br/>[![](https://img.shields.io/github/workflow/status/winfsp/cgofuse/test)](https://github.com/winfsp/cgofuse/actions/workflows/test.yml)|macOS<br/>[![](https://img.shields.io/github/workflow/status/winfsp/cgofuse/test)](https://github.com/winfsp/cgofuse/actions/workflows/test.yml)|Linux<br/>[![](https://img.shields.io/github/workflow/status/winfsp/cgofuse/test)](https://github.com/winfsp/cgofuse/actions/workflows/test.yml)|FreeBSD<br/>[![no CI](https://img.shields.io/badge/build-none-lightgrey.svg)](https://cirrus-ci.com/github/billziss-gh/cgofuse)|NetBSD<sup>*</sup><br/>![no CI](https://img.shields.io/badge/build-none-lightgrey.svg)|OpenBSD<sup>*</sup><br/>![no CI](https://img.shields.io/badge/build-none-lightgrey.svg)|
|:-----:|:------:|:------:|:------:|:------:|:------:|:------:|
|  cgo  |&#x2713;|&#x2713;|&#x2713;|&#x2713;|&#x2713;|&#x2713;|
| !cgo  |&#x2713;|        |        |        |        |        |

**\*** NetBSD and OpenBSD support is experimental. There are known issues that stem from the differences in the NetBSD [librefuse](https://github.com/NetBSD/src/tree/bbc46b99bff565d75f55fb23b51eff511068b183/lib/librefuse) and OpenBSD [libfuse](https://github.com/openbsd/src/tree/dae5ffec5618b0b660e9064e3b0991bb4ab1b1e8/lib/libfuse) implementations from the reference [libfuse](https://github.com/libfuse/libfuse) implementation

## How to build

**Windows cgo**
- Prerequisites: [WinFsp](https://github.com/winfsp/winfsp), gcc (e.g. from [Mingw-builds](http://mingw-w64.org/doku.php/download))
- Build:
    ```
    > cd cgofuse
    > set CPATH=C:\Program Files (x86)\WinFsp\inc\fuse
    > go install -v ./fuse ./examples/memfs
    ```

**Windows !cgo**
- Prerequisites: [WinFsp](https://github.com/winfsp/winfsp)
- Build:
    ```
    > cd cgofuse
    > set CGO_ENABLED=0
    > go install -v ./fuse ./examples/memfs
    ```

**macOS**
- Prerequisites: [FUSE for macOS](https://osxfuse.github.io), [command line tools](https://developer.apple.com/library/content/technotes/tn2339/_index.html)
- Build:
    ```
    $ cd cgofuse
    $ go install -v ./fuse ./examples/memfs ./examples/passthrough
    ```

**Linux**
- Prerequisites: libfuse-dev, gcc
- Build:
    ```
    $ cd cgofuse
    $ go install -v ./fuse ./examples/memfs ./examples/passthrough
    ```

**FreeBSD**
- Prerequisites: fusefs-libs
- Build:
    ```
    $ cd cgofuse
    $ go install -v ./fuse ./examples/memfs ./examples/passthrough

    # You may also need the following in order to run FUSE file systems.
    # Commands must be run as root.
    $ vi /boot/loader.conf                      # add: fuse_load="YES"
    $ sysctl vfs.usermount=1                    # allow user mounts
    $ pw usermod USERNAME -G operator           # allow user to open /dev/fuse
    ```

**NetBSD**
- Prerequisites: NONE
- Build:
    ```
    $ cd cgofuse
    $ go install -v ./fuse ./examples/memfs ./examples/passthrough

    # You may also need the following in order to run FUSE file systems.
    # Commands must be run as root.
    $ chmod go+rw /dev/puffs
    $ sysctl -w vfs.generic.usermount=1
    ```

**OpenBSD**
- Prerequisites: NONE
- Build:
    ```
    $ cd cgofuse
    $ go install -v ./fuse ./examples/memfs ./examples/passthrough
    ```
- **NOTE**: OpenBSD 6 removed the `kern.usermount` option, which allowed non-root users to mount file systems [[link](https://undeadly.org/cgi?action=article&sid=20160715125022&mode=expanded&count=0)]. Therefore you must be root in order to use FUSE and cgofuse.

## How to use

User mode file systems are expected to implement `fuse.FileSystemInterface`. To make implementation simpler a file system can embed ("inherit") a `fuse.FileSystemBase` which provides default implementations for all operations. To mount a file system one must instantiate a `fuse.FileSystemHost` using `fuse.NewFileSystemHost`.

The full documentation is available at GoDoc.org: [package fuse](https://godoc.org/github.com/winfsp/cgofuse/fuse)

There are currently three example file systems:

- [Hellofs](examples/hellofs/hellofs.go) is an extremely simple file system. Runs on all OS'es.
- [Memfs](examples/memfs/memfs.go) is an in memory file system. Runs on all OS'es.
- [Passthrough](examples/passthrough/passthrough.go) is a file system that passes all operations to the underlying file system. Runs on all OS'es except Windows.
- [Notifyfs](examples/notifyfs/notifyfs.go) is a file system that can issue file change notifications. Runs on Windows only.

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

## Contributors

- Bill Zissimopoulos \<billziss at navimatics.com>
- Nick Craig-Wood \<nick at craig-wood.com>
- Fredrik Medley <fredrik.medley at veoneer.com>

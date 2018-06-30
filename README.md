# Cross-platform FUSE library for Go

[![CircleCI](https://img.shields.io/circleci/project/github/billziss-gh/cgofuse.svg?label=cross-build)](https://circleci.com/gh/billziss-gh/cgofuse)
[![GoDoc](https://godoc.org/github.com/billziss-gh/cgofuse/fuse?status.svg)](https://godoc.org/github.com/billziss-gh/cgofuse/fuse)

Cgofuse is a cross-platform FUSE library for Go. It is supported on multiple platforms and can be ported to any platform that has a FUSE implementation. It has [cgo](https://golang.org/cmd/cgo/) and [!cgo](https://github.com/golang/go/wiki/WindowsDLLs) ("nocgo") variants depending on the platform.

|       |macOS<br/>[![Travis CI](https://img.shields.io/travis/billziss-gh/cgofuse.svg)](https://travis-ci.org/billziss-gh/cgofuse)|FreeBSD<br/>[![PMCI](https://storage.googleapis.com/pmci-logs/github.com/billziss-gh/cgofuse/freebsd/badge.svg)](https://storage.googleapis.com/pmci-logs/github.com/billziss-gh/cgofuse/freebsd/build.html)|NetBSD<br/>![no CI](https://img.shields.io/badge/build-none-lightgrey.svg)|OpenBSD<br/>![no CI](https://img.shields.io/badge/build-none-lightgrey.svg)|Linux<br/>[![Travis CI](https://img.shields.io/travis/billziss-gh/cgofuse.svg)](https://travis-ci.org/billziss-gh/cgofuse)|Windows<br/>[![AppVeyor](https://img.shields.io/appveyor/ci/billziss-gh/cgofuse.svg)](https://ci.appveyor.com/project/billziss-gh/cgofuse)|
|:-----:|:----------------:|:----------------:|:----------------:|:----------------:|:----------------:|:----------------:|
|  cgo  |:heavy_check_mark:|:heavy_check_mark:<sup>1</sup>|:heavy_check_mark:<sup>2</sup>|:heavy_check_mark:<sup>2</sup>|:heavy_check_mark:|:heavy_check_mark:|
| !cgo  |                  |                  |                  |                  |                  |:heavy_check_mark:<sup>1</sup>|

- **1**: Requires Go 1.11.
- **2**: NetBSD and OpenBSD support is experimental. There are known issues that stem from the differences in the NetBSD [librefuse](https://github.com/NetBSD/src/tree/bbc46b99bff565d75f55fb23b51eff511068b183/lib/librefuse) and OpenBSD [libfuse](https://github.com/openbsd/src/tree/dae5ffec5618b0b660e9064e3b0991bb4ab1b1e8/lib/libfuse) implementations from the reference [libfuse](https://github.com/libfuse/libfuse) implementation.
    - NetBSD and OpenBSD: Option parsing may fail because the `fuse_opt_parse` function is not fully compatible with the one in libfuse.
    - OpenBSD only: Signal handling is broken due to a bug in the OpenBSD implementation of [`fuse_set_signal_handlers`](https://github.com/openbsd/src/blob/dae5ffec5618b0b660e9064e3b0991bb4ab1b1e8/lib/libfuse/fuse.c#L485-L493).

## How to build

**macOS**
- Prerequisites: [FUSE for macOS](https://osxfuse.github.io), [command line tools](https://developer.apple.com/library/content/technotes/tn2339/_index.html)
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

**Linux**
- Prerequisites: libfuse-dev, gcc
- Build:
    ```
    $ cd cgofuse
    $ go install -v ./fuse ./examples/memfs ./examples/passthrough
    ```

**Windows cgo**
- Prerequisites: [WinFsp](https://github.com/billziss-gh/winfsp), gcc (e.g. from [Mingw-builds](http://mingw-w64.org/doku.php/download))
- Build:
    ```
    > cd cgofuse
    > set CPATH=C:\Program Files (x86)\WinFsp\inc\fuse
    > go install -v ./fuse ./examples/memfs
    ```

**Windows !cgo**
- Prerequisites: [WinFsp](https://github.com/billziss-gh/winfsp)
- Build:
    ```
    > cd cgofuse
    > set CGO_ENABLED=0
    > go install -v ./fuse ./examples/memfs
    ```

## How to cross-compile your project using xgo

You can easily cross-compile your project using [xgo](https://github.com/karalabe/xgo) and the [billziss/xgo-cgofuse](https://hub.docker.com/r/billziss/xgo-cgofuse/) docker image.

- Prerequisites: [docker](https://www.docker.com), [xgo](https://github.com/karalabe/xgo)
- Build:
    ```
    $ docker pull billziss/xgo-cgofuse
    $ go get -u github.com/karalabe/xgo
    $ cd YOUR-PROJECT-THAT-USES-CGOFUSE
    $ xgo --image=billziss/xgo-cgofuse \
        --targets=darwin/386,darwin/amd64,linux/386,linux/amd64,windows/386,windows/amd64 .
    ```

Cross-compilation only works for macOS, Linux and Windows.

## How to use

User mode file systems are expected to implement `fuse.FileSystemInterface`. To make implementation simpler a file system can embed ("inherit") a `fuse.FileSystemBase` which provides default implementations for all operations. To mount a file system one must instantiate a `fuse.FileSystemHost` using `fuse.NewFileSystemHost`.

The full documentation is available at GoDoc.org: [package fuse](https://godoc.org/github.com/billziss-gh/cgofuse/fuse)

There are currently three example file systems:

- [Hellofs](examples/hellofs/hellofs.go) is an extremely simple file system. Runs on all OS'es.
- [Memfs](examples/memfs/memfs.go) is an in memory file system. Runs on all OS'es.
- [Passthrough](examples/passthrough/passthrough.go) is a file system that passes all operations to the underlying file system. Runs on all OS'es except Windows.

## How it is tested

Cgofuse is regularly built and tested on [Travis CI](https://travis-ci.org/billziss-gh/cgofuse), [Poor Man's CI](https://github.com/billziss-gh/pmci) and [AppVeyor](https://ci.appveyor.com/project/billziss-gh/cgofuse). The following software is being used to test cgofuse.

**macOS**
- [fstest](https://github.com/billziss-gh/secfs.test/tree/master/fstest/ntfs-3g-pjd-fstest-8af5670)
- [fsx](https://github.com/billziss-gh/secfs.test/tree/master/fstools/src/fsx)

**FreeBSD**
- [fsx](https://github.com/billziss-gh/secfs.test/tree/master/fstools/src/fsx)

**Linux**
- [fstest](https://github.com/billziss-gh/secfs.test/tree/master/fstest/ntfs-3g-pjd-fstest-8af5670)
- [fsx](https://github.com/billziss-gh/secfs.test/tree/master/fstools/src/fsx)

**Windows (cgo and !cgo)**
- [winfsp-tests](https://github.com/billziss-gh/winfsp/tree/master/tst/winfsp-tests)
- [fsx](https://github.com/billziss-gh/secfs.test/tree/master/fstools/src/fsx)

## Contributors

- Bill Zissimopoulos \<billziss at navimatics.com>
- Nick Craig-Wood \<nick at craig-wood.com>

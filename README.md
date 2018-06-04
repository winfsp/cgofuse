# Cross-platform FUSE library for Go

[![Travis CI](https://img.shields.io/travis/billziss-gh/cgofuse.svg?label=macOS/linux)](https://travis-ci.org/billziss-gh/cgofuse)
[![AppVeyor](https://img.shields.io/appveyor/ci/billziss-gh/cgofuse.svg?label=windows)](https://ci.appveyor.com/project/billziss-gh/cgofuse)
[![CircleCI](https://img.shields.io/circleci/project/github/billziss-gh/cgofuse.svg?label=cross-build)](https://circleci.com/gh/billziss-gh/cgofuse)
[![GoDoc](https://godoc.org/github.com/billziss-gh/cgofuse/fuse?status.svg)](https://godoc.org/github.com/billziss-gh/cgofuse/fuse)

Cgofuse is a cross-platform FUSE library for Go. It is supported on multiple platforms and can be ported to any platform that has a FUSE implementation. It has [cgo](https://golang.org/cmd/cgo/) and [!cgo](https://github.com/golang/go/wiki/WindowsDLLs) ("nocgo") variants depending on the platform.


|       |macOS             |FreeBSD<sup>2</sup>|Linux            |Windows<sup>1</sup>|
|:-----:|:----------------:|:----------------:|:----------------:|:----------------:|
|  cgo  |:heavy_check_mark:|:heavy_check_mark:|:heavy_check_mark:|:heavy_check_mark:|
| !cgo  |                  |                  |                  |:heavy_check_mark:|

- **1**: Windows requires [WinFsp](https://github.com/billziss-gh/winfsp).
- **2**: FreeBSD signal handling is currently broken. Please see this [discussion](https://github.com/billziss-gh/cgofuse/issues/18#issuecomment-390446362) for details.

## How to build

**macOS**
- Prerequisites: [FUSE for macOS](https://osxfuse.github.io), [command line tools](https://developer.apple.com/library/content/technotes/tn2339/_index.html)
- Build:
    ```
    $ cd cgofuse
    $ go install -v ./fuse ./examples/memfs ./examples/passthrough
    ```

**FreeBSD**
- Prerequisites: fusefs-libs, clang
- Build:
    ```
    $ cd cgofuse
    $ go install -v ./fuse ./examples/memfs ./examples/passthrough

    # you may also need the following in order to run FUSE file systems:
    $ sudo vi /boot/loader.conf                 # add: fuse_load="YES"
    $ sudo sysctl vfs.usermount=1               # allow user mounts
    $ sudo pw usermod USERNAME -G operator      # allow user to open /dev/fuse
    ```

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

- [Hellofs](examples/hellofs/hellofs.go) is an extremely simple file system. Runs on macOS, FreeBSD, Linux and Windows.
- [Memfs](examples/memfs/memfs.go) is an in memory file system. Runs on macOS, FreeBSD, Linux and Windows.
- [Passthrough](examples/passthrough/passthrough.go) is a file system that passes all operations to the underlying file system. Runs on macOS, FreeBSD, Linux.

## How it is tested

Cgofuse is regularly built and tested on [Travis CI](https://travis-ci.org/billziss-gh/cgofuse) and [AppVeyor](https://ci.appveyor.com/project/billziss-gh/cgofuse). The following software is being used to test cgofuse.

**macOS/Linux**
- [fstest](https://github.com/billziss-gh/secfs.test/tree/master/fstest/ntfs-3g-pjd-fstest-8af5670)
- [fsx](https://github.com/billziss-gh/secfs.test/tree/master/fstools/src/fsx)

**Windows**
- [winfsp-tests](https://github.com/billziss-gh/winfsp/tree/master/tst/winfsp-tests)

## Contributors

- Bill Zissimopoulos \<billziss at navimatics.com>
- Nick Craig-Wood \<nick at craig-wood.com>

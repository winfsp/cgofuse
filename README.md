# FUSE library for Go

Cgofuse is a cross-platform FUSE library for Go. It is implemented using [cgo](https://golang.org/cmd/cgo/) and can be ported to any platform that has a FUSE implementation.

Cgofuse currently runs on OSX, Linux and Windows (using [WinFsp](https://github.com/billziss-gh/winfsp)).

## How to build

Cgofuse is regularly built on [Travis CI](https://travis-ci.org/billziss-gh/cgofuse) and [AppVeyor](https://ci.appveyor.com/project/billziss-gh/cgofuse).

To build locally you will need a Go installation and a C compiler.

OSX:

- Prerequisites: [OSXFUSE](https://osxfuse.github.io), [command line tools](https://developer.apple.com/library/content/technotes/tn2339/_index.html)
- Build:
    ```
    $ cd cgofuse
    $ go install -v ./...
    ```

Linux:

- Prerequisites: libfuse-dev, gcc
- Build:
    ```
    $ cd cgofuse
    $ go install -v ./...
    ```

Windows:

- [WinFsp](https://github.com/billziss-gh/winfsp), gcc (e.g. from [Mingw-builds](http://mingw-w64.org/doku.php/download))
- Build:
    ```
    > cd cgofuse
    > set CPATH=C:\Program Files (x86)\WinFsp\inc\fuse
    > go install -v ./examples/memfs
    ```

## How to use

User mode file systems are expected to implement `fuse.FileSystemInterface`. To make implementation simpler a file system can embed ("inherit") a `fuse.FileSystemBase` which provides default implementations for all operations. To mount a file system one must instantiate a `fuse.FileSystemHost` using `fuse.NewFileSystemHost`.

There are currently two example file systems:

- [Memfs](examples/memfs/memfs.go) is a simple in memory file system. Runs on OSX, Linux and Windows.
- [Passthrough](examples/passthrough/passthrough.go) is a file system that passes all operations to the underlying file system. Runs on OSX, Linux.

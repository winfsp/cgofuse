# Changelog


**v1.5.0**

- Add `FileSystemHost.Notify` API which allows file change notification to be issued from the user mode file system [Windows only].
- Add `notifyfs` file system to showcase the new API functionality [Windows only].


**v1.4.0**

- The FUSE library is demand-loaded on all platforms.
    - Prior to this version only the Windows FUSE library (WinFsp) was demand loaded.
    - This introduces a behavior change for programs that use FUSE when FUSE is not present. Prior to this version such programs would not start at all. In this new version such programs will start, but will fail with a `panic` in `FileSystemHost.Mount` if FUSE is not present.


**v1.3.0**

- Add FileSystemOpenEx interface.
- Miscellaneous other fixes and improvements.


**v1.2.0**

- Cgo and !cgo variants for Windows port.
- FreeBSD, NetBSD and OpenBSD ports.


**v1.1.0**

- `OptParse` function parses FUSE options.
- `fmt.Stringer` and `fmt.GoStringer` implementation for `fuse.Error`.


**v1.0.4**

- Implement BSD `flags`, `chflags`, `setcrtime`, `setchgtime`.
- Improve documentation.


**v1.0.3**

- Windows XP compatibility (eliminate `RegGetValueW`).


**v1.0.2**

- Windows XP compatibility (eliminate slim R/W lock).


**v1.0.1**

- Cross-compilation `Dockerfile`.
- CircleCI integration.
- Do not catch `SIGHUP`.
- Improve documentation.


**v1.0**

- Initial cgofuse release.
- The API is now **FROZEN**. Breaking API changes will receive a major version update (`2.0`). Incremental API changes will receive a minor version update (`1.x`).

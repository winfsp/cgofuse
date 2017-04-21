# FUSE library for Go

Cgofuse is a FUSE library for Go using cgo. The benefit is that this library can be ported to all platforms that have a FUSE implementation. This includes Windows with my own [WinFsp](https://github.com/billziss-gh/winfsp).

Please note that this library is written by an extreme Go novice. So keep the laughs to a minimum!

## How to test on OSX

Clone and build cgofuse:

```
$ cd $GOPATH/src
$ git clone https://github.com/billziss-gh/cgofuse.git
$ go install cgofuse/examples/passthrough
```

Clone and build fstest:

```
$ cd ~/Projects
$ git clone https://github.com/billziss-gh/secfs.test.git
$ cd secfs.test/fstest/fstest
$ make
```

Apply the following patch to secfs.test:

```diff
diff --git a/fstest/ntfs-3g-pjd-fstest-8af5670/tests/conf b/fstest/ntfs-3g-pjd-fstest-8af5670/tests/conf
index 18cd344..6c64086 100644
--- a/fstest/ntfs-3g-pjd-fstest-8af5670/tests/conf
+++ b/fstest/ntfs-3g-pjd-fstest-8af5670/tests/conf
@@ -5,4 +5,4 @@
 os=`uname`
 
 # Known file systems: UFS, ZFS, ext3, ext4, ntfs-3g, xfs, btrfs, glusterfs, HFS+, secfs
-fs="secfs"
+fs="HFS+"
diff --git a/fstest/ntfs-3g-pjd-fstest-8af5670/tests/misc.sh b/fstest/ntfs-3g-pjd-fstest-8af5670/tests/misc.sh
index 6714f8f..5d5c18f 100644
--- a/fstest/ntfs-3g-pjd-fstest-8af5670/tests/misc.sh
+++ b/fstest/ntfs-3g-pjd-fstest-8af5670/tests/misc.sh
@@ -137,7 +137,7 @@ supported()
                        return 0
                fi
                if [ ${os} = "Darwin" ]; then
-                       return 0
+                       return 1
                fi
         return 1
                ;;
```

The passthrough file system uses a system directory as its underlying storage (`/tmp/t/p`). We will also need a mount directory (`/tmp/t/m`) and we can then launch the file system:

```
$ mkdir -p /tmp/t/{p,m}
$ chmod 777 /tmp/t/p
sudo $GOPATH/bin/passthrough -f -o attr_timeout=0,use_ino,allow_other /tmp/t/p /tmp/t/m
```

From a different command prompt run fstest against the file system:

```
$ cd /tmp/t/m
$ sudo prove -fr ~/Projects/secfs.test/fstest/fstest/tests
```

To unmount the file system use:

```
$ cd
$ sudo umount /tmp/t/m
```

### Failing tests

Most fstest tests will pass, but there are a few failing tests:

```
/Users/billziss/Projects/secfs.test/fstest/fstest/tests/chown/00.t ............. 48/171 
not ok 48 - expect 06555 lstat fstest_61de0b9fd8a698e01a34a150f0534f9b mode - got 0555
not ok 55 - expect 06555 lstat fstest_61de0b9fd8a698e01a34a150f0534f9b mode - got 0555
not ok 62 - expect 06555 lstat fstest_61de0b9fd8a698e01a34a150f0534f9b mode - got 0555

/Users/billziss/Projects/secfs.test/fstest/fstest/tests/open/17.t .............. 1/3 
not ok 2 - expect ENXIO open fstest_df12f85343217433ce784a53447310aa O_WRONLY,O_NONBLOCK - got EPERM

/Users/billziss/Projects/secfs.test/fstest/fstest/tests/rmdir/12.t ............. 1/6 
not ok 4 - expect ENOTEMPTY rmdir fstest_486cf7e2f6146f54f985f5816da63a24/fstest_85b2baba0a95d210d3092e33b04e02d4/.. - got EINVAL

/Users/billziss/Projects/secfs.test/fstest/fstest/tests/zzz_ResourceFork/00.t .. 1/6 xattr: fstest_65b295c43f8cae84ea8a5b7b51e0f9be: No such xattr: com.apple.ResourceFork
```

Most of these failures are OSX API or OSXFUSE quirks. The last test (`zzz_ResourceFork`) does not pass because we do not have xattr support yet.

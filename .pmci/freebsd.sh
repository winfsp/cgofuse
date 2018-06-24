set -ex

# FUSE
kldload fuse
pkg install -y fusefs-libs

# cgofuse: build and test
export GOPATH=/tmp/go
mkdir -p /tmp/go/src/github.com/billziss-gh
cp -R /tmp/repo/cgofuse /tmp/go/src/github.com/billziss-gh
cd /tmp/go/src/github.com/billziss-gh/cgofuse
go build ./examples/memfs
go build ./examples/passthrough
go test -v ./fuse

# secfs.test
pkg install -y gmake
git clone -q https://github.com/billziss-gh/secfs.test.git /tmp/repo/secfs.test
git -C /tmp/repo/secfs.test checkout -q 105e0fe6280631d5077f950589de1b2c44b0faad
gmake -C /tmp/repo/secfs.test/fstools/src/fsx
mkdir -p /tmp/t/m /tmp/t/p

# cgofuse/memfs
./memfs /tmp/t/m &
(cd /tmp/t/m && /tmp/repo/secfs.test/fstools/src/fsx/fsx -N 10000 test xxxxxx)
umount /tmp/t/m

# cgofuse/passthrough
./passthrough /tmp/t/p /tmp/t/m &
(cd /tmp/t/m && /tmp/repo/secfs.test/fstools/src/fsx/fsx -N 10000 test xxxxxx)
umount /tmp/t/m

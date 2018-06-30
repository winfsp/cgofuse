set -ex

# cgofuse: build and test
export GOPATH=/tmp/go
mkdir -p /tmp/go/src/github.com/billziss-gh
cp -R /tmp/repo/cgofuse /tmp/go/src/github.com/billziss-gh
cd /tmp/go/src/github.com/billziss-gh/cgofuse
go build -v ./...
# go test -v ./fuse -run 'TestUnmount|TestSignal'

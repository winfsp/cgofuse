set -ex

# cgofuse: build and test
export GOPATH=/tmp/go
mkdir -p /tmp/go/src/github.com/winfsp
cp -R /tmp/repo/cgofuse /tmp/go/src/github.com/winfsp
cd /tmp/go/src/github.com/winfsp/cgofuse
go build -v ./...
# go test -v ./fuse -run 'TestUnmount|TestSignal'

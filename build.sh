export PATH=$PATH:/root/goroot/bin
export GOPATH=/root/goroot/bin
cd /root/download-delegator
git pull
go get
go build

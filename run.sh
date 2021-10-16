pkill download-delegator
export PATH=$PATH:/root/goroot/bin
cd /root/download-delegator
git pull
go build
nohup ./download-delegator config.toml &> /var/log/download-delegator.log &

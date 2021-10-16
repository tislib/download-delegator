pkill download-delegator
go build
nohup ./download-delegator config.toml &> /var/log/download-delegator.log
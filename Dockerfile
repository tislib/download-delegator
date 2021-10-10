FROM ubuntu

COPY download-delegator /

COPY server.crt /
COPY server.key /

COPY container.config.toml /config.toml

# Command to run
ENTRYPOINT ["/download-delegator", "/config.toml"]

FROM ubuntu

COPY download-delegator /

COPY server.crt /
COPY server.key /

# Command to run
ENTRYPOINT ["/download-delegator", "/server.crt", "/server.key", "/proxy.conf"]

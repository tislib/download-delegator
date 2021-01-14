FROM ubuntu

COPY download-delegator /

# Command to run
ENTRYPOINT ["/download-delegator"]

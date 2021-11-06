FROM golang:1.16
WORKDIR /app

COPY . /app

RUN go get
RUN go get download-delegator/awslambda
RUN go get download-delegator/app

RUN curl -sSL "https://github.com/gotestyourself/gotestsum/releases/download/v0.3.1/gotestsum_0.3.1_linux_amd64.tar.gz" | sudo tar -xz -C /usr/local/bin gotestsum

RUN gotestsum
FROM golang:1.16

WORKDIR /app

COPY . /app

RUN go get
RUN go get download-delegator/awslambda
RUN go get download-delegator/app
RUN go get github.com/stretchr/testify/assert

RUN go get gotest.tools/gotestsum

RUN gotestsum
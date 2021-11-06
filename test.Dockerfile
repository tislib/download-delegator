FROM golang:1.16

RUN go get
RUN go get download-delegator/awslambda
RUN go get download-delegator/app

RUN go get gotest.tools/gotestsum

WORKDIR /app

COPY . /app

RUN gotestsum
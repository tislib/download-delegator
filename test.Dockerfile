FROM golang:1.16
WORKDIR /app

COPY . /app

RUN go get
RUN go get gotest.tools/gotestsum

RUN gotestsum
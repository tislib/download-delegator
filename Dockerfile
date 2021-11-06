FROM golang:1.16

RUN go get
RUN go get download-delegator/awslambda
RUN go get download-delegator/app

WORKDIR /app

COPY . /app

RUN CGO_ENABLED=0 GOOS=linux go build -o download-delegator .


FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /app/download-delegator ./
CMD ["./download-delegator"]
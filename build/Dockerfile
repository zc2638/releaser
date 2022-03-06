FROM golang:1.17 as builder
ENV GOPROXY=https://goproxy.io,direct
ENV GO111MODULE=on

WORKDIR /go/cache
ADD go.mod .
ADD go.sum .
RUN go mod download

WORKDIR /work
ADD . .
RUN go run github.com/zc2638/releaser/cmd set --git &&\
 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w -X $(go run github.com/zc2638/releaser/cmd get)" -o releaser github.com/zc2638/releaser/cmd

FROM alpine:3.6 as alpine
RUN apk add -U --no-cache ca-certificates

FROM alpine:3.6
MAINTAINER zc2638
LABEL maintainer="zc2638" \
    email="zc2638@qq.com" \
    version="0.0.1"

COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /work/releaser /bin/releaser

WORKDIR /work
CMD ["releaser", "-v"]

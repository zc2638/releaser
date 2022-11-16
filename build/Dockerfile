FROM golang:1.19 as builder
ENV GOPROXY=https://goproxy.io,direct
ENV GO111MODULE=on

WORKDIR /go/cache
ADD go.mod .
ADD go.sum .
RUN go mod download

WORKDIR /work
ADD . .
RUN make build

FROM alpine:3.6 as alpine
RUN apk add -U --no-cache ca-certificates

FROM alpine:3.6
MAINTAINER zc2638
LABEL maintainer="zc2638" \
    email="zc2638@qq.com" \
    version="0.0.3"

COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/local/bin/releaser /usr/local/bin/releaser

WORKDIR /work
CMD ["releaser", "-v"]

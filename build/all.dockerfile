FROM golang:1.17 as builder
ENV GOPROXY=https://goproxy.io,direct
ENV GO111MODULE=on

WORKDIR /go/cache
ADD go.mod .
ADD go.sum .
RUN go mod download

WORKDIR /work
ADD . .
RUN go run github.com/zc2638/releaser/cmd/releaser set --git &&\
 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w -X $(go run github.com/zc2638/releaser/cmd/releaser get)" -o releaser github.com/zc2638/releaser/cmd/releaser

FROM node:18.5.0-alpine3.16

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && \
    apk --update upgrade && \
    apk add git ca-certificates tzdata && \
    rm -rf /var/cache/apk/* && \
    npm install -g -registry=https://registry.npm.taobao.org/ \
  conventional-changelog \
  conventional-changelog-cli \
  cz-conventional-changelog \
  commitizen \
  standard-version

COPY --from=mikefarah/yq:4.25.3 /usr/bin/yq /usr/bin/yq
COPY --from=builder /work/releaser /usr/bin/releaser

ENV TZ=Asia/Shanghai
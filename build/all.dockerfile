FROM golang:1.19 as builder
ENV GOPROXY=https://goproxy.io,direct
ENV GO111MODULE=on

WORKDIR /go/cache
ADD go.mod .
ADD go.sum .
RUN go mod download

WORKDIR /work
ADD . .
RUN go run github.com/zc2638/releaser/cmd/releaser set --git && \
    CGO_ENABLED=0 go build -ldflags="-s -w -X $(go run github.com/zc2638/releaser/cmd/releaser get)" -o /usr/local/bin/releaser github.com/zc2638/releaser/cmd/releaser


FROM node:18.5.0-alpine3.16

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories \
    && apk --update upgrade \
    && apk add git ca-certificates tzdata tar \
    # Install envsubst command for replacing config files in system startup
    # - it needs libintl package
    # - only weights 100KB combined with it's libraries
    && apk add gettext libintl \
    && rm -rf /var/cache/apk/* \
    && npm install -g -registry=https://registry.npm.taobao.org/ \
  conventional-changelog \
  conventional-changelog-cli \
  cz-conventional-changelog \
  commitizen \
  standard-version

COPY --from=mikefarah/yq:4.25.3 /usr/bin/yq /usr/local/bin/yq
COPY --from=builder /usr/local/bin/releaser /usr/local/bin/releaser

ENV TZ=Asia/Shanghai
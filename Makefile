build:
	@go run github.com/zc2638/releaser/cmd/releaser set --git
    @CGO_ENABLED=0 go build -ldflags="-s -w -X $(go run github.com/zc2638/releaser/cmd/releaser get)" -o /usr/local/bin/releaser github.com/zc2638/releaser/cmd/releaser

docker:
	@docker build -t  zc2638/releaser -f build/Dockerfile .

docker-helper:
	@docker build -t  zc2638/releaser:helper -f build/all.dockerfile .

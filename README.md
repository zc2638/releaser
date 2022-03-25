# releaser

A tool for Golang programs that automatically generate versions at build time.

## Quick Start

### 1. Code import

```go
package main

import (
	"fmt"

	"github.com/zc2638/releaser"
)

func main() {
	fmt.Printf("version: %s\n", releaser.Version.String())
	fmt.Printf("gitCommit: %s\n", releaser.Version.Git.Commit)
	fmt.Printf("goVersion: %s\n", releaser.Version.GoVersion)
	fmt.Printf("compiler: %s\n", releaser.Version.Compiler)
	fmt.Printf("platform: %s\n", releaser.Version.Platform)
	fmt.Printf("buildDate: %s\n", releaser.Version.BuildDate)
}
```

### 2. Run Build

```shell
go build -ldflags="-X $(releaser get)" main.go
```

## Install

### Install from source

```shell
go install -v github.com/zc2638/releaser/cmd/releaser@latest
```

### Install from Docker

```shell
docker run --rm -it zc2638/releaser:0.0.1 
```

### Build from source

#### 1. Clone

```shell
git clone https://github.com/zc2638/releaser.git releaser && cd "$_"
```

#### 2. Build

```shell
go build -ldflags="-X $(go run github.com/zc2638/releaser/cmd get)" -o releaser github.com/zc2638/releaser/cmd
```

## Commands

### init project

```shell
releaser init <project name>
```

### create service

```shell
releaser create <service name>
```

### set

#### use for project

```shell
releaser set --version 1.0.0 --meta status=ok --meta repo=github
```

#### use for service

```shell
releaser set <service name> --version 1.0.0 --meta status=ok --meta repo=github
```

### get

#### use for project

```shell
# get build info, output format `gobuild` is default
releaser get [-o gobuild|json]
# get version
releaser get --filter version
# get metadata field value
releaser get --filter meta.status
```

#### use for service

```shell
# get build info, output format `gobuild` is default
releaser get <service name> [-o gobuild|json]
# get version
releaser get <service name> --filter version
# get metadata field value
releaser get <service name> --filter meta.status
```

### delete

```shell
releaser delete <service name>
```

### walk

```shell
releaser walk --command "echo $version && echo $meta.status"
```

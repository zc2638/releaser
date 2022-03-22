// Copyright © 2022 zc2638 <zc2638@qq.com>.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package git

import (
	"bytes"
	"fmt"

	"github.com/zc2638/releaser/pkg/util"
)

func run(command string) (string, error) {
	output, err := util.Exec(command)
	if err != nil {
		return "", err
	}
	output = bytes.TrimSuffix(output, []byte("\n"))
	return string(output), nil
}

func CurrentCommit() (string, error) {
	return run("git rev-parse HEAD")
}

func CurrentBranch() (string, error) {
	return run("git rev-parse --abbrev-ref HEAD")
}

func GetTagByCommit(commit string) (string, error) {
	return run(fmt.Sprintf("git tag --points-at %s", commit))
}

func GetRecentTag() (string, error) {
	return run("git describe --abbrev=0 --tags")
}

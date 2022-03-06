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

package releaser

import (
	"fmt"
	"runtime"
	"time"

	"github.com/blang/semver/v4"
)

var data string

var Version *VersionEntry

func init() {
	Version = &VersionEntry{}
	if data != "" {
		_ = Decode(data, Version)
	}
	Version.Complete()
}

type VersionGit struct {
	Branch string `json:"branch,omitempty" yaml:"branch,omitempty"`
	Commit string `json:"commit,omitempty" yaml:"commit,omitempty"`
	Tag    string `json:"tag,omitempty" yaml:"tag,omitempty"`
}

type VersionEntry struct {
	Version   *semver.Version   `json:"version,omitempty" yaml:"version,omitempty"`
	Git       *VersionGit       `json:"git,omitempty" yaml:"git,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	BuildDate string            `json:"buildDate,omitempty" yaml:"buildDate,omitempty"`
	GoVersion string            `json:"goVersion,omitempty" yaml:"goVersion,omitempty"`
	Compiler  string            `json:"compiler,omitempty" yaml:"compiler,omitempty"`
	Platform  string            `json:"platform,omitempty" yaml:"platform,omitempty"`
}

func (ve VersionEntry) String() string {
	if ve.Version != nil {
		version := ve.Version.String()
		if version != "" {
			return version
		}
	}
	version := "0.0.0"
	branch := "unknown"
	if ve.Git == nil {
		return version + "-" + branch
	}
	if ve.Git.Tag != "" {
		return ve.Git.Tag
	}
	if ve.Git.Branch != "" {
		branch = ve.Git.Branch
	}
	return version + "-" + branch + "." + ve.Git.Commit
}

func (ve *VersionEntry) Complete() {
	ve.BuildDate = time.Now().Format("2006-01-02 15:04:05")
	ve.GoVersion = runtime.Version()
	ve.Compiler = runtime.Compiler
	ve.Platform = fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
}

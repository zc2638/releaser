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
)

var version string

var Version *VersionEntry

func init() {
	ve, err := Decode(version)
	if err != nil {
		ve = &VersionEntry{}
	}
	ve.Complete()
	Version = ve
}

type VersionEntry struct {
	GitTag    string `json:"gitTag" yaml:"gitTag"`
	GitCommit string `json:"gitCommit" yaml:"gitCommit"`
	GitBranch string `json:"gitBranch" yaml:"gitBranch"`
	BuildDate string `json:"buildDate,omitempty" yaml:"buildDate,omitempty"`
	GoVersion string `json:"goVersion,omitempty" yaml:"goVersion,omitempty"`
	Compiler  string `json:"compiler,omitempty" yaml:"compiler,omitempty"`
	Platform  string `json:"platform,omitempty" yaml:"platform,omitempty"`
}

func (ve VersionEntry) String() string {
	if ve.GitTag != "" {
		return ve.GitTag
	}
	return ve.GitCommit
}

func (ve *VersionEntry) Complete() {
	ve.BuildDate = time.Now().Format("2006-01-02 15:04:05")
	ve.GoVersion = runtime.Version()
	ve.Compiler = runtime.Compiler
	ve.Platform = fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
}

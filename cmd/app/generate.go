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

package app

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"

	"github.com/zc2638/releaser/pkg/git"

	"github.com/zc2638/releaser"

	"github.com/spf13/cobra"
	"github.com/zc2638/releaser/pkg/definition"
)

func NewGenerateCommand() *cobra.Command {
	cfg := &definition.CommonConfig{}
	cmd := &cobra.Command{
		Use:          "generate",
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, _ []string) error {
			_, err := os.Stat(cfg.Config)
			if err != nil {
				if err := os.MkdirAll(cfg.Config, os.ModePerm); err != nil {
					return err
				}
			}

			manifestPath := filepath.Join(cfg.Config, definition.Manifest)
			branch, err := git.CurrentBranch()
			if err != nil {
				return err
			}
			commit, err := git.CurrentCommit()
			if err != nil {
				return err
			}
			tag, _ := git.GetTagByCommit(commit)

			ve := &releaser.VersionEntry{
				GitTag:    tag,
				GitCommit: commit,
				GitBranch: branch,
			}

			data, err := releaser.BuildYaml(ve)
			if err != nil {
				return err
			}
			if err := os.WriteFile(manifestPath, data, os.ModePerm); err != nil {
				return err
			}

			switch cfg.Output {
			case definition.OutputFormatJSON:
				jsonData, err := json.Marshal(ve)
				if err != nil {
					return err
				}
				fmt.Println(string(jsonData))
			case definition.OutputFormatGoBuild:
				encodeData, err := releaser.Encode(ve)
				if err != nil {
					return err
				}
				module := "github.com/zc2638/releaser"
				if info, ok := debug.ReadBuildInfo(); ok {
					module = info.Main.Path
				}
				fmt.Println(fmt.Sprintf("%s.version=%s", module, encodeData))
			default:
				return fmt.Errorf("unsupport output format: %s", cfg.Output)
			}
			return nil
		},
	}

	cfg.AddConfigFlag(cmd)
	cfg.AddOutputFormatFlag(cmd)
	return cmd
}

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
	"fmt"
	"runtime/debug"

	"github.com/zc2638/releaser/pkg/util"

	"github.com/spf13/cobra"

	"github.com/zc2638/releaser/pkg/storage"
)

type WalkOption struct {
	Command string
}

func newWalkCommand() *cobra.Command {
	opt := &WalkOption{}
	cmd := &cobra.Command{
		Use:          "walk",
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, args []string) error {
			info, ok := debug.ReadBuildInfo()
			if ok {
				moduleName = info.Main.Path
			}

			data, err := storage.Read(storage.ManifestName)
			if err != nil {
				return fmt.Errorf("read manifest failed: %v", err)
			}
			manifest := &storage.Manifest{}
			if err := manifest.Unmarshal(data); err != nil {
				return fmt.Errorf("unmarshal manifest failed: %v", err)
			}

			entries := make([]*storage.Entry, 0, len(manifest.Services))
			for _, s := range manifest.Services {
				entry := &storage.Entry{}
				sd, err := storage.Read(s + ".yaml")
				if err != nil {
					entry.Name = s
				} else if err := entry.Unmarshal(sd); err != nil {
					return err
				}
				entries = append(entries, entry)
			}

			for _, entry := range entries {
				m := entry.ToMap()
				output, err := util.ExecWithEnv(opt.Command, m)
				if err != nil {
					return err
				}
				fmt.Print(string(output))
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&opt.Command, "command", "c", opt.Command, "Set the command to be executed")
	return cmd
}

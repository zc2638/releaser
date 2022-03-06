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
	"runtime/debug"
	"strings"

	"github.com/99nil/go/sets"
	"github.com/blang/semver/v4"
	"github.com/spf13/cobra"

	"github.com/zc2638/releaser"
	"github.com/zc2638/releaser/pkg/storage"
	"github.com/zc2638/releaser/pkg/util"
)

var moduleName = "github.com/zc2638/releaser"

type GetOption struct {
	Filter string
	Output string
}

func newGetCommand() *cobra.Command {
	opt := &GetOption{
		Output: storage.OutputFormatGoBuild,
	}
	cmd := &cobra.Command{
		Use:          "get",
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, args []string) error {
			if err := util.CheckArgsGT(args, 1); err != nil {
				return err
			}

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

			entry := &storage.Entry{}
			if len(args) == 0 {
				entry = &manifest.Entry
			} else {
				serviceName := args[0]
				servicePath := serviceName + ".yaml"
				ss := sets.NewString(manifest.Services...)
				if !ss.Has(serviceName) {
					return fmt.Errorf("service %s not exist", serviceName)
				}

				// Regenerate information if read failed
				data, err = storage.Read(servicePath)
				if err != nil {
					return fmt.Errorf("read service failed: %v", err)
				}
				if err := entry.Unmarshal(data); err != nil {
					return fmt.Errorf("unmarshal service failed: %v", err)
				}
			}

			if opt.Filter == "version" {
				fmt.Println(entry.Version)
			} else if strings.HasPrefix(opt.Filter, "meta.") {
				key := strings.TrimPrefix(opt.Filter, "meta.")
				fmt.Println(entry.Metadata[key])
			} else {
				switch opt.Output {
				case storage.OutputFormatJSON:
					marshal, err := json.Marshal(entry)
					if err != nil {
						return fmt.Errorf("marshal service failed: %v", err)
					}
					fmt.Println(string(marshal))
				default:
					version, err := semver.New(strings.TrimPrefix(entry.Version, "v"))
					if err != nil {
						return fmt.Errorf("check service version failed: %v", err)
					}

					ve := &releaser.VersionEntry{}
					ve.Version = version
					if entry.Git != nil {
						ve.Git = &releaser.VersionGit{
							Branch: entry.Git.Branch,
							Commit: entry.Git.Commit,
							Tag:    entry.Git.Tag,
						}
					}
					if entry.Metadata != nil {
						ve.Metadata = entry.Metadata
					}
					encode, err := releaser.Encode(ve)
					if err != nil {
						return fmt.Errorf("encode service failed: %v", err)
					}
					fmt.Println(moduleName + ".data=" + encode)
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&opt.Output, "output", "o", opt.Output, "Specify the format of the output. e.g. gobuild,json")
	cmd.Flags().StringVar(&opt.Filter, "filter", "", "Get content according to filter rules. e.g. version,meta.status")
	return cmd
}

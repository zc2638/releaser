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
	"strings"

	"github.com/99nil/go/sets"
	"github.com/spf13/cobra"

	"github.com/zc2638/releaser/pkg/git"
	"github.com/zc2638/releaser/pkg/storage"
	"github.com/zc2638/releaser/pkg/util"
)

type SetOption struct {
	Git      bool
	Version  string
	Metadata []string
}

func newSetCommand() *cobra.Command {
	opt := &SetOption{}
	cmd := &cobra.Command{
		Use:          "set",
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, args []string) error {
			if err := util.CheckArgsGT(args, 1); err != nil {
				return err
			}

			data, err := storage.Read(storage.ManifestName)
			if err != nil {
				return fmt.Errorf("read manifest failed: %v", err)
			}
			manifest := &storage.Manifest{}
			if err := manifest.Unmarshal(data); err != nil {
				return fmt.Errorf("unmarshal manifest failed: %v", err)
			}

			gitEntry := &storage.GitEntry{}
			if opt.Git {
				branch, err := git.CurrentBranch()
				if err != nil {
					return err
				}
				commit, err := git.CurrentCommit()
				if err != nil {
					return err
				}
				tag, _ := git.GetTagByCommit(commit)
				gitEntry.Branch = branch
				gitEntry.Commit = commit
				gitEntry.Tag = tag
			}

			meta := make(map[string]string)
			for _, v := range opt.Metadata {
				arr := strings.SplitN(v, "=", 2)
				if len(arr) != 2 {
					continue
				}
				meta[arr[0]] = arr[1]
			}

			if len(args) == 0 {
				if opt.Git {
					manifest.Git = gitEntry
				}
				if opt.Version != "" {
					manifest.Version = opt.Version
				}
				if manifest.Metadata == nil {
					manifest.Metadata = make(map[string]string)
				}
				for k, v := range meta {
					manifest.Metadata[k] = v
				}
				return storage.Save(manifest, storage.ManifestName)
			}

			serviceName := args[0]
			servicePath := serviceName + ".yaml"
			ss := sets.NewString(manifest.Services...)
			if !ss.Has(serviceName) {
				return fmt.Errorf("service %s not exist", serviceName)
			}

			entry := &storage.Entry{}
			// Regenerate information if read failed
			data, err = storage.Read(servicePath)
			if err == nil {
				_ = entry.Unmarshal(data)
			}
			if opt.Git {
				entry.Git = gitEntry
			}
			if opt.Version != "" {
				entry.Version = opt.Version
			}
			if entry.Metadata == nil {
				entry.Metadata = make(map[string]string)
			}
			for k, v := range meta {
				entry.Metadata[k] = v
			}
			entry.Name = serviceName
			entry.Kind = storage.KindService
			return storage.Save(entry, servicePath)
		},
	}

	cmd.Flags().BoolVar(&opt.Git, "git", opt.Git, "Whether to automatically parse git information")
	cmd.Flags().StringVar(&opt.Version, "version", opt.Version, "Set service version. e.g. 1.0.0、1.0.0-beta.1")
	cmd.Flags().StringArrayVar(&opt.Metadata, "meta", nil, "Set custom key-value pair parameters. e.g. status=ok")
	return cmd
}

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
	"errors"
	"fmt"
	"strings"

	"github.com/99nil/gopkg/sets"
	"github.com/spf13/cobra"

	"github.com/zc2638/releaser/pkg/git"
	"github.com/zc2638/releaser/pkg/storage"
	"github.com/zc2638/releaser/pkg/util"
)

type SetOption struct {
	Git         bool
	ProjectPath string
	Version     string
	Metadata    []string
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

			meta := make(map[string]string)
			for _, v := range opt.Metadata {
				arr := strings.SplitN(v, "=", 2)
				if len(arr) != 2 {
					continue
				}
				meta[arr[0]] = arr[1]
			}

			if len(args) == 0 {
				return setManifest(opt, manifest, meta)
			}
			return setService(opt, manifest, meta, args[0])
		},
	}

	cmd.Flags().BoolVar(&opt.Git, "git", opt.Git, "Whether to automatically parse git information")
	cmd.Flags().StringVarP(&opt.ProjectPath, "project-path", "p", opt.ProjectPath,
		"The root path of the project that needs to be imported as a service, the system will automatically find and resolve the storage directory under it")
	cmd.Flags().StringVar(&opt.Version, "version", opt.Version, "Set service version. e.g. 1.0.0、1.0.0-beta.1")
	cmd.Flags().StringArrayVar(&opt.Metadata, "meta", nil, "Set custom key-value pair parameters. e.g. status=ok")
	return cmd
}

func setManifest(opt *SetOption, manifest *storage.Manifest, meta map[string]string) error {
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

		manifest.Git = &storage.GitEntry{
			Branch: branch,
			Commit: commit,
			Tag:    tag,
		}
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

func setService(opt *SetOption, manifest *storage.Manifest, meta map[string]string, serviceName string) error {
	servicePath := serviceName + ".yaml"
	ss := sets.NewString(manifest.Services...)
	if !ss.Has(serviceName) {
		return fmt.Errorf("service %s not exist", serviceName)
	}

	entry := &storage.Entry{}
	if opt.ProjectPath != "" {
		data, err := storage.ReadWithRoot(opt.ProjectPath, storage.ManifestName)
		if err != nil {
			return fmt.Errorf("read [project as service] manifest failed: %v", err)
		}
		svcManifest := &storage.Manifest{}
		if err := svcManifest.Unmarshal(data); err != nil {
			return fmt.Errorf("unmarshal [project as service] manifest failed: %v", err)
		}
		entry = &svcManifest.Entry
	}

	// Regenerate information if read failed
	data, err := storage.Read(servicePath)
	if err == nil {
		_ = entry.Unmarshal(data)
	}
	if opt.Git {
		return errors.New("does not support the service to automatically generate git information")
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
}

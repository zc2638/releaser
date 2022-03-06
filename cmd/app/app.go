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
	"os"

	"github.com/99nil/go/sets"
	"github.com/spf13/cobra"

	"github.com/zc2638/releaser"
	"github.com/zc2638/releaser/pkg/storage"
	"github.com/zc2638/releaser/pkg/util"
)

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "releaser",
		Version: releaser.Version.String(),
	}
	cmd.AddCommand(
		newInitCommand(),
		newCreateCommand(),
		newSetCommand(),
		newGetCommand(),
	)
	return cmd
}

func newInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "init",
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("please enter the project name")
			}
			if err := util.CheckArgs(args, 1); err != nil {
				return err
			}
			if err := initStorage(); err != nil {
				return fmt.Errorf("init storage failed: %v", err)
			}

			manifest := &storage.Manifest{
				Entry: storage.Entry{
					Name: args[0],
					Kind: storage.KindProject,
				},
			}
			if err := storage.Save(manifest, storage.ManifestName); err != nil {
				return fmt.Errorf("init manifest failed: %v", err)
			}
			return nil
		},
	}

	return cmd
}

func newCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "create",
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("please enter the service name")
			}
			if err := util.CheckArgs(args, 1); err != nil {
				return err
			}
			serviceName := args[0]

			data, err := storage.Read(storage.ManifestName)
			if err != nil {
				return fmt.Errorf("read manifest failed: %v", err)
			}
			manifest := &storage.Manifest{}
			if err := manifest.Unmarshal(data); err != nil {
				return fmt.Errorf("unmarshal manifest failed: %v", err)
			}

			ss := sets.NewString(manifest.Services...)
			if ss.Has(serviceName) {
				return fmt.Errorf("service %s already exists", serviceName)
			}
			ss.Add(serviceName)
			manifest.Services = ss.List()
			if err := storage.Save(manifest, storage.ManifestName); err != nil {
				return fmt.Errorf("save manifest failed: %v", err)
			}
			return nil
		},
	}

	return cmd
}

func initStorage() error {
	_, err := os.Stat(storage.DefaultPath)
	if err == nil {
		return nil
	}
	if !os.IsNotExist(err) {
		return err
	}
	return os.MkdirAll(storage.DefaultPath, os.ModePerm)
}

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

package storage

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// DefaultPath defines the default storage dir path
const DefaultPath = ".releaser"

const ManifestName = "manifest"

const (
	KindProject = "Project"
	KindService = "Service"
)

const (
	OutputFormatGoBuild = "gobuild"
	OutputFormatJSON    = "json"
)

type Marshaler interface {
	Marshal() ([]byte, error)
}

func Read(name string) ([]byte, error) {
	return ReadWithRoot("", name)
}

func Save(m Marshaler, name string) error {
	return SaveWithRoot(m, "", name)
}

func ReadWithRoot(root, name string) ([]byte, error) {
	return os.ReadFile(filepath.Join(root, DefaultPath, name))
}

func SaveWithRoot(m Marshaler, root, name string) error {
	b, err := m.Marshal()
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(root, DefaultPath, name), b, os.ModePerm)
}

type Manifest struct {
	Entry    `json:",inline" yaml:",inline"`
	Services []string `json:"services,omitempty" yaml:"services,omitempty"`
}

func (m *Manifest) Marshal() ([]byte, error) {
	return yaml.Marshal(&m)
}

func (m *Manifest) Unmarshal(in []byte) error {
	return yaml.Unmarshal(in, &m)
}

type GitEntry struct {
	Branch string `json:"branch" yaml:"branch"`
	Commit string `json:"commit" yaml:"commit"`
	Tag    string `json:"tag" yaml:"tag"`
}

type Entry struct {
	Name     string            `json:"name" yaml:"name"`
	Kind     string            `json:"kind" yaml:"kind"`
	Version  string            `json:"version" yaml:"version"`
	Metadata map[string]string `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	Git      *GitEntry         `json:"git,omitempty" yaml:"git,omitempty"`
}

func (e *Entry) Marshal() ([]byte, error) {
	return yaml.Marshal(&e)
}

func (e *Entry) Unmarshal(in []byte) error {
	return yaml.Unmarshal(in, &e)
}

func (e *Entry) ToMap() map[string]string {
	m := make(map[string]string)
	m[buildWithPrefix("name")] = e.Name
	m[buildWithPrefix("kind")] = e.Kind
	m[buildWithPrefix("version")] = e.Version
	for k, v := range e.Metadata {
		m[buildWithPrefix("meta."+k)] = v
	}
	return m
}

func buildWithPrefix(s string) string {
	return "releaser." + s
}

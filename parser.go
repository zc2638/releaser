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
	"encoding/base64"

	"gopkg.in/yaml.v3"
)

func Encode(ve *VersionEntry) (string, error) {
	b, err := BuildYaml(ve)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func Decode(data string) (*VersionEntry, error) {
	b, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	return ParseYaml(b)
}

func ParseYaml(b []byte) (*VersionEntry, error) {
	var ve VersionEntry
	if err := yaml.Unmarshal(b, &ve); err != nil {
		return nil, err
	}
	return &ve, nil
}

func BuildYaml(ve *VersionEntry) ([]byte, error) {
	return yaml.Marshal(ve)
}

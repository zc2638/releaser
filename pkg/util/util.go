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

package util

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func CheckArgs(args []string, num int) error {
	l := len(args)
	if l != num {
		return fmt.Errorf("unusual number of parameters, expect %d, take %d", num, l)
	}
	return nil
}

func CheckArgsGT(args []string, num int) error {
	l := len(args)
	if l > num {
		return fmt.Errorf("unusual number of parameters, max expect %d, take %d", num, l)
	}
	return nil
}

func Exec(command string) ([]byte, error) {
	return ExecWithEnv(command, nil)
}

func ExecWithEnv(command string, envMap map[string]string) ([]byte, error) {
	args := strings.Split(command, " ")
	var cmd *exec.Cmd
	if len(args) > 1 {
		execPath, err := exec.LookPath("bash")
		if err == nil {
			cmd = exec.Command(execPath, "-c", command)
		} else if execPath, err = exec.LookPath("sh"); err == nil {
			cmd = exec.Command(execPath, "-c", command)
		} else {
			cmd = exec.Command(args[0], args[1:]...)
		}
	} else {
		cmd = exec.Command(args[0])
	}

	currentEnvSet := os.Environ()
	envSet := make([]string, 0, len(envMap)+len(currentEnvSet))
	envSet = append(envSet, currentEnvSet...)
	for k, v := range envMap {
		envSet = append(envSet, k+"="+v)
	}
	cmd.Env = envSet

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("%v: %s", err, output)
	}
	return output, nil
}

// Copyright 2017 Xiaomi, Inc.
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

package g

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

/*
	以下是open-falcon的工具模块，主要是对配置文件、模块顺序、日志文件、模块进程
	提供相应的工具函数
*/

/* 判断当前日志文件是否存在 */
func HasLogfile(name string) bool {
	if _, err := os.Stat(LogPath(name)); err != nil {
		return false
	}
	return true
}

/* 对模块排序 */
func PreqOrder(moduleArgs []string) []string {
	/* 模块为空，返回空 */
	if len(moduleArgs) == 0 {
		return []string{}
	}

	var modulesInOrder []string

	// get arguments which are found in the order
	/* 按顺序查找出模块 */
	for _, nameOrder := range AllModulesInOrder {
		for _, nameArg := range moduleArgs {
			if nameOrder == nameArg {
				modulesInOrder = append(modulesInOrder, nameOrder)
			}
		}
	}

	// get arguments which are not found in the order
	/* 把没有在AllModulesInOrder的模块添加到modulesInOrder中 */
	for _, nameArg := range moduleArgs {
		end := 0
		for _, nameOrder := range modulesInOrder {
			if nameOrder == nameArg {
				break
			}
			end++
		}
		if end == len(modulesInOrder) {
			modulesInOrder = append(modulesInOrder, nameArg)
		}
	}
	return modulesInOrder
}

/*
	返回给当前路劲和给出路劲的相对路劲
	可以查看tool_test.go的Rel单元测试
*/
func Rel(p string) string {
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}

	// filepath.Abs() returns an error only when os.Getwd() returns an error;
	abs, _ := filepath.Abs(p)

	r, err := filepath.Rel(wd, abs)
	if err != nil {
		return ""
	}

	return r
}

/* 判断是否有配置 */
func HasCfg(name string) bool {
	if _, err := os.Stat(Cfg(name)); err != nil {
		return false
	}
	return true
}

/* 判断是否有name模块 */
func HasModule(name string) bool {
	return Modules[name]
}

/* 设置进程pid, 通过os/exec模块获取进程pid, 并放在PidOf, pid可以是多个 */
func setPid(name string) {
	output, _ := exec.Command("pgrep", "-f", ModuleApps[name]).Output()
	pidStr := strings.TrimSpace(string(output))
	PidOf[name] = pidStr
}

/* 获取模块的pid,如果没有，则先设置，然后返回 */
func Pid(name string) string {
	if PidOf[name] == "<NOT SET>" {
		setPid(name)
	}
	return PidOf[name]
}

/* 如果当前的pid不是空，那么返回true,否则返回false */
func IsRunning(name string) bool {
	setPid(name)
	return Pid(name) != ""
}

/* 参数去重 */
func RmDup(args []string) []string {
	if len(args) == 0 {
		return []string{}
	}
	if len(args) == 1 {
		return args
	}

	ret := []string{}
	isDup := make(map[string]bool)
	for _, arg := range args {
		if isDup[arg] == true {
			continue
		}
		ret = append(ret, arg)
		isDup[arg] = true
	}
	return ret
}

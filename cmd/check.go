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

package cmd

import (
	"fmt"
	"strings"

	"github.com/open-falcon/falcon-plus/g"
	"github.com/spf13/cobra"
)

/*
	Check命令，RunE是运行后返回错误
	g.AllModulesInOrder是字符串数组，包含open-falcon的模块名称，并且按照一定的顺序，
*/
var Check = &cobra.Command{
	Use:   "check [Module ...]",
	Short: "Check the status of Open-Falcon modules",
	Long: `
Check if the specified Open-Falcon modules are running.
Modules:
  ` + "all " + strings.Join(g.AllModulesInOrder, " "),
	RunE: check,
}

/* Check命令: 主要检查哪些模块是在运行的，哪些没有运行 */
func check(c *cobra.Command, args []string) error {
	args = g.RmDup(args) /* 复制参数并重名 */

	/*
		如果参数为空，则默认所有的模块(agent/aggregator/graph/hbs/judge...)[在g/g.go文件]
		为什么是g/g.go文件，名字好怪，为什么不可以是utility/utils.go?
	*/
	if len(args) == 0 {
		args = g.AllModulesInOrder
	}

	for _, moduleName := range args {
		if !g.HasModule(moduleName) {
			return fmt.Errorf("%s doesn't exist", moduleName)
		}

		/* 如果模块在运行，返回模块App名称（比较正经）, UP, 模块的进程PID */
		if g.IsRunning(moduleName) {
			fmt.Printf("%20s %10s %15s \n", g.ModuleApps[moduleName], "UP", g.Pid(moduleName))
		} else {
			fmt.Printf("%20s %10s %15s \n", g.ModuleApps[moduleName], "DOWN", "-")
		}
	}

	return nil
}

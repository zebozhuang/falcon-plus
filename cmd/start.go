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
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/open-falcon/falcon-plus/g"
	"github.com/spf13/cobra"
)

/*
	Open-Falcon启动模块: start [Module ...]
*/
var Start = &cobra.Command{
	Use:   "start [Module ...]",
	Short: "Start Open-Falcon modules",
	Long: `
Start the specified Open-Falcon modules and run until a stop command is received.
A module represents a single node in a cluster.
Modules:
	` + "all " + strings.Join(g.AllModulesInOrder, " "), /* g.AllModulesInOrder是一个字符串数组，包含所有模块的名称 */
	RunE:          start,
	SilenceUsage:  true, /* TODO: 这个两个使用的效果? 试了下，用跟没用没区别。 */
	SilenceErrors: true,
}

/*
	PreOrderFlag在falcon-plus\main.go文件中初始化，默认false, 表示模块是否需要排序
	ConsoleOutputFlag在falcon-plus\main.go文件中初始化，默认false，表示模块是否需要在console输出日志
*/
var PreqOrderFlag bool
var ConsoleOutputFlag bool

/* 获取模块参数，即获取配置文件路径 */
func cmdArgs(name string) []string {
	return []string{"-c", g.Cfg(name)}
}

/* 打开日志文件 */
func openLogFile(name string) (*os.File, error) {
	logDir := g.LogDir(name)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	logPath := g.LogPath(name)
	logOutput, err := os.OpenFile(logPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return logOutput, nil
}

/* 执行模块, co参数表示打印log日志到console, name表示模块名字 */
func execModule(co bool, name string) error {
	/* 执行： 命令+参数（配置）*/
	cmd := exec.Command(g.Bin(name), cmdArgs(name)...)

	if co {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	/* 日志 */
	logOutput, err := openLogFile(name)
	if err != nil {
		return err
	}
	defer logOutput.Close()
	cmd.Stdout = logOutput
	cmd.Stderr = logOutput
	return cmd.Start()
}

/* 模块检查:(1)模块存在 (2)模块配置*/
func checkStartReq(name string) error {
	if !g.HasModule(name) {
		return fmt.Errorf("%s doesn't exist", name)
	}

	if !g.HasCfg(name) {
		r := g.Rel(g.Cfg(name))
		return fmt.Errorf("expect config file: %s", r)
	}

	return nil
}

/* 模块是否启动 */
func isStarted(name string) bool {
	ticker := time.NewTicker(time.Millisecond * 100)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C: /* 一直在检测, 不会跑到time.After，是否是BUG */
			if g.IsRunning(name) {
				return true
			}
		case <-time.After(time.Second): /* 不会运行到这里，BUG? */
			return false
		}
	}
}

func start(c *cobra.Command, args []string) error {
	args = g.RmDup(args)

	if PreqOrderFlag {
		args = g.PreqOrder(args)
	}

	if len(args) == 0 {
		args = g.AllModulesInOrder
	}

	for _, moduleName := range args {
		if err := checkStartReq(moduleName); err != nil {
			return err
		}

		// Skip starting if the module is already running
		if g.IsRunning(moduleName) {
			fmt.Print("[", g.ModuleApps[moduleName], "] ", g.Pid(moduleName), "\n")
			continue
		}

		/* 执行模块 */
		if err := execModule(ConsoleOutputFlag, moduleName); err != nil {
			return err
		}
		/* 再次检测是否在运行 */
		if isStarted(moduleName) {
			fmt.Print("[", g.ModuleApps[moduleName], "] ", g.Pid(moduleName), "\n")
			continue
		}

		return fmt.Errorf("[%s] failed to start", g.ModuleApps[moduleName])
	}
	return nil
}

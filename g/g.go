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

import "path/filepath"

/* 模块的名称、二进制路劲、配置路径、app名称、日志路劲、进程、部署顺序 */
var Modules map[string]bool
var BinOf map[string]string
var cfgOf map[string]string
var ModuleApps map[string]string
var logpathOf map[string]string
var PidOf map[string]string
var AllModulesInOrder []string

func init() {
	/* 当前已经实现的模块 */
	Modules = map[string]bool{
		"agent":      true,
		"aggregator": true,
		"graph":      true,
		"hbs":        true,
		"judge":      true,
		"nodata":     true,
		"transfer":   true,
		"gateway":    true,
		"api":        true,
		"alarm":      true,
	}

	/* 模块的二进制路劲 BinOf */
	BinOf = map[string]string{
		"agent":      "./agent/bin/falcon-agent",
		"aggregator": "./aggregator/bin/falcon-aggregator",
		"graph":      "./graph/bin/falcon-graph",
		"hbs":        "./hbs/bin/falcon-hbs",
		"judge":      "./judge/bin/falcon-judge",
		"nodata":     "./nodata/bin/falcon-nodata",
		"transfer":   "./transfer/bin/falcon-transfer",
		"gateway":    "./gateway/bin/falcon-gateway",
		"api":        "./api/bin/falcon-api",
		"alarm":      "./alarm/bin/falcon-alarm",
	}

	/* 各个服务的配置路劲 cfgOf */
	cfgOf = map[string]string{
		"agent":      "./agent/config/cfg.json",
		"aggregator": "./aggregator/config/cfg.json",
		"graph":      "./graph/config/cfg.json",
		"hbs":        "./hbs/config/cfg.json",
		"judge":      "./judge/config/cfg.json",
		"nodata":     "./nodata/config/cfg.json",
		"transfer":   "./transfer/config/cfg.json",
		"gateway":    "./gateway/config/cfg.json",
		"api":        "./api/config/cfg.json",
		"alarm":      "./alarm/config/cfg.json",
	}

	/* 模块对应app名称 */
	ModuleApps = map[string]string{
		"agent":      "falcon-agent",
		"aggregator": "falcon-aggregator",
		"graph":      "falcon-graph",
		"hbs":        "falcon-hbs",
		"judge":      "falcon-judge",
		"nodata":     "falcon-nodata",
		"transfer":   "falcon-transfer",
		"gateway":    "falcon-gateway",
		"api":        "falcon-api",
		"alarm":      "falcon-alarm",
	}

	/* 模块的日志路劲 */
	logpathOf = map[string]string{
		"agent":      "./agent/logs/agent.log",
		"aggregator": "./aggregator/logs/aggregator.log",
		"graph":      "./graph/logs/graph.log",
		"hbs":        "./hbs/logs/hbs.log",
		"judge":      "./judge/logs/judge.log",
		"nodata":     "./nodata/logs/nodata.log",
		"transfer":   "./transfer/logs/transfer.log",
		"gateway":    "./gateway/logs/gateway.log",
		"api":        "./api/logs/api.log",
		"alarm":      "./alarm/logs/alarm.log",
	}

	/* 模块的进程pid,初始化为<NOT SET> */
	PidOf = map[string]string{
		"agent":      "<NOT SET>",
		"aggregator": "<NOT SET>",
		"graph":      "<NOT SET>",
		"hbs":        "<NOT SET>",
		"judge":      "<NOT SET>",
		"nodata":     "<NOT SET>",
		"transfer":   "<NOT SET>",
		"gateway":    "<NOT SET>",
		"api":        "<NOT SET>",
		"alarm":      "<NOT SET>",
	}

	// Modules are deployed in this order
	/* 按顺序部署模块 */
	AllModulesInOrder = []string{
		"graph",
		"hbs",
		"judge",
		"transfer",
		"nodata",
		"aggregator",
		"agent",
		"gateway",
		"api",
		"alarm",
	}
}

/* 获取模块二进制的路劲 */
func Bin(name string) string {
	p, _ := filepath.Abs(BinOf[name])
	return p
}

/* 获取模块的配置路径 */
func Cfg(name string) string {
	p, _ := filepath.Abs(cfgOf[name])
	return p
}

/* 获取日志的路劲 */
func LogPath(name string) string {
	p, _ := filepath.Abs(logpathOf[name])
	return p
}

/* 获取日志的目录 */
func LogDir(name string) string {
	d, _ := filepath.Abs(filepath.Dir(logpathOf[name]))
	return d
}

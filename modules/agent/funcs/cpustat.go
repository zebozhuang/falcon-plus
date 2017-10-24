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

package funcs

import (
	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/toolkits/nux"
	"sync"
)

const (
	historyCount int = 2
)

var (
	procStatHistory [historyCount]*nux.ProcStat
	psLock          = new(sync.RWMutex)
)

/* CPU数据 */
func UpdateCpuStat() error {
	ps, err := nux.CurrentProcStat()
	if err != nil {
		return err
	}

	psLock.Lock()
	defer psLock.Unlock()
	for i := historyCount - 1; i > 0; i-- {
		procStatHistory[i] = procStatHistory[i-1]
	}

	procStatHistory[0] = ps
	return nil
}

/* CPU差值 */
func deltaTotal() uint64 {
	if procStatHistory[1] == nil {
		return 0
	}
	return procStatHistory[0].Cpu.Total - procStatHistory[1].Cpu.Total
}

/* CPU闲置 */
func CpuIdle() float64 {
	psLock.RLock()
	defer psLock.RUnlock()
	dt := deltaTotal()
	if dt == 0 {
		return 0.0
	}
	invQuotient := 100.00 / float64(dt)
	return float64(procStatHistory[0].Cpu.Idle-procStatHistory[1].Cpu.Idle) * invQuotient
}

/* CPU用户使用 */
func CpuUser() float64 {
	psLock.RLock()
	defer psLock.RUnlock()
	dt := deltaTotal()
	if dt == 0 {
		return 0.0
	}
	invQuotient := 100.00 / float64(dt)
	return float64(procStatHistory[0].Cpu.User-procStatHistory[1].Cpu.User) * invQuotient
}

/* CPU Nice值 */
func CpuNice() float64 {
	psLock.RLock()
	defer psLock.RUnlock()
	dt := deltaTotal()
	if dt == 0 {
		return 0.0
	}
	invQuotient := 100.00 / float64(dt)
	return float64(procStatHistory[0].Cpu.Nice-procStatHistory[1].Cpu.Nice) * invQuotient
}

/* CPU 核心时间 */
func CpuSystem() float64 {
	psLock.RLock()
	defer psLock.RUnlock()
	dt := deltaTotal()
	if dt == 0 {
		return 0.0
	}
	invQuotient := 100.00 / float64(dt)
	return float64(procStatHistory[0].Cpu.System-procStatHistory[1].Cpu.System) * invQuotient
}

/*CPU IO等待时间 */
func CpuIowait() float64 {
	psLock.RLock()
	defer psLock.RUnlock()
	dt := deltaTotal()
	if dt == 0 {
		return 0.0
	}
	invQuotient := 100.00 / float64(dt)
	return float64(procStatHistory[0].Cpu.Iowait-procStatHistory[1].Cpu.Iowait) * invQuotient
}

/* CPU 硬中断时间 */
func CpuIrq() float64 {
	psLock.RLock()
	defer psLock.RUnlock()
	dt := deltaTotal()
	if dt == 0 {
		return 0.0
	}
	invQuotient := 100.00 / float64(dt)
	return float64(procStatHistory[0].Cpu.Irq-procStatHistory[1].Cpu.Irq) * invQuotient
}

/* CPU 软中断时间 */
func CpuSoftIrq() float64 {
	psLock.RLock()
	defer psLock.RUnlock()
	dt := deltaTotal()
	if dt == 0 {
		return 0.0
	}
	invQuotient := 100.00 / float64(dt)
	return float64(procStatHistory[0].Cpu.SoftIrq-procStatHistory[1].Cpu.SoftIrq) * invQuotient
}

/* CPU 虚拟机的时间之比 */
func CpuSteal() float64 {
	psLock.RLock()
	defer psLock.RUnlock()
	dt := deltaTotal()
	if dt == 0 {
		return 0.0
	}
	invQuotient := 100.00 / float64(dt)
	return float64(procStatHistory[0].Cpu.Steal-procStatHistory[1].Cpu.Steal) * invQuotient
}

/* 虚拟CPU所占的比重 */
func CpuGuest() float64 {
	psLock.RLock()
	defer psLock.RUnlock()
	dt := deltaTotal()
	if dt == 0 {
		return 0.0
	}
	invQuotient := 100.00 / float64(dt)
	return float64(procStatHistory[0].Cpu.Guest-procStatHistory[1].Cpu.Guest) * invQuotient
}

/* 上下文切换数 */
func CurrentCpuSwitches() uint64 {
	psLock.RLock()
	defer psLock.RUnlock()
	return procStatHistory[0].Ctxt
}

func CpuPrepared() bool {
	psLock.RLock()
	defer psLock.RUnlock()
	return procStatHistory[1] != nil
}

/* CPU 测量数据 */
func CpuMetrics() []*model.MetricValue {
	if !CpuPrepared() {
		return []*model.MetricValue{}
	}

	cpuIdleVal := CpuIdle()
	idle := GaugeValue("cpu.idle", cpuIdleVal)
	busy := GaugeValue("cpu.busy", 100.0-cpuIdleVal)
	user := GaugeValue("cpu.user", CpuUser())
	nice := GaugeValue("cpu.nice", CpuNice())
	system := GaugeValue("cpu.system", CpuSystem())
	iowait := GaugeValue("cpu.iowait", CpuIowait())
	irq := GaugeValue("cpu.irq", CpuIrq())
	softirq := GaugeValue("cpu.softirq", CpuSoftIrq())
	steal := GaugeValue("cpu.steal", CpuSteal())
	guest := GaugeValue("cpu.guest", CpuGuest())
	switches := CounterValue("cpu.switches", CurrentCpuSwitches())
	return []*model.MetricValue{idle, busy, user, nice, system, iowait, irq, softirq, steal, guest, switches}
}

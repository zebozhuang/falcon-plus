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

package sender

import (
	"time"

	"github.com/open-falcon/falcon-plus/common/model"
)

/* 生成元数据对象 */
func MakeMetaData(endpoint, metric, tags string, val interface{}, counterType string, step_and_ts ...int64) *model.JsonMetaData {
	md := model.JsonMetaData{
		Endpoint:    endpoint,
		Metric:      metric,
		Tags:        tags,
		Value:       val,
		CounterType: counterType,
	}

	argc := len(step_and_ts)
	if argc == 0 {
		md.Step = 60
		md.Timestamp = time.Now().Unix()
	} else if argc == 1 {
		md.Step = step_and_ts[0]
		md.Timestamp = time.Now().Unix()
	} else if argc == 2 {
		md.Step = step_and_ts[0]
		md.Timestamp = step_and_ts[1]
	}

	return &md
}

/* 生成测量数据 */
func MakeGaugeValue(endpoint, metric, tags string, val interface{}, step_and_ts ...int64) *model.JsonMetaData {
	return MakeMetaData(endpoint, metric, tags, val, "GAUGE", step_and_ts...)
}

/* 生成计数数据 */
func MakeCounterValue(endpoint, metric, tags string, val interface{}, step_and_ts ...int64) *model.JsonMetaData {
	return MakeMetaData(endpoint, metric, tags, val, "COUNTER", step_and_ts...)
}

/* 添加测量数据到队列，MetaDataQueue在sender.go*/
func PushGauge(endpoint, metric, tags string, val interface{}, step_and_ts ...int64) {
	md := MakeGaugeValue(endpoint, metric, tags, val, step_and_ts...)
	MetaDataQueue.PushFront(md)
}

/*添加计数数据到队列 */
func PushCounter(endpoint, metric, tags string, val interface{}, step_and_ts ...int64) {
	md := MakeCounterValue(endpoint, metric, tags, val, step_and_ts...)
	MetaDataQueue.PushFront(md)
}

/*添加[测量|计数]数据到队列 */
func Push(endpoint, metric, tags string, val interface{}, counterType string, step_and_ts ...int64) {
	md := MakeMetaData(endpoint, metric, tags, val, counterType, step_and_ts...)
	MetaDataQueue.PushFront(md)
}

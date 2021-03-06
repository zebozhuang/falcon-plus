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

package utils

import (
	"fmt"
)

/* 是指生成PRIMARY KEY? */
func PK(endpoint, metric string, tags map[string]string) string {
	if tags == nil || len(tags) == 0 {
		return fmt.Sprintf("%s/%s", endpoint, metric)
	}
	return fmt.Sprintf("%s/%s/%s", endpoint, metric, SortedTags(tags))
}

func PK2(endpoint, counter string) string {
	return fmt.Sprintf("%s/%s", endpoint, counter)
}

/* 按照一定格式生成uuid,不是真正意义上的uuid */
func UUID(endpoint, metric string, tags map[string]string, dstype string, step int) string {
	if tags == nil || len(tags) == 0 {
		return fmt.Sprintf("%s/%s/%s/%d", endpoint, metric, dstype, step)
	}
	return fmt.Sprintf("%s/%s/%s/%s/%d", endpoint, metric, SortedTags(tags), dstype, step)
}

/* 对数据生成checksum md5校验值 */
func Checksum(endpoint string, metric string, tags map[string]string) string {
	pk := PK(endpoint, metric, tags)
	return Md5(pk)
}

/* 对数据生成uuid,并生成checksum校验值 */
func ChecksumOfUUID(endpoint, metric string, tags map[string]string, dstype string, step int64) string {
	return Md5(UUID(endpoint, metric, tags, dstype, int(step)))
}

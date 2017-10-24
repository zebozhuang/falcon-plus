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

package backend_pool

import (
	"fmt"
	"net"
	"time"

	connp "github.com/toolkits/conn_pool"
)

/*
	TSDB客户端连接对象
	实现 Name, Closed, Close函数
*/

// TSDB
type TsdbClient struct {
	cli  net.Conn
	name string
}

func (t TsdbClient) Name() string {
	return t.name
}

func (t TsdbClient) Closed() bool {
	return t.cli == nil
}

func (t TsdbClient) Close() error {
	if t.cli != nil {
		err := t.cli.Close()
		t.cli = nil
		return err
	}
	return nil
}

/* 创建TsdbConnPool 连接池, 参数：连接地址，最大连接数，最大闲置数，超时时间 */
func newTsdbConnPool(address string, maxConns int, maxIdle int, connTimeout int) *connp.ConnPool {
	/* ConnPool提供连接池管理, 返回的NConn是一个包含(Name, Closed, Close函数的interface) */
	pool := connp.NewConnPool("tsdb", address, int32(maxConns), int32(maxIdle))

	/* 实现自定义创建连接对象 */
	pool.New = func(name string) (connp.NConn, error) {
		_, err := net.ResolveTCPAddr("tcp", address)
		if err != nil {
			return nil, err
		}

		conn, err := net.DialTimeout("tcp", address, time.Duration(connTimeout)*time.Millisecond)
		if err != nil {
			return nil, err
		}

		return TsdbClient{conn, name}, nil
	}

	return pool
}

/* TSDB连接池管理（Helper） */
type TsdbConnPoolHelper struct {
	p           *connp.ConnPool
	maxConns    int
	maxIdle     int
	connTimeout int
	callTimeout int
	address     string
}

/*创建TSDB连接池管理对象 */
func NewTsdbConnPoolHelper(address string, maxConns, maxIdle, connTimeout, callTimeout int) *TsdbConnPoolHelper {
	return &TsdbConnPoolHelper{
		p:           newTsdbConnPool(address, maxConns, maxIdle, connTimeout),
		maxConns:    maxConns,
		maxIdle:     maxIdle,
		connTimeout: connTimeout,
		callTimeout: callTimeout,
		address:     address,
	}
}

/* 发送数据 */
func (t *TsdbConnPoolHelper) Send(data []byte) (err error) {
	conn, err := t.p.Fetch()
	if err != nil {
		return fmt.Errorf("get connection fail: err %v. proc: %s", err, t.p.Proc())
	}

	cli := conn.(TsdbClient).cli

	done := make(chan error, 1)
	go func() {
		_, err = cli.Write(data)
		done <- err
	}()

	select {
	case <-time.After(time.Duration(t.callTimeout) * time.Millisecond):
		t.p.ForceClose(conn)
		return fmt.Errorf("%s, call timeout", t.address)
	case err = <-done:
		if err != nil {
			t.p.ForceClose(conn)
			err = fmt.Errorf("%s, call failed, err %v. proc: %s", t.address, err, t.p.Proc())
		} else {
			t.p.Release(conn)
		}
		return err
	}
}

/* 销毁连接池 */
func (t *TsdbConnPoolHelper) Destroy() {
	if t.p != nil {
		t.p.Destroy()
	}
}

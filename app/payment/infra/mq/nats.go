// Copyright 2024 CloudWeGo Authors
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

package mq

import (
	"github.com/nats-io/nats.go"
)

var (
	Nc  *nats.Conn
	err error
)

func Init() {
	Nc, err = nats.Connect("nats-svc")
	if err != nil {
		panic(err)
	}
}

const (
	// Stream 名称
	OrderStreamName = "ORDERS"

	// 主题
	OrderCancelSubject   = "order.cancel"   // 订单取消
	PaymentCancelSubject = "payment.cancel" // 支付取消
)

// 订单取消消息
type OrderCancelMessage struct {
	OrderID     string `json:"order_id"`
	UserID      int64  `json:"user_id"`
	ExpireTime  int64  `json:"expire_time"`  // 过期时间戳
	OrderStatus string `json:"order_status"` // 当前订单状态
}

// 支付取消消息
type PaymentCancelMessage struct {
	OrderID    string `json:"order_id"`
	UserID     int64  `json:"user_id"`
	PaymentID  string `json:"payment_id"`
	ExpireTime int64  `json:"expire_time"` // 过期时间戳
}

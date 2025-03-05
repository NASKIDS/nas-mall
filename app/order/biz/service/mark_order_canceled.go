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

package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/naskids/nas-mall/app/order/biz/dal/mysql"
	"github.com/naskids/nas-mall/app/order/biz/model"
	order "github.com/naskids/nas-mall/rpc_gen/kitex_gen/order"
)

const (
	// 锁前缀和过期时间
	OrderCancelLockPrefix = "lock:order:cancel:"
	LockExpiration        = 10 * time.Second
)

type MarkOrderCanceledService struct {
	ctx context.Context
}

// NewMarkOrderCanceledService new MarkOrderCanceledService
func NewMarkOrderCanceledService(ctx context.Context) *MarkOrderCanceledService {
	return &MarkOrderCanceledService{ctx: ctx}
}

// Run create note info
func (s *MarkOrderCanceledService) Run(req *order.MarkOrderCanceledReq) (resp *order.MarkOrderCanceledResp, err error) {
	// Finish your business logic.
	if req.UserId == 0 || req.OrderId == "" {
		err = fmt.Errorf("user_id or order_id can not be empty")
		return
	}

	orderResult, err := model.GetOrder(mysql.DB, s.ctx, req.UserId, req.OrderId)
	if err != nil {
		return nil, err
	}

	// 检查订单状态，只有未支付的订单才能取消
	if orderResult.OrderState != model.OrderStatePlaced {
		log.Printf("订单[%s]当前状态为[%s]，无需取消", req.OrderId, orderResult.OrderState)
		return &order.MarkOrderCanceledResp{}, nil // 非未支付状态，无需取消，视为成功
	}

	err = model.UpdateOrderState(mysql.DB, s.ctx, req.UserId, req.OrderId, model.OrderStateCanceled)
	if err != nil {

		return nil, err
	}
	log.Printf("订单[%s]已成功取消", req.OrderId)
	resp = &order.MarkOrderCanceledResp{}
	return
}

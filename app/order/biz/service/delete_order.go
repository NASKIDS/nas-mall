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

    "github.com/cloudwego/kitex/pkg/klog"

    "github.com/naskids/nas-mall/app/order/biz/dal/mysql"
    "github.com/naskids/nas-mall/app/order/biz/model"
    order "github.com/naskids/nas-mall/rpc_gen/kitex_gen/order"
)

type DeleteOrderService struct {
    ctx context.Context
}

// NewDeleteOrderService new DeleteOrderService
func NewDeleteOrderService(ctx context.Context) *DeleteOrderService {
    return &DeleteOrderService{ctx: ctx}
}

// Run delete order
func (s *DeleteOrderService) Run(req *order.DeleteOrderReq) (resp *order.DeleteOrderResp, err error) {
    // 参数验证
    if req.UserId == 0 || req.OrderId == "" {
        err = fmt.Errorf("user_id or order_id can not be empty")
        return
    }

    // 检查订单是否存在
    _, err = model.GetOrder(mysql.DB, s.ctx, req.UserId, req.OrderId)
    if err != nil {
        klog.Errorf("model.GetOrder error:%v", err)
        return nil, err
    }

    // 删除订单
    err = model.DeleteOrder(mysql.DB, s.ctx, req.UserId, req.OrderId)
    if err != nil {
        klog.Errorf("model.DeleteOrder error:%v", err)
        return nil, err
    }

    resp = &order.DeleteOrderResp{}
    return
} 
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

package main

import (
	"context"

	"github.com/naskids/nas-mall/app/order/biz/service"
	order "github.com/naskids/nas-mall/rpc_gen/kitex_gen/order"
)

// OrderServiceImpl implements the last service interface defined in the IDL.
type OrderServiceImpl struct{}

// PlaceOrder implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) PlaceOrder(ctx context.Context, req *order.PlaceOrderReq) (resp *order.PlaceOrderResp, err error) {
	resp, err = service.NewPlaceOrderService(ctx).Run(req)

	return resp, err
}

// ListOrder implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) ListOrder(ctx context.Context, req *order.ListOrderReq) (resp *order.ListOrderResp, err error) {
	resp, err = service.NewListOrderService(ctx).Run(req)

	return resp, err
}

// MarkOrderPaid implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) MarkOrderPaid(ctx context.Context, req *order.MarkOrderPaidReq) (resp *order.MarkOrderPaidResp, err error) {
	resp, err = service.NewMarkOrderPaidService(ctx).Run(req)

	return resp, err
}

// MarkOrderCanceled implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) MarkOrderCanceled(ctx context.Context, req *order.MarkOrderCanceledReq) (resp *order.MarkOrderCanceledResp, err error) {
	resp, err = service.NewMarkOrderCanceledService(ctx).Run(req)

	return resp, err
}

// DeleteOrder implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) DeleteOrder(ctx context.Context, req *order.DeleteOrderReq) (resp *order.DeleteOrderResp, err error) {
	resp, err = service.NewDeleteOrderService(ctx).Run(req)
	return resp, err
}

// GetOrderStatus implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) GetOrderStatus(ctx context.Context, req *order.GetOrderStatusReq) (resp *order.GetOrderStatusResp, err error) {
	resp, err = service.NewGetOrderStatusService(ctx).Run(req)

	return resp, err
}

// GetOrderByID implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) GetOrderByID(ctx context.Context, req *order.GetOrderReq) (resp *order.GetOrderResp, err error) {
	resp, err = service.NewGetOrderByIDService(ctx).Run(req)

	return resp, err
}

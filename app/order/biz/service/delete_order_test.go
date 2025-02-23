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
	"testing"

	order "github.com/naskids/nas-mall/rpc_gen/kitex_gen/order"
)

func TestDeleteOrder_Run(t *testing.T) {
	ctx := context.Background()
	s := NewDeleteOrderService(ctx)

	tests := []struct {
		name    string
		req     *order.DeleteOrderReq
		wantErr bool
	}{
		{
			name: "invalid params - empty user_id",
			req: &order.DeleteOrderReq{
				UserId:  0,
				OrderId: "test-order-id",
			},
			wantErr: true,
		},
		{
			name: "invalid params - empty order_id",
			req: &order.DeleteOrderReq{
				UserId:  1,
				OrderId: "",
			},
			wantErr: true,
		},
		{
			name: "order not found",
			req: &order.DeleteOrderReq{
				UserId:  1,
				OrderId: "non-existent-order",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := s.Run(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteOrderService.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("resp: %v", resp)
		})
	}
} 
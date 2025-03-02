package service

import (
	"context"
	"testing"
	order "github.com/naskids/nas-mall/rpc_gen/kitex_gen/order"
)

func TestGetOrderByID_Run(t *testing.T) {
	ctx := context.Background()
	s := NewGetOrderByIDService(ctx)
	// init req and assert value

	req := &order.GetOrderReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}

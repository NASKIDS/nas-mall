package service

import (
	"context"
	"testing"
	order "github.com/naskids/nas-mall/rpc_gen/kitex_gen/order"
)

func TestScheduledOrderCancel_Run(t *testing.T) {
	ctx := context.Background()
	s := NewScheduledOrderCancelService(ctx)
	// init req and assert value

	req := &order.ScheduledOrderCancelReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}

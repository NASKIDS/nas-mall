package service

import (
	"context"
	"testing"

	payment "github.com/naskids/nas-mall/rpc_gen/kitex_gen/payment"
)

func TestCancelCharge_Run(t *testing.T) {
	ctx := context.Background()
	s := NewCancelChargeService(ctx)
	// init req and assert value

	req := &payment.CancelChargeReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}

package service

import (
	"context"
	"testing"

	payment "github.com/naskids/nas-mall/rpc_gen/kitex_gen/payment"
)

func TestCreatePaymentLog_Run(t *testing.T) {
	ctx := context.Background()
	s := NewCreatePaymentLogService(ctx)
	// init req and assert value

	req := &payment.CreatePaymentLogReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}

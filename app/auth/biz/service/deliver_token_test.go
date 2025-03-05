package service

import (
	"context"
	"testing"

	auth "github.com/naskids/nas-mall/rpc_gen/kitex_gen/auth"
)

func TestDeliverToken_Run(t *testing.T) {
	ctx := context.Background()
	s := NewDeliverTokenService(ctx)
	// init req and assert value

	req := &auth.DeliverTokenReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test
}

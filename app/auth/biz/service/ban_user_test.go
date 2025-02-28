package service

import (
	"context"
	"testing"

	auth "github.com/naskids/nas-mall/rpc_gen/kitex_gen/auth"
)

func TestBanUser_Run(t *testing.T) {
	ctx := context.Background()
	s := NewBanUserService(ctx)
	// init req and assert value

	req := &auth.BanUserReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}

package service

import (
	"context"
	auth "github.com/naskids/nas-mall/rpc_gen/kitex_gen/auth"
	"testing"
)

func TestRefreshToken_Run(t *testing.T) {
	ctx := context.Background()
	s := NewRefreshTokenService(ctx)
	// init req and assert value

	req := &auth.RefreshTokenReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}

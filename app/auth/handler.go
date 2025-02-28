package main

import (
	"context"

	"github.com/naskids/nas-mall/app/auth/biz/service"
	auth "github.com/naskids/nas-mall/rpc_gen/kitex_gen/auth"
)

// AuthServiceImpl implements the last service interface defined in the IDL.
type AuthServiceImpl struct{}

// DeliverToken implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) DeliverToken(ctx context.Context, req *auth.DeliverTokenReq) (resp *auth.DeliveryTokenResp, err error) {
	resp, err = service.NewDeliverTokenService(ctx).Run(req)

	return resp, err
}

// RefreshToken implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) RefreshToken(ctx context.Context, req *auth.RefreshTokenReq) (resp *auth.RefreshTokenResp, err error) {
	resp, err = service.NewRefreshTokenService(ctx).Run(req)

	return resp, err
}

// BanUser implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) BanUser(ctx context.Context, req *auth.BanUserReq) (resp *auth.BanUserResp, err error) {
	resp, err = service.NewBanUserService(ctx).Run(req)

	return resp, err
}

// VerifyTokenByRPC implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) VerifyTokenByRPC(ctx context.Context, req *auth.VerifyTokenReq) (resp *auth.VerifyTokenResp, err error) {
	resp, err = service.NewVerifyTokenByRPCService(ctx).Run(req)

	return resp, err
}

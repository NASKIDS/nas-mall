package service

import (
	"context"
	"fmt"

	"github.com/naskids/nas-mall/common/token"
	auth "github.com/naskids/nas-mall/rpc_gen/kitex_gen/auth"
)

type VerifyTokenByRPCService struct {
	ctx        context.Context
	tokenMaker *token.Maker
} // NewVerifyTokenByRPCService new VerifyTokenByRPCService
func NewVerifyTokenByRPCService(ctx context.Context) *VerifyTokenByRPCService {

	return &VerifyTokenByRPCService{ctx: ctx, tokenMaker: token.DefaultMaker()}
}

// Run create note info
func (s *VerifyTokenByRPCService) Run(req *auth.VerifyTokenReq) (resp *auth.VerifyTokenResp, err error) {
	claims, err := s.tokenMaker.ParseRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("remote failed to verify token: [%w]", err)
	}

	userId := claims["uid"].(uint64)
	role := claims["rol"].(string)

	return &auth.VerifyTokenResp{
		UserId: userId,
		Role:   role,
	}, nil
}

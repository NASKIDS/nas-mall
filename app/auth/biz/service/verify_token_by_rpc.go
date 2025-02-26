package service

import (
	"context"

	"github.com/naskids/nas-mall/app/auth/utils/token"
	auth "github.com/naskids/nas-mall/rpc_gen/kitex_gen/auth"
)

type VerifyTokenByRPCService struct {
	ctx        context.Context
	tokenMaker token.Maker
} // NewVerifyTokenByRPCService new VerifyTokenByRPCService
func NewVerifyTokenByRPCService(ctx context.Context) *VerifyTokenByRPCService {
	return &VerifyTokenByRPCService{ctx: ctx}
}

// Run create note info
func (s *VerifyTokenByRPCService) Run(req *auth.VerifyTokenReq) (resp *auth.VerifyResp, err error) {
	userID, err := s.tokenMaker.ParseAccessToken(req.AccessToken)
	if err != nil {
		return &auth.VerifyResp{Valid: false}, nil
	}
	return &auth.VerifyResp{
		Valid:  true,
		UserId: userID,
	}, nil
}

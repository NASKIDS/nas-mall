package auth

import (
	"context"
	"reflect"
	"testing"

	"github.com/bytedance/gopkg/cloud/metainfo"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/callopt"

	"github.com/naskids/nas-mall/rpc_gen/kitex_gen/auth"
)

func TestBanUser(t *testing.T) {
	type args struct {
		ctx         context.Context
		req         *auth.BanUserReq
		callOptions []callopt.Option
	}
	tests := []struct {
		name     string
		args     args
		wantResp *auth.BanUserResp
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := BanUser(tt.args.ctx, tt.args.req, tt.args.callOptions...)
			if (err != nil) != tt.wantErr {
				t.Errorf("BanUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("BanUser() gotResp = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}

func TestDeliverToken(t *testing.T) {
	InitClient("auth", client.WithHostPorts("192.168.31.223:8888"))
	type args struct {
		ctx         context.Context
		req         *auth.DeliverTokenReq
		callOptions []callopt.Option
	}
	tests := []struct {
		name     string
		args     args
		wantResp *auth.DeliveryTokenResp
		wantErr  bool
	}{
		{
			name: "test success",
			args: args{
				ctx: metainfo.WithValue(context.Background(), "access_token", "ak"),
				req: &auth.DeliverTokenReq{
					UserId: 1,
				},
			},
			wantResp: &auth.DeliveryTokenResp{
				AccessToken:  "",
				RefreshToken: "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := DeliverToken(tt.args.ctx, tt.args.req, tt.args.callOptions...)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeliverToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(gotResp)
		})
	}
}

func TestRefreshToken(t *testing.T) {
	type args struct {
		ctx         context.Context
		req         *auth.RefreshTokenReq
		callOptions []callopt.Option
	}
	tests := []struct {
		name     string
		args     args
		wantResp *auth.RefreshTokenResp
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := RefreshToken(tt.args.ctx, tt.args.req, tt.args.callOptions...)
			if (err != nil) != tt.wantErr {
				t.Errorf("RefreshToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("RefreshToken() gotResp = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}

func TestVerifyTokenByRPC(t *testing.T) {
	type args struct {
		ctx         context.Context
		req         *auth.VerifyTokenReq
		callOptions []callopt.Option
	}
	tests := []struct {
		name     string
		args     args
		wantResp *auth.VerifyTokenResp
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := VerifyTokenByRPC(tt.args.ctx, tt.args.req, tt.args.callOptions...)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyTokenByRPC() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("VerifyTokenByRPC() gotResp = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}

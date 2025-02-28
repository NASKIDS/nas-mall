package token

import (
	"time"

	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/samber/lo"

	"github.com/hertz-contrib/paseto"

	"github.com/naskids/nas-mall/app/auth/biz/model"
)

const (
	defaultSymmetricKey = "707172737475767778797a7b7c7d7e7f808182838485868788898a8b8c8d8e8f"
	DefaultPublicKey    = "1eb9dbbbbc047c03fd70604e0071f0987e16b28b757225c11f00415d0e20b1a2"
	defaultPrivateKey   = "b4cbfb43df4ce210727d953e4a713307fa19bb7d9f85041438d9e11b942a37741eb9dbbbbc047c03fd70604e0071f0987e16b28b757225c11f00415d0e20b1a2"
	defaultImplicit     = "2025-nasmall-implicit"
	issuer              = "nasmall-issuer"
)

type Maker struct {
	accessKeyDuration     time.Duration
	refreshKeyDuration    time.Duration
	genAccessTokenFunc    paseto.GenTokenFunc
	genRefreshTokenFunc   paseto.GenTokenFunc
	parseAccessTokenFunc  paseto.ParseFunc
	parseRefreshTokenFunc paseto.ParseFunc
	userStore             model.AuthUser
}

func NewMaker() *Maker {
	genAccessTokenFunc := lo.Must(paseto.NewV4SignFunc(defaultPrivateKey, []byte(defaultImplicit)))
	genRefreshTokenFunc := lo.Must(paseto.NewV4EncryptFunc(defaultSymmetricKey, []byte(defaultImplicit)))
	parseAccessTokenFunc := lo.Must(paseto.NewV4PublicParseFunc(DefaultPublicKey, []byte(paseto.DefaultImplicit), paseto.WithIssuer(issuer)))
	parseRefreshTokenFunc := lo.Must(paseto.NewV4LocalParseFunc(defaultSymmetricKey, []byte(paseto.DefaultImplicit), paseto.WithIssuer(issuer)))
	return &Maker{
		15 * time.Minute,
		24 * time.Hour,
		genAccessTokenFunc,
		genRefreshTokenFunc,
		parseAccessTokenFunc,
		parseRefreshTokenFunc,
		nil,
	}
}

func (m *Maker) GenerateAccessToken(customClaims utils.H) (string, error) {
	now := time.Now()
	token, err := m.genAccessTokenFunc(&paseto.StandardClaims{
		Issuer:    "nasmall-issuer",
		ExpiredAt: now.Add(m.accessKeyDuration),
		NotBefore: now,
		IssuedAt:  now,
	}, customClaims, nil)
	if err != nil {
		klog.Error("generate token failed")
	}
	return token, nil
}

func (m *Maker) GenerateRefreshToken(customClaims utils.H) (string, error) {
	now := time.Now()
	token, err := m.genAccessTokenFunc(&paseto.StandardClaims{
		Issuer:    issuer,
		ExpiredAt: now.Add(m.refreshKeyDuration),
		NotBefore: now,
		IssuedAt:  now,
	}, customClaims, nil)
	if err != nil {
		klog.Error("generate token failed")
	}
	return token, nil
}

func (m *Maker) ParseAccessToken(tokenStr string) (claims utils.H, err error) {
	token, err := m.parseAccessTokenFunc(tokenStr)
	if err != nil {
		return nil, err
	}
	return token.Claims(), nil
}

func (m *Maker) ParseRefreshToken(tokenStr string) (claims utils.H, err error) {
	token, err := m.parseRefreshTokenFunc(tokenStr)
	if err != nil {
		return
	}
	return token.Claims(), nil
}

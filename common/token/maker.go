package token

import (
	"os"
	"sync"
	"time"

	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/hertz-contrib/paseto"
	"github.com/samber/lo"
)

const (
	defaultSymmetricKey = "707172737475767778797a7b7c7d7e7f808182838485868788898a8b8c8d8e8f"
	defaultPublicKey    = "1eb9dbbbbc047c03fd70604e0071f0987e16b28b757225c11f00415d0e20b1a2"
	defaultPrivateKey   = "b4cbfb43df4ce210727d953e4a713307fa19bb7d9f85041438d9e11b942a37741eb9dbbbbc047c03fd70604e0071f0987e16b28b757225c11f00415d0e20b1a2"
	defaultImplicit     = "2025-nasmall-implicit"
	issuer              = "nasmall-issuer"
)

type TokenMaker struct {
	AccessKeyDuration     time.Duration
	RefreshKeyDuration    time.Duration
	genAccessTokenFunc    paseto.GenTokenFunc
	genRefreshTokenFunc   paseto.GenTokenFunc
	parseAccessTokenFunc  paseto.ParseFunc
	parseRefreshTokenFunc paseto.ParseFunc
}

var (
	once  sync.Once
	Maker *TokenMaker
)

func InitTokenMaker() *TokenMaker {
	once.Do(func() {
		publicKey := os.Getenv("PUBLIC_KEY")
		privateKey := os.Getenv("PRIVATE_KEY")
		symmetricKey := os.Getenv("SYMMETRIC_KEY")
		genAccessTokenFunc := lo.Must(paseto.NewV4SignFunc(privateKey, []byte(defaultImplicit)))
		genRefreshTokenFunc := lo.Must(paseto.NewV4EncryptFunc(symmetricKey, []byte(defaultImplicit)))
		parseAccessTokenFunc := lo.Must(paseto.NewV4PublicParseFunc(publicKey, []byte(defaultImplicit), paseto.WithIssuer(issuer)))
		parseRefreshTokenFunc := lo.Must(paseto.NewV4LocalParseFunc(symmetricKey, []byte(defaultImplicit), paseto.WithIssuer(issuer)))
		Maker = &TokenMaker{
			15 * time.Minute,
			24 * time.Hour,
			genAccessTokenFunc,
			genRefreshTokenFunc,
			parseAccessTokenFunc,
			parseRefreshTokenFunc,
		}
	})
	return Maker
}

func (m *TokenMaker) GenerateAccessToken(customClaims utils.H) (string, error) {
	now := time.Now()
	token, err := m.genAccessTokenFunc(&paseto.StandardClaims{
		Issuer:    issuer,
		ExpiredAt: now.Add(m.AccessKeyDuration),
		NotBefore: now,
		IssuedAt:  now,
	}, customClaims, nil)
	if err != nil {
		klog.Error("generate token failed")
	}
	return token, nil
}

func (m *TokenMaker) GenerateRefreshToken(customClaims utils.H) (string, error) {
	now := time.Now()
	token, err := m.genRefreshTokenFunc(&paseto.StandardClaims{
		Issuer:    issuer,
		ExpiredAt: now.Add(m.RefreshKeyDuration),
		NotBefore: now,
		IssuedAt:  now,
	}, customClaims, nil)
	if err != nil {
		klog.Error("generate token failed")
	}
	return token, nil
}

func (m *TokenMaker) ParseAccessToken(tokenStr string) (claims utils.H, err error) {
	token, err := m.parseAccessTokenFunc(tokenStr)
	if err != nil {
		return nil, err
	}
	return token.Claims(), nil
}

func (m *TokenMaker) ParseRefreshToken(tokenStr string) (claims utils.H, err error) {
	token, err := m.parseRefreshTokenFunc(tokenStr)
	if err != nil {
		return
	}
	return token.Claims(), nil
}

package token

import (
	"aidanwo
	"errors"
paseto"
	"
	"strconv"
	"golang.org/x/crypto/ed25519"

	"github.com/naskids/nas-mall/app/auth/biz/model"
)

type Maker struct {
	tokenDuration   time.Duration
	refreshDuration time.Duration
	userStore       model.AuthUser
}

func (m *Maker) GenerateAccessToken(userID uint64, expiry time.Duration) (string, error) {
	accessToken := paseto.NewToken()
	accessToken.SetIssuedAt(time.Now())
	accessToken.SetNotBefore(time.Now())
	accessToken.SetExpiration(time.Now().Add(2 * time.Hour))
	accessToken.SetString()
	accessToken.SetAudience()
	key := paseto.NewV4SymmetricKey() // don't share this!!

	encrypted := accessToken.V4Encrypt(key, nil)
	token := paseto.JSONToken{
		Expiration: time.Now().Add(expiry),
	}
	token.Set("user_id", strconv.FormatUint(userID, 10))
	return m.paseto.Sign(m.accessPrivateKey, token, nil)
}

func (m *Maker) GenerateRefreshToken(userID, version uint64, expiry time.Duration) (string, error) {
	token := paseto.JSONToken{
		Expiration: time.Now().Add(expiry),
	}
	token.Set("user_id", strconv.FormatUint(userID, 10))
	token.Set("refresh_version", strconv.FormatUint(version, 10))
	return m.paseto.Encrypt(m.refreshKey, token, nil)
}

func (m *Maker) ParseAccessToken(token string) (userID uint64, err error) {
	var parsedToken paseto.JSONToken
	err = m.paseto.Verify(token, m.accessPublicKey, &parsedToken, nil)
	if err != nil {
		return 0, err
	}
	userID, err = strconv.ParseUint(parsedToken.Get("user_id"), 10, 64)
	return
}

func (m *Maker) ParseRefreshToken(token string) (userID, version uint64, err error) {
	var parsedToken paseto.JSONToken
	err = m.paseto.Decrypt(token, m.refreshKey, &parsedToken, nil)
	if err != nil {
		return 0, 0, err
	}
	userID, err = strconv.ParseUint(parsedToken.Get("user_id"), 10, 64)
	var err2 error
	version, err2 = strconv.ParseUint(parsedToken.Get("refresh_version"), 10, 64)
	err = errors.Join(err2, err)
	return
}

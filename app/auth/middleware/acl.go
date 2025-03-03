package middleware

import (
	"context"
	"errors"
	"math/rand"

	"github.com/cloudwego/kitex/pkg/acl"
)

var errRejected = errors.New("casbin rejected")

func reject(ctx context.Context, request interface{}) (reason error) {
	if rand.Intn(100) == 0 {
		return errRejected // 拒绝请求时，需要返回一个错误
	}
	return nil
}

var ACL = acl.NewACLMiddleware([]acl.RejectFunc{reject})

package middleware

import (
	"context"
	"errors"
	"fmt"

	"github.com/bytedance/gopkg/cloud/metainfo"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/redis-watcher/v2"
	"github.com/cloudwego/kitex/pkg/acl"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/utils/kitexutil"
	"github.com/redis/go-redis/v9"

	"github.com/naskids/nas-mall/app/auth/conf"
	"github.com/naskids/nas-mall/common/token"
)

var (
	errRejected = errors.New("casbin rejected")
	errNoMethod = errors.New("no method")
	enforcer    *casbin.Enforcer
)

func reject(ctx context.Context, request interface{}) (reason error) {
	accessToken, _ := metainfo.GetValue(ctx, "access_token")
	var role string
	if accessToken != "" {
		claims, err := token.VerifyAccessToken(accessToken)
		if err != nil {
			// access token 无效
			return err
		}
		role = claims["rol"].(string)
	} else {
		role = "anonymous"
	}

	method, ok1 := kitexutil.GetMethod(ctx)
	svcName, ok2 := kitexutil.GetIDLServiceName(ctx)
	if !ok1 || !ok2 {
		return errNoMethod
	}

	obj := fmt.Sprintf("%s/%s", svcName, method)
	if enforce, err := enforcer.Enforce(role, obj, "CALL"); err != nil || !enforce {
		return errRejected
	}

	return nil
}

var ACL = acl.NewACLMiddleware([]acl.RejectFunc{reject})

func InitACL() {
	w, _ := rediswatcher.NewWatcher(conf.GetConf().Redis.Address, rediswatcher.WatcherOptions{
		Options: redis.Options{
			Network:  "tcp",
			Password: "",
		},
		Channel: "/casbin",
	})

	m := model.NewModel()
	m.AddDef("r", "r", "sub, obj, act")
	m.AddDef("p", "p", "sub, obj, act")
	m.AddDef("g", "g", "_, _")
	m.AddDef("e", "e", "some(where (p.eft == allow))")
	m.AddDef("m", "m", "g(r.sub, p.sub) && keyMatch(r.obj, p.obj) && regexMatch(r.act, p.act)")

	enforcer, _ = casbin.NewEnforcer(m)

	err := enforcer.SetWatcher(w)
	if err != nil {
		klog.Fatalf("init casbin enforcer failed: %s", err)
	}
	err = w.SetUpdateCallback(func(s string) {
		klog.Info(s)
		rediswatcher.DefaultUpdateCallback(enforcer)(s)
	})
	if err != nil {
		klog.Fatalf("init enforcer failed : %s", err)
	}
	ok, err := enforcer.AddPolicy("anonymous", "AuthService/*", "CALL")
	if !ok || err != nil {
		klog.Fatalf("init policy failed : %s", err)
	}
}

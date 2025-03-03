package main

import (
	"net"
	"strings"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/joho/godotenv"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/naskids/nas-mall/app/auth/biz/dal"
	"github.com/naskids/nas-mall/app/auth/conf"
	"github.com/naskids/nas-mall/app/auth/middleware"
	"github.com/naskids/nas-mall/common/mtl"
	"github.com/naskids/nas-mall/common/serversuite"
	"github.com/naskids/nas-mall/common/token"
	"github.com/naskids/nas-mall/common/utils"
	"github.com/naskids/nas-mall/rpc_gen/kitex_gen/auth/authservice"
)

var serviceName = conf.GetConf().Kitex.Service

func main() {
	_ = godotenv.Load()
	mtl.InitLog(&lumberjack.Logger{
		Filename:   conf.GetConf().Kitex.LogFileName,
		MaxSize:    conf.GetConf().Kitex.LogMaxSize,
		MaxBackups: conf.GetConf().Kitex.LogMaxBackups,
		MaxAge:     conf.GetConf().Kitex.LogMaxAge,
	})
	mtl.InitTracing(serviceName)
	mtl.InitMetric(serviceName, conf.GetConf().Kitex.MetricsPort, conf.GetConf().Registry.RegistryAddress[0])
	dal.Init()
	token.InitTokenMaker()
	middleware.InitACL()

	_ = middleware.E.SavePolicy()
	opts := kitexInit()

	svr := authservice.NewServer(new(AuthServiceImpl), opts...)

	err := svr.Run()
	if err != nil {
		klog.Error(err.Error())
	}
}

func kitexInit() (opts []server.Option) {
	// address
	address := conf.GetConf().Kitex.Address
	if strings.HasPrefix(address, ":") {
		localIp := utils.MustGetLocalIPv4()
		address = localIp + address
	}
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		panic(err)
	}
	opts = append(opts,
		server.WithServiceAddr(addr),
		server.WithMiddleware(middleware.ACL),
		server.WithSuite(serversuite.CommonServerSuite{CurrentServiceName: serviceName, RegistryAddr: conf.GetConf().Registry.RegistryAddress[0]}))

	// service info
	opts = append(opts, server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
		ServiceName: conf.GetConf().Kitex.Service,
	}))
	return
}

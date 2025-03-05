// Copyright 2024 CloudWeGo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rpc

import (
	"sync"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/transport"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"

	"github.com/naskids/nas-mall/app/payment/conf"
	"github.com/naskids/nas-mall/rpc_gen/kitex_gen/order/orderservice"
	rpcorder "github.com/naskids/nas-mall/rpc_gen/rpc/order"
)

var (
	OrderClient orderservice.Client
	once        sync.Once
	serviceName string
)

func InitClient() {
	once.Do(func() {
		serviceName = conf.GetConf().Kitex.Service
		initOrderClient()
	})
}

type OrderClientSuite struct {
	CurrentServiceName string
}

func (s OrderClientSuite) Options() []client.Option {
	opts := []client.Option{
		client.WithHostPorts("10.1.2.133:8886"),
		client.WithMetaHandler(transmeta.ClientHTTP2Handler),
		client.WithTransportProtocol(transport.GRPC),
	}

	opts = append(opts,
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: s.CurrentServiceName,
		}),
		client.WithSuite(tracing.NewClientSuite()),
	)

	return opts
}
func initOrderClient() {
	opts := []client.Option{
		client.WithSuite(OrderClientSuite{
			CurrentServiceName: serviceName,
		}),
	}
	rpcorder.InitClient("order", opts...)
	OrderClient = rpcorder.DefaultClient()
}

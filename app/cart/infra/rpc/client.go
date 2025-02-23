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

	"github.com/naskids/nas-mall/common/clientsuite"
	rpcproduct "github.com/naskids/nas-mall/rpc_gen/rpc/product"

	"github.com/naskids/nas-mall/app/cart/conf"
	"github.com/naskids/nas-mall/rpc_gen/kitex_gen/product/productcatalogservice"
)

var (
	ProductClient productcatalogservice.Client
	once          sync.Once
	serviceName   string
)

func InitClient() {
	once.Do(func() {
		serviceName = conf.GetConf().Kitex.Service
		initProductClient()
	})
}

func initProductClient() {
	opts := []client.Option{
		client.WithSuite(clientsuite.CommonGrpcClientSuite{
			CurrentServiceName: serviceName,
		}),
	}

	rpcproduct.InitClient("product", opts...)
	ProductClient = rpcproduct.DefaultClient()
}

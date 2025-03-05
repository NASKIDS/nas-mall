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

	"github.com/naskids/nas-mall/app/checkout/conf"
	"github.com/naskids/nas-mall/rpc_gen/kitex_gen/cart/cartservice"
	"github.com/naskids/nas-mall/rpc_gen/kitex_gen/order/orderservice"
	"github.com/naskids/nas-mall/rpc_gen/kitex_gen/payment/paymentservice"
	"github.com/naskids/nas-mall/rpc_gen/kitex_gen/product/productcatalogservice"
	rpccart "github.com/naskids/nas-mall/rpc_gen/rpc/cart"
	rpcorder "github.com/naskids/nas-mall/rpc_gen/rpc/order"
	rpcpayment "github.com/naskids/nas-mall/rpc_gen/rpc/payment"
	rpcproduct "github.com/naskids/nas-mall/rpc_gen/rpc/product"
)

var (
	CartClient    cartservice.Client
	ProductClient productcatalogservice.Client
	PaymentClient paymentservice.Client
	OrderClient   orderservice.Client
	once          sync.Once
	serviceName   string
	commonSuite   client.Option
)

func InitClient() {
	once.Do(func() {
		serviceName = conf.GetConf().Kitex.Service
		// commonSuite = client.WithSuite(clientsuite.CommonGrpcClientSuite{
		// 	CurrentServiceName: serviceName,
		// })
		initCartClient()
		initProductClient()
		initPaymentClient()
		initOrderClient()
	})
}

func initProductClient() {
	// rpcproduct.InitClient("product", commonSuite)
	ProductClient = rpcproduct.DefaultClient()
}

func initCartClient() {
	// rpccart.InitClient("cart", commonSuite)
	CartClient = rpccart.DefaultClient()
}

func initPaymentClient() {
	// rpcpayment.InitClient("payment", commonSuite)
	PaymentClient = rpcpayment.DefaultClient()
}

func initOrderClient() {
	// rpcorder.InitClient("order", commonSuite)
	OrderClient = rpcorder.DefaultClient()
}

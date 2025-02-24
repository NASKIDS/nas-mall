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
	"context"
	"sync"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/circuitbreak"
	"github.com/cloudwego/kitex/pkg/fallback"
	"github.com/cloudwego/kitex/pkg/rpcinfo"

	"github.com/naskids/nas-mall/app/frontend/infra/mtl"
	"github.com/naskids/nas-mall/rpc_gen/kitex_gen/cart/cartservice"
	"github.com/naskids/nas-mall/rpc_gen/kitex_gen/checkout/checkoutservice"
	"github.com/naskids/nas-mall/rpc_gen/kitex_gen/order/orderservice"
	"github.com/naskids/nas-mall/rpc_gen/kitex_gen/product"
	"github.com/naskids/nas-mall/rpc_gen/kitex_gen/product/productcatalogservice"
	"github.com/naskids/nas-mall/rpc_gen/kitex_gen/user/userservice"
	rpccart "github.com/naskids/nas-mall/rpc_gen/rpc/cart"
	rpccheckout "github.com/naskids/nas-mall/rpc_gen/rpc/checkout"
	rpcorder "github.com/naskids/nas-mall/rpc_gen/rpc/order"
	rpcproduct "github.com/naskids/nas-mall/rpc_gen/rpc/product"
	rpcuser "github.com/naskids/nas-mall/rpc_gen/rpc/user"

	prometheus "github.com/kitex-contrib/monitor-prometheus"
)

var (
	ProductClient  productcatalogservice.Client
	UserClient     userservice.Client
	CartClient     cartservice.Client
	CheckoutClient checkoutservice.Client
	OrderClient    orderservice.Client
	once           sync.Once
	commonSuite    client.Option
)

func InitClient() {
	once.Do(func() {
		initProductClient()
		initUserClient()
		initCartClient()
		initCheckoutClient()
		initOrderClient()
	})
}

func initProductClient() {
	var opts []client.Option

	cbs := circuitbreak.NewCBSuite(func(ri rpcinfo.RPCInfo) string {
		return circuitbreak.RPCInfo2Key(ri)
	})
	cbs.UpdateServiceCBConfig("shop-frontend/product/GetProduct", circuitbreak.CBConfig{Enable: true, ErrRate: 0.5, MinSample: 2})

	opts = append(
		opts,
		client.WithCircuitBreaker(cbs),
		client.WithFallback(fallback.NewFallbackPolicy(fallback.UnwrapHelper(func(ctx context.Context, req, resp interface{}, err error) (fbResp interface{}, fbErr error) {
			methodName := rpcinfo.GetRPCInfo(ctx).To().Method()
			if err == nil {
				return resp, err
			}
			if methodName != "ListProducts" {
				return resp, err
			}
			return &product.ListProductsResp{
				Products: []*product.Product{
					{
						Price:       6.6,
						Id:          3,
						Picture:     "/static/image/t-shirt.jpeg",
						Name:        "T-Shirt",
						Description: "CloudWeGo T-Shirt",
					},
				},
			}, nil
		}))),
	)
	opts = append(opts, client.WithTracer(prometheus.NewClientTracer("", "", prometheus.WithDisableServer(true), prometheus.WithRegistry(mtl.Registry))))
	rpcproduct.InitClient("product", opts...)
	ProductClient = rpcproduct.DefaultClient()
}

func initUserClient() {
	rpcuser.InitClient("user", commonSuite)
	UserClient = rpcuser.DefaultClient()
}

func initCartClient() {
	rpccart.InitClient("cart", commonSuite)
	CartClient = rpccart.DefaultClient()
}

func initCheckoutClient() {
	rpccheckout.InitClient("checkout", commonSuite)
	CheckoutClient = rpccheckout.DefaultClient()
}

func initOrderClient() {
	rpcorder.InitClient("order", commonSuite)
	OrderClient = rpcorder.DefaultClient()
}

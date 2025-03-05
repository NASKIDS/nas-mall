package service

import (
	"context"
	"fmt"

	"github.com/cloudwego/kitex/pkg/klog"

	"github.com/naskids/nas-mall/app/order/biz/dal/mysql"
	"github.com/naskids/nas-mall/app/order/biz/model"
	"github.com/naskids/nas-mall/rpc_gen/kitex_gen/cart"
	order "github.com/naskids/nas-mall/rpc_gen/kitex_gen/order"
)

type GetOrderByIDService struct {
	ctx context.Context
}

// NewGetOrderByIDService 创建GetOrderByIDService实例
func NewGetOrderByIDService(ctx context.Context) *GetOrderByIDService {
	return &GetOrderByIDService{ctx: ctx}
}

// Run 根据userId和orderId查询订单
func (s *GetOrderByIDService) Run(req *order.GetOrderReq) (resp *order.GetOrderResp, err error) {
	// 参数验证
	if req.UserId == 0 || req.OrderId == "" {
		err = fmt.Errorf("user_id或order_id不能为空")
		return nil, err
	}

	// 从数据库查询订单
	orderData, err := model.GetOrder(mysql.DB, s.ctx, req.UserId, req.OrderId)
	if err != nil {
		klog.Errorf("model.GetOrder error: %v", err)
		return nil, err
	}

	// 构建响应
	var orderItems []*order.OrderItem
	for _, item := range orderData.OrderItems {
		orderItems = append(orderItems, &order.OrderItem{
			Cost: item.Cost,
			Item: &cart.CartItem{
				ProductId: item.ProductId,
				Quantity:  item.Quantity,
			},
		})
	}
	// 封装订单信息
	orderResp := &order.Order{
		OrderId:      orderData.OrderId,
		UserId:       orderData.UserId,
		UserCurrency: orderData.UserCurrency,
		Email:        orderData.Consignee.Email,
		CreatedAt:    int32(orderData.CreatedAt.Unix()),
		Address: &order.Address{
			Country:       orderData.Consignee.Country,
			City:          orderData.Consignee.City,
			StreetAddress: orderData.Consignee.StreetAddress,
			ZipCode:       orderData.Consignee.ZipCode,
		},
		OrderItems: orderItems,
	}

	resp = &order.GetOrderResp{
		Order: orderResp,
	}

	return resp, nil
}

package service

import (
	"context"

	"github.com/cloudwego/kitex/pkg/klog"

	"github.com/naskids/nas-mall/app/ai/biz/dal/mysql"
	"github.com/naskids/nas-mall/app/ai/biz/graph/rag_sql"
	"github.com/naskids/nas-mall/app/ai/biz/model"
	ai "github.com/naskids/nas-mall/rpc_gen/kitex_gen/ai"
	"github.com/naskids/nas-mall/rpc_gen/kitex_gen/cart"
	order "github.com/naskids/nas-mall/rpc_gen/kitex_gen/order"
)

type QueryOrderService struct {
	ctx context.Context
} // NewQueryOrderService new QueryOrderService
func NewQueryOrderService(ctx context.Context) *QueryOrderService {
	return &QueryOrderService{ctx: ctx}
}

// Run create note info
func (s *QueryOrderService) Run(req *ai.QueryOrderRequest) (resp *order.ListOrderResp, err error) {
	sql, err := rag_sql.Text2SQL.Invoke(s.ctx, &req.UserMessage)
	if err != nil {
		klog.Error(err)
	}

	orders, err := model.GetOrderFromRawSQL(s.ctx, mysql.DB, sql)
	if err != nil {
		klog.Errorf("model.ListOrder.err:%v", err)
		return nil, err
	}
	var list []*order.Order
	for _, v := range orders {
		var items []*order.OrderItem
		for _, v := range v.OrderItems {
			items = append(items, &order.OrderItem{
				Cost: v.Cost,
				Item: &cart.CartItem{
					ProductId: v.ProductId,
					Quantity:  v.Quantity,
				},
			})
		}
		o := &order.Order{
			OrderId:      v.OrderId,
			UserId:       v.UserId,
			UserCurrency: v.UserCurrency,
			Email:        v.Consignee.Email,
			CreatedAt:    int32(v.CreatedAt.Unix()),
			Address: &order.Address{
				Country:       v.Consignee.Country,
				City:          v.Consignee.City,
				StreetAddress: v.Consignee.StreetAddress,
				ZipCode:       v.Consignee.ZipCode,
			},
			OrderItems: items,
		}
		list = append(list, o)
	}
	resp = &order.ListOrderResp{
		Orders: list,
	}
	return
}

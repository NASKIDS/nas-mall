func checkout(Req req){
	从请求中获取cartId
	查询得到购物车具体内容
	计算购物车总价
	生成orderId
	插入一条新的订单并将其状态设置为placed
	往消息队里中发送一条消息，传入orderID和超时时间
	生成交易ID
	调用支付函数
	if(支付失败){
		return
	}
	if(成功抢到order_redis_key的锁){
		if(订单状态不是placed){
			取消支付结果
			return
		}
			插入一条支付log
			将订单标记为已支付
			清空购物车
		
	}
	
}
func cancelOrderMessageCal(){
	if(未超时间){
		放回消息队列
	}
	if(成功抢到order_redis_key的锁){
		if(订单状态是已支付){
			return
		}
		将订单标记为已取消
	}
}

func cancelPayMessageCal(){
	if(未超时间){
		放回消息队列
	}
	if(成功抢到order_redis_key的锁){
		if(订单状态是已取消){
			return
		}
		将订单标记为已取消
	}
}
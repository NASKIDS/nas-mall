package dal

import (
	"github.com/naskids/nas-mall/app/ai/biz/dal/mysql"
	"github.com/naskids/nas-mall/app/ai/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}

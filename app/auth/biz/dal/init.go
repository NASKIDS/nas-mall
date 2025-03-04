package dal

import (
	"github.com/naskids/nas-mall/app/auth/biz/dal/mysql"
	"github.com/naskids/nas-mall/app/auth/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}

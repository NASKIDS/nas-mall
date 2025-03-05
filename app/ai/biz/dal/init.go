package dal

import (
	"github.com/naskids/nas-mall/app/ai/biz/dal/mysql"
)

func Init() {
	// redis.Init()
	mysql.Init()
}

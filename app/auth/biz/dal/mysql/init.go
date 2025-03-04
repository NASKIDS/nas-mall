package mysql

import (
	"fmt"
	"os"

	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/plugin/opentelemetry/tracing"

	"github.com/naskids/nas-mall/app/auth/biz/model"
	"github.com/naskids/nas-mall/app/auth/conf"
	"github.com/naskids/nas-mall/common/mtl"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error
)

func Init() {
	dsn := fmt.Sprintf(conf.GetConf().MySQL.DSN, os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_HOST"))
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		panic(err)
	}
	if os.Getenv("GO_ENV") != "online" {
		needDemoData := !DB.Migrator().HasTable(&model.AuthUser{})
		if needDemoData {
			err = DB.AutoMigrate(&model.AuthUser{})
			if err != nil {
				klog.Warnf("auto migrate tables failed, err:%v", err)
			}

			DB.Create(&model.AuthUser{UserID: 1, Role: "admin", RefreshVersion: 2})
			DB.Create(&model.AuthUser{UserID: 2, Role: "user", RefreshVersion: 1239})
		}
	}
	if err = DB.Use(tracing.NewPlugin(tracing.WithoutMetrics(), tracing.WithTracerProvider(mtl.TracerProvider))); err != nil {
		panic(err)
	}
}

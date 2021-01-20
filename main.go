package main

import (
	"fmt"
	"myblog/dao"
	"myblog/models"
	"myblog/routers"
)

func main() {
	err := dao.InitMysql()
	if err != nil {
		fmt.Println("init mysql failed!", err)
		return
	}
	defer dao.Clsoe()
	dao.DB.AutoMigrate(&models.User{})
	dao.DB.AutoMigrate(&models.Record{})
	dao.DB.AutoMigrate(&models.Comment{})

	err = dao.InitRedis()
	if err != nil {
		fmt.Println("init redis failed!", err)
		return
	}
	defer dao.RClose()

	r := routers.SetupRouter()
	r.Run()
}

package main

import (
	_ "gitee.com/King_of_Universe_Sailing/MircoCenter/ServiceCenter/apps"
	"gitee.com/King_of_Universe_Sailing/MircoCenter/ServiceCenter/ioc"
	"gitee.com/King_of_Universe_Sailing/MircoCenter/ServiceCenter/middlewares/http"
	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()
	group := engine.Group("/wcenter/ServiceCenter/api/v1")
	group.Use(http.DefaultCors())
	err := ioc.RegisterGinApi(group)
	if err != nil {
		panic(err)
	}
	err = ioc.Init()
	if err != nil {
		panic(err)
	}
	err = engine.Run(":8080")
	if err != nil {
		panic(err)
	}
}

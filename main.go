package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"main/controller"
	"main/db"
	"main/middle"
	"main/service"
)

func main() {

	//获取配置文件中的mysql地址
	gconf, err := service.Conf.ReadConfFunc()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//初始化mysql
	db.InitMysql(gconf)

	//创建路由引擎
	router := gin.Default()
	//加载中间件
	router.Use(middle.Cors())    //跨域
	router.Use(middle.JWTAuth()) //加载jwt中间件，用于token验证
	//初始化路由
	controller.Router.RouterInit(router)

	//携程开启定时器
	go func() {
		service.ClusTers.CronFunc()
	}()

	//启动websocket
	//go func() {
	//	http.HandleFunc("/ws", service.Terminal.WsHandler)
	//	http.ListenAndServe(":8083", nil)
	//	fmt.Println("ws服务已启动。。。")
	//}()

	//_, _ = service.Imfile.ImportFile()
	router.Run("0.0.0.0:8999")
	//关闭GORM
	defer db.Close()
	//关闭cron
	defer service.ClusTers.CloseCron()
}

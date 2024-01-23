package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"main/service"
)

var Keepalive keepalive

type keepalive struct {
}

func (k *keepalive) KeepaliveFunc(c *gin.Context) {
	//去数据库更新状态
	name := c.Query("clusterName")
	fmt.Println("name= ", name)
	err := service.Keep.UpdateStatus(name, "active")
	if err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.Status(200)
}

package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"main/middle"
	"main/service"
)

var Namespace namespace

type namespace struct {
}

// 获取namespace列表
func (n *namespace) GetNamespaces(c *gin.Context) {
	//获取token
	token := middle.Jwt.GetToken()
	fmt.Println("准备携带的token为：", token)
	clusterName := c.Query("clusterName")
	fmt.Println("集群名：", clusterName)
	data, err := service.Namespace.GetNamespaces(token, clusterName)
	if err != nil {
		if err.Error() == "err" {
			c.JSON(400, gin.H{
				"err": data,
			})
		} else {
			c.JSON(400, gin.H{
				"err": err.Error(),
			})
		}
		return
	}
	c.JSON(200, data)
}

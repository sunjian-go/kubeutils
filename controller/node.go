package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"main/middle"
	"main/service"
)

var Node node

type node struct {
}

// 获取namespace列表
func (n *node) GetNodes(c *gin.Context) {
	//获取token
	token := middle.Jwt.GetToken()
	fmt.Println("准备携带的token为：", token)
	clusterName := c.Query("clusterName")
	data, err := service.Node.GetNodes(token, clusterName)
	if err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.JSON(200, data)
}

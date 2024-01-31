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

// 获取node列表
func (n *node) GetNodes(c *gin.Context) {
	//获取token
	token := middle.Jwt.GetToken()
	fmt.Println("准备携带的token为：", token)
	clusterName := c.Query("clusterName")
	node := new(service.NodeInfo)
	if err := c.Bind(node); err != nil {
		c.JSON(400, gin.H{
			"err":  "绑定node参数失败" + err.Error(),
			"data": nil,
		})
		return
	}
	fmt.Println("客户端传过来的为：", *node)
	data, err := service.Node.GetNodes(token, clusterName, node)
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

// 获取node详情
func (n *node) GetNodeDetail(c *gin.Context) {
	//获取token
	token := middle.Jwt.GetToken()
	fmt.Println("准备携带的token为：", token)
	clusterName := c.Query("clusterName")
	nodeName := c.Query("nodeName")
	fmt.Println("客户端传过来的为：", nodeName)
	data, err := service.Node.GetNodeDetail(token, clusterName, nodeName)
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

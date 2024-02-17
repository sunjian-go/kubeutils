package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"main/middle"
	"main/service"
)

var Listpath listpath

type listpath struct {
}

func (l *listpath) ListContainerPath(c *gin.Context) {
	clusterName := c.Query("clusterName")
	podinfo := new(service.PodPath)
	if err := c.Bind(podinfo); err != nil {
		c.JSON(400, gin.H{
			"err": "绑定数据失败" + err.Error(),
		})
		return
	}
	fmt.Println("pod信息：", podinfo)

	//获取token
	token := middle.Jwt.GetToken()
	fmt.Println("准备携带的token为：", token)

	out, err := service.Listpath.ListContainerPath(podinfo, token, clusterName)
	if err != nil {
		if err.Error() == "err" {
			c.JSON(400, gin.H{
				"err": out,
			})
		} else {
			c.JSON(400, gin.H{
				"err": err.Error(),
			})
		}
		return
	}
	c.JSON(200, gin.H{
		"msg":  "列出容器路径成功",
		"data": out,
	})
}

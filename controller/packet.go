package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"main/service"
)

var Pack packet

type packet struct {
}

// 启动抓包进程
func (p *packet) StartPacket(c *gin.Context) {
	url := c.Query("url")
	clusterName := c.Query("clusterName")
	packinfo := new(service.PackInfo)
	if err := c.ShouldBind(packinfo); err != nil {
		c.JSON(400, gin.H{
			"err": "绑定数据失败",
		})
		return
	}
	fmt.Println("需要抓包的数据为：", packinfo)
	info, err := service.Pack.StartPacket(packinfo, clusterName, url)
	if err != nil {
		if err.Error() == "err" {
			c.JSON(400, gin.H{
				"err": info,
			})
		} else {
			c.JSON(400, gin.H{
				"err": err.Error(),
			})
		}
		return
	}
	c.JSON(200, gin.H{
		"msg": info,
	})
}

func (p *packet) StopPacket(c *gin.Context) {
	url := c.Query("url")
	clusterName := c.Query("clusterName")
	err := service.Pack.StopPacket(c, clusterName, url)
	if err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
}

// 获取所有网卡信息
func (p *packet) GetAllInterface(c *gin.Context) {
	url := c.Query("url")
	clusterName := c.Query("clusterName")
	interfaces, err := service.Pack.GetAllInterface(c, clusterName, url)
	if err != nil {
		if err.Error() == "err" {
			c.JSON(400, gin.H{
				"err": interfaces,
			})
		} else {
			c.JSON(400, gin.H{
				"err": err.Error(),
			})
		}
	}
	c.JSON(200, interfaces)

}

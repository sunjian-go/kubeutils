package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"main/dao"
)

var Ipaddr ipaddr

type ipaddr struct {
}

func (i *ipaddr) GetClusterIP(c *gin.Context) {
	cname := c.Query("cluster_name")
	ip, err := dao.RegCluster.GetClusterIP(cname)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"ip": ip,
	})
}

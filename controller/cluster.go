package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"main/model"
	"main/service"
)

var Clus clus

type clus struct {
}

// 集群注册
func (cl *clus) RegisterFunc(ctx *gin.Context) {
	clu := new(model.Cluster)
	if err := ctx.ShouldBindJSON(clu); err != nil {
		ctx.JSON(400, gin.H{
			"err": "数据绑定失败：" + err.Error(),
		})
		return
	}
	fmt.Println("获取到：", clu)
	if err := service.ClusTers.RegisterFunc(clu); err != nil {
		fmt.Println(err.Error())
		ctx.JSON(400, gin.H{
			"err": "注册失败" + err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"msg": "注册成功",
	})
}

// 获取所有集群
func (cl *clus) GetAllClusters(c *gin.Context) {
	cluInfo := new(service.ClusterInfo)
	if err := c.Bind(cluInfo); err != nil {
		c.JSON(400, gin.H{
			"err": "绑定数据失败：" + err.Error(),
		})
		return
	}
	clus, total, err := service.ClusTers.GetAllClusters(cluInfo)
	if err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	//fmt.Println("获取所有集群：", clus)
	c.JSON(200, gin.H{
		"msg":   "获取所有集群成功",
		"data":  clus,
		"total": total,
	})
}

// 删除目标集群
func (cl *clus) DeleteCluster(c *gin.Context) {
	cluName := c.Query("clusterName")
	err := service.ClusTers.DeleteCluster(cluName)
	if err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"msg": "删除目标集群成功",
	})
}

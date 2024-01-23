package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"main/middle"
	"main/service"
)

var Pod pod

type pod struct {
}

// 此结构体用于内部，用来绑定客户端传过来的pod信息

func (p *pod) GetPods(c *gin.Context) {
	clusterName := c.Query("clusterName")
	pod := new(service.PodInfo)
	if err := c.Bind(pod); err != nil {
		c.JSON(400, gin.H{
			"err":  "绑定pod参数失败" + err.Error(),
			"data": nil,
		})
		return
	}
	fmt.Println("客户端传过来的为：", *pod)

	//获取token
	token := middle.Jwt.GetToken()
	fmt.Println("准备携带的token为：", token)

	data, err := service.Pod.GetObjs(token, clusterName, "getpods", pod)
	if err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.JSON(200, data)

}

// 获取容器信息
func (p *pod) GetContainer(c *gin.Context) {
	clusterName := c.Query("clusterName")
	pod := new(service.PodDetail)
	if err := c.Bind(pod); err != nil {
		c.JSON(400, gin.H{
			"err":  "绑定数据失败" + err.Error(),
			"data": nil,
		})
		return
	}
	//获取token
	token := middle.Jwt.GetToken()
	fmt.Println("准备携带的token为：", token)

	containers, err := service.Pod.GetObjs(token, clusterName, "getContainer", pod)
	if err != nil {
		c.JSON(400, gin.H{
			"err":  "获取容器信息失败" + err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "获取容器信息成功",
		"data": containers,
	})
}

func (p *pod) WsFunc(c *gin.Context) {
	namespace := c.Query("namespace")
	podName := c.Query("pod_name")
	containerName := c.Query("container_name")
	bashType := c.Query("bashType")
	clusterName := c.Query("clusterName")
	err := service.Terminal.WsHandler(namespace, podName, containerName, bashType, clusterName, c)
	if err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"msg": "建立ws连接成功",
	})
}

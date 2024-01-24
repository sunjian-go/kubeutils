package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"main/service"
)

var Login login

type login struct {
}

func (l *login) Login(c *gin.Context) {
	user := new(service.User)
	if err := c.ShouldBind(user); err != nil {
		c.JSON(400, gin.H{
			"err": "数据绑定失败：" + err.Error(),
		})
		return
	}
	fmt.Println("客户端登录：", user.Username, user.Password)

	//返回token和配置文件中的webSocket地址
	token, kubeconf, err := service.Login.Login(user)
	if err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	fmt.Println("获取ws主机地址：", kubeconf["host"])
	c.JSON(200, gin.H{
		"msg":   "登录成功",
		"token": token,
		"host":  kubeconf["host"],
		"port":  kubeconf["port"],
	})
}

// 获取验证公式
func (l *login) GetAuthCode(c *gin.Context) {
	//获取验证公式
	formula := service.Auth.GetAuthCode()
	c.JSON(200, gin.H{
		"msg":  "获取验证码成功",
		"data": formula,
	})
}

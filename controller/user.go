package controller

import (
	"github.com/gin-gonic/gin"
	"main/model"
	"main/service"
)

var Users user

type user struct {
}

func (u *user) GetUsers(c *gin.Context) {
	user := new(service.UserInfo)
	if err := c.Bind(user); err != nil {
		c.JSON(400, gin.H{

			"err": "绑定数据失败：" + err.Error(),
		})
		return
	}
	userItem, err := service.Users.GetFilterUsers(user)
	if err != nil {
		c.JSON(400, gin.H{
			"err": "获取用户列表失败",
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "获取用户列表成功",
		"data": userItem,
	})
}

// 查询单条用户信息
func (u *user) GetUser(c *gin.Context) {
	user := c.Query("user")
	userItem, err := service.Users.GetUser(user)
	if err != nil {
		c.JSON(400, gin.H{
			"err": "获取用户信息失败",
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "获取用户信息成功",
		"data": userItem,
	})
}

// 添加user
func (u *user) AddUser(c *gin.Context) {
	user := new(model.User)
	if err := c.ShouldBind(user); err != nil {
		c.JSON(400, gin.H{
			"err": "绑定数据失败：" + err.Error(),
		})
		return
	}
	err := service.Users.AddUser(user)
	if err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"msg": "新增用户成功",
	})
}

// 更新用户信息
func (g *user) UpdateUser(c *gin.Context) {
	user := new(model.User)
	if err := c.ShouldBind(user); err != nil {
		c.JSON(400, gin.H{
			"err": "绑定数据失败：" + err.Error(),
		})
		return
	}

	err := service.Users.UpdateUserAuth(user)
	if err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"msg": "更新用户信息成功",
	})
}

// 删除用户
func (u *user) DeleteUser(c *gin.Context) {
	user := c.Query("user")
	err := service.Users.DeleteUser(user)
	if err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"msg": "删除用户成功",
	})
}

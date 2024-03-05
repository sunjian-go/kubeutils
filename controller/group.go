package controller

import (
	"github.com/gin-gonic/gin"
	"main/model"
	"main/service"
)

var Group group

type group struct {
}

func (g *group) Getgroups(c *gin.Context) {
	group := new(service.GroupInfo)
	if err := c.Bind(group); err != nil {
		c.JSON(400, gin.H{

			"err": "绑定数据失败：" + err.Error(),
		})
		return
	}
	groupItem, err := service.Group.GetFilterGroups(group)
	if err != nil {
		c.JSON(400, gin.H{
			"err": "获取用户组列表失败",
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "获取用户组列表成功",
		"data": groupItem,
	})
}

// 查询单条group
func (g *group) Getgroup(c *gin.Context) {
	group := c.Query("group")
	groupItem, err := service.Group.GetGroup(group)
	if err != nil {
		c.JSON(400, gin.H{
			"err": "获取用户组失败",
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "获取用户组成功",
		"data": groupItem,
	})
}

// 添加group
func (g *group) Addgroup(c *gin.Context) {
	group := new(model.Group)
	if err := c.ShouldBind(group); err != nil {
		c.JSON(400, gin.H{
			"err": "绑定数据失败：" + err.Error(),
		})
		return
	}
	err := service.Group.AddGroup(group)
	if err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"msg": "新增用户组成功",
	})
}

// 更新用户组信息
func (g *group) UpdategroupAuth(c *gin.Context) {
	group := new(model.Group)
	if err := c.ShouldBind(group); err != nil {
		c.JSON(400, gin.H{
			"err": "绑定数据失败：" + err.Error(),
		})
		return
	}
	err := service.Group.UpdateGroupAuth(group)
	if err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"msg": "更新角色成功",
	})
}

// 删除用户组
func (g *group) Deletegroup(c *gin.Context) {
	group := c.Query("group")
	err := service.Group.DeleteGroup(group)
	if err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"msg": "删除用户组成功",
	})
}

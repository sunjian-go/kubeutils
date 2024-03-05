package controller

import (
	"github.com/gin-gonic/gin"
	"main/model"
	"main/service"
)

var Role role

type role struct {
}

func (r *role) GetRoles(c *gin.Context) {
	role := new(service.RoleInfo)
	if err := c.Bind(role); err != nil {
		c.JSON(400, gin.H{

			"err": "绑定数据失败：" + err.Error(),
		})
		return
	}
	roleItem, err := service.Role.GetFilterRoles(role)
	if err != nil {
		c.JSON(400, gin.H{
			"err": "获取role列表失败",
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "获取role列表成功",
		"data": roleItem,
	})
}

// 获取单条role
func (r *role) GetRole(c *gin.Context) {
	role := c.Query("role")
	roleItem, err := service.Role.GetRole(role)
	if err != nil {
		c.JSON(400, gin.H{
			"err": "获取role失败",
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "获取role成功",
		"data": roleItem,
	})
}

// 添加role
func (r *role) AddRole(c *gin.Context) {
	role := new(model.Role)
	if err := c.ShouldBind(role); err != nil {
		c.JSON(400, gin.H{
			"err": "绑定数据失败：" + err.Error(),
		})
		return
	}
	err := service.Role.AddRole(role)
	if err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"msg": "新增role成功",
	})
}

// 更新权限
func (r *role) UpdateRoleAuth(c *gin.Context) {
	role := new(model.Role)
	if err := c.ShouldBind(role); err != nil {
		c.JSON(400, gin.H{
			"err": "绑定数据失败：" + err.Error(),
		})
		return
	}
	err := service.Role.UpdateRoleAuth(role)
	if err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"msg": "更新权限成功",
	})
}

// 删除角色
func (r *role) DeleteRole(c *gin.Context) {
	role := c.Query("role")
	err := service.Role.DeleteRole(role)
	if err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"msg": "删除角色成功",
	})
}

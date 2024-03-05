package dao

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"main/db"
	"main/model"
)

var Group group

type group struct {
}

// 获取表中所有用户组列表
func (g *group) GetAllGroups(page, limit int) ([]model.Group, int, error) {
	//定义分页的起始位置
	startSet := (page - 1) * limit
	//不指定容量定义切片
	var group []model.Group

	//先查出所有角色的数量
	tx := db.GORM.Find(&group)
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		return group, 0, errors.New("获取用户组列表失败，" + tx.Error.Error())
	}
	total := len(group)

	//查询出所有数据存入切片
	tx = db.GORM.Limit(limit).
		Offset(startSet).
		Order("id desc").
		Find(&group)
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		return group, 0, errors.New("获取用户组列表失败，" + tx.Error.Error())
	}
	return group, total, nil
}

// 根据条件查询用户组信息
func (g *group) GetFilterGroups(fiterName string, page, limit int) ([]model.Group, error) {
	//定义分页的起始位置
	startSet := (page - 1) * limit
	var group []model.Group
	var tx *gorm.DB

	//按需查
	tx = db.GORM.Model(group).Where("group_name like ?", "%"+fiterName+"%").
		Limit(limit).
		Offset(startSet).
		Order("id desc").
		Find(&group)
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		return nil, errors.New("查询用户组失败，" + tx.Error.Error())
	}

	fmt.Println("查询用户组信息：", group)
	return group, nil
}

// 查询单条group
func (g *group) GetGroup(group string) (*model.Group, error) {
	gro := new(model.Group)
	//查询出所有数据存入切片
	tx := db.GORM.Where("group_name = ?", group).First(gro)
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		return gro, errors.New("获取用户组失败，" + tx.Error.Error())
	}
	return gro, nil
}

// 新增用户组
func (g *group) AddGroup(group *model.Group) error {
	tx := db.GORM.Create(group)
	if tx.Error != nil {
		return errors.New("新增用户组失败: " + tx.Error.Error())
	}
	return nil
}

// 更新用户组
func (g *group) UpdateGroup(group *model.Group) error {
	//gro := &model.Group{}
	if group.Role == "" {
		tx := db.GORM.Model(group).Where("group_name = ?", group.GroupName).Update(model.Group{Role: "null"})
		if tx.Error != nil && tx.Error.Error() != "record not found" {
			return errors.New("更新角色失败：" + tx.Error.Error())
		}
	} else {
		tx := db.GORM.Model(group).Where("group_name = ?", group.GroupName).Update(model.Group{Role: group.Role})
		if tx.Error != nil && tx.Error.Error() != "record not found" {
			return errors.New("更新角色失败：" + tx.Error.Error())
		}
	}
	return nil
}

// 删除用户组
func (g *group) DeleteGroup(name string) error {
	group := &model.Group{}
	//开始删除
	tx := db.GORM.Where("group_name = ?", name).Delete(group)
	if tx.Error != nil {
		return errors.New("删除用户组失败: " + tx.Error.Error())
	}
	return nil
}

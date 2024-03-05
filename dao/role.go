package dao

import (
	"errors"
	"github.com/jinzhu/gorm"
	"main/db"
	"main/model"
)

var Role role

type role struct {
}

// 获取表中所有role列表
func (r *role) GetAllRoles(page, limit int) ([]model.Role, int, error) {
	//定义分页的起始位置
	startSet := (page - 1) * limit
	//不指定容量定义切片
	var rol []model.Role

	//先查出所有角色的数量
	tx := db.GORM.Find(&rol)
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		return rol, 0, errors.New("获取role列表失败，" + tx.Error.Error())
	}
	total := len(rol)
	//再进行分页查询
	tx = db.GORM.Limit(limit).
		Offset(startSet).
		Order("id desc").
		Find(&rol)
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		return rol, 0, errors.New("获取role列表失败，" + tx.Error.Error())
	}
	return rol, total, nil
}

// 查询单条role
func (r *role) GetRole(name string) (*model.Role, error) {
	rol := new(model.Role)
	//查询出所有数据存入切片
	tx := db.GORM.Where("role = ?", name).First(rol)
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		return rol, errors.New("获取role失败，" + tx.Error.Error())
	}
	return rol, nil
}

// 根据条件查询role信息
func (r *role) GetFilterRoles(fiterName string, page, limit int) ([]model.Role, error) {
	//定义分页的起始位置
	startSet := (page - 1) * limit
	var rols []model.Role
	var tx *gorm.DB

	//按需查
	tx = db.GORM.Model(rols).Where("role like ?", "%"+fiterName+"%").
		Limit(limit).
		Offset(startSet).
		Order("id desc").
		Find(&rols)
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		return nil, errors.New("查询role失败，" + tx.Error.Error())
	}

	//fmt.Println("查询role信息：", rols)
	return rols, nil
}

// 新增role
func (r *role) AddRole(role *model.Role) error {
	tx := db.GORM.Create(role)
	if tx.Error != nil {
		return errors.New("新增role失败: " + tx.Error.Error())
	}
	return nil
}

// 更新role权限
func (r *role) UpdateRole(role *model.Role) error {
	//rol := &model.Role{}
	if role.Auth == "" {
		tx := db.GORM.Model(role).Where("role = ?", role.Role).Update(model.Role{Auth: "null"})
		if tx.Error != nil && tx.Error.Error() != "record not found" {
			return errors.New("更新权限失败：" + tx.Error.Error())
		}
	} else {
		tx := db.GORM.Model(role).Where("role = ?", role.Role).Update(model.Role{Auth: role.Auth})
		if tx.Error != nil && tx.Error.Error() != "record not found" {
			return errors.New("更新权限失败：" + tx.Error.Error())
		}
	}

	return nil
}

// 删除role
func (r *role) DeleteRole(name string) error {
	role := &model.Role{}
	//开始删除
	tx := db.GORM.Where("role = ?", name).Delete(role)
	if tx.Error != nil {
		return errors.New("删除角色失败: " + tx.Error.Error())
	}
	return nil
}

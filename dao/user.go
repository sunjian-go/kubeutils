package dao

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"main/db"
	"main/model"
)

var User user

type user struct {
}

// 获取表中所有用户列表
func (u *user) GetAllUsers(page, limit int) ([]model.User, int, error) {
	//定义分页的起始位置
	startSet := (page - 1) * limit
	//不指定容量定义切片
	var user []model.User

	//先查出所有角色的数量
	tx := db.GORM.Find(&user)
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		return user, 0, errors.New("获取用户列表失败，" + tx.Error.Error())
	}
	total := len(user)

	//查询出所有数据存入切片
	tx = db.GORM.Limit(limit).
		Offset(startSet).
		Order("id desc").
		Find(&user)
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		return user, 0, errors.New("获取用户列表失败: " + tx.Error.Error())
	}
	return user, total, nil
}

// 根据条件查询用户信息
func (u *user) GetFilterUsers(fiterName string, page, limit int) ([]model.User, error) {
	//定义分页的起始位置
	startSet := (page - 1) * limit
	var user []model.User
	var tx *gorm.DB

	//按需查
	tx = db.GORM.Model(user).Where("username like ?", "%"+fiterName+"%").
		Limit(limit).
		Offset(startSet).
		Order("id desc").
		Find(&user)
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		return nil, errors.New("查询用户失败，" + tx.Error.Error())
	}

	fmt.Println("查询用户组信息：", user)
	return user, nil
}

// 查询单条用户信息
func (u *user) GetUser(name string) (*model.User, error) {
	user := new(model.User)
	//查询出所有数据存入切片
	tx := db.GORM.Where("username = ?", name).First(user)
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		return user, errors.New("获取用户失败，" + tx.Error.Error())
	}
	return user, nil
}

// 新增用户
func (u *user) AddUser(user *model.User) error {
	tx := db.GORM.Create(user)
	if tx.Error != nil {
		return errors.New("新增用户失败: " + tx.Error.Error())
	}
	return nil
}

// 更新用户密码和用户组
func (u *user) UpdateUser(user *model.User) error {
	//先修改密码
	tx := db.GORM.Model(user).Where("username = ?", user.Username).Update(model.User{Password: user.Password})
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		return errors.New("更新用户密码失败：" + tx.Error.Error())
	}
	//再修改用户组
	tx = db.GORM.Model(user).Where("username = ?", user.Username).Update(model.User{GroupName: user.GroupName})
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		return errors.New("更新用户组失败：" + tx.Error.Error())
	}

	return nil
}

// 删除用户
func (u *user) DeleteUser(name string) error {
	user := &model.User{}
	//开始删除
	tx := db.GORM.Where("username = ?", name).Delete(user)
	if tx.Error != nil {
		return errors.New("删除用户失败: " + tx.Error.Error())
	}
	return nil
}

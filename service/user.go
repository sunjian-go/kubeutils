package service

import (
	"errors"
	"main/dao"
	"main/model"
	"main/utils"
	"strconv"
)

var Users user

type user struct {
}

type Userlist struct {
	Items []model.User `json:"items"`
	Total int          `json:"total"`
}

type UserInfo struct {
	FilterName string `form:"filterName"`
	Page       string `form:"page"`
	Limit      string `form:"limit"`
}

// 查询用户列表
func (u *user) GetFilterUsers(userInfo *UserInfo) (*Userlist, error) {
	newpage, _ := strconv.Atoi(userInfo.Page)
	newlimit, _ := strconv.Atoi(userInfo.Limit)
	var Users []model.User
	var err error
	total := 0
	if userInfo.FilterName == "" {
		//查所有
		Users, total, err = dao.User.GetAllUsers(newpage, newlimit)
		if err != nil {
			utils.Logg.Error(err.Error())
			return nil, err
		}
	} else {
		//按需查
		Users, err = dao.User.GetFilterUsers(userInfo.FilterName, newpage, newlimit)
		if err != nil {
			utils.Logg.Error(err.Error())
			return nil, err
		}
		total = len(Users)
	}
	return &Userlist{
		Items: Users,
		Total: total,
	}, nil
}

// 查询单条用户信息
func (u *user) GetUser(name string) (*model.User, error) {
	user, err := dao.User.GetUser(name)
	if err != nil {
		utils.Logg.Error(err.Error())
		return nil, err
	}
	return user, nil
}

// 新增用户
func (u *user) AddUser(addUser *model.User) error {
	//先去重
	Userlist, err := dao.User.GetFilterUsers(addUser.Username, 1, 10)
	if err != nil {
		return err
	}
	if len(Userlist) > 0 {
		return errors.New("该用户已存在")
	}

	err = dao.User.AddUser(addUser)
	if err != nil {
		utils.Logg.Error(err.Error())
		return err
	}
	return nil
}

// 更新用户信息
func (u *user) UpdateUserAuth(User *model.User) error {
	err := dao.User.UpdateUser(User)
	if err != nil {
		utils.Logg.Error(err.Error())
		return err
	}
	return nil
}

// 删除用户
func (u *user) DeleteUser(name string) error {
	//先查询是否删除
	users, err := dao.User.GetFilterUsers(name, 1, 10)
	if err != nil {
		utils.Logg.Error(err.Error())
		return err
	}
	if len(users) == 0 {
		return errors.New("用户不存在或已被删除")
	}
	err = dao.User.DeleteUser(name)
	if err != nil {
		utils.Logg.Error(err.Error())
		return err
	}
	return nil
}

package service

import (
	"errors"
	"main/dao"
	"main/model"
	"main/utils"
	"strconv"
)

var Role role

type role struct {
}

type Roles struct {
	Items []model.Role `json:"items"`
	Total int          `json:"total"`
}

type RoleInfo struct {
	FilterName string `form:"filterName"`
	Page       string `form:"page"`
	Limit      string `form:"limit"`
}

// 查询role列表
func (r *role) GetFilterRoles(roleinfo *RoleInfo) (*Roles, error) {
	newpage, _ := strconv.Atoi(roleinfo.Page)
	newlimit, _ := strconv.Atoi(roleinfo.Limit)
	var roles []model.Role
	var err error
	total := 0
	if roleinfo.FilterName == "" {
		//查所有
		roles, total, err = dao.Role.GetAllRoles(newpage, newlimit)
		if err != nil {
			utils.Logg.Error(err.Error())
			return nil, err
		}
	} else {
		//按需查
		roles, err = dao.Role.GetFilterRoles(roleinfo.FilterName, newpage, newlimit)
		if err != nil {
			utils.Logg.Error(err.Error())
			return nil, err
		}
		total = len(roles)
	}
	return &Roles{
		Items: roles,
		Total: total,
	}, nil
}

// 获取单条role
func (r *role) GetRole(name string) (*model.Role, error) {
	role, err := dao.Role.GetRole(name)
	if err != nil {
		utils.Logg.Error(err.Error())
		return nil, err
	}
	return role, nil
}

// 新增role
func (r *role) AddRole(addRole *model.Role) error {
	//先去重
	rolelist, err := dao.Role.GetFilterRoles(addRole.Role, 1, 10)
	if err != nil {
		return err
	}
	if len(rolelist) > 0 {
		return errors.New("该角色已存在")
	}

	err = dao.Role.AddRole(addRole)
	if err != nil {
		utils.Logg.Error(err.Error())
		return err
	}
	return nil
}

// 更新role权限
func (r *role) UpdateRoleAuth(role *model.Role) error {
	err := dao.Role.UpdateRole(role)
	if err != nil {
		utils.Logg.Error(err.Error())
		return err
	}
	return nil
}

// 删除角色
func (r *role) DeleteRole(name string) error {
	//先查询是否已删除
	roles, err := dao.Role.GetFilterRoles(name, 1, 10)
	if err != nil {
		utils.Logg.Error(err.Error())
		return err
	}
	if len(roles) == 0 {
		return errors.New("该角色不存在或已被删除")
	}
	err = dao.Role.DeleteRole(name)
	if err != nil {
		utils.Logg.Error(err.Error())
		return err
	}
	return nil
}

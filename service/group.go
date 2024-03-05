package service

import (
	"errors"
	"main/dao"
	"main/model"
	"main/utils"
	"strconv"
)

var Group group

type group struct {
}

type Grouplist struct {
	Items []model.Group `json:"items"`
	Total int           `json:"total"`
}

type GroupInfo struct {
	FilterName string `form:"filterName"`
	Page       string `form:"page"`
	Limit      string `form:"limit"`
}

// 查询Group列表
func (g *group) GetFilterGroups(Groupinfo *GroupInfo) (*Grouplist, error) {
	newpage, _ := strconv.Atoi(Groupinfo.Page)
	newlimit, _ := strconv.Atoi(Groupinfo.Limit)
	var Groups []model.Group
	total := 0
	var err error
	if Groupinfo.FilterName == "" {
		//查所有
		Groups, total, err = dao.Group.GetAllGroups(newpage, newlimit)
		if err != nil {
			utils.Logg.Error(err.Error())
			return nil, err
		}
	} else {
		//按需查
		Groups, err = dao.Group.GetFilterGroups(Groupinfo.FilterName, newpage, newlimit)
		if err != nil {
			utils.Logg.Error(err.Error())
			return nil, err
		}
		total = len(Groups)
	}
	return &Grouplist{
		Items: Groups,
		Total: total,
	}, nil
}

// 查询单条group
func (g *group) GetGroup(name string) (*model.Group, error) {
	group, err := dao.Group.GetGroup(name)
	if err != nil {
		utils.Logg.Error(err.Error())
		return nil, err
	}
	return group, nil
}

// 新增Group
func (g *group) AddGroup(addGroup *model.Group) error {
	//先去重
	Grouplist, err := dao.Group.GetFilterGroups(addGroup.GroupName, 1, 10)
	if err != nil {
		return err
	}
	if len(Grouplist) > 0 {
		return errors.New("该用户组已存在")
	}

	err = dao.Group.AddGroup(addGroup)
	if err != nil {
		utils.Logg.Error(err.Error())
		return err
	}
	return nil
}

// 更新Group权限
func (g *group) UpdateGroupAuth(Group *model.Group) error {
	err := dao.Group.UpdateGroup(Group)
	if err != nil {
		utils.Logg.Error(err.Error())
		return err
	}
	return nil
}

// 删除角色
func (g *group) DeleteGroup(name string) error {
	//先查询是否删除
	group, err := dao.Group.GetFilterGroups(name, 1, 10)
	if err != nil {
		utils.Logg.Error(err.Error())
		return err
	}
	if len(group) == 0 {
		return errors.New("用户组不存在或已被删除")
	}
	err = dao.Group.DeleteGroup(name)
	if err != nil {
		utils.Logg.Error(err.Error())
		return err
	}
	return nil
}

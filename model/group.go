package model

import "time"

// 定义集群注册表结构
type Group struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	GroupName string     `json:"groupName"`
	Role      string     `json:"role"`
}

// 返回mysql表名，以此来定义mysql中的表名
func (*Group) TableName() string {
	return "group_tab"
}

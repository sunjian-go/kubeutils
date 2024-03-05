package model

import "time"

// 定义集群注册表结构
type User struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	Username  string     `json:"username"`
	Password  string     `json:"password"`
	GroupName string     `json:"groupName"`
}

// 返回mysql表名，以此来定义mysql中的表名
func (*User) TableName() string {
	return "user"
}

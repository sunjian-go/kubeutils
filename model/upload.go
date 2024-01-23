package model

import "time"

// 定义集群注册表结构
type Upload_History struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
	ClusterName string     `json:"cluster_name"`
	PodName     string     `json:"podName"`
	Namespace   string     `json:"namespace"`
	Path        string     `json:"path"`
	File        string     `json:"file"`
	Status      string     `json:"status"`
	Code        string     `json:"code"`
}

// 返回mysql表名，以此来定义mysql中的表名
func (*Upload_History) TableName() string {
	return "upload_history"
}

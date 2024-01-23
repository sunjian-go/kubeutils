package model

import "time"

// 定义集群注册表结构
type Cluster struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
	ClusterName string     `json:"cluster_name"`
	Ipaddr      string     `json:"ipaddr"`
	Port        string     `json:"port"`
	K8sVersion  string     `json:"k8s_version"`
	Status      string     `json:"status"`
}

// 返回mysql表名，以此来定义mysql中的表名
func (*Cluster) TableName() string {
	return "cluster"
}

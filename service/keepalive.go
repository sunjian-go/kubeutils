package service

import (
	"errors"
	"fmt"
	"main/dao"
	"main/model"
)

var Keep keepalive

type keepalive struct {
}

// 更新集群在线状态
func (k *keepalive) UpdateStatus(name, status string) error {
	//更新之前先查询该集群是否被软删除
	clu := new(model.Cluster)
	clu.ClusterName = name
	clus, err := dao.RegCluster.GetClusetrInfo(clu)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	if clus.ClusterName == "" {
		return errors.New("集群已被删除")
	}

	//如果集群没被删除就继续更新状态
	err = dao.RegCluster.UpdateClusterStatus(name, status)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

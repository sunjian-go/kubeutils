package dao

import (
	"errors"
	"fmt"
	"main/db"
	"main/model"
	"main/utils"
	"strconv"
)

var RegCluster reqcluster

type reqcluster struct {
}

// 获取表中单个集群信息
func (r *reqcluster) GetClusetrInfo(cluster *model.Cluster) (*model.Cluster, error) {
	cl := new(model.Cluster)
	tx := db.GORM.Where("cluster_name = ?", cluster.ClusterName).First(cl)
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		return nil, errors.New("获取单个集群信息失败，" + tx.Error.Error())
	}
	return cl, nil
}

// 获取表中所有集群信息
func (r *reqcluster) GetAllClusetrInfo() ([]model.Cluster, error) {
	//不指定容量定义切片
	var cls []model.Cluster
	//查询出所有数据存入切片
	tx := db.GORM.Find(&cls)
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		return cls, errors.New("获取集群信息失败，" + tx.Error.Error())
	}
	return cls, nil
}

// 根据条件查询集群信息
func (r *reqcluster) GetAllClusetr(fiterName string, page, limit int) ([]model.Cluster, int, error) {
	//定义分页的起始位置
	startSet := (page - 1) * limit
	var clusters []model.Cluster

	//先查出符合条件的总数量
	tx := db.GORM.Model(clusters).Find(&clusters)
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		return nil, 0, errors.New("查询所有集群失败，" + tx.Error.Error())
	}
	utils.Logg.Info("total= " + strconv.Itoa(len(clusters)))
	total := len(clusters)
	clusters = nil
	tx = db.GORM.Model(clusters).Where("cluster_name like ?", "%"+fiterName+"%").
		Limit(limit).
		Offset(startSet).
		Order("id desc").
		Find(&clusters)
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		return nil, 0, errors.New("按需查询集群失败，" + tx.Error.Error())
	}
	fmt.Println("查询集群信息：", clusters)
	return clusters, total, nil
}

// 注册集群
func (r *reqcluster) Register(cluster *model.Cluster) error {
	tx := db.GORM.Create(cluster)
	if tx.Error != nil {
		return errors.New("初始化cluster表失败: " + tx.Error.Error())
	}
	return nil
}

// 更新集群IP信息
func (r *reqcluster) UpdateClusterInfo(opt string, cluster *model.Cluster) error {
	clu := &model.Cluster{}
	switch opt {
	case "ip":
		tx := db.GORM.Model(clu).Where("cluster_name = ?", cluster.ClusterName).Update(model.Cluster{Ipaddr: cluster.Ipaddr})
		if tx.Error != nil && tx.Error.Error() != "record not found" {
			return errors.New("更新集群IP失败，" + tx.Error.Error())
		}
		break
	case "port":
		tx := db.GORM.Model(clu).Where("cluster_name = ?", cluster.ClusterName).Update(model.Cluster{Port: cluster.Port})
		if tx.Error != nil && tx.Error.Error() != "record not found" {
			return errors.New("更新集群port失败，" + tx.Error.Error())
		}
		break
	}
	return nil
}

// 更新集群活跃信息
func (r *reqcluster) UpdateClusterStatus(name, status string) error {
	utils.Logg.Info("将" + name + "状态改为：" + status)
	clu := &model.Cluster{}
	tx := db.GORM.Model(clu).Where("cluster_name = ?", name).Update(model.Cluster{Status: status})
	if tx.Error != nil {
		if tx.Error.Error() == "record not found" {
			return errors.New("集群不存在：" + tx.Error.Error())
		} else {
			return errors.New("更新集群状态失败，" + tx.Error.Error())
		}

	}
	return nil
}

// 通过集群name获取ip端口等信息
func (r *reqcluster) GetClusterIP(name string) (*model.Cluster, error) {
	clu := &model.Cluster{}
	tx := db.GORM.Where("cluster_name = ?", name).First(clu)
	if tx.Error != nil {
		if tx.Error.Error() == "record not found" {
			return nil, errors.New("该集群不存在")
		} else {
			return nil, errors.New("获取该集群ip失败，" + tx.Error.Error())
		}
	}

	return clu, nil
}

// 删除集群
func (r *reqcluster) DeleteCluster(name string) error {
	clu := &model.Cluster{}
	//先查询是否已删除
	tx := db.GORM.Where("cluster_name = ?", name).First(clu)
	if tx.Error != nil {
		if tx.Error.Error() == "record not found" {
			return errors.New("目标集群不存在或已被删除")
		} else {
			return errors.New("查找目标集群失败，" + tx.Error.Error())
		}
	}
	//开始删除
	tx = db.GORM.Where("cluster_name = ?", name).Delete(clu)
	if tx.Error != nil {
		return errors.New("删除目标集群失败: " + tx.Error.Error())
	}
	return nil
}

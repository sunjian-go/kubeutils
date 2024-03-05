package service

import (
	"fmt"
	"github.com/robfig/cron"
	"main/dao"
	"main/model"
	"main/utils"
	"time"
)

var ClusTers cluss

type cluss struct {
}

// 注册
func (c *cluss) RegisterFunc(cluster *model.Cluster) error {
	//先查询表里是否存在该集群
	clu, err := dao.RegCluster.GetClusetrInfo(cluster)
	if err != nil {
		utils.Logg.Error(err.Error())
		return err
	}
	if clu.ClusterName != "" {
		if cluster.Ipaddr != clu.Ipaddr {
			utils.Logg.Info("该集群信息已变更")
			//代表表中集群ip变了,就去更新ip即可
			err := dao.RegCluster.UpdateClusterInfo("ip", cluster)
			if err != nil {
				return err
			}
		}
		if cluster.Port != clu.Port {
			//代表表中集群port变了,就去更新ip即可
			err := dao.RegCluster.UpdateClusterInfo("port", cluster)
			if err != nil {
				return err
			}
		}
	} else {
		utils.Logg.Info("该集群信息不存在,去创建")
		//代表表里没有该集群，泽去创建
		err := dao.RegCluster.Register(cluster)
		if err != nil {
			return err
		}
	}
	return nil
}

var cronObj *cron.Cron

// 定时器
func (c *cluss) CronFunc() {
	cronObj = cron.New()
	err := cronObj.AddFunc("*/60 * * * * *", func() {
		//每60秒执行一次检查集群健康状态
		CheckHelth()
	})
	if err != nil {
		utils.Logg.Error("定时器执行报错：" + err.Error())
	}
	//启动/关闭
	cronObj.Start()
}
func (c *cluss) CloseCron() {
	cronObj.Stop()
}

// 启动go携程来定时检查集群健康状态
func CheckHelth() {
	//先查出所有agent的数据存入一个结构体切片
	cls, err := dao.RegCluster.GetAllClusetrInfo()
	if err != nil {
		utils.Logg.Error(err.Error())
	}
	//fmt.Println("取出数据为：", cls)
	for _, data := range cls {
		//根据查出来的name查询更新时间
		currentDate := time.Now()
		fmt.Println("打印数据：", data.ClusterName, data.UpdatedAt.Unix())
		fmt.Println("当前时间戳：", currentDate.Unix())
		//如果集群更新时间戳与当前时间戳相差大于60秒，则代表该集群失联（agent每60秒发一次心跳）
		if currentDate.Unix()-data.UpdatedAt.Unix() > 60 {
			//失联的集群将该集群状态置为inactive
			utils.Logg.Info(data.ClusterName + "已失联。。。")
			err := dao.RegCluster.UpdateClusterStatus(data.ClusterName, "inactive")
			if err != nil {
				utils.Logg.Error(err.Error())
			}
		}
	}
}

type ClusterInfo struct {
	FilterName string `form:"filterName"`
	Page       int    `form:"page"`
	Limit      int    `form:"limit"`
}

// 获取所有集群
func (c *cluss) GetAllClusters(clu *ClusterInfo) ([]model.Cluster, int, error) {
	//clus, err := dao.RegCluster.GetAllClusetrInfo()
	clus, total, err := dao.RegCluster.GetAllClusetr(clu.FilterName, clu.Page, clu.Limit)
	if err != nil {
		utils.Logg.Error(err.Error())
		return nil, 0, err
	}
	return clus, total, nil
}

// 删除目标集群
func (c *cluss) DeleteCluster(name string) error {
	err := dao.RegCluster.DeleteCluster(name)
	if err != nil {
		utils.Logg.Error(err.Error())
		return err
	}
	//删除本集群上传记录
	err = dao.Uploadhist.DeleteAllFilesInfo(name)
	if err != nil {
		utils.Logg.Error(err.Error())
		return err
	}
	return nil
}

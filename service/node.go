package service

import (
	"encoding/json"
	"errors"
	"main/dao"
	"main/utils"
)

var Node node

type node struct {
}

type NodeInfo struct {
	FilterName string `form:"filter_name"`
	Limit      string `form:"limit"`
	Page       string `form:"page"`
}

// 获取node列表
func (n *node) GetNodes(token, clusterName string, nodeinfo *NodeInfo) (interface{}, error) {
	//根据集群名获取IP
	clu, err := dao.RegCluster.GetClusterIP(clusterName)
	if err != nil {
		utils.Logg.Error(err.Error())
		return "", err
	}

	urls := "http://" + clu.Ipaddr + ":" + clu.Port + "/api/corev1/getnodes?filter_name=" + nodeinfo.FilterName + "&limit=" + nodeinfo.Limit + "&page=" + nodeinfo.Page
	body, code, err := utils.Http.HttpSend(urls, nil, "GET")
	if err != nil {
		return nil, err
	}

	// 解析 body 内容为 JSON 格式
	var data map[string]interface{}
	//解码到data中
	err = json.Unmarshal(body, &data)
	if err != nil {
		utils.Logg.Error("解析 JSON 数据时出错:" + err.Error())
		return "", errors.New("解析 JSON 数据时出错:" + err.Error())
	}

	if code == 200 {
		return data, nil
	} else {
		return data["err"], errors.New("err")
	}
}

// 获取node详情
func (n *node) GetNodeDetail(token, clusterName, nodeName string) (interface{}, error) {
	//根据集群名获取IP
	clu, err := dao.RegCluster.GetClusterIP(clusterName)
	if err != nil {
		utils.Logg.Error(err.Error())
		return "", err
	}

	urls := "http://" + clu.Ipaddr + ":" + clu.Port + "/api/corev1/getnodedetail?nodeName=" + nodeName
	body, code, err := utils.Http.HttpSend(urls, nil, "GET")
	if err != nil {
		return nil, err
	}

	// 解析 body 内容为 JSON 格式
	var data map[string]interface{}
	//解码到data中
	err = json.Unmarshal(body, &data)
	if err != nil {
		utils.Logg.Error("解析 JSON 数据时出错:" + err.Error())
		return "", errors.New("解析 JSON 数据时出错:" + err.Error())
	}

	if code == 200 {
		return data, nil
	} else {
		return data["err"], errors.New("err")
	}
}

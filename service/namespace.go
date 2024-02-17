package service

import (
	"encoding/json"
	"errors"
	"github.com/wonderivan/logger"
	"main/dao"
	"main/utils"
)

var Namespace namespace

type namespace struct {
}

func (n *namespace) GetNamespaces(token, clusterName string) (interface{}, error) {
	//根据集群名获取IP
	clu, err := dao.RegCluster.GetClusterIP(clusterName)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	urls := "http://" + clu.Ipaddr + ":" + clu.Port + "/api/corev1/getnamespaces"
	body, code, err := utils.Http.HttpSend(urls, nil, "GET")
	if err != nil {
		return nil, err
	}

	// 解析 body 内容为 JSON 格式
	var data map[string]interface{}
	//解码到data中
	err = json.Unmarshal(body, &data)
	if err != nil {
		logger.Error("解析 JSON 数据时出错:" + err.Error())
		return "", errors.New("解析 JSON 数据时出错:" + err.Error())
	}

	if code == 200 {
		return data, nil
	} else {
		return data["err"], errors.New("err")
	}
}

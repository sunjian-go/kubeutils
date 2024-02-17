package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wonderivan/logger"
	"main/dao"
	"main/utils"
)

var Listpath listpath

type listpath struct {
}

type PodPath struct {
	PodName       string `form:"podName"`
	Namespace     string `form:"namespace"`
	ContainerName string `form:"containerName"`
	Path          string `form:"path"`
}

func (l *listpath) ListContainerPath(podinfo *PodPath, token, clusterName string) (interface{}, error) {
	fmt.Println("podinfo:", podinfo)
	// 将结构体编码为 JSON
	podData, err := json.Marshal(podinfo)
	if err != nil {
		logger.Error("编码结构体为 JSON 时出错：" + err.Error())
		return "", errors.New("编码结构体为 JSON 时出错：" + err.Error())
	}
	// 创建一个包含 JSON 数据的 io.Reader
	//jsonReader := bytes.NewReader(podData)
	jsonReader := bytes.NewBuffer(podData)

	//根据集群名获取IP
	clu, err := dao.RegCluster.GetClusterIP(clusterName)
	if err != nil {
		logger.Error(err.Error())
	}
	urls := "http://" + clu.Ipaddr + ":" + clu.Port + "/api/listPath"

	body, code, err := utils.Http.HttpSend(urls, jsonReader, "GET")
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

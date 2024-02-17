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

var Pod pod

type pod struct {
}
type PodInfo struct {
	FilterName string `form:"filter_name"`
	NameSpace  string `form:"namespace"`
	Limit      int    `form:"limit"`
	Page       int    `form:"page"`
}

type PodDetail struct {
	Name      string `form:"name"`
	Namespace string `form:"namespace"`
}

// 根据需求获取pod或container并返回前端
func (p *pod) GetObjs(token, clusterName, opt string, podinfo interface{}) (interface{}, error) {
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
		return nil, err
	}

	// 创建 HTTP 请求
	var urls string
	switch opt {
	case "getpods":
		urls = "http://" + clu.Ipaddr + ":" + clu.Port + "/api/corev1/getpods"
		break
	case "getContainer":
		urls = "http://" + clu.Ipaddr + ":" + clu.Port + "/api/corev1/getcontainers"
		break
	}

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

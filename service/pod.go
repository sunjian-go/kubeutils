package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wonderivan/logger"
	"io/ioutil"
	"main/dao"
	"net/http"
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
	req, err := http.NewRequest("GET", urls, jsonReader) //后端需要用ShouldBindJSON来接收参数
	if err != nil {
		logger.Error("创建 HTTP 请求报错：" + err.Error())
		return "", errors.New("创建 HTTP 请求报错：" + err.Error())
	}

	// 在请求头中添加Authorization头，携带Token
	req.Header.Set("Authorization", token)
	// 设置请求头的 Content-Type 为 application/json
	req.Header.Set("Content-Type", "application/json")
	fmt.Println("发送：", req)
	// 发送 HTTP 请求
	var resp *http.Response
	// 创建 HTTP 客户端
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		logger.Error("发送 HTTP 请求报错：" + err.Error())
		return "", errors.New("发送 HTTP 请求报错，请检查后端agent服务是否正常运行")
	}
	defer resp.Body.Close()

	fmt.Println("状态信息：", resp.Status)
	if resp.Status == "200 OK" {
		// 读取响应的 body 内容
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.Error("读取响应 body 时出错:" + err.Error())
			return "", errors.New("读取响应 body 时出错:" + err.Error())
		}
		// 解析 body 内容为 JSON 格式
		var data map[string]interface{}
		//解码到data中
		err = json.Unmarshal(body, &data)
		if err != nil {
			logger.Error("解析 JSON 数据时出错:" + err.Error())
			return "", errors.New("解析 JSON 数据时出错:" + err.Error())
		}
		//fmt.Println("获取到data: ", data)
		return data, nil
	} else {
		logger.Error("获取数据失败。。。")
		return "", errors.New("获取数据失败。。。")
	}
}

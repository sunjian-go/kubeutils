package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"main/dao"
	"net/http"
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
		fmt.Println("编码结构体为 JSON 时出错：" + err.Error())
		return "", errors.New("编码结构体为 JSON 时出错：" + err.Error())
	}
	// 创建一个包含 JSON 数据的 io.Reader
	//jsonReader := bytes.NewReader(podData)
	jsonReader := bytes.NewBuffer(podData)

	//根据集群名获取IP
	clu, err := dao.RegCluster.GetClusterIP(clusterName)
	if err != nil {
		fmt.Println(err.Error())
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest("GET", "http://"+clu.Ipaddr+":"+clu.Port+"/api/listPath", jsonReader) //后端需要用ShouldBindJSON来接收参数
	if err != nil {
		fmt.Println("创建 HTTP 请求报错：" + err.Error())
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
		fmt.Println("发送 HTTP 请求报错：" + err.Error())
		return "", errors.New("发送 HTTP 请求报错：" + err.Error())
	}
	defer resp.Body.Close()

	//for {
	//	fmt.Println("状态：", resp.Status)
	//}
	fmt.Println("状态信息：", resp.Status)
	if resp.Status == "200 OK" {
		// 读取响应的 body 内容
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("读取响应 body 时出错:" + err.Error())
			return "", errors.New("读取响应 body 时出错:" + err.Error())
		}
		// 解析 body 内容为 JSON 格式
		var data map[string]interface{}
		//解码到data中
		err = json.Unmarshal(body, &data)
		if err != nil {
			fmt.Println("解析 JSON 数据时出错:" + err.Error())
			return "", errors.New("解析 JSON 数据时出错:" + err.Error())
		}
		return data, nil
	} else {
		fmt.Println("获取数据失败。。。")
		return "", errors.New("获取数据失败。。。")
	}
}

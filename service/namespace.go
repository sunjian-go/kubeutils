package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wonderivan/logger"
	"io/ioutil"
	"main/dao"
	"net/http"
)

var Namespace namespace

type namespace struct {
}

func (n *namespace) GetNamespaces(token, clusterName string) (interface{}, error) {
	// 创建 HTTP 客户端
	client := &http.Client{}

	//根据集群名获取IP
	clu, err := dao.RegCluster.GetClusterIP(clusterName)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest(http.MethodGet, "http://"+clu.Ipaddr+":"+clu.Port+"/api/corev1/getnamespaces", nil)
	if err != nil {
		logger.Error("创建 HTTP 请求报错：" + err.Error())
		return "", errors.New("创建 HTTP 请求报错：" + err.Error())
	}

	// 在请求头中添加Authorization头，携带Token
	req.Header.Set("Authorization", token)

	// 发送 HTTP 请求
	var resp *http.Response

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
		return data, nil
	} else {
		logger.Error("获取namespace列表失败。。。")
		return "", errors.New("获取namespace列表失败。。。")
	}
}

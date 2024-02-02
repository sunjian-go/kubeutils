package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wonderivan/logger"
	"io"
	"main/dao"
	"main/utils"
	"net/http"
)

var Port portt

type portt struct {
}

type PortData struct {
	Ip      string `json:"ip"`
	TcpPort string `json:"tcpPort"`
}

func (p *portt) TCPTelnet(portdata *PortData, clusterName, url string) (interface{}, error) {
	//根据集群名获取IP
	clu, err := dao.RegCluster.GetClusterIP(clusterName)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	urls := "http://" + clu.Ipaddr + ":" + clu.Port + "/api/port?url=" + url

	//将结构体转为json格式
	jsonReader, err := utils.Stj.StructToJson(portdata)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	//创建http请求
	req, err := http.NewRequest("POST", urls, jsonReader) //后端需要用ShouldBindJSON来接收参数
	if err != nil {
		logger.Error("创建 HTTP 请求报错：" + err.Error())
		return "", errors.New("创建 HTTP 请求报错：" + err.Error())
	}
	fmt.Println("发送：", req)

	// 发送 HTTP 请求
	// 创建 HTTP 客户端
	var resp *http.Response
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		logger.Error("发送 HTTP 请求报错：" + err.Error())
		return "", errors.New("发送 HTTP 请求报错，请检查后端agent服务是否正常运行") //代表后端agent挂了或者连不上了
	}
	defer resp.Body.Close()

	fmt.Println("状态信息：", resp.Status)

	// 读取响应的 body 内容
	body, err := io.ReadAll(resp.Body)
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
	if resp.StatusCode == 200 {
		return data["msg"], nil
	} else {
		return data["err"], errors.New("err")
	}

}
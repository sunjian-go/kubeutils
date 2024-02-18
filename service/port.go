package service

import (
	"encoding/json"
	"errors"
	"main/dao"
	"main/utils"
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
		utils.Logg.Error(err.Error())
		return nil, err
	}

	//将结构体转为json格式
	jsonReader, err := utils.Stj.StructToJson(portdata)
	if err != nil {
		utils.Logg.Error(err.Error())
		return nil, err
	}

	urls := "http://" + clu.Ipaddr + ":" + clu.Port + "/api/port?url=" + url
	body, code, err := utils.Http.HttpSend(urls, jsonReader, "POST")
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
		return data["msg"], nil
	} else {
		return data["err"], errors.New("err")
	}

}

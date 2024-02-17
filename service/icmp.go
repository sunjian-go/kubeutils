package service

import (
	"encoding/json"
	"errors"
	"github.com/wonderivan/logger"
	"main/dao"
	"main/utils"
)

var Icmp icmp

type icmp struct {
}

type Icmpdata struct {
	Ip      string `json:"ip"`
	TimeOut string `json:"timeOut"` //超时秒
	Count   string `json:"count"`   //数据包数量
}

// ping方法
func (i *icmp) PingFunc(icmpdata *Icmpdata, clusterName, url string) (interface{}, error) {
	//根据集群名获取IP
	clu, err := dao.RegCluster.GetClusterIP(clusterName)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	urls := "http://" + clu.Ipaddr + ":" + clu.Port + "/api/icmp?url=" + url

	//将结构体转为json格式
	//将结构体转为json格式
	jsonReader, err := utils.Stj.StructToJson(icmpdata)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	body, code, err := utils.Http.HttpSend(urls, jsonReader, "POST")
	if err != nil {
		return nil, err
	}

	// 解析 body 内容为 JSON 格式
	var data map[string]interface{}
	//解码到data中
	err = json.Unmarshal(body, &data)
	if err != nil {
		logger.Error("解析 JSON 数据时出错:" + err.Error())
		return nil, errors.New("解析 JSON 数据时出错:" + err.Error())
	}
	if code == 200 {
		return data, nil
	} else {
		return data["err"], errors.New("err")
	}

}

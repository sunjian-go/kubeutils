package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"io"
	"main/dao"
	"main/utils"
	"net/http"
	"strings"
)

var Pack packet

type packet struct {
}

type PackInfo struct {
	Ip      string `json:"ip"`
	Port    string `json:"port"`
	NetName string `json:"netName"`
}

// 开始抓包
func (p *packet) StartPacket(pcakinfo *PackInfo, clusterName, url string) (interface{}, error) {
	//根据集群名获取IP
	clu, err := dao.RegCluster.GetClusterIP(clusterName)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	// 将结构体编码为 JSON
	//将结构体转为json格式
	jsonReader, err := utils.Stj.StructToJson(pcakinfo)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	//创建http请求
	urls := "http://" + clu.Ipaddr + ":" + clu.Port + "/api/startPacket?url=" + url //去往agent
	req, err := http.NewRequest("POST", urls, jsonReader)                           //后端需要用ShouldBindJSON来接收参数
	if err != nil {
		fmt.Println("创建 HTTP 请求报错：" + err.Error())
		return nil, errors.New("创建 HTTP 请求报错：" + err.Error())
	}

	fmt.Println("发送：", req)

	// 发送 HTTP 请求
	// 创建 HTTP 客户端
	var resp *http.Response
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		fmt.Println("发送 HTTP 请求报错：" + err.Error())
		return nil, errors.New("发送 HTTP 请求报错：" + err.Error())
	}
	defer resp.Body.Close()

	fmt.Println("状态信息：", resp.Status)

	// 读取响应的 body 内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("读取响应 body 时出错:" + err.Error())
		return nil, errors.New("读取响应 body 时出错:" + err.Error())
	}
	// 解析 body 内容为 JSON 格式
	var data map[string]interface{}
	//解码到data中
	err = json.Unmarshal(body, &data)
	if err != nil {
		logger.Error("解析 JSON 数据时出错:" + err.Error())
		return nil, errors.New("解析 JSON 数据时出错:" + err.Error())
	}
	if resp.StatusCode == 200 {
		return data["msg"], nil
	} else {
		return data["err"], errors.New("err")
	}
}

// 停止抓包并获取pcap文件
func (p *packet) StopPacket(cont *gin.Context, clusterName, url string) error {
	//根据集群名获取IP
	clu, err := dao.RegCluster.GetClusterIP(clusterName)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	urls := "http://" + clu.Ipaddr + ":" + clu.Port + "/api/stopPacket?url=" + url
	req, err := http.NewRequest("POST", urls, nil) //后端需要用ShouldBindJSON来接收参数
	if err != nil {
		fmt.Println("创建 HTTP 请求报错：" + err.Error())
		return errors.New("创建 HTTP 请求报错：" + err.Error())
	}
	fmt.Println("发送：", req)

	// 发送 HTTP 请求
	var resp *http.Response
	// 创建 HTTP 客户端
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		fmt.Println("发送 HTTP 请求报错：" + err.Error())
		return errors.New("发送 HTTP 请求报错：" + err.Error())
	}
	defer resp.Body.Close()

	fmt.Println("状态信息：", resp.Status)
	fmt.Println("长度为: ", resp.ContentLength)
	DisStr := strings.Split(resp.Header.Values("Content-Disposition")[0], "=")
	pcapname := strings.ReplaceAll(DisStr[1], "\"", "")

	if resp.StatusCode == 200 {
		//设置响应头，告诉浏览器这是一个要下载的文件
		cont.Header("Content-Type", "application/vnd.tcpdump.pcap")
		cont.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", pcapname))
		cont.Header("Content-Transfer-Encoding", "binary")
		cont.Header("Content-Length", fmt.Sprintf("%d", resp.ContentLength))

		//将管道中的数据写入响应体
		n, err := io.Copy(cont.Writer, resp.Body)
		fmt.Println("写入字节：", n)
		if err != nil {
			fmt.Println("写入流失败：" + err.Error())
			return errors.New("写入流失败：" + err.Error())
		}
		return nil
	} else {
		return errors.New("停止抓包失败")
	}
}

// 获取网卡列表
func (p *packet) GetAllInterface(cont *gin.Context, clusterName, url string) (interface{}, error) {
	//根据集群名获取IP
	clu, err := dao.RegCluster.GetClusterIP(clusterName)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	urls := "http://" + clu.Ipaddr + ":" + clu.Port + "/api/interfaces?url=" + url
	req, err := http.NewRequest("POST", urls, nil) //后端需要用ShouldBindJSON来接收参数
	if err != nil {
		fmt.Println("创建 HTTP 请求报错：" + err.Error())
		return nil, errors.New("创建 HTTP 请求报错：" + err.Error())
	}
	fmt.Println("发送：", req)

	// 发送 HTTP 请求
	var resp *http.Response
	// 创建 HTTP 客户端
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		fmt.Println("发送 HTTP 请求报错：" + err.Error())
		return nil, errors.New("发送 HTTP 请求报错：" + err.Error())
	}
	defer resp.Body.Close()

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
		return data, nil
	} else {
		return data["err"], errors.New("err")
	}
}

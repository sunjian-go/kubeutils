package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"main/dao"
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
func (p *packet) StartPacket(pcakinfo *PackInfo, clusterName, url string) error {
	//根据集群名获取IP
	clu, err := dao.RegCluster.GetClusterIP(clusterName)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// 将结构体编码为 JSON
	packData, err := json.Marshal(pcakinfo)
	if err != nil {
		fmt.Println("编码结构体为 JSON 时出错：" + err.Error())
		return errors.New("编码结构体为 JSON 时出错：" + err.Error())
	}

	// 创建一个包含 JSON 数据的 io.Reader
	jsonReader := bytes.NewBuffer(packData)
	//创建http请求
	urls := "http://" + clu.Ipaddr + ":" + clu.Port + "/api/startPacket?url=" + url //去往agent
	req, err := http.NewRequest("POST", urls, jsonReader)                           //后端需要用ShouldBindJSON来接收参数
	if err != nil {
		fmt.Println("创建 HTTP 请求报错：" + err.Error())
		return errors.New("创建 HTTP 请求报错：" + err.Error())
	}

	fmt.Println("发送：", req)

	// 发送 HTTP 请求
	// 创建 HTTP 客户端
	var resp *http.Response
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		fmt.Println("发送 HTTP 请求报错：" + err.Error())
		return errors.New("发送 HTTP 请求报错：" + err.Error())
	}
	defer resp.Body.Close()

	fmt.Println("状态信息：", resp.Status)
	if resp.StatusCode != 200 {
		if resp.StatusCode == 442 {
			return errors.New("当前已有抓包程序运行，请稍后重试")
		}
		return errors.New("启动抓包程序失败，请检查抓包程序")
	}
	return nil
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

		// 将管道中的数据写入响应体
		n, err := io.Copy(cont.Writer, resp.Body)
		fmt.Println("写入字节：", n)
		if err != nil {
			fmt.Println("写入流失败：" + err.Error())
			return errors.New("写入流失败：" + err.Error())
		}
	} else {
		return errors.New("关闭抓包程序失败，请检查抓包程序")
	}
	return nil
}

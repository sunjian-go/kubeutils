package utils

import (
	"bytes"
	"errors"
	"io"
	"net/http"
)

var Http httpTools

type httpTools struct {
}

func (h *httpTools) HttpSend(url string, jsondata *bytes.Buffer, method string) ([]byte, int, error) {
	//创建http请求
	var req *http.Request
	var err error
	if jsondata == nil {
		req, err = http.NewRequest(method, url, nil) //后端需要用ShouldBindJSON来接收参数
	} else {
		req, err = http.NewRequest(method, url, jsondata) //后端需要用ShouldBindJSON来接收参数
	}

	if err != nil {
		Logg.Error("创建 HTTP 请求报错：" + err.Error())
		return nil, 0, errors.New("创建 HTTP 请求报错：" + err.Error())
	}
	//fmt.Println("发送：", req)

	// 创建 HTTP 客户端
	var resp *http.Response
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		Logg.Error("发送 HTTP 请求报错：" + err.Error())
		return nil, 0, errors.New("发送 HTTP 请求报错，请检查后端agent服务是否正常运行")
	}
	defer resp.Body.Close()
	//fmt.Println("状态信息：", resp.Status)

	// 读取响应的 body 内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		Logg.Error("读取响应 body 时出错:" + err.Error())
		return nil, 0, errors.New("读取响应 body 时出错:" + err.Error())
	}
	return body, resp.StatusCode, nil
}

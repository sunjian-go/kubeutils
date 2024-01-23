package service

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"main/dao"
	"main/model"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var File file

type file struct {
}

type QueryStr struct {
	PodName       string `json:"podName"`
	Namespace     string `json:"namespace"`
	ContainerName string `json:"containerName"`
	Path          string `json:"path"`
	ClusterName   string `json:"clusterName"`
}

type PodInFo struct {
	PodName       string `json:"podName"`
	Namespace     string `json:"namespace"`
	ContainerName string `json:"containerName"`
	Path          string `json:"path"`
}

// 拷贝文件到pod转发方法
func (f *file) CopyToPod(files []*multipart.FileHeader, podinfo *QueryStr) error {
	// 用于保存二进制文件流的缓冲区
	var buf bytes.Buffer

	// 创建multipart写入器,将multipart/form-data数据写入到buf缓冲区中
	writer := multipart.NewWriter(&buf)

	//创建要插入数据的表结构体
	uploadData := new(model.Upload_History)
	for _, file := range files {
		//将城市对应的用户对应的文件存入表中
		fmt.Println("文件名：", file.Filename)
		uploadData.File = file.Filename
		// 打开文件
		srcFile, err := file.Open()
		if err != nil {
			fmt.Println("打开文件失败：", err.Error())
		}
		defer srcFile.Close()

		// 创建multipart表单字段
		part, err := writer.CreateFormFile("file", file.Filename)
		if err != nil {
			fmt.Println("创建multipart表单字段失败：", err.Error())
		}

		// 将文件内容复制到multipart表单字段中
		if _, err := io.Copy(part, srcFile); err != nil {
			fmt.Println("将文件内容复制到multipart表单字段中报错：", err.Error())
		}
	}
	fmt.Println("上传pod: ", podinfo.PodName)
	fmt.Println("namespace：", podinfo.Namespace)
	fmt.Println("上传路径：", podinfo.Path)
	fmt.Println("集群名：", podinfo.ClusterName)

	//根据时间戳和文件名生成md5值
	// 获取当前时间
	currentTime := time.Now()
	// 获取当前时间的秒数
	seconds := currentTime.Unix()
	// 使用 strconv.FormatInt 进行转换
	secondsStr := strconv.FormatInt(seconds, 10)
	//获取md5值
	md5str := createMd5(secondsStr, uploadData.File)

	//组装表数据
	uploadData.ClusterName = podinfo.ClusterName
	uploadData.PodName = podinfo.PodName
	uploadData.Namespace = podinfo.Namespace
	uploadData.Path = podinfo.Path
	uploadData.Status = "pending"
	uploadData.Code = md5str

	//开始插入表数据
	err := dao.Uploadhist.UploadData(uploadData)
	if err != nil {
		fmt.Println(err.Error())
	}

	// 添加其他文本字段（如果有的话）
	//writer.WriteField("key", "value")

	// 关闭写入器
	writer.Close()

	// 创建一个 url.Values 对象
	params := url.Values{}
	params.Add("podName", podinfo.PodName)
	params.Add("namespace", podinfo.Namespace)
	params.Add("containerName", podinfo.ContainerName)
	params.Add("path", podinfo.Path)

	//根据集群名获取IP
	clu, err := dao.RegCluster.GetClusterIP(podinfo.ClusterName)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("ip= ", clu.Ipaddr)
	// 构建 URL
	baseURL := "http://" + clu.Ipaddr + ":" + clu.Port + "/api/upload"
	fullURL := baseURL + "?" + params.Encode()

	// 创建HTTP请求
	req, err := http.NewRequest("POST", fullURL, &buf)
	if err != nil {
		fmt.Println("创建HTTP请求失败：", err.Error())
	}

	// 设置请求头，指定Content-Type为multipart/form-data
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 发送HTTP请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("发送HTTP请求失败：", err.Error())
	}
	defer resp.Body.Close()

	//到这一步就说明上传成功了，修改状态即可
	err = dao.Uploadhist.UpdateUploadDataStatus(md5str, "success")
	if err != nil {
		fmt.Println("更新文件上传状态失败：", err.Error())
	}
	return nil
}

// 生成md5哈希值
func createMd5(time, file string) string {
	// 创建 MD5 哈希对象
	hash := md5.New()

	// 将字符串转换为字节数组并写入哈希对象
	hash.Write([]byte(time + file))

	// 计算 MD5 哈希值
	hashInBytes := hash.Sum(nil)

	// 将哈希值转换为十六进制字符串
	md5String := hex.EncodeToString(hashInBytes)

	return md5String
}

// 从pod拷贝文件转发方法
func (f *file) CopyFromPod(podinfo *QueryStr, cont *gin.Context, clusterName string) error {
	//根据集群名获取IP
	clu, err := dao.RegCluster.GetClusterIP(clusterName)
	if err != nil {
		fmt.Println(err.Error())
	}

	// 将结构体编码为 JSON
	podData, err := json.Marshal(podinfo)
	if err != nil {
		fmt.Println("编码结构体为 JSON 时出错：" + err.Error())
		return errors.New("编码结构体为 JSON 时出错：" + err.Error())
	}

	// 创建一个包含 JSON 数据的 io.Reader
	jsonReader := bytes.NewBuffer(podData)
	//创建http请求
	urls := "http://" + clu.Ipaddr + ":" + clu.Port + "/api/download"
	req, err := http.NewRequest("POST", urls, jsonReader) //后端需要用ShouldBindJSON来接收参数
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
	if resp.Status == "200 OK" {
		//设置响应头，告诉浏览器这是一个要下载的文件
		respFilename := ""
		if strings.Contains(podinfo.Path, ".") {
			respFilename = strings.Split(podinfo.Path, ".")[0]
		} else {
			respFilename = podinfo.Path
		}
		cont.Header("Content-Disposition", "attachment; filename="+respFilename+".tar")
		cont.Header("Content-Type", "application/octet-stream")

		// 将管道中的数据写入响应体
		n, err := io.Copy(cont.Writer, resp.Body)
		fmt.Println("写入字节：", n)
		if err != nil {
			fmt.Println("写入流失败：" + err.Error())
			return errors.New("写入流失败：" + err.Error())
		}
	}
	return nil
}

type UploadInfo struct {
	CLusterName string `form:"clusterName"`
	FilterName  string `form:"filterName"`
	Page        int    `form:"page"`
	Limit       int    `form:"limit"`
}

// 获取当前集群所有上传记录
func (f *file) GetUploadHistory(uploadinfo *UploadInfo) ([]model.Upload_History, int, error) {
	uploadh, total, err := dao.Uploadhist.GetUploadHistory(uploadinfo.CLusterName, uploadinfo.FilterName, uploadinfo.Page, uploadinfo.Limit)
	if err != nil {
		fmt.Println(err.Error())
		return nil, 0, err
	}
	return uploadh, total, nil
}

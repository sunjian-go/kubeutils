package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"main/service"
	"net/http"
)

var File file

type file struct {
}

// 上传文件
func (f *file) UploadFile(c *gin.Context) {
	podinfo := new(service.QueryStr)
	podinfo.PodName = c.Query("podName")
	podinfo.Namespace = c.Query("namespace")
	podinfo.ContainerName = c.Query("containerName")
	podinfo.Path = c.Query("path")
	podinfo.ClusterName = c.Query("clusterName")
	//fmt.Println("aaaaaaaaaaaaaaaaa", podinfo)
	if podinfo.PodName == "" || podinfo.Namespace == "" || podinfo.Path == "" {
		fmt.Println("pod信息不完善，请设置完再上传")
		c.JSON(400, gin.H{
			"err": "pod信息不完善，请设置完再上传",
		})
		return
	}

	// 从请求中获取文件
	formFiles, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, "获取上传文件失败: %v", err)
		return
	}
	files := formFiles.File["file"]
	err = service.File.CopyToPod(files, podinfo)
	if err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"msg": "上传文件成功",
	})
}

// 下载文件
func (f *file) DownLoadFile(c *gin.Context) {
	clusterName := c.Query("clusterName")
	podinfo := new(service.QueryStr)
	if err := c.ShouldBind(podinfo); err != nil {
		c.JSON(400, gin.H{
			"err": "数据绑定失败：" + err.Error(),
		})
		return
	}
	fmt.Println("获取数据为：", podinfo)
	err := service.File.CopyFromPod(podinfo, c, clusterName)
	if err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.Status(200)
}

// 获取所有上传记录信息
func (f *file) GetUploadHistory(c *gin.Context) {
	cluInfo := new(service.UploadInfo)
	if err := c.Bind(cluInfo); err != nil {
		c.JSON(400, gin.H{
			"err": "绑定数据失败：" + err.Error(),
		})
		return
	}
	uploadh, total, err := service.File.GetUploadHistory(cluInfo)
	if err != nil {
		c.JSON(400, gin.H{
			"err":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":   "获取上传记录成功",
		"data":  uploadh,
		"total": total,
	})
}

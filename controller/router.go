package controller

import (
	"github.com/gin-gonic/gin"
)

var Router router

type router struct {
}

// 初始化路由
func (r *router) RouterInit(router *gin.Engine) {
	router.POST("/api/register", Clus.RegisterFunc).
		GET("/api/getAllClus", Clus.GetAllClusters).
		DELETE("/api/delClus", Clus.DeleteCluster).
		POST("/api/keepalive", Keepalive.KeepaliveFunc).
		GET("/api/getIP", Ipaddr.GetClusterIP).
		POST("/api/upload", File.UploadFile).
		POST("/api/login", Login.Login).
		GET("/api/corev1/getnamespaces", Namespace.GetNamespaces).
		GET("/api/corev1/getpods", Pod.GetPods).
		GET("/api/corev1/getnodes", Node.GetNodes).
		GET("/api/corev1/getcontainers", Pod.GetContainer).
		GET("/api/listPath", Listpath.ListContainerPath).
		POST("/api/download", File.DownLoadFile).
		GET("/api/uploadHistory", File.GetUploadHistory).
		GET("/api/importfile", Imfile.ImportFile).
		//ws获取日志
		GET("/api/ws", Pod.WsFunc)
}

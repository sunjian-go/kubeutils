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
		GET("/api/authformula", Login.GetAuthCode).
		GET("/api/corev1/getnamespaces", Namespace.GetNamespaces).
		GET("/api/corev1/getpods", Pod.GetPods).
		GET("/api/corev1/getnodes", Node.GetNodes).
		GET("/api/corev1/getnodedetail", Node.GetNodeDetail).
		GET("/api/corev1/getcontainers", Pod.GetContainer).
		GET("/api/listPath", Listpath.ListContainerPath).
		POST("/api/download", File.DownLoadFile).
		GET("/api/uploadHistory", File.GetUploadHistory).
		POST("/api/startPacket", Pack.StartPacket).
		POST("/api/stopPacket", Pack.StopPacket).
		GET("/api/interfaces", Pack.GetAllInterface).
		POST("/api/icmp", Icmp.PingFunc).
		POST("/api/port", Port.PortTel).
		//角色
		GET("/api/getrole", Role.GetRoles).
		GET("/api/getonerole", Role.GetRole).
		POST("/api/addrole", Role.AddRole).
		PUT("/api/updaterole", Role.UpdateRoleAuth).
		DELETE("/api/deleterole", Role.DeleteRole).
		//用户组
		GET("/api/getgroup", Group.Getgroups).
		GET("/api/getonegroup", Group.Getgroup).
		POST("/api/addgroup", Group.Addgroup).
		PUT("/api/updategroup", Group.UpdategroupAuth).
		DELETE("/api/deletegroup", Group.Deletegroup).
		//用户
		GET("/api/getuser", Users.GetUsers).
		GET("/api/getoneuser", Users.GetUser).
		POST("/api/adduser", Users.AddUser).
		PUT("/api/updateuser", Users.UpdateUser).
		DELETE("/api/deleteuser", Users.DeleteUser).
		//ws获取日志
		GET("/api/ws", Pod.WsFunc)
}

package middle

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		//获取请求方法
		method := c.Request.Method
		//添加跨域响应头
		c.Header("Access-Control-Allow-Origin", "*")                                       // 设置允许跨域请求的源，使用*表示允许所有域名访问
		c.Header("Content-Type", "application/json, text/plain, */*, charset=UTF-8")       // 设置响应头部的Content-Type为application/json
		c.Header("Access-Control-Max-Age", "86400")                                        // 设置OPTIONS预检请求的缓存时间，单位为秒
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE") // 设置允许的请求方法
		c.Header("Access-Control-Allow-Headers", "X-Token, Content-Type, ContentLength,"+
			"Accept-Encoding, X-CSRF-Token, Authorization, X-Max") // 设置允许的请求头
		c.Header("Access-Control-Allow-Credentials", "false") // 设置是否允许发送Cookie
		//放行所有OPTIONS方法: 检查如果是OPTIONS预检请求，直接返回状态码204，表示服务器允许该请求
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)

		}
		//处理请求
		c.Next()

	}
}

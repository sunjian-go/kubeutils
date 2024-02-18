package middle

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"main/utils"
	"strings"
)

var Jwt jwt

type jwt struct {
}

var Token string

func (j *jwt) GetToken() string {
	return Token
}

// JWTAuth 中间件，检查token
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		//对登录接口放行,如果获取的URL是登陆路由就直接放行
		fmt.Println("获取的URL为：", c.Request.URL.String())
		if c.Request.URL.String() == "/api/login" ||
			strings.Contains(c.Request.URL.String(), "register") ||
			strings.Contains(c.Request.URL.String(), "keepalive") ||
			strings.Contains(c.Request.URL.String(), "upload") ||
			strings.Contains(c.Request.URL.String(), "authformula") ||
			strings.Contains(c.Request.URL.String(), "ws") {
			c.Next()
		} else {
			//如果不是登录路由就需要进行token验证
			//1.获取Header中的Authorization
			Token = c.Request.Header.Get("Authorization")
			utils.Logg.Info("获取到token为：" + Token)
			if Token == "" {
				utils.Logg.Info("请求未携带token，无权限访问")
				c.JSON(400, gin.H{
					"msg":  "请求未携带token，无权限访问",
					"data": nil,
				})
				c.Abort() //中止当前请求的执行并立即返回响应
				return
			}
			// 2.parseToken 解析token包含的信息
			claims, err := utils.JWTToken.ParseToken(Token)
			if err != nil {
				rdata := struct {
					Code int    `json:"code"`
					Data string `json:"data"`
				}{}
				rdata.Code = 10086
				rdata.Data = ""
				//token延期错误
				if err.Error() == "TokenExpired" {
					c.JSON(400, gin.H{
						"msg":  "授权已过期",
						"data": rdata,
					})
					c.Abort()
					return
				}
				//其他解析错误
				c.JSON(400, gin.H{
					"msg":  err.Error(),
					"data": rdata,
				})
				c.Abort()
				return
			}
			// 3.token验证没问题了就继续交由下一个路由处理,并将解析出的信息传递下去
			c.Set("claims", claims)
			c.Next()
		}
	}
}

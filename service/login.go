package service

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"main/dao"
	"main/utils"
	"time"
	//"google.golang.org/genproto/googleapis/ads/googleads/v3/errors"
)

var Login login

type login struct {
}
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Formula  string `json:"formula"`
	Result   string `json:"result"`
}

func (l *login) Login(adminuser *User) (string, map[string]string, error) {
	//首先验证码校验
	err := Auth.CountResult(adminuser.Formula, adminuser.Result)
	if err != nil {
		utils.Logg.Error(err.Error())
		return "", nil, err
	}
	//校验账密
	gconf, err := Conf.ReadConfFunc()
	if err != nil {
		return "", nil, err
	}

	//获取用户信息
	user, err := dao.User.GetUser(adminuser.Username)
	if err != nil {
		utils.Logg.Error(err.Error())
		return "", nil, err
	}

	if adminuser.Username != "" && adminuser.Password != "" {
		if adminuser.Username == "admin" {
			if adminuser.Username != gconf["AdminUser"] || adminuser.Password != gconf["AdminPasswd"] {
				utils.Logg.Error("账号或密码错误，请重试")
				return "", nil, errors.New("账号或密码错误，请重试")
			}
		} else {
			if adminuser.Username != user.Username || adminuser.Password != user.Password {
				utils.Logg.Error("账号或密码错误，请重试")
				return "", nil, errors.New("账号或密码错误，请重试")
			}
		}

	} else {
		utils.Logg.Error("账号或密码不能为空")
		return "", nil, errors.New("账号或密码不能为空")
	}

	//验证账密通过后，生成token
	// 定义加密因子
	secret := "sunjiandevops"
	// 创建一个新的Token对象
	token := jwt.New(jwt.SigningMethodHS256)
	// 设置Token的Claim(声明)，这是您自定义的数据
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Minute * 40).Unix() // 设置Token过期时间（40分钟）
	//claims["exp"] = time.Now().Add(time.Second * 5).Unix()
	claims["user_id"] = "1234567"
	claims["username"] = adminuser.Username

	// 使用加密因子进行签名，并获取最终的Token字符串
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		utils.Logg.Error("生成Token失败: " + err.Error())
		return "", nil, errors.New("生成Token失败: " + err.Error())
	}
	utils.Logg.Info("生成的Token: " + tokenString)
	kubeconf, err := Conf.ReadConfFunc()
	if err != nil {
		utils.Logg.Error(err.Error())
		return "", nil, errors.New("获取配置文件失败: " + err.Error())
	}
	return tokenString, kubeconf, nil
}

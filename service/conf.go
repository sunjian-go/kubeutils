package service

import (
	"github.com/Unknwon/goconfig"
	"main/utils"
	"os"
)

var Conf conf

type conf struct {
}

// 读取整个session配置
func (c *conf) ReadConfFunc() (map[string]string, error) {
	currentPath, _ := os.Getwd()
	confPath := currentPath + "/config/conf.ini"
	_, err := os.Stat(confPath)
	if err != nil {
		utils.Logg.Error("目标文件未找到")
		return nil, err
	}
	// 加载配置
	config, err := goconfig.LoadConfigFile(confPath)
	if err != nil {
		utils.Logg.Error("读取配置文件出错:" + err.Error())
		return nil, err
	}
	// 获取 section
	gconf, err := config.GetSection("server")
	if err != nil {
		utils.Logg.Error("获取配置文件内容失败：" + err.Error())
		return nil, err
	}
	return gconf, nil
}

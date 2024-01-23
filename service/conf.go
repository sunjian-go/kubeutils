package service

import (
	"github.com/Unknwon/goconfig"
	"github.com/wonderivan/logger"
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
		logger.Error("file is not found %s")
		return nil, err
	}
	// 加载配置
	config, err := goconfig.LoadConfigFile(confPath)
	if err != nil {
		logger.Error("读取配置文件出错:", err)
		return nil, err
	}
	// 获取 section
	gconf, err := config.GetSection("server")
	if err != nil {
		logger.Error("获取配置文件内容失败：", err.Error())
		return nil, err
	}
	return gconf, nil
}

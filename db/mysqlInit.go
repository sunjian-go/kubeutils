package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/wonderivan/logger"
	"main/model"
	"strconv"
	"time"
)

var GORM *gorm.DB

// 初始化数据库
func InitMysql(gconf map[string]string) {
	port, err := strconv.Atoi(gconf["dbPort"])
	if err != nil {
		fmt.Println("gconf[\"dbPort\"]转换失败:", err)
		return
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		gconf["dbUser"],
		gconf["dbPwd"],
		gconf["dbhost"],
		port,
		gconf["dbName"])
	//与数据库建立连接，生成一个*gorm.DB类型的对象
	GORM, err = gorm.Open("mysql", dsn)
	fmt.Println("数据库连接信息为：", dsn)
	if err != nil {
		fmt.Println("数据库连接失败: " + err.Error())
		return
	}

	//打印sql语句
	GORM.LogMode(true)
	//开启连接池
	//连接池最大允许的空闲连接数，如果没有sql任务需要执行的连接数大于20，超过的连接会被连接池关闭
	GORM.DB().SetMaxIdleConns(10)
	// 设置了并发连接数
	GORM.DB().SetMaxOpenConns(100)
	//设置了连接可复用的最大时间
	GORM.DB().SetConnMaxLifetime(time.Duration(30 * time.Second))

	logger.Info("连接数据库成功")
	//建表
	GORM.AutoMigrate(&model.Cluster{})        //集群表
	GORM.AutoMigrate(&model.Upload_History{}) //上传文件历史表
	GORM.AutoMigrate(&model.Role{})           //角色表
	GORM.AutoMigrate(&model.Group{})          //用户组表
	GORM.AutoMigrate(&model.User{})           //用户表

}

func Close() {
	GORM.Close()
}

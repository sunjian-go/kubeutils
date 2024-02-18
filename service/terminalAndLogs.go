package service

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"main/dao"
	"main/utils"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"
)

var Terminal terminal

type terminal struct {
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// 允许所有来源的WebSocket连接
		return true
	},
}

var (
	connFromSrv   = make(chan *websocket.Conn)
	connFromAgent = make(chan *websocket.Conn)
	exitConn      = make(chan bool)
)

func (t *terminal) WsHandler(namespace, podName, containerName, bashType, clusterName string, c *gin.Context) error {
	fmt.Println("接入客户端")

	// 将HTTP 请求升级为 WebSocket 连接
	serverWsConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		utils.Logg.Error("升级为 WebSocket 连接失败：" + err.Error())
		return err
	}

	//携程去建立agent链接(必须)
	go func() {
		t.ConAgentws(namespace, podName, containerName, bashType, clusterName)
	}()

	//将server套接字写入通道，为了使agent端可以使用server端套接字
	connFromSrv <- serverWsConn
	utils.Logg.Info("serverWsConn已写入")

	//应该主程序结束时关闭的
	defer func() {
		utils.Logg.Info("WsHandler方法结束")
		serverWsConn.Close()
	}()

	//先去检测agent端ws是否连接成功
	if !<-exitConn {
		utils.Logg.Error("agent端ws请求失败")
		return errors.New("agent端ws请求失败")
	}
	//获取agent套接字
	agentConn := <-connFromAgent
	utils.Logg.Info("agentConn已读取")

	// 设置连接关闭的回调函数
	serverWsConn.SetCloseHandler(func(code int, text string) error {
		fmt.Printf("ws server 连接关闭，状态码：%d，原因：%s\n", code, text)
		agentConn.Close()
		return nil
	})

	utils.Logg.Info("ws已连接。。。")

	for {
		//utils.Logg.Info("开始读取消息")
		// 读取客户端消息
		_, p, err := serverWsConn.ReadMessage()
		if err != nil {
			utils.Logg.Error(err.Error())
			return err
		}

		// 发送消息给agent端
		err = agentConn.WriteMessage(websocket.TextMessage, p)
		if err != nil {
			utils.Logg.Error("发送消息给agent端报错: " + err.Error())
			return err
		}
	}
}

// 连接agent
func (t *terminal) ConAgentws(namespace, podname, containerName, bash, clusterName string) {
	//首先获取server的套接字,用于下面给前端发消息
	srvConn := <-connFromSrv
	utils.Logg.Info("srvConn已读取")

	utils.Logg.Info("连接angent")
	// 创建一个用于接收系统中断信号的通道
	interrupt := make(chan os.Signal, 1)
	// 注册系统中断信号处理函数
	signal.Notify(interrupt, os.Interrupt)

	//根据集群名获取IP
	clu, err := dao.RegCluster.GetClusterIP(clusterName)
	if err != nil {
		utils.Logg.Error(err.Error())
	}

	// 定义agent端的 URL
	var agentUrl *url.URL
	if bash != "log" {
		//bash不为log说明是terminal连接
		//agentUrl, err = url.Parse("ws://" + clu.Ipaddr + ":" + clu.WsPort + "/ws")
		agentUrl, err = url.Parse("ws://" + clu.Ipaddr + ":" + clu.Port + "/api/terminal")
	} else {
		//bash为log说明是获取日志连接
		agentUrl, err = url.Parse("ws://" + clu.Ipaddr + ":" + clu.Port + "/api/corev1/getlog")
	}
	if err != nil {
		utils.Logg.Error("无法解析后端 URL: " + err.Error())
	}
	// 设置查询参数
	query := url.Values{}
	if bash != "log" {
		query.Set("namespace", namespace)
		query.Set("pod_name", podname)
		query.Set("container_name", containerName)
		query.Set("bashType", bash)
	} else {
		query.Set("namespace", namespace)
		query.Set("podname", podname)
		query.Set("container", containerName)
	}
	agentUrl.RawQuery = query.Encode()

	// 使用默认的 Dialer 发起 WebSocket 连接
	utils.Logg.Info("发起agent WebSocket 连接")
	agentWsConn, _, err := websocket.DefaultDialer.Dial(agentUrl.String(), nil)
	if err != nil {
		utils.Logg.Error("发起 WebSocket 连接失败: " + err.Error())
		exitConn <- false
		return
	}
	defer func() {
		utils.Logg.Info("ConAgentws方法关退出")
		agentWsConn.Close()
	}()
	//到这一步说明已经连接上agent的ws了
	exitConn <- true

	//将angent套接字写入通道，为了让server端可以使用agent套接字
	connFromAgent <- agentWsConn
	utils.Logg.Info("agentWsConn已写入")
	defer func() {
		utils.Logg.Info("ConAgentws方法结束")
	}()

	// 创建一个完成信号的通道
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			//一直读取agent端发过来的数据
			readNum, message, err := agentWsConn.ReadMessage()
			if err != nil {
				utils.Logg.Error("读取agent端数据失败: " + err.Error())
				return
			}

			//当读取到agent端的数据之后写入前端
			err = srvConn.WriteMessage(readNum, message)
			if err != nil {
				utils.Logg.Error("写入前端数据报错：" + err.Error())
			}
		}
	}()
	//循环保持该函数一直运行
	for {
		select {
		case <-done:
			utils.Logg.Info("agent ws关闭。。。")
			return
		case <-interrupt:
			utils.Logg.Info("终端信号")
			// 关闭连接
			err := agentWsConn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				utils.Logg.Error("写入关闭消息失败:" + err.Error())
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

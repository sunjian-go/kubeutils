# k8s实用性效率工具后端
### 说明：
此项目是一个可以在页面对pod以及node进行远程操作的平台后端
### 功能：
pod: 页面上传文件到pod内部，从pod内部下载文件到本地，终端、日志功能。

node: 远程抓包、icmp测试、端口测试（仅支持TCP），先决条件：需配合kubepacket插件完成

支持添加多个集群

# 构建可执行程序
go build -o kubeutils main.go

# Dockerfile
```
#使用 alpine 作为基础镜像
FROM alpine:latest
#创建所需目录
RUN mkdir -p /home/kubeUtils/config && mkdir /home/kubeUtils/yaml 
#将可执行程序复制到镜像中
COPY kubeutils /home/kubeUtils
COPY kubeagent.yaml /home/kubeUtils/yaml
RUN chmod +x /home/kubeUtils/kubeutils
#将配置文件复制到镜像中的 conf 目录
COPY conf.ini /home/kubeUtils/config/
#设置工作目录为 /home/kubeServer
WORKDIR /home/kubeUtils
#暴露端口
EXPOSE 8999
#指定容器启动时要运行的命令
CMD ["./kubeutils"]
```

# 配置文件conf.ini
```
[server]
AdminUser   = admin
AdminPasswd = admin123
dbhost = 1.1.1.1    #mysql地址
dbPort = 3306       #mysql端口
dbName = clusterinfo  #mysql库
dbUser = root       #mysql用户
dbPwd  = 123    #mysql密码
; 实际暴漏出去的ip
host = 2.2.2.2       #本服务实际暴漏出去的地址
; 实际暴露出去的端口   
port = 8999           #本服务实际暴漏出去的端口
```

# 部署方式
```
docker部署：
docker run -d --name kubeutils \
        -v conf/conf.ini:/home/kubeUtils/config/conf.ini \
        -p 8999:8999 \
        kubeutils:v1.0
```

# k8s部署：
```
根据docker部署修改
```

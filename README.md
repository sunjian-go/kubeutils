#构建可执行程序
go build -o kubeutils main.go

Dockerfile
# 使用 alpine 作为基础镜像
FROM alpine:latest
#创建所需目录
RUN mkdir -p /home/kubeUtils/config && mkdir /home/kubeUtils/yaml 
# 将可执行程序复制到镜像中
COPY kubeutils /home/kubeUtils
COPY kubeagent.yaml /home/kubeUtils/yaml
RUN chmod +x /home/kubeUtils/kubeutils
# 将配置文件复制到镜像中的 conf 目录
COPY conf.ini /home/kubeUtils/config/
# 设置工作目录为 /home/kubeServer
WORKDIR /home/kubeUtils
#暴露端口
EXPOSE 8999
# 指定容器启动时要运行的命令
CMD ["./kubeutils"]

配置文件conf.ini
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

部署方式
docker部署：
docker run -d --name kubeutils \
        -v conf/conf.ini:/home/kubeUtils/config/conf.ini \
        -p 8999:8999 \
        kubeutils:v1.0

k8s部署：
根据docker部署修改

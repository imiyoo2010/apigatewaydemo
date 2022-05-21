FROM centos:centos7 

#ADD会对目标文件进行覆盖
ADD ./gateway /opt/ 

#安装网络包(需要使用ifconfig和ping命令) 
RUN yum install net-tools -y


WORKDIR /opt

EXPOSE 8080

ENTRYPOINT ["./api-gateway_linux","-c","config.json"]

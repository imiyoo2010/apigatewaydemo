FROM centos:centos7

#ADD会对目标文件进行覆盖
ADD aliyun-mirror.repo /etc/yum.repos.d/CentOS-Base.repo
ADD aliyun-epel.repo /etc/yum.repos.d/epel.repo

ADD server/ /opt/ 

#安装网络包(需要使用ifconfig和ping命令)
RUN yum install net-tools -y

EXPOSE 8088

ENTRYPOINT ["/opt/api-server_linux"]

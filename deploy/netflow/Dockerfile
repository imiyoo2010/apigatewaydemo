FROM centos:centos7

#ADD会对目标文件进行覆盖
ADD aliyun-mirror.repo /etc/yum.repos.d/CentOS-Base.repo
ADD aliyun-epel.repo /etc/yum.repos.d/epel.repo

ADD entrypoint.sh /
RUN chmod +x /entrypoint.sh

#RUN yum upgrade ca-certificates --disablerepo=epel -y
RUN yum install nginx -y 

#安装网络包(需要使用ifconfig和ping命令)
RUN yum install net-tools -y


RUN yum install keepalived -y
#ADD keepalived-master.conf /etc/keepalived/keepalived.conf
#RUN systemctl start keepalived

ADD nginx_check.sh /opt/
ADD nginx-upstream.conf /etc/nginx/conf.d/nginx-upstream.conf
ADD nginx-location.conf /etc/nginx/default.d/nginx-location.conf


EXPOSE 80

ENTRYPOINT ["/entrypoint.sh"]


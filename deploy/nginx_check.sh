#! /bin/bash
#检测nginx是否启动了
A=`ps -C nginx -no-header | wc - 1`
if [ $A -eq 0];then    			#如果nginx没有启动就启动nginx 
    /usr/local/nginx/sbin/nginx    	#通过Nginx的启动脚本来重启nginx
    sleep 2
    if [`ps -C nginx --no-header| wc -1` -eq 0 ];then #如果nginx重启失败，则下面就会停掉keepalived服务，进行VIP转移
        killall keepalived
    fi
fi

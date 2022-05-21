#!/bin/bash
#############################################
## control.sh启动脚本, 至少实现start和stop两个方法
#############################################
workspace=$(cd $(dirname $0) && pwd -P)
cd ${workspace}
module=api-gateway
app=${module}
conf=config.json
logfile=var/app.log
pidfile=var/app.pid
## function
function start() {
    # 创建日志目录
    mkdir -p var &>/dev/null

    # 以后台方式 启动程序
    nohup ./${app} -c ${conf} >>${logfile} 2>&1 &

    local pid=$(get_pid)

    echo "${app} start ok, pid=${pid}"
    # 启动成功, 退出码为 0
    exit 0
}

function stop() {
    local pid=$(get_pid)
    # 停止该服务
    kill ${pid} &>/dev/null

    sleep 1
    # stop服务失败, 返回码为 非0
    echo "stop timeout(60s)"
    echo "${app} stop ok"
    exit 1
}
function status(){
    exit 0
}
function restart() {
  local pid=$(get_pid)
  kill -9 ${pid} &>/dev/null
  start
  exit 0
}

## internals
function get_pid() {
    real_pid=`ps -ef | grep api-gateway | grep -v grep | awk '{print $2}'`
    echo "${real_pid}"
}

action=$1
case ${action} in
    "start" )
        # 启动服务
        start
        ;;
    "stop" )
        # 停止服务
        stop
        ;;
    "status" )
        # 检查服务
        status
        ;;
    "restart" )
        #重启服务
        restart
        ;;
    * )
        echo "unknown command"
        exit 1
        ;;
esac
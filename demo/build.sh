#!/bin/bash

if [ $# -lt 1 ]; then
    echo "===== Usage: ./build.sh mac or ./build.sh linux ====="
    exit 1
fi

set -e
workspace=$(cd $(dirname $0) && pwd -P)

module_C="api-gateway"

module_S="api-server"

if [ "$module_C" == "" ] || [ "$module_S" == "" ]; then
    echo "===== please uncomment variable 'module' ====="
    exit 1
fi

#go build -o $module main.go    #编译目标文件
if [ $1 == "linux" ]; then
    make build-linux
fi

if [ $1 == "mac" ]; then
    make build-mac
fi

ret=$?
if [ $ret -ne 0 ];then
    echo "===== $module build failure ====="
    exit $ret
else
    echo -n "===== $module build successfully! ====="
fi

#将程序和脚本进行打包
output="output"
rm -rf $output
mkdir -p $output/gateway
mkdir -p $output/server

# 填充output目录, output的内容即为待部署内容
    (
        cp -f control.sh ${output}/gateway &&     # 拷贝部署脚本control.sh至output/gateway目录
        cp -rf config.json ${output}/gateway &&
        cp -rf storage ${output}/gateway &&
        cp -rf seelog.xml ${output}/gateway &&

        mv ${module_C} ${output}/gateway &&

        #测试文件
        cp -rf add_api_to_server.sh ${output}/ &&

        cp -rf abc.db ${output}/server &&
        cp -f control.sh ${output}/server &&     # 拷贝部署脚本control.sh至output/server目录
        mv ${module_S} ${output}/server &&        # 移动需要部署的文件到output目录下

        echo -e "===== Generate output ok ====="
    ) || { echo -e "===== Generate output failure ====="; exit 2; } # 填充output目录失败后, 退出码为 非0

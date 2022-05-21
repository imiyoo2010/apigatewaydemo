#!/usr/bin/env bash

#一、添加测试接口

##1、添加后端服务地址
curl http://127.0.0.1:8088/gateway/upstream -d '{
    "upstream_name":"测试",
    "upstream_sign":"test",
    "upstream_ip":"127.0.0.1:8080",
    "protocol":"http"
}'

##2、添加后端服务接口映射
curl http://127.0.0.1:8088/gateway/apimap -d '{
"name":"测试接口",
"method":"GET",
"protocol":"http",
"upstream_sign":"test",
"gate_path":"/map/ping",
"gate_params":"",
"back_path":"/ping",
"back_params":"",
"back_position":"querystring"}'


#二、添加外部接口

##1、添加后端服务地址
curl http://127.0.0.1:8088/gateway/upstream -d '{
    "upstream_name":"查询",
    "upstream_sign":"query",
    "upstream_ip":"httpbin.org",
    "protocol":"http"
}'

##2、添加后端服务接口映射
curl http://127.0.0.1:8088/gateway/apimap -d '{
"name":"IP查询",
"method":"GET",
"protocol":"http",
"upstream_sign":"query",
"gate_path":"/query/ip",
"gate_params":"",
"back_path":"/ip",
"back_params":"",
"back_position":"querystring"}'
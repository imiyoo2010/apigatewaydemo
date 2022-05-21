#!/bin/sh
keepalived -D
/usr/sbin/nginx -g "daemon off;"

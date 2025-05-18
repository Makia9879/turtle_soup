#!/bin/bash

# 定义程序路径和参数
PROGRAM="./turtle_soup_go"
ARGS="server --config etc/etc.yaml"
PID_FILE="turtle_soup.pid"

# 检查程序是否存在
if [ ! -f "$PROGRAM" ]; then
    echo "错误: 程序文件 $PROGRAM 不存在"
    exit 1
fi

# 检查配置文件是否存在
if [ ! -f "etc/etc.yaml" ]; then
    echo "错误: 配置文件 etc/etc.yaml 不存在"
    exit 1
fi

# 启动程序到后台
nohup $PROGRAM $ARGS > /dev/null 2>&1 &

# 保存PID到文件
echo $! > $PID_FILE
echo "程序已启动，PID: $!, 已保存到 $PID_FILE"
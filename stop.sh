#!/bin/bash

PID_FILE="turtle_soup.pid"

# 检查PID文件是否存在
if [ ! -f "$PID_FILE" ]; then
    echo "错误: PID文件 $PID_FILE 不存在"
    exit 0
fi

# 读取PID并停止进程
PID=$(cat $PID_FILE)
if ps -p $PID > /dev/null; then
    kill $PID
    echo "已停止进程(PID: $PID)"
    rm $PID_FILE
else
    echo "错误: 进程(PID: $PID)未运行"
    rm $PID_FILE
    exit 0
fi
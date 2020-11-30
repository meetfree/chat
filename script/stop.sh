#!/bin/bash

pid=$(ps -ef | grep 'im-service' | grep -v grep | awk '{print $2}')
if [ -z "$pid" ]; then
  echo "进程不存在"
else
  kill -9 "$pid"
  echo "进程 $pid 已结束"
fi
exit 0

#!/bin/bash

run(){
    chmod u+x ./nats-server
    nohup ./nats-server >> nats.log 2>&1 &
    chmod u+x ./im-service
    nohup ./im-service >> im.log 2>&1 &
}

# shellcheck disable=SC2028
echo "1.run"
# shellcheck disable=SC2028
echo "2.reload"
echo "3.kill"
# shellcheck disable=SC2162
read -p "请选择：" op
echo "$op"
# 获取进程号
if "$op"==1
#if [ -z "$pid" ]
  then run
  echo "已启动"
  elif "$op"==2
   then
pid=$(ps -ef | grep 'im-service' | grep -v grep | awk '{print $2}')
  kill -9 "$pid"
  run
  echo "应用已重启"
else
pid=$(ps -ef | grep 'im-service' | grep -v grep | awk '{print $2}')
  kill -9 "$pid"
  echo "进程 $pid 已结束"
fi
exit 0
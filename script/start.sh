#!/bin/bash

mode=$1
run() {
  chmod u+x ./nats-server
  nohup ./nats-server >>nats.log 2>&1 &
  chmod u+x ./im-service
  nohup ./im-service "$mode" >>im.log 2>&1 &
}

pid=$(ps -ef | grep 'im-service' | grep -v grep | awk '{print $2}')
if [ -z "$pid" ]; then
  run
  echo "应用已启动"
else
  kill -9 "$pid"
  run
  echo "应用已重启"
fi
exit 0

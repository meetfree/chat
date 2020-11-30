# 多房间聊天室
1. 消息发送采用接口发送的方式需要对接[nats](https://nats.io) golang的一个队列
2. 配置文件使用config.toml
3. 消息存储采用mongo db
4. websocket 握手时携带房间号
5. script 使用脚本打包和启动
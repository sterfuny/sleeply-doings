# sleeply-alive

一个轻量的设备在线状态同步工具，通过 WebSocket 保持长连接，使用 HTTP API 推送和拉取设备信息。

## 功能

- WebSocket 长连接，自动重连
- 设备注册与状态跟踪
- 通用键值对消息（`map[string]any`），自动识别数字、布尔、字符串
- SQLite 持久化存储设备信息
- 支持服务端/客户端两种模式
- 配置文件自动生成

## 配置文件

默认路径 `~/.config/alive/config.yaml`，首次运行自动生成：

```yaml
mode: server
port: 9180
id: 自动生成的 UUID
addrs:
  - ws://localhost:9180/ws
```

#### 服务端模式

```yaml
mode: server
port: 9180
```

服务端监听 WebSocket（/ws）和 HTTP 拉取接口（/pull）。

#### 客户端模式

```yaml
mode: client
port: 9191
id: 自动生成的 UUID
addrs:
  - ws://localhost:9180/ws
```

客户端连接 WebSocket 后自动注册，并监听 HTTP 推送接口。

## API

- push

*GET/POST* /push?any ...

参数任意 any 自动识别 int、float、bool、string

- pull

*GET* /pull?uuid=<设备UUID>

获取设备 JSON

*GET* /pull?uuid=<设备UUID>&setname=<名称>

返回结果并设置名称

***设备信息结构***

```json
{
  "name": "",
  "lastmsg": {
    "uuid": "",
    "sendkey": {}
  },
  "status": true
}
```

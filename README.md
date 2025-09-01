# Gateway Go Server

## 项目概述

这是一个基于Go语言的Gateway服务器项目，提供WebSocket连接管理、消息发布/订阅功能和HTTP代理服务。

## ⚠️ 配套客户端

**此服务器需要配合 TypeScript SDK 客户端使用！**

完整的 Gateway 解决方案包括：
1. **Gateway Go Server** (本项目): WebSocket 服务器，处理连接和消息路由
2. **Gateway TypeScript SDK**: 客户端 SDK，参考 [gateway-ts-sdk](https://github.com/langgexyz/gateway-ts-sdk)

```bash
# 典型使用流程
# 1. 启动服务器
make debug

# 2. 使用客户端 SDK (需要单独克隆 gateway-ts-sdk 仓库)
cd path/to/gateway-ts-sdk
npm install
npm run examples:node
```

## 项目结构

```
gateway-go-server/
├── bin/                          # 配置文件和二进制文件
│   ├── config.debug.json        # 调试配置
│   ├── config.json.default      # 默认配置模板
│   └── gateway-go-server         # 编译后的服务器程序
├── src/                         # Go服务端源码
│   ├── api/                     # API接口层
│   │   ├── subscribe.go         # 订阅接口
│   │   ├── publish.go           # 发布接口
│   │   ├── ping.go             # Ping接口
│   │   ├── unsubscribe.go      # 取消订阅接口
│   │   └── suite.go            # API套件
│   ├── main/
│   │   └── main.go             # 服务器入口
│   ├── proxy/
│   │   ├── proxy.go            # 代理实现
│   │   └── suite.go            # 代理套件
│   └── utils/
│       └── compression.go      # 工具函数
├── go.mod                      # Go模块定义
├── go.sum                      # Go依赖锁定
├── makefile                    # 构建脚本
└── README.md                   # 本文档
```

## 核心功能

- **WebSocket连接管理**: 支持大量并发WebSocket连接
- **订阅/发布系统**: 实现消息的发布和订阅功能
- **API接口**: 提供基于WebSocket的API接口
- **HTTP代理**: 支持HTTP请求的代理和转发
- **TypeScript SDK**: 提供完整的TypeScript客户端SDK
- **跨平台构建**: 支持Linux、Windows、macOS等多平台编译

## 快速开始

### 前置条件

由于这是一个独立的服务器项目，你还需要克隆配套的客户端 SDK：

```bash
# 克隆服务器仓库
git clone https://github.com/langgexyz/gateway-go-server.git

# 克隆客户端 SDK 仓库  
git clone https://github.com/langgexyz/gateway-ts-sdk.git
```

### 步骤 1: 启动 Gateway 服务器

```bash
# 编译服务端
make build

# 启动服务器（生产模式）
./bin/gateway-go-server

# 或者调试模式启动（推荐开发时使用）
make debug

# 使用自定义配置文件启动
./bin/gateway-go-server -config=bin/config.debug.json
```

服务器启动后会监听端口 `18443`（默认），等待客户端连接。

### 步骤 2: 使用客户端 SDK

服务器启动后，需要使用 TypeScript SDK 进行连接和测试：

```bash
# 切换到 SDK 目录 (需要单独克隆 gateway-ts-sdk 仓库)
cd path/to/gateway-ts-sdk

# 安装依赖
npm install

# 运行基础测试示例
npm run examples:node

# 运行完整功能测试（Hook + Proxy）
npm run examples:hook-proxy
```

### 步骤 3: 验证连接

确认以下功能正常工作：
1. ✅ WebSocket 连接到 `ws://localhost:18443`
2. ✅ 订阅/发布消息功能
3. ✅ Ping 连接测试
4. ✅ HTTP 代理功能

## API文档

### 服务端API

#### Subscribe API
- **路径**: `API/Subscribe`
- **方法**: WebSocket
- **功能**: 订阅消息推送
- **请求**:
```json
{
  "time": 1640995200000,
  "clientNo": 1
}
```
- **响应**:
```json
{
  "time": 1640995200000,
  "clientNo": 1
}
```

#### Publish API
- **路径**: `API/Publish`
- **方法**: WebSocket
- **功能**: 发布消息到所有订阅者
- **请求**:
```json
{
  "dataSize": 1024,
  "clientTime": 1640995200000
}
```
- **响应**:
```json
{
  "success": 95,
  "failed": 5,
  "totalTime": 120,
  "clientTime": 1640995200000
}
```

#### Ping API
- **路径**: `API/Ping`
- **方法**: WebSocket
- **功能**: 测试连接延迟
- **请求**: `{}`
- **响应**: `{}`

#### Unsubscribe API
- **路径**: `API/Unsubscribe`
- **方法**: WebSocket
- **功能**: 取消订阅消息推送
- **请求**:
```json
{
  "channel": "demo-channel"
}
```
- **响应**: `{}`

#### Proxy API
- **路径**: `API/Proxy`
- **方法**: WebSocket
- **功能**: HTTP请求代理
- **Headers**: 
  - `x-proxy-url`: 目标URL
  - `x-proxy-method`: HTTP方法 (GET/POST/PUT/DELETE)
- **请求**: HTTP请求体内容
- **响应**: 代理的HTTP响应内容

### TypeScript SDK

详见 [`gateway-ts-sdk/README.md`](gateway-ts-sdk/README.md)

## 配置说明

### 服务端配置

配置文件示例（JSON格式）：
```json
{
  "server": {
    "host": "0.0.0.0",
    "port": 8080,
    "readTimeout": 30,
    "writeTimeout": 30
  },
  "websocket": {
    "maxConnections": 10000,
    "bufferSize": 1024,
    "enableCompression": true
  },
  "log": {
    "level": "info",
    "output": "stdout"
  }
}
```

## 性能特性

- **高并发**: 支持数千个并发WebSocket连接
- **低延迟**: 优化的消息路由和转发机制
- **内存效率**: 高效的连接池和内存管理
- **可扩展**: 模块化设计，易于扩展新功能
- **跨平台**: 支持Linux、Windows、macOS等多平台部署
- **HTTP代理**: 高性能的HTTP请求代理转发


## 开发指南

### 服务端开发

```bash
# 运行测试
make test

# 代码格式化
go fmt ./...

# 静态检查
go vet ./...

# 查看依赖
make dep

# 编译 (本地)
make build

# 编译 (Linux)
make build-linux

# 编译 (所有平台)
make build-all
```

### TypeScript SDK 开发

```bash
cd gateway-ts-sdk

# 安装依赖
npm install

# 构建SDK
npm run build

# 运行测试
npm test

# 发布测试
npm run release:dry
```

## 部署指南

### 生产环境部署

1. **编译服务端**:
```bash
make build
```

2. **配置文件**:
- 复制 `bin/config.json.default` 为 `bin/config.json`
- 根据生产环境修改配置

3. **启动服务**:
```bash
./gateway-go-server -config=bin/config.json
```

4. **系统服务**:
可以配置为systemd服务或Docker容器运行。

### 监控和日志

- **性能监控**: 通过API获取实时统计信息
- **日志记录**: 支持不同级别的日志输出
- **健康检查**: 内置健康检查端点

## 故障排除

### 常见问题

1. **连接失败**:
   - 检查服务器是否正常运行
   - 确认端口是否被占用
   - 检查防火墙设置

2. **性能问题**:
   - 调整连接池大小
   - 优化缓冲区配置
   - 监控内存使用

3. **SDK使用问题**:
   - 检查TypeScript版本兼容性
   - 确认WebSocket URL格式
   - 查看SDK示例代码

### 调试技巧

1. **启用调试日志**:
```bash
./gateway-go-server -debug
```

2. **使用SDK测试**:
TypeScript SDK提供了完整的示例和测试功能。

3. **性能分析**:
```bash
go tool pprof http://localhost:8080/debug/pprof/profile
```

## 贡献指南

### 代码规范

- **Go代码**: 遵循Go官方代码规范
- **TypeScript代码**: 使用ESLint和Prettier
- **提交信息**: 使用约定式提交格式

### 测试要求

- 新功能必须包含单元测试
- 确保测试覆盖率不低于80%
- 集成测试覆盖主要业务流程

### 提交流程

1. Fork项目
2. 创建特性分支
3. 提交代码和测试
4. 创建Pull Request
5. 代码审查和合并

## 许可证

MIT License

## 联系方式

如有问题或建议，请提交Issue或联系开发团队。

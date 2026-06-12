# simonStu

基于 Gin 的 Web 服务脚手架，集成 Redis。

## 项目结构

```
├── cmd/server/main.go        # 服务入口
├── config/config.go          # 配置（加载 .env + 环境变量）
├── internal/
│   ├── database/redis.go     # Redis 客户端初始化
│   ├── handler/              # HTTP handler
│   ├── middleware/           # 中间件（CORS 等）
│   └── router/               # 路由注册
├── .env                      # 本地配置（不提交）
├── .env.example              # 配置模板
└── go.mod
```

## 快速开始

```bash
# 复制配置模板，按需修改
cp .env.example .env

# 启动
go run ./cmd/server/
```

## 配置

通过 `.env` 文件或环境变量配置：

| 变量 | 默认值 | 说明 |
|---|---|---|
| `PORT` | `8080` | 监听端口 |
| `GIN_MODE` | `debug` | Gin 运行模式（debug/release/test） |
| `REDIS_HOST` | `127.0.0.1` | Redis 地址 |
| `REDIS_PORT` | `6379` | Redis 端口 |
| `REDIS_PASSWORD` | (空) | Redis 密码 |

`.env` 已加入 `.gitignore`，不会提交到仓库。

## API

| 方法 | 路径 | 说明 |
|---|---|---|
| `GET` | `/api/ping` | 健康检查 |
| `GET` | `/api/hello` | 示例接口 |
| `POST` | `/api/redis/set` | 写入 Redis key-value |
| `GET` | `/api/redis/get?key=xxx` | 读取 Redis key |

### Redis 示例

```bash
# 写入
curl -X POST http://localhost:8080/api/redis/set \
  -H "Content-Type: application/json" \
  -d '{"key":"name","value":"simon"}'

# 读取
curl http://localhost:8080/api/redis/get?key=name
```

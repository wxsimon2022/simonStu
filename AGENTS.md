# Repository Guidelines

## Project Structure & Module Organization

```
goStu/
├── cmd/server/            # Go 服务入口
├── cmd/example/           # 独立示例程序
├── config/                # 配置（.env + 环境变量）
├── internal/
│   ├── database/          # Redis / MySQL 连接初始化
│   ├── handler/           # HTTP handler（请求→响应）
│   ├── logger/            # 终端 + 文件日志
│   ├── middleware/        # JWT 认证 / 权限校验 / CORS
│   ├── model/             # GORM 数据模型，嵌入 BaseModel
│   ├── repository/        # 通用泛型仓储 BaseRepo[T]
│   ├── response/          # 统一响应格式 {code, data, message}
│   ├── router/            # Gin 路由注册
│   └── service/           # JWT / bcrypt 认证服务
├── frontend/              # Vue 3 + Element Plus 前端
│   └── src/
│       ├── views/         # 页面组件
│       ├── router/        # 前端路由（含登录守卫）
│       └── utils/         # 请求封装（自动带 JWT）
└── sql/                   # 数据库建表 + 初始数据 SQL
```

Go 模块名 `github.com/wxsimon8888/simonStu`，内部包按职责分层：`handler → service → repository → model`。

## Build, Test, and Development Commands

```bash
# 启动后端（开发模式）
GOCACHE=/tmp/gocache go run ./cmd/server/

# 编译后端
GOCACHE=/tmp/gocache go build ./cmd/server/

# 初始化数据库表和数据
mysql -u root -p <db> < sql/init.sql
mysql -u root -p <db> < sql/seed.sql

# 启动前端（开发模式，端口 3000）
cd frontend && npm run dev

# 构建前端
cd frontend && npm run build

# Docker 部署
docker compose build --no-cache
docker compose up -d
```

## Coding Style & Naming Conventions

- **Go**: `gofmt` 缩进（tab），标识符使用 CamelCase。导出符号大写，私有符号小写。error 值用 `errors.Is()` 而非 `==`。
- **Vue**: Composition API + `<script setup>`，组件名多单词 PascalCase，文件名 kebab-case（如 `user-list.vue`）。
- **SQL**: 表名以 `c_` 开头（如 `c_users`, `c_admin_roles`），字段统一包含 `id`, `is_delete`, `create_time`, `update_time`，不使用外键。
- **API 响应**: 统一 `{ code: 200, data: {...}, message: "success" }` 格式，通过 `response.Success()` / `response.Error()` 构建。
- **日志**: handler 内用 `logger.Infof(c, ...)` 记录操作日志，`logger.Errorf(c, ...)` 记录错误日志，错误日志附请求路径和关键参数。
- **认证**: 需登录的路由在 `auth` group 下，权限校验通过 `middleware.PermissionRequired("perm:name")` 声明。请求携带 `Authorization: Bearer <token>`。

## Testing Guidelines

当前项目尚未建立自动化测试。新增 handler 或 service 后建议添加：

- **单元测试** 放在与被测文件同包的文件中，命名为 `*_test.go`，使用标准 `testing` 包。
- **测试函数** 命名 `TestXxx`。
- 运行测试：`go test ./internal/...`。

## Commit & Pull Request Guidelines

提交信息使用中文，格式为 `<type>: <description>`：

```
add: 用户管理 CRUD 接口
fix: 权限列表返回空数据的 nil 指针问题
refactor: UserList 改用 BaseRepo 分页查询
chore: 更新依赖版本
```

- 一个提交专注一个改动，不混杂无关修改。
- 提交前运行 `go build ./cmd/server/` 和 `npm run build` 确保编译通过。
- PR 描述写明改动内容、影响范围，涉及 UI 时附截图。

## Security & Configuration Tips

- `.env` 文件包含数据库密码和 JWT 密钥，已加入 `.gitignore`，不提交到仓库。参考 `.env.example` 配置。
- 密码使用 bcrypt 哈希存储，不在日志或响应中输出明文。
- JWT 令牌生成时同时存入 Redis，退出登录时从 Redis 撤销，使令牌立即失效。
- 管理员账号通过 `is_admin: true` 跳过所有权限校验，确保该字段不被外部篡改。
- `Response` 结构体中 `json:"-"` 标记的字段（如 `password_hash`、`is_delete`）不会出现在 API 响应中。

# simonStu 项目技能

## 项目概述

Gin + Vue 3 + Element Plus 全栈后台管理模板，集成 RBAC 权限管理、Redis Token 存储、JWT 认证。

## 项目结构

```
goStu/
├── cmd/server/main.go       # 入口：初始化配置 → 日志 → Redis/MySQL → 认证 → HTTP
├── config/config.go         # 配置：.env + 环境变量，getEnv/getEnvInt 封装
├── internal/
│   ├── database/            # MySQL (database.DB) + Redis (database.RedisClient) 初始化
│   ├── handler/             # HTTP handler：验证参数 → 调用下层 → 返回统一响应
│   ├── service/             # 认证服务：JWT + bcrypt + Redis Token 存储
│   ├── repository/          # 泛型仓储 BaseRepo[T]，所有 CRUD 走这里
│   ├── middleware/          # JWT 认证 + 权限校验 + CORS
│   ├── model/               # GORM 数据模型，嵌入 BaseModel
│   ├── response/            # 统一响应格式 {code, data, message}
│   ├── logger/              # 终端 + 文件日志，按天拆分
│   └── router/              # 路由注册，按功能拆分到独立文件
├── frontend/                # Vue 3 + Element Plus + Vite
│   └── src/
│       ├── views/           # 页面组件
│       ├── router/          # 前端路由 + 登录守卫
│       └── utils/           # api 封装（自动带 JWT Token，401 自动跳登录）
├── sql/                     # 建表 SQL + 种子数据
├── Dockerfile + docker-compose.yml
└── .env / .env.example
```

## 核心模式

### 模型层 (model)

所有业务模型嵌入 `BaseModel`，自动获得 `id / is_delete / create_time / update_time`：

```go
type Xxx struct {
    BaseModel
    Name string `gorm:"column:name;type:varchar(64);not null" json:"name"`
}
func (Xxx) TableName() string { return "c_xxx" }
```

- 表名以 `c_` 开头（如 `c_users`、`c_admin`）
- `json:"-"` 标记敏感字段（如 `password_hash`、`is_delete`），确保不暴露到 API 响应
- `*int` 用于可空字段（如 `ParentID *int` → JSON 中为 `null` 或数字）

### 仓储层 (repository)

泛型 `BaseRepo[T any]`，所有业务模型复用同一套 CRUD：

```go
var XxxRepo = NewBaseRepo[model.Xxx](database.DB)

// 可用方法
XxxRepo.GetByID(id)     // → *model.Xxx, error（自动过滤 is_delete=0）
XxxRepo.List(page, size) // → []model.Xxx, int64 total, error
XxxRepo.Create(data)     // → error
XxxRepo.Update(id, map)  // → error（自动填充 update_time）
XxxRepo.Delete(id)       // → error（软删除，设置 is_delete=1）
```

- 所有查询自动追加 `WHERE is_delete = 0`
- `Update` 和 `Delete` 返回 `ErrNotFound`，handler 可用 `errors.Is(err, repository.ErrNotFound)` 判断

### Handler 层

四步走模式：

```go
func XxxList(c *gin.Context) {
    // 1. 校验参数（ShouldBindJSON / DefaultQuery）
    // 2. 调用下层（repository / database.DB）
    // 3. 处理错误 + 打日志（logger.Errorf）
    // 4. 返回统一响应（response.Success / response.Error）
}
```

- 所有 handler 使用 `response.Success(c, data)` / `response.Error(c, httpStatus, msg)`
- 返回格式统一：`{"code":200, "data":{...}, "message":"success"}`
- 对于列表接口，`data` 中包含 `list` 和 `total`

### 响应格式 (response)

```go
response.Success(c, gin.H{"list": items, "total": n})  // 200 + data
response.Error(c, 400, "参数错误")                        // 400 + data=nil
response.Error(c, 500, "服务器错误")                      // 500 + data=nil
```

### 路由组织 (router)

按功能拆分到独立文件，主入口 `router.go` 只负责组装：

```go
api := r.Group("/api")
{
    RegisterPublicRoutes(api)          // → public.go（公开接口）

    auth := api.Group("")
    auth.Use(middleware.AuthRequired())
    {
        // 普通业务接口
        auth.PUT("/user", middleware.PermissionRequired("perm:name"), handler.XxxUpdate)
        RegisterSystemRoutes(auth)     // → system.go（系统管理）
    }
}
```

- `RegisterXxxRoutes(auth *gin.RouterGroup)` 模式，每个子文件一个注册函数
- 权限通过 `middleware.PermissionRequired("perm:name")` 声明式控制

### 认证体系

```
请求 → AuthRequired 中间件
  ├─ 解析 JWT → 验证签名
  ├─ 查 Redis token → 确认未撤销
  └─ 注入 user_id / username / is_admin 到 context

PermissionRequired("perm:name") 中间件
  ├─ is_admin=true → 跳过（管理员通吃）
  └─ 否则联表查 c_admin_roles + c_role_permissions + c_permissions
```

- JWT 24 小时过期，签名算法 HS256
- 登录成功后 token 同步存入 Redis（key: `token:<full-jwt>`，TTL 24h）
- 退出登录时从 Redis 删除 token，使之立即失效
- Redis 不可用时静默降级，不阻塞登录

### 数据库规范

- 所有表：`id` + `is_delete` + `create_time` + `update_time`
- 禁止外键约束
- 使用软删除（`is_delete = 1`），查询默认过滤已删记录
- RBAC 五表：`c_admin` - `c_roles` - `c_permissions` - `c_admin_roles` - `c_role_permissions`
- 权限类型：`dir`(目录) / `menu`(菜单) / `btn`(按钮)，通过 `type` 字段区分

### 前端规范

- Vue 3 Composition API + `<script setup>`
- 所有 API 请求走 `api()` 函数（自动带 JWT Token，401 自动跳登录）
- Element Plus 中文版，`el-tree` / `el-tree-select` / `el-table` / `el-pagination` 为主要组件
- 路由守卫：未登录自动跳 `/login`，已登录访问 `/login` 自动跳首页
- 右键菜单 + 键盘快捷键（Enter 编辑 / Delete 删除）+ 搜索过滤

### 错误日志

- 所有 handler 在 `response.Error` 前调用 `logger.Errorf(c, "描述 key=%s err=%v", k, e)`
- 日志格式：`[ERROR] file.go:42: [GET /api/xxx] 描述 key=val err=xxx`
- 成功操作记 `logger.Infof`，失败记 `logger.Errorf`

## 安全配置

- `.env` 包含数据库密码和 JWT 密钥，已加入 `.gitignore`
- `password_hash` 仅存 bcrypt 哈希，不存明文
- 管理员 JWT 声明 `is_admin: true`，跳过所有权限校验（必须在服务端严格控制此字段）
- `json:"-"` 确保敏感字段不被序列化到 API 响应中

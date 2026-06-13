# simonStu

基于 Gin + Vue 3 + Element Plus 的全栈后台管理模板，集成 RBAC 权限管理、Redis Token 存储。

## 项目结构

```
├── cmd/
│   └── server/main.go          # 服务入口
├── config/config.go            # 配置（.env + 环境变量）
├── internal/
│   ├── database/               # MySQL + Redis 连接初始化
│   ├── handler/                # HTTP handler
│   ├── logger/                 # 终端 + 文件日志（按天拆分）
│   ├── middleware/             # JWT 认证 / 权限校验 / CORS
│   ├── model/                  # GORM 数据模型，嵌入 BaseModel
│   ├── repository/             # 通用泛型仓储 BaseRepo[T]
│   ├── response/               # 统一响应 {code, data, message}
│   ├── router/                 # 路由注册（按文件拆分）
│   └── service/                # JWT / bcrypt 认证服务
├── frontend/                   # Vue 3 + Element Plus 前端
│   └── src/
│       ├── views/              # 页面组件（Dashboard, Admin, Role, Permission, Profile）
│       ├── router/             # 前端路由 + 登录守卫
│       └── utils/              # 请求封装（自动带 JWT）
├── sql/                        # 数据库建表 + 初始数据 SQL
├── Dockerfile
├── docker-compose.yml
├── .env / .env.example
└── README.md
```

## 快速开始

```bash
# 1. 初始化数据库
mysql -u root -p <db> < sql/init.sql
mysql -u root -p <db> < sql/seed.sql

# 2. 启动后端
GOCACHE=/tmp/gocache go run ./cmd/server/

# 3. 启动前端（另一终端）
cd frontend && npm run dev

# 4. 浏览器打开 http://localhost:3000
#    默认管理员账号 admin / admin123
```

## 配置

通过 `.env` 文件配置，参考 `.env.example`。

| 变量 | 默认值 | 说明 |
|---|---|---|
| `PORT` | `8080` | HTTP 端口 |
| `GIN_MODE` | `debug` | Gin 运行模式 |
| `TZ` | `Asia/Shanghai` | 时区 |
| `LOG_DIR` | `storage/logs` | 日志目录 |
| `JWT_SECRET` | `simon-stu-secret-key` | JWT 签名密钥 |
| `DB_HOST` | `127.0.0.1` | MySQL 地址 |
| `DB_PORT` | `3306` | MySQL 端口 |
| `DB_USER` | `root` | MySQL 用户名 |
| `DB_PASSWORD` | (空) | MySQL 密码 |
| `DB_NAME` | `test` | MySQL 数据库 |
| `REDIS_HOST` | `127.0.0.1` | Redis 地址 |
| `REDIS_PORT` | `6379` | Redis 端口 |
| `REDIS_PASSWORD` | (空) | Redis 密码 |
| `REDIS_DB` | `0` | Redis 库号 |

## API 一览

### 认证（公开）

| 方法 | 路径 | 说明 |
|---|---|---|
| `POST` | `/api/auth/login` | 普通用户登录 |
| `POST` | `/api/auth/admin/login` | 管理员登录 |
| `POST` | `/api/auth/logout` | 退出登录（撤销 Redis 中的令牌） |
| `GET` | `/api/auth/userinfo` | 当前用户信息 + 角色 + 权限 |

### 个人中心（需 JWT 登录）

| 方法 | 路径 | 说明 |
|---|---|---|
| `GET` | `/api/auth/profile` | 获取当前管理员个人信息 |
| `PUT` | `/api/auth/profile` | 修改个人信息（姓名/手机/邮箱/密码） |

### 用户管理（需 `user:*` 权限）

| 方法 | 路径 | 权限 | 说明 |
|---|---|---|---|
| `PUT` | `/api/user` | `user:update` | 修改用户 |

### 系统管理（需对应权限）

| 方法 | 路径 | 权限 | 说明 |
|---|---|---|---|
| `GET` | `/api/admin` | `admin:list` | 管理员列表（含角色） |
| `POST` | `/api/admin` | `admin:create` | 新建管理员 |
| `PUT` | `/api/admin` | `admin:update` | 修改管理员（含角色） |
| `DELETE` | `/api/admin/:id` | `admin:delete` | 删除管理员 |
| `GET` | `/api/role` | `role:list` | 角色列表（含权限） |
| `POST` | `/api/role` | `role:create` | 新建角色 |
| `PUT` | `/api/role` | `role:update` | 修改角色 |
| `DELETE` | `/api/role/:id` | `role:delete` | 删除角色 |
| `GET` | `/api/permission` | `perm:list` | 权限树 + 平铺列表 |
| `POST` | `/api/permission` | `perm:create` | 新建权限/菜单 |
| `PUT` | `/api/permission` | `perm:update` | 修改权限/菜单 |
| `DELETE` | `/api/permission/:id` | `perm:delete` | 删除权限/菜单 |
| `PUT` | `/api/permission/reorder` | `perm:update` | 拖拽排序批量更新（id / parent_id / sort_order） |
| `GET` | `/api/menus` | — | 当前用户菜单树（侧边栏导航，基于角色过滤） |

### 其他

| 方法 | 路径 | 说明 |
|---|---|---|
| `GET` | `/api/hello` | 基础问候 |
| `GET` | `/api/ping` | 健康检查 |
| `POST` | `/api/redis/set` | 写入 Redis |
| `GET` | `/api/redis/get` | 读取 Redis |
| `POST` | `/api/redis/stock/deduct` | 原子扣减库存（Lua） |
| `POST` | `/api/concurrent/process` | 并发 Worker Pool 示例 |

## 前端功能

| 页面 | 路由 | 功能 |
|---|---|---|
| 仪表盘 | `/` | 系统概览 |
| 管理员管理 | `/admin` | 管理员 CRUD，关联角色 |
| 角色管理 | `/role` | 角色 CRUD，关联权限（树形勾选） |
| 菜单管理 | `/permission` | 菜单树拖拽排序、搜索高亮、层级连线、右键菜单 |
| 个人中心 | `/profile` | 个人信息维护、修改密码 |

### 菜单管理特色

- **拖拽排序**：直接拖拽树节点重新排列，自动保存排序号
- **搜索高亮**：实时过滤匹配节点，自动展开祖先，匹配字符黄色高亮
- **层级连线**：子节点左侧显示虚线缩进指示线，结构一目了然
- **快捷键**：回车编辑、Delete 删除、集群展开/折叠
- **三种类型**：目录、菜单、按钮，各有独立图标和标签色
- **图标选择**：集成 Element Plus 图标库，预览选择

## Docker 部署

```bash
docker compose build --no-cache
docker compose up -d
```

## 数据库规范

- 所有表含 `id` / `is_delete` / `create_time` / `update_time` 四个公共字段
- 不使用外键约束
- 表名以 `c_` 开头（如 `c_users`、`c_admin`）
- RBAC 五表：`c_admin` / `c_roles` / `c_permissions` / `c_admin_roles` / `c_role_permissions`
- 菜单排序：`c_permissions.sort_order` 字段，同级按此升序排列

## 技术栈

| 层 | 技术 |
|---|---|
| 后端框架 | Gin |
| ORM | GORM |
| 认证 | JWT + Redis Token 存储 |
| 密码 | bcrypt |
| 缓存 | Redis（Lua 脚本） |
| 前端 | Vue 3 + Vite |
| UI | Element Plus |
| 数据库 | MySQL 8.0 |
| 部署 | Docker / Docker Compose |

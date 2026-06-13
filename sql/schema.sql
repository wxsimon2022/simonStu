-- =========================================================
-- 完整数据库表结构（不含数据）
-- 执行: mysql -u root -p <db> < sql/schema.sql
-- 规范: 所有表含 id / is_delete / create_time / update_time，无外键
-- =========================================================

CREATE TABLE IF NOT EXISTS c_users (
  id            INT           AUTO_INCREMENT PRIMARY KEY,
  username      VARCHAR(64)   NOT NULL COMMENT '用户名',
  password_hash VARCHAR(255)  NOT NULL COMMENT 'bcrypt 密码哈希',
  is_admin      TINYINT(1)    NOT NULL DEFAULT 0 COMMENT '1=管理员 0=普通用户',
  is_delete     TINYINT       NOT NULL DEFAULT 0 COMMENT '0=正常 1=删除',
  create_time   DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time   DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY ix_c_users_username (username)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='普通用户表';

CREATE TABLE IF NOT EXISTS c_admin (
  id            INT           AUTO_INCREMENT PRIMARY KEY,
  username      VARCHAR(64)   NOT NULL COMMENT '登录用户名',
  password_hash VARCHAR(255)  NOT NULL COMMENT 'bcrypt 密码哈希',
  real_name     VARCHAR(64)   DEFAULT '' COMMENT '真实姓名',
  phone         VARCHAR(20)   DEFAULT '' COMMENT '手机号',
  email         VARCHAR(128)  DEFAULT '' COMMENT '邮箱',
  status        TINYINT       NOT NULL DEFAULT 1 COMMENT '1=启用 0=禁用',
  is_delete     TINYINT       NOT NULL DEFAULT 0 COMMENT '0=正常 1=删除',
  create_time   DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time   DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_admin_username (username)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='管理后台用户表';

CREATE TABLE IF NOT EXISTS c_roles (
  id          INT           AUTO_INCREMENT PRIMARY KEY,
  name        VARCHAR(64)   NOT NULL UNIQUE COMMENT '角色标识，如 admin / editor',
  description VARCHAR(255)  DEFAULT '' COMMENT '角色描述',
  status      TINYINT       NOT NULL DEFAULT 1 COMMENT '1=启用 0=禁用',
  is_delete   TINYINT       NOT NULL DEFAULT 0 COMMENT '0=正常 1=删除',
  create_time DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色表';

CREATE TABLE IF NOT EXISTS c_permissions (
  id          INT           AUTO_INCREMENT PRIMARY KEY,
  name        VARCHAR(128)  NOT NULL UNIQUE COMMENT '权限标识，如 user:list',
  description VARCHAR(255)  DEFAULT '' COMMENT '权限描述',
  type        VARCHAR(16)   NOT NULL DEFAULT 'menu' COMMENT 'dir=目录 menu=菜单 btn=按钮',
  icon        VARCHAR(64)   DEFAULT '' COMMENT 'Element Plus 图标名，如 User / Setting / Management',
  parent_id   INT           DEFAULT NULL COMMENT '父权限 ID',
  is_delete   TINYINT       NOT NULL DEFAULT 0 COMMENT '0=正常 1=删除',
  create_time DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='权限（菜单）表';

CREATE TABLE IF NOT EXISTS c_admin_roles (
  id          INT       AUTO_INCREMENT PRIMARY KEY,
  user_id     INT       NOT NULL COMMENT '管理员 ID',
  role_id     INT       NOT NULL COMMENT '角色 ID',
  is_delete   TINYINT   NOT NULL DEFAULT 0 COMMENT '0=正常 1=删除',
  create_time DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_admin_role (user_id, role_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='管理员-角色关联表';

CREATE TABLE IF NOT EXISTS c_role_permissions (
  id             INT       AUTO_INCREMENT PRIMARY KEY,
  role_id        INT       NOT NULL COMMENT '角色 ID',
  permission_id  INT       NOT NULL COMMENT '权限 ID',
  is_delete      TINYINT   NOT NULL DEFAULT 0 COMMENT '0=正常 1=删除',
  create_time    DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time    DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_role_perm (role_id, permission_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色-权限关联表';

-- =========================================================
-- RBAC 初始数据
-- =========================================================

-- 角色
INSERT INTO c_roles (name, description) VALUES
('admin', '超级管理员'),
('editor', '编辑员')
ON DUPLICATE KEY UPDATE description = VALUES(description);

-- 权限
INSERT INTO c_permissions (name, description) VALUES
('user:list',   '用户列表'),
('user:create', '创建用户'),
('user:update', '修改用户'),
('user:delete', '删除用户'),
('role:list',   '角色列表'),
('role:create', '创建角色'),
('role:update', '修改角色'),
('role:delete', '删除角色')
ON DUPLICATE KEY UPDATE description = VALUES(description);

-- admin 角色分配所有权限
INSERT INTO c_role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM c_roles r, c_permissions p WHERE r.name = 'admin'
ON DUPLICATE KEY UPDATE role_id = role_id;

-- editor 角色分配部分权限
INSERT INTO c_role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM c_roles r, c_permissions p
WHERE r.name = 'editor' AND p.name IN ('user:list', 'user:create', 'user:update')
ON DUPLICATE KEY UPDATE role_id = role_id;

-- 管理后台初始用户（密码在服务端用 bcrypt 生成后更新）
INSERT INTO c_admin (username, password_hash, real_name, status) VALUES
('admin', '$2a$10$gyRC9Lo9BxcKgDr9deqOu.mw2GzIpTNfQOkOPZCxzAJUBZ16LvtYy', '超级管理员', 1)
ON DUPLICATE KEY UPDATE real_name = VALUES(real_name);

-- 新增系统管理权限
INSERT IGNORE INTO c_permissions (name, description) VALUES
('admin:list',   '管理员列表'),
('admin:create', '创建管理员'),
('admin:update', '修改管理员'),
('admin:delete', '删除管理员'),
('perm:list',    '权限列表');

-- admin 角色获取所有新权限
INSERT IGNORE INTO c_role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM c_roles r, c_permissions p WHERE r.name = 'admin';

-- 新增权限管理权限
INSERT IGNORE INTO c_permissions (name, description) VALUES
('perm:create', '创建权限'),
('perm:update', '修改权限'),
('perm:delete', '删除权限');

-- admin 角色获取所有新权限
INSERT IGNORE INTO c_role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM c_roles r, c_permissions p WHERE r.name = 'admin';

-- =========================================================
-- 权限树：添加父级权限分组
-- =========================================================
INSERT IGNORE INTO c_permissions (name, description) VALUES
('user', '用户管理'),
('role', '角色管理'),
('admin', '管理员管理'),
('perm', '权限管理');

-- 设置子权限的 parent_id（使用临时表绕过 MySQL 限制）
UPDATE c_permissions SET parent_id = (SELECT id FROM (SELECT id FROM c_permissions WHERE name = 'user') AS tmp) WHERE name IN ('user:list','user:create','user:update','user:delete');
UPDATE c_permissions SET parent_id = (SELECT id FROM (SELECT id FROM c_permissions WHERE name = 'role') AS tmp) WHERE name IN ('role:list','role:create','role:update','role:delete');
UPDATE c_permissions SET parent_id = (SELECT id FROM (SELECT id FROM c_permissions WHERE name = 'admin') AS tmp) WHERE name IN ('admin:list','admin:create','admin:update','admin:delete');
UPDATE c_permissions SET parent_id = (SELECT id FROM (SELECT id FROM c_permissions WHERE name = 'perm') AS tmp) WHERE name LIKE 'perm:%';

-- admin 角色获取所有权限（含新加的父级权限）
INSERT IGNORE INTO c_role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM c_roles r, c_permissions p WHERE r.name = 'admin';

-- =========================================================
-- 权限（菜单）类型：dir=目录 menu=菜单 btn=按钮
-- =========================================================
SET @cnt = (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'c_permissions' AND COLUMN_NAME = 'type');
SET @q = IF(@cnt = 0, 'ALTER TABLE c_permissions ADD COLUMN type VARCHAR(16) NOT NULL DEFAULT ''menu'' COMMENT ''dir=目录 menu=菜单 btn=按钮'' AFTER description', 'SELECT 1 AS done');
PREPARE s FROM @q;
EXECUTE s;
DEALLOCATE PREPARE s;

UPDATE c_permissions SET type = 'dir' WHERE name IN ('user', 'role', 'admin', 'perm');
UPDATE c_permissions SET type = 'menu' WHERE name LIKE '%:list';
UPDATE c_permissions SET type = 'btn' WHERE name LIKE '%:create' OR name LIKE '%:update' OR name LIKE '%:delete';

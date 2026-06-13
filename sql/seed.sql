-- =========================================================
-- RBAC 初始数据
-- 可重复执行（INSERT IGNORE / ON DUPLICATE KEY UPDATE）
-- =========================================================

-- 角色
INSERT INTO c_roles (name, description, status) VALUES
('admin', '超级管理员', 1),
('editor', '编辑员', 1)
ON DUPLICATE KEY UPDATE description = VALUES(description);

-- =========================================================
-- 权限（菜单）
-- =========================================================

-- 父级目录（type = dir）
INSERT IGNORE INTO c_permissions (name, description) VALUES
('user',  '用户管理'),
('role',  '角色管理'),
('admin', '管理员管理'),
('perm',  '权限管理');

-- 子权限 - 菜单页（type = menu，命名规范 xxx:list）
INSERT IGNORE INTO c_permissions (name, description) VALUES
('user:list',   '用户列表'),
('role:list',   '角色列表'),
('admin:list',  '管理员列表'),
('perm:list',   '权限列表');

-- 子权限 - 按钮操作（type = btn，命名规范 xxx:create / update / delete）
INSERT IGNORE INTO c_permissions (name, description) VALUES
('user:create', '创建用户'),
('user:update', '修改用户'),
('user:delete', '删除用户'),
('role:create', '创建角色'),
('role:update', '修改角色'),
('role:delete', '删除角色'),
('admin:create', '创建管理员'),
('admin:update', '修改管理员'),
('admin:delete', '删除管理员'),
('perm:create', '创建权限'),
('perm:update', '修改权限'),
('perm:delete', '删除权限');

-- =========================================================
-- 设置父子关系
-- =========================================================
UPDATE c_permissions SET parent_id = (SELECT id FROM (SELECT id FROM c_permissions WHERE name = 'user') AS tmp)  WHERE name IN ('user:list','user:create','user:update','user:delete');
UPDATE c_permissions SET parent_id = (SELECT id FROM (SELECT id FROM c_permissions WHERE name = 'role') AS tmp)  WHERE name IN ('role:list','role:create','role:update','role:delete');
UPDATE c_permissions SET parent_id = (SELECT id FROM (SELECT id FROM c_permissions WHERE name = 'admin') AS tmp) WHERE name IN ('admin:list','admin:create','admin:update','admin:delete');
UPDATE c_permissions SET parent_id = (SELECT id FROM (SELECT id FROM c_permissions WHERE name = 'perm') AS tmp)  WHERE name LIKE 'perm:%' AND name != 'perm';

-- =========================================================
-- 设置类型（兼容旧表升级）
-- =========================================================
SET @col_exists = (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'c_permissions' AND COLUMN_NAME = 'type');
SET @ddl = IF(@col_exists = 0, 'ALTER TABLE c_permissions ADD COLUMN type VARCHAR(16) NOT NULL DEFAULT ''menu'' COMMENT ''dir=目录 menu=菜单 btn=按钮'' AFTER description', 'SELECT 1');
PREPARE s FROM @ddl; EXECUTE s; DEALLOCATE PREPARE s;

UPDATE c_permissions SET type = 'dir'  WHERE name IN ('user','role','admin','perm');
UPDATE c_permissions SET type = 'menu' WHERE name LIKE '%:list';
UPDATE c_permissions SET type = 'btn'  WHERE name LIKE '%:create' OR name LIKE '%:update' OR name LIKE '%:delete';

-- =========================================================
-- 角色 - 权限分配
-- =========================================================
INSERT IGNORE INTO c_role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM c_roles r, c_permissions p WHERE r.name = 'admin' AND p.is_delete = 0;

INSERT IGNORE INTO c_role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM c_roles r, c_permissions p
WHERE r.name = 'editor' AND p.name IN ('user:list','user:create','user:update');

-- =========================================================
-- 管理后台管理员
-- =========================================================
INSERT INTO c_admin (username, password_hash, real_name, status) VALUES
('admin', '$2a$10$gyRC9Lo9BxcKgDr9deqOu.mw2GzIpTNfQOkOPZCxzAJUBZ16LvtYy', '超级管理员', 1)
ON DUPLICATE KEY UPDATE password_hash = VALUES(password_hash), real_name = VALUES(real_name), status = VALUES(status);

INSERT IGNORE INTO c_admin_roles (user_id, role_id)
SELECT a.id, r.id FROM c_admin a, c_roles r
WHERE a.username = 'admin' AND r.name = 'admin';

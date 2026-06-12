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

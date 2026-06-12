-- =========================================================
-- 初始化管理员账号 + 角色 + 权限
-- 执行: mysql -u root -p <db> < sql/init_admin.sql
-- 默认密码: admin123（如需修改，换掉 @admin_password 的值）
-- =========================================================

-- bcrypt 哈希（对应明文 admin123）
SET @admin_password = '$2a$10$gyRC9Lo9BxcKgDr9deqOu.mw2GzIpTNfQOkOPZCxzAJUBZ16LvtYy';

-- 1. 确保 admin 角色存在
INSERT IGNORE INTO c_roles (name, description, status) VALUES ('admin', '超级管理员', 1);

-- 2. 创建或更新管理员账号
INSERT INTO c_admin (username, password_hash, real_name, status)
VALUES ('admin', @admin_password, '超级管理员', 1)
ON DUPLICATE KEY UPDATE
  password_hash = @admin_password,
  real_name     = '超级管理员',
  status        = 1;

-- 3. 将 admin 角色授予管理员
INSERT IGNORE INTO c_admin_roles (user_id, role_id)
SELECT a.id, r.id FROM c_admin a, c_roles r
WHERE a.username = 'admin' AND r.name = 'admin';

-- 4. admin 角色获取全部权限
INSERT IGNORE INTO c_role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM c_roles r, c_permissions p WHERE r.name = 'admin';

-- 5. admin 角色获取所有新增的权限（包括目录、菜单、按钮）
INSERT IGNORE INTO c_role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM c_roles r, c_permissions p
WHERE r.name = 'admin' AND p.is_delete = 0
ON DUPLICATE KEY UPDATE role_id = role_id;

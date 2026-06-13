-- =========================================================
-- RBAC 初始数据（角色、权限/菜单、管理员账号）
-- 可重复执行（INSERT IGNORE / ON DUPLICATE KEY UPDATE）
-- 默认密码: admin123（bcrypt 已预生成）
-- =========================================================

-- =========================================================
-- 1. 角色
-- =========================================================
INSERT INTO c_roles (name, description, status) VALUES
('admin',  '超级管理员', 1),
('editor', '编辑员', 1)
ON DUPLICATE KEY UPDATE description = VALUES(description);

-- =========================================================
-- 2. 权限（菜单） — 三级结构
--    第一级: type=dir （目录，侧边栏分组）
--    第二级: type=menu（菜单页，含路由）
--    第三级: type=btn （按钮操作，仅权限校验）
-- =========================================================

-- 2.1 目录（顶级分组）
INSERT IGNORE INTO c_permissions (name, description, type, icon) VALUES
('user',   '用户管理', 'dir', 'User'),
('system', '系统管理', 'dir', 'Setting');

-- 2.2 菜单页（二级节点，parent → dir）
--     命名规范: xxx:list，路由自动解析为 /{xxx}
INSERT INTO c_permissions (name, description, type, icon, parent_id) VALUES
('user:list',  '用户列表', 'menu', 'UserFilled',
  (SELECT id FROM (SELECT id FROM c_permissions WHERE name = 'user') AS tmp)),
('admin:list', '管理员管理', 'menu', 'Management',
  (SELECT id FROM (SELECT id FROM c_permissions WHERE name = 'system') AS tmp)),
('role:list',  '角色管理', 'menu', 'UserFilled',
  (SELECT id FROM (SELECT id FROM c_permissions WHERE name = 'system') AS tmp)),
('perm:list',  '菜单管理', 'menu', 'Menu',
  (SELECT id FROM (SELECT id FROM c_permissions WHERE name = 'system') AS tmp))
ON DUPLICATE KEY UPDATE description = VALUES(description), type = VALUES(type), icon = VALUES(icon), parent_id = VALUES(parent_id);

-- 2.3 按钮（三级节点，parent → 所属菜单）
--     命名规范: xxx:create / :update / :delete，无路由
INSERT INTO c_permissions (name, description, type, icon, parent_id) VALUES
('user:create',  '创建用户',   'btn', '',
  (SELECT id FROM (SELECT id FROM c_permissions WHERE name = 'user:list') AS tmp)),
('user:update',  '修改用户',   'btn', '',
  (SELECT id FROM (SELECT id FROM c_permissions WHERE name = 'user:list') AS tmp)),
('user:delete',  '删除用户',   'btn', '',
  (SELECT id FROM (SELECT id FROM c_permissions WHERE name = 'user:list') AS tmp)),
('admin:create', '创建管理员', 'btn', '',
  (SELECT id FROM (SELECT id FROM c_permissions WHERE name = 'admin:list') AS tmp)),
('admin:update', '修改管理员', 'btn', '',
  (SELECT id FROM (SELECT id FROM c_permissions WHERE name = 'admin:list') AS tmp)),
('admin:delete', '删除管理员', 'btn', '',
  (SELECT id FROM (SELECT id FROM c_permissions WHERE name = 'admin:list') AS tmp)),
('role:create',  '创建角色',   'btn', '',
  (SELECT id FROM (SELECT id FROM c_permissions WHERE name = 'role:list') AS tmp)),
('role:update',  '修改角色',   'btn', '',
  (SELECT id FROM (SELECT id FROM c_permissions WHERE name = 'role:list') AS tmp)),
('role:delete',  '删除角色',   'btn', '',
  (SELECT id FROM (SELECT id FROM c_permissions WHERE name = 'role:list') AS tmp)),
('perm:create',  '创建菜单',   'btn', '',
  (SELECT id FROM (SELECT id FROM c_permissions WHERE name = 'perm:list') AS tmp)),
('perm:update',  '修改菜单',   'btn', '',
  (SELECT id FROM (SELECT id FROM c_permissions WHERE name = 'perm:list') AS tmp)),
('perm:delete',  '删除菜单',   'btn', '',
  (SELECT id FROM (SELECT id FROM c_permissions WHERE name = 'perm:list') AS tmp))
ON DUPLICATE KEY UPDATE description = VALUES(description), type = VALUES(type), icon = VALUES(icon), parent_id = VALUES(parent_id);

-- =========================================================
-- 3. 角色-权限分配
-- =========================================================
INSERT IGNORE INTO c_role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM c_roles r, c_permissions p WHERE r.name = 'admin' AND p.is_delete = 0;

INSERT IGNORE INTO c_role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM c_roles r, c_permissions p
WHERE r.name = 'editor' AND p.name IN ('user:list','user:create','user:update');

-- =========================================================
-- 4. 管理后台管理员
--    密码: admin123（bcrypt 预生成哈希）
-- =========================================================
INSERT INTO c_admin (username, password_hash, real_name, status) VALUES
('admin', '$2a$10$gyRC9Lo9BxcKgDr9deqOu.mw2GzIpTNfQOkOPZCxzAJUBZ16LvtYy', '超级管理员', 1)
ON DUPLICATE KEY UPDATE
  password_hash = VALUES(password_hash),
  real_name     = VALUES(real_name),
  status        = VALUES(status);

INSERT IGNORE INTO c_admin_roles (user_id, role_id)
SELECT a.id, r.id FROM c_admin a, c_roles r
WHERE a.username = 'admin' AND r.name = 'admin';

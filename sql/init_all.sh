#!/bin/bash
# =========================================================
# 一键初始化：建库 → 建表 → 种子数据 → 管理员账号
# 用法: bash sql/init_all.sh
# =========================================================

set -e

DB_HOST="${DB_HOST:-124.223.72.223}"
DB_PORT="${DB_PORT:-13306}"
DB_USER="${DB_USER:-root}"
DB_NAME="${DB_NAME:-simon_admin}"

echo "==> 连接数据库 $DB_HOST:$DB_PORT"
echo "==> 数据库名 $DB_NAME"

# 创建数据库
mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p \
  -e "CREATE DATABASE IF NOT EXISTS \`$DB_NAME\` DEFAULT CHARSET utf8mb4"

# 导入表结构
echo "==> 导入表结构 schema.sql"
mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p "$DB_NAME" < "$(dirname "$0")/schema.sql"

# 导入种子数据
echo "==> 导入种子数据 seed.sql"
mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p "$DB_NAME" < "$(dirname "$0")/seed.sql"

# 初始化管理员
echo "==> 初始化管理员账号 init_admin.sql"
mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p "$DB_NAME" < "$(dirname "$0")/init_admin.sql"

echo "==> 初始化完成！"
echo "    管理员账号: admin / admin123"

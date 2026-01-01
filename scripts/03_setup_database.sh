#!/bin/bash

# iNarbit 数据库初始化脚本

set -e

echo "=========================================="
echo "iNarbit 数据库初始化"
echo "=========================================="

PROJECT_PATH="/root/inarbit"
DB_INIT_SQL="$PROJECT_PATH/database/init.sql"

# 检查SQL文件是否存在
if [ ! -f "$DB_INIT_SQL" ]; then
    echo "错误：找不到数据库初始化脚本: $DB_INIT_SQL"
    exit 1
fi

echo "1. 检查PostgreSQL服务..."
if ! systemctl is-active --quiet postgresql; then
    echo "启动PostgreSQL服务..."
    systemctl start postgresql
fi

echo "2. 等待PostgreSQL就绪..."
sleep 2

echo "3. 初始化数据库..."
sudo -u postgres psql -d inarbit_db -f "$DB_INIT_SQL"

echo "4. 验证数据库..."
sudo -u postgres psql -d inarbit_db -c "\dt"

echo "5. 设置数据库权限..."
sudo -u postgres psql << 'EOF'
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO inarbit;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO inarbit;
GRANT EXECUTE ON ALL FUNCTIONS IN SCHEMA public TO inarbit;
EOF

echo "=========================================="
echo "数据库初始化完成！"
echo "=========================================="
echo ""
echo "数据库信息："
echo "  主机: localhost"
echo "  端口: 5432"
echo "  数据库: inarbit_db"
echo "  用户: inarbit"
echo "  密码: inarbit_password"
echo ""
echo "默认用户："
echo "  用户名: admin"
echo "  密码: password (需要修改)"
echo "=========================================="

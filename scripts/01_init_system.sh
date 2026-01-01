#!/bin/bash

# iNarbit 系统初始化脚本
# 用途：初始化系统环境、创建必要的目录和用户

set -e

echo "=========================================="
echo "iNarbit 系统初始化"
echo "=========================================="

# 检查是否为root用户
if [[ $EUID -ne 0 ]]; then
   echo "此脚本必须以root用户身份运行"
   exit 1
fi

# 定义项目路径
PROJECT_PATH="/root/inarbit"
BACKEND_PATH="$PROJECT_PATH/backend"
FRONTEND_PATH="$PROJECT_PATH/frontend"
SCRIPTS_PATH="$PROJECT_PATH/scripts"
NGINX_PATH="$PROJECT_PATH/nginx"

echo "1. 创建项目目录结构..."
mkdir -p "$PROJECT_PATH"
mkdir -p "$BACKEND_PATH"/{config,models,handlers,services,websocket,database/migrations,utils}
mkdir -p "$FRONTEND_PATH"/{public,src/{components,pages,services,styles}}
mkdir -p "$SCRIPTS_PATH"
mkdir -p "$NGINX_PATH/ssl"
mkdir -p "$PROJECT_PATH/logs"

echo "2. 更新系统包..."
apt-get update
apt-get upgrade -y

echo "3. 安装系统依赖..."
apt-get install -y \
    curl \
    wget \
    git \
    build-essential \
    pkg-config \
    libssl-dev \
    postgresql \
    postgresql-contrib \
    nginx \
    certbot \
    python3-certbot-nginx \
    supervisor \
    htop \
    vim \
    nano

echo "4. 安装Go 1.21..."
if ! command -v go &> /dev/null; then
    cd /tmp
    wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
    tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
    rm go1.21.0.linux-amd64.tar.gz
    
    # 添加到PATH
    echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
    source /etc/profile
fi

echo "5. 安装Node.js和npm..."
if ! command -v node &> /dev/null; then
    curl -fsSL https://deb.nodesource.com/setup_20.x | bash -
    apt-get install -y nodejs
fi

echo "6. 创建应用用户..."
if ! id -u "inarbit" &>/dev/null; then
    useradd -m -s /bin/bash inarbit
fi

echo "7. 设置目录权限..."
chown -R inarbit:inarbit "$PROJECT_PATH"
chmod -R 755 "$PROJECT_PATH"
chmod -R 755 "$SCRIPTS_PATH"

echo "8. 创建环境变量文件..."
cat > "$PROJECT_PATH/.env" << 'EOF'
# 应用配置
PORT=8080
ENVIRONMENT=production

# 数据库配置
DB_HOST=localhost
DB_PORT=5432
DB_USER=inarbit
DB_PASSWORD=inarbit_password
DB_NAME=inarbit_db
DB_SSL_MODE=disable

# JWT配置
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production-12345678

# Binance配置
BINANCE_API_KEY=your_binance_api_key_here
BINANCE_API_SECRET=your_binance_api_secret_here
BINANCE_TESTNET=true

# 日志配置
LOG_LEVEL=info
LOG_FILE=/root/inarbit/logs/app.log

# WebSocket配置
WS_READ_BUFFER_SIZE=1024
WS_WRITE_BUFFER_SIZE=1024
EOF

echo "9. 启动PostgreSQL服务..."
systemctl start postgresql
systemctl enable postgresql

echo "10. 配置PostgreSQL..."
sudo -u postgres psql << 'EOF'
-- 创建数据库用户
CREATE USER inarbit WITH PASSWORD 'inarbit_password';
ALTER ROLE inarbit CREATEDB;

-- 创建数据库
CREATE DATABASE inarbit_db OWNER inarbit;

-- 授予权限
GRANT ALL PRIVILEGES ON DATABASE inarbit_db TO inarbit;
EOF

echo "11. 启动Nginx服务..."
systemctl start nginx
systemctl enable nginx

echo "12. 创建日志目录..."
mkdir -p "$PROJECT_PATH/logs"
chown -R inarbit:inarbit "$PROJECT_PATH/logs"
chmod -R 755 "$PROJECT_PATH/logs"

echo "=========================================="
echo "系统初始化完成！"
echo "=========================================="
echo ""
echo "下一步："
echo "1. 运行: bash $SCRIPTS_PATH/02_install_deps.sh"
echo "2. 运行: bash $SCRIPTS_PATH/03_setup_database.sh"
echo "3. 运行: bash $SCRIPTS_PATH/04_build_backend.sh"
echo "4. 运行: bash $SCRIPTS_PATH/05_build_frontend.sh"
echo "5. 运行: bash $SCRIPTS_PATH/06_setup_ssl.sh"
echo "6. 运行: bash $SCRIPTS_PATH/07_start_services.sh"
echo ""
echo "或者直接运行: bash $SCRIPTS_PATH/deploy.sh"
echo "=========================================="

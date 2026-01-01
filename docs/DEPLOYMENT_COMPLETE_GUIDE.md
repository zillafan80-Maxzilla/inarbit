# iNarbit 完整部署指南

## 目录

1. [系统要求](#系统要求)
2. [部署前准备](#部署前准备)
3. [快速部署](#快速部署)
4. [详细步骤](#详细步骤)
5. [服务管理](#服务管理)
6. [监控和维护](#监控和维护)
7. [故障排除](#故障排除)
8. [性能优化](#性能优化)
9. [安全加固](#安全加固)
10. [常见问题](#常见问题)

---

## 系统要求

### 硬件要求

- **CPU**: 2核或以上
- **内存**: 4GB或以上
- **磁盘**: 50GB或以上（SSD推荐）
- **网络**: 稳定的公网连接

### 软件要求

- **操作系统**: Ubuntu 22.04 LTS
- **Go**: 1.19或以上
- **PostgreSQL**: 12或以上
- **Nginx**: 1.18或以上
- **Supervisor**: 4.0或以上
- **Node.js**: 16或以上（用于前端编译）

### 域名和SSL

- **域名**: inarbit.work, www.inarbit.work
- **SSL证书**: Let's Encrypt（免费，自动续期）

---

## 部署前准备

### 1. 更新系统

```bash
sudo apt-get update
sudo apt-get upgrade -y
```

### 2. 安装基础依赖

```bash
sudo apt-get install -y \
    curl \
    wget \
    git \
    build-essential \
    libssl-dev \
    libffi-dev \
    python3-dev \
    python3-pip \
    supervisor \
    nginx \
    postgresql \
    postgresql-contrib \
    certbot \
    python3-certbot-nginx
```

### 3. 安装Go

```bash
# 下载Go
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz

# 解压
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz

# 添加到PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 验证
go version
```

### 4. 安装Node.js

```bash
# 使用NodeSource仓库
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt-get install -y nodejs

# 验证
node --version
npm --version
```

### 5. 克隆项目

```bash
cd /root
git clone https://github.com/zillafan80-Maxzilla/inarbit.git
cd inarbit
```

---

## 快速部署

### 一键部署（推荐）

```bash
# 进入项目目录
cd /root/inarbit

# 赋予脚本执行权限
chmod +x build.sh deploy.sh scripts/*.sh

# 执行部署脚本
sudo bash deploy.sh install
```

### 验证部署

```bash
# 查看服务状态
sudo supervisorctl status inarbit-backend

# 检查服务健康
curl http://localhost:8080/api/health

# 访问应用
# HTTP: http://inarbit.work
# HTTPS: https://inarbit.work
```

---

## 详细步骤

### 步骤1：系统初始化

```bash
# 创建项目目录
sudo mkdir -p /root/inarbit/{logs,data,backups,config,scripts,supervisor,nginx}

# 设置权限
sudo chown -R root:root /root/inarbit
sudo chmod -R 755 /root/inarbit
```

### 步骤2：数据库初始化

```bash
# 启动PostgreSQL
sudo systemctl start postgresql
sudo systemctl enable postgresql

# 创建数据库用户
sudo -u postgres psql << EOF
CREATE USER inarbit WITH PASSWORD 'inarbit_password';
ALTER USER inarbit CREATEDB;
EOF

# 创建数据库
sudo -u postgres createdb -O inarbit inarbit

# 初始化表结构
sudo -u postgres psql -d inarbit -f /root/inarbit/backend/scripts/init_db.sql
```

### 步骤3：编译后端

```bash
cd /root/inarbit
bash build.sh backend
```

### 步骤4：编译前端

```bash
cd /root/inarbit/frontend
npm install
npm run build
```

### 步骤5：配置环境变量

```bash
# 创建.env文件
cat > /root/inarbit/.env << EOF
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
ENV=production

DB_HOST=localhost
DB_PORT=5432
DB_USER=inarbit
DB_PASSWORD=inarbit_password
DB_NAME=inarbit

JWT_SECRET=$(openssl rand -base64 32)
CORS_ALLOWED_ORIGINS=https://inarbit.work,https://www.inarbit.work

LOG_LEVEL=info
EOF

chmod 600 /root/inarbit/.env
```

### 步骤6：配置Supervisor

```bash
# 复制Supervisor配置
sudo cp /root/inarbit/supervisor/inarbit.conf /etc/supervisor/conf.d/

# 重新加载Supervisor
sudo supervisorctl reread
sudo supervisorctl update

# 启动服务
sudo supervisorctl start inarbit-backend
```

### 步骤7：配置Nginx

```bash
# 复制Nginx配置
sudo cp /root/inarbit/nginx/inarbit.conf /etc/nginx/sites-available/

# 启用站点
sudo ln -s /etc/nginx/sites-available/inarbit /etc/nginx/sites-enabled/

# 禁用默认站点
sudo rm -f /etc/nginx/sites-enabled/default

# 测试配置
sudo nginx -t

# 启动Nginx
sudo systemctl start nginx
sudo systemctl enable nginx
```

### 步骤8：配置SSL证书

```bash
# 申请证书
sudo certbot certonly --nginx \
    -d inarbit.work \
    -d www.inarbit.work \
    --non-interactive \
    --agree-tos \
    --email admin@inarbit.work

# 更新Nginx配置为HTTPS
# 编辑 /etc/nginx/sites-available/inarbit，添加SSL配置

# 重启Nginx
sudo systemctl restart nginx

# 配置自动续期
sudo systemctl enable certbot.timer
sudo systemctl start certbot.timer
```

---

## 服务管理

### 启动服务

```bash
# 启动后端
sudo supervisorctl start inarbit-backend

# 启动Nginx
sudo systemctl start nginx

# 启动PostgreSQL
sudo systemctl start postgresql
```

### 停止服务

```bash
# 停止后端
sudo supervisorctl stop inarbit-backend

# 停止Nginx
sudo systemctl stop nginx
```

### 重启服务

```bash
# 重启后端
sudo supervisorctl restart inarbit-backend

# 重启Nginx
sudo systemctl restart nginx
```

### 查看服务状态

```bash
# 查看后端状态
sudo supervisorctl status inarbit-backend

# 查看Nginx状态
sudo systemctl status nginx

# 查看PostgreSQL状态
sudo systemctl status postgresql
```

### 查看日志

```bash
# 后端日志
tail -f /root/inarbit/logs/inarbit.log

# Nginx访问日志
tail -f /root/inarbit/logs/nginx_access.log

# Nginx错误日志
tail -f /root/inarbit/logs/nginx_error.log

# Supervisor日志
tail -f /var/log/supervisor/supervisord.log
```

---

## 监控和维护

### 健康检查

```bash
# 检查后端API
curl http://localhost:8080/api/health

# 检查前端
curl http://localhost/

# 检查WebSocket
curl -i -N -H "Connection: Upgrade" -H "Upgrade: websocket" http://localhost:8080/ws
```

### 性能监控

```bash
# 查看进程信息
ps aux | grep inarbit-backend

# 查看内存使用
free -h

# 查看磁盘使用
df -h

# 查看网络连接
netstat -an | grep ESTABLISHED | wc -l
```

### 数据库维护

```bash
# 备份数据库
sudo -u postgres pg_dump inarbit > /root/inarbit/backups/inarbit_$(date +%Y%m%d_%H%M%S).sql

# 恢复数据库
sudo -u postgres psql inarbit < /root/inarbit/backups/inarbit_YYYYMMDD_HHMMSS.sql

# 清理数据库日志
sudo -u postgres psql -d inarbit -c "VACUUM FULL;"
```

### 日志清理

```bash
# 清理7天前的日志
find /root/inarbit/logs -name "*.log" -mtime +7 -delete

# 压缩日志
gzip /root/inarbit/logs/*.log
```

---

## 故障排除

### 问题1：后端服务无法启动

**症状**: `supervisorctl status inarbit-backend` 显示 FATAL

**解决方案**:

```bash
# 查看错误日志
sudo supervisorctl tail inarbit-backend

# 检查可执行文件权限
ls -l /root/inarbit/dist/backend/inarbit-backend

# 检查环境变量
cat /root/inarbit/.env

# 检查数据库连接
psql -h localhost -U inarbit -d inarbit -c "SELECT 1"

# 重新编译
cd /root/inarbit && bash build.sh backend
```

### 问题2：Nginx无法反向代理

**症状**: 访问网站显示 502 Bad Gateway

**解决方案**:

```bash
# 检查Nginx配置
sudo nginx -t

# 检查后端服务是否运行
sudo supervisorctl status inarbit-backend

# 检查后端端口是否开放
netstat -tuln | grep 8080

# 检查Nginx日志
tail -f /root/inarbit/logs/nginx_error.log

# 重启Nginx
sudo systemctl restart nginx
```

### 问题3：数据库连接失败

**症状**: 日志显示 "database connection refused"

**解决方案**:

```bash
# 检查PostgreSQL是否运行
sudo systemctl status postgresql

# 检查数据库用户和密码
sudo -u postgres psql -d inarbit -c "SELECT 1"

# 检查.env文件中的数据库配置
cat /root/inarbit/.env | grep DB_

# 重启PostgreSQL
sudo systemctl restart postgresql
```

### 问题4：SSL证书过期

**症状**: 浏览器显示证书过期警告

**解决方案**:

```bash
# 手动续期
sudo certbot renew

# 检查证书有效期
sudo certbot certificates

# 自动续期已配置，可以等待自动续期
```

### 问题5：内存或磁盘不足

**症状**: 服务缓慢或无响应

**解决方案**:

```bash
# 查看磁盘使用
df -h

# 清理日志
find /root/inarbit/logs -name "*.log" -mtime +7 -delete

# 清理备份
find /root/inarbit/backups -name "*.sql.gz" -mtime +30 -delete

# 查看内存使用
free -h

# 重启服务释放内存
sudo supervisorctl restart inarbit-backend
```

---

## 性能优化

### 1. Nginx优化

```nginx
# 增加worker进程数
worker_processes auto;

# 增加连接数
events {
    worker_connections 2048;
}

# 启用gzip压缩
gzip on;
gzip_types text/plain text/css text/javascript application/json;
gzip_min_length 1000;

# 启用缓存
proxy_cache_path /var/cache/nginx levels=1:2 keys_zone=api_cache:10m;
```

### 2. PostgreSQL优化

```bash
# 编辑 /etc/postgresql/14/main/postgresql.conf

# 增加共享缓冲区
shared_buffers = 256MB

# 增加有效缓存大小
effective_cache_size = 1GB

# 增加工作内存
work_mem = 16MB

# 增加维护工作内存
maintenance_work_mem = 64MB
```

### 3. Go应用优化

```go
// 在main.go中添加

// 增加连接池大小
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(5)

// 启用HTTP/2
server.HTTP2.Enabled = true

// 启用Keep-Alive
server.KeepAlive = 75 * time.Second
```

### 4. 系统优化

```bash
# 增加文件描述符限制
echo "* soft nofile 65535" >> /etc/security/limits.conf
echo "* hard nofile 65535" >> /etc/security/limits.conf

# 增加网络缓冲区
echo "net.core.rmem_max = 134217728" >> /etc/sysctl.conf
echo "net.core.wmem_max = 134217728" >> /etc/sysctl.conf

# 应用更改
sudo sysctl -p
```

---

## 安全加固

### 1. 防火墙配置

```bash
# 启用UFW防火墙
sudo ufw enable

# 允许SSH
sudo ufw allow 22/tcp

# 允许HTTP
sudo ufw allow 80/tcp

# 允许HTTPS
sudo ufw allow 443/tcp

# 拒绝其他端口
sudo ufw default deny incoming
```

### 2. 更新默认密码

```bash
# 更改数据库密码
sudo -u postgres psql << EOF
ALTER USER inarbit WITH PASSWORD 'your_secure_password';
EOF

# 更新.env文件
sed -i 's/inarbit_password/your_secure_password/g' /root/inarbit/.env
```

### 3. 启用HTTPS

```bash
# 已在步骤8中配置
# 确保所有HTTP请求重定向到HTTPS
```

### 4. 定期备份

```bash
# 创建每日备份脚本
cat > /root/inarbit/scripts/daily_backup.sh << 'EOF'
#!/bin/bash
BACKUP_DIR="/root/inarbit/backups"
DATE=$(date +%Y%m%d_%H%M%S)
sudo -u postgres pg_dump inarbit | gzip > $BACKUP_DIR/inarbit_$DATE.sql.gz
find $BACKUP_DIR -name "*.sql.gz" -mtime +30 -delete
EOF

# 添加到crontab
echo "0 2 * * * /root/inarbit/scripts/daily_backup.sh" | sudo crontab -
```

### 5. 监控和告警

```bash
# 安装监控工具
sudo apt-get install -y htop iotop nethogs

# 设置日志告警
# 监控错误日志中的异常
tail -f /root/inarbit/logs/inarbit.log | grep ERROR
```

---

## 常见问题

### Q1: 如何更新应用？

```bash
cd /root/inarbit
git pull origin main
bash build.sh backend
sudo supervisorctl restart inarbit-backend
```

### Q2: 如何查看实时日志？

```bash
tail -f /root/inarbit/logs/inarbit.log
```

### Q3: 如何修改应用端口？

```bash
# 编辑.env文件
sed -i 's/SERVER_PORT=8080/SERVER_PORT=9000/g' /root/inarbit/.env

# 编辑Nginx配置
sed -i 's/127.0.0.1:8080/127.0.0.1:9000/g' /etc/nginx/sites-available/inarbit

# 重启服务
sudo supervisorctl restart inarbit-backend
sudo systemctl restart nginx
```

### Q4: 如何备份数据库？

```bash
sudo -u postgres pg_dump inarbit > /root/inarbit/backups/inarbit_backup.sql
gzip /root/inarbit/backups/inarbit_backup.sql
```

### Q5: 如何恢复数据库？

```bash
gunzip /root/inarbit/backups/inarbit_backup.sql.gz
sudo -u postgres psql inarbit < /root/inarbit/backups/inarbit_backup.sql
```

### Q6: 如何查看服务器资源使用情况？

```bash
# CPU和内存
top

# 磁盘
df -h

# 网络
netstat -an

# 进程
ps aux | grep inarbit
```

### Q7: 如何处理高并发？

```bash
# 增加Nginx worker进程
# 增加PostgreSQL连接池
# 启用Redis缓存
# 使用CDN加速静态资源
```

### Q8: 如何配置自动备份？

```bash
# 创建备份脚本
cat > /root/inarbit/scripts/backup.sh << 'EOF'
#!/bin/bash
BACKUP_DIR="/root/inarbit/backups"
DATE=$(date +%Y%m%d_%H%M%S)
sudo -u postgres pg_dump inarbit | gzip > $BACKUP_DIR/inarbit_$DATE.sql.gz
find $BACKUP_DIR -name "*.sql.gz" -mtime +30 -delete
EOF

# 添加到crontab（每天凌晨2点执行）
echo "0 2 * * * /root/inarbit/scripts/backup.sh" | sudo crontab -
```

---

## 总结

部署完成后，您的iNarbit应用应该：

- ✅ 在 https://inarbit.work 上可访问
- ✅ 后端API运行在 http://localhost:8080
- ✅ 数据库正常运行
- ✅ SSL证书已配置
- ✅ 服务由Supervisor守护
- ✅ 日志正常记录
- ✅ 定期备份已配置

如有问题，请查看日志文件或联系技术支持。

---

**最后更新**: 2024年1月
**版本**: 1.0.0

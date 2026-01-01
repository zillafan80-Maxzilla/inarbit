# iNarbit 三角套利机器人 - 部署指南

## 项目概述

iNarbit是一个基于Go和React的专业级三角套利机器人Web控制系统，支持：
- 三角/四角/五角套利策略
- 虚拟盘和实盘切换
- 实时数据推送（WebSocket）
- 多交易所支持（默认Binance）
- HTTPS安全连接

## 系统要求

- **操作系统**：Ubuntu 22.04 LTS
- **CPU**：2vCPU 或更高
- **内存**：4GB 或更高
- **存储**：50GB ESSD云盘
- **网络**：公网IP（8.211.158.208）
- **域名**：inarbit.work, www.inarbit.work

## 预装软件

- Go 1.21+
- Node.js 20+
- PostgreSQL 14+
- Nginx
- Certbot (Let's Encrypt)

## 快速开始

### 方式一：一键部署（推荐）

```bash
# 1. 登录服务器
ssh root@8.211.158.208

# 2. 进入项目目录
cd /root/inarbit

# 3. 运行一键部署脚本
bash scripts/deploy.sh

# 选择"1"执行完整部署
```

### 方式二：分步部署

```bash
# 1. 系统初始化
bash scripts/01_init_system.sh

# 2. 安装依赖
bash scripts/02_install_deps.sh

# 3. 数据库初始化
bash scripts/03_setup_database.sh

# 4. 编译后端
bash scripts/04_build_backend.sh

# 5. 编译前端
bash scripts/05_build_frontend.sh

# 6. 配置SSL证书
bash scripts/06_setup_ssl.sh

# 7. 启动服务
bash scripts/07_start_services.sh
```

## 项目结构

```
/root/inarbit/
├── backend/              # Go后端应用
│   ├── main.go          # 主程序
│   ├── go.mod           # 模块定义
│   ├── config/          # 配置
│   ├── models/          # 数据模型
│   ├── handlers/        # HTTP处理器
│   ├── services/        # 业务逻辑
│   ├── websocket/       # WebSocket
│   ├── database/        # 数据库
│   └── utils/           # 工具函数
├── frontend/            # React前端应用
│   ├── public/          # 静态资源
│   ├── src/             # 源代码
│   ├── package.json     # 依赖定义
│   └── vite.config.js   # Vite配置
├── nginx/               # Nginx配置
│   ├── nginx.conf       # 配置文件
│   └── ssl/             # SSL证书
├── scripts/             # 部署脚本
├── logs/                # 日志文件
└── .env                 # 环境变量
```

## 配置说明

### 环境变量 (.env)

```env
# 应用配置
PORT=8080
ENVIRONMENT=production

# 数据库
DB_HOST=localhost
DB_PORT=5432
DB_USER=inarbit
DB_PASSWORD=inarbit_password
DB_NAME=inarbit_db

# JWT
JWT_SECRET=your-secret-key

# Binance
BINANCE_API_KEY=your_api_key
BINANCE_API_SECRET=your_api_secret
BINANCE_TESTNET=true
```

### 数据库配置

- **主机**：localhost
- **端口**：5432
- **数据库**：inarbit_db
- **用户**：inarbit
- **密码**：inarbit_password

### Nginx配置

Nginx配置文件位于：`/etc/nginx/sites-available/inarbit`

主要功能：
- HTTP → HTTPS 重定向
- 前端静态文件服务
- API代理到后端
- WebSocket支持
- SSL/TLS配置

## 访问应用

### 前端应用

```
https://inarbit.work
https://www.inarbit.work
```

### API端点

```
https://inarbit.work/api/
```

### WebSocket连接

```
wss://inarbit.work/ws
```

### 默认用户

- **用户名**：admin
- **密码**：password（请立即修改）

## 服务管理

### 后端服务

```bash
# 启动
systemctl start inarbit

# 停止
systemctl stop inarbit

# 重启
systemctl restart inarbit

# 查看状态
systemctl status inarbit

# 查看日志
journalctl -u inarbit -f
```

### Nginx服务

```bash
# 启动
systemctl start nginx

# 停止
systemctl stop nginx

# 重启
systemctl restart nginx

# 测试配置
nginx -t

# 查看日志
tail -f /root/inarbit/logs/nginx_access.log
tail -f /root/inarbit/logs/nginx_error.log
```

### PostgreSQL服务

```bash
# 启动
systemctl start postgresql

# 停止
systemctl stop postgresql

# 重启
systemctl restart postgresql

# 连接数据库
psql -U inarbit -d inarbit_db -h localhost
```

## SSL证书管理

### 证书信息

- **签发机构**：Let's Encrypt
- **证书路径**：`/root/inarbit/nginx/ssl/cert.pem`
- **密钥路径**：`/root/inarbit/nginx/ssl/key.pem`
- **有效期**：90天
- **自动续期**：每天凌晨3点

### 手动续期

```bash
certbot renew --force-renewal
systemctl reload nginx
```

### 验证证书

```bash
openssl x509 -in /root/inarbit/nginx/ssl/cert.pem -text -noout
```

## 性能优化

### 数据库优化

```sql
-- 创建索引
CREATE INDEX idx_trades_created_at ON trades(created_at);
CREATE INDEX idx_trades_status ON trades(status);

-- 分析表
ANALYZE trades;
```

### 应用优化

1. **启用Gzip压缩**：已在Nginx中配置
2. **缓存策略**：静态文件缓存30天
3. **连接池**：PostgreSQL连接池配置
4. **WebSocket优化**：心跳检测和自动重连

## 监控和日志

### 应用日志

```bash
# 后端日志
journalctl -u inarbit -f

# 前端日志（浏览器控制台）
# 按F12打开开发者工具
```

### Nginx日志

```bash
# 访问日志
tail -f /root/inarbit/logs/nginx_access.log

# 错误日志
tail -f /root/inarbit/logs/nginx_error.log
```

### 系统监控

```bash
# CPU和内存使用
top

# 磁盘使用
df -h

# 网络连接
netstat -an | grep ESTABLISHED
```

## 常见问题

### Q1：无法访问应用

**解决方案**：
1. 检查防火墙：`sudo ufw status`
2. 检查Nginx：`systemctl status nginx`
3. 检查后端：`systemctl status inarbit`
4. 查看日志：`journalctl -u inarbit -f`

### Q2：SSL证书过期

**解决方案**：
```bash
certbot renew --force-renewal
systemctl reload nginx
```

### Q3：数据库连接失败

**解决方案**：
```bash
# 检查PostgreSQL
systemctl status postgresql

# 测试连接
psql -U inarbit -d inarbit_db -h localhost

# 检查密码
# 编辑 /root/inarbit/.env
```

### Q4：WebSocket连接失败

**解决方案**：
1. 检查Nginx配置：`nginx -t`
2. 检查防火墙是否允许443端口
3. 查看浏览器控制台错误信息

## 备份和恢复

### 数据库备份

```bash
# 完整备份
pg_dump -U inarbit -d inarbit_db > backup.sql

# 压缩备份
pg_dump -U inarbit -d inarbit_db | gzip > backup.sql.gz

# 恢复
psql -U inarbit -d inarbit_db < backup.sql
```

### 应用备份

```bash
# 备份整个项目
tar -czf inarbit_backup.tar.gz /root/inarbit/

# 恢复
tar -xzf inarbit_backup.tar.gz -C /
```

## 安全建议

1. **修改默认密码**：登录后立即修改admin密码
2. **配置防火墙**：
   ```bash
   ufw allow 22/tcp
   ufw allow 80/tcp
   ufw allow 443/tcp
   ufw enable
   ```
3. **启用SSH密钥认证**：禁用密码登录
4. **定期更新**：`apt-get update && apt-get upgrade`
5. **备份数据**：定期备份数据库和配置
6. **监控日志**：定期检查系统和应用日志

## 更新和维护

### 更新应用

```bash
# 1. 停止服务
systemctl stop inarbit

# 2. 更新代码
cd /root/inarbit
git pull origin main

# 3. 重新编译
bash scripts/04_build_backend.sh
bash scripts/05_build_frontend.sh

# 4. 启动服务
systemctl start inarbit
```

### 更新依赖

```bash
# Go依赖
cd /root/inarbit/backend
go get -u ./...

# Node依赖
cd /root/inarbit/frontend
npm update
```

## 技术支持

如有问题，请：
1. 查看日志文件
2. 检查系统资源
3. 验证配置文件
4. 参考本文档的常见问题部分

## 许可证

iNarbit © 2024 All Rights Reserved

---

**最后更新**：2024年1月
**版本**：1.0.0

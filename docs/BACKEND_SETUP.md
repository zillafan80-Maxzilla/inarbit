# iNarbit 后端设置指南

## 项目结构

```
/root/inarbit/
├── backend/
│   ├── main.go                    # 主程序入口（inarbit_main_updated.go）
│   ├── go.mod                     # Go模块定义（inarbit_go_mod_updated.mod）
│   ├── go.sum                     # Go依赖锁定文件（自动生成）
│   ├── config.go                  # 配置管理（inarbit_config.go）
│   ├── models/
│   │   └── types.go               # 数据模型（inarbit_models_types.go）
│   ├── database/
│   │   └── db.go                  # 数据库操作（inarbit_database_db.go）
│   ├── services/
│   │   └── auth.go                # 认证服务（inarbit_services_auth.go）
│   ├── handlers/
│   │   └── api.go                 # API处理器（inarbit_handlers_api.go）
│   ├── websocket/
│   │   └── manager.go             # WebSocket管理器（inarbit_websocket_manager.go）
│   └── scripts/
│       ├── init_db.sql            # 数据库初始化脚本（inarbit_init_database.sql）
│       └── seed_db.sql            # 数据库种子数据（待创建）
└── frontend/
    └── ...（前端代码）
```

## 文件对应关系

| 生成的文件 | 目标位置 |
|-----------|--------|
| inarbit_main_updated.go | backend/main.go |
| inarbit_go_mod_updated.mod | backend/go.mod |
| inarbit_config.go | backend/config.go |
| inarbit_models_types.go | backend/models/types.go |
| inarbit_database_db.go | backend/database/db.go |
| inarbit_services_auth.go | backend/services/auth.go |
| inarbit_handlers_api.go | backend/handlers/api.go |
| inarbit_websocket_manager.go | backend/websocket/manager.go |
| inarbit_init_database.sql | backend/scripts/init_db.sql |

## 安装步骤

### 1. 创建目录结构

```bash
cd /root/inarbit
mkdir -p backend/{models,database,services,handlers,websocket,scripts}
```

### 2. 复制文件

将所有生成的Go文件复制到相应的目录。

### 3. 安装Go依赖

```bash
cd /root/inarbit/backend
go mod download
go mod tidy
```

### 4. 初始化数据库

```bash
# 连接到PostgreSQL
psql -U postgres -h localhost

# 创建数据库和用户
CREATE USER inarbit WITH PASSWORD 'inarbit_password';
CREATE DATABASE inarbit OWNER inarbit;
GRANT ALL PRIVILEGES ON DATABASE inarbit TO inarbit;

# 运行初始化脚本
psql -U inarbit -d inarbit -h localhost -f backend/scripts/init_db.sql
```

### 5. 创建初始用户

```bash
# 使用以下SQL创建默认管理员用户
psql -U inarbit -d inarbit -h localhost

INSERT INTO users (username, password_hash, email, is_active) VALUES (
  'admin',
  '$2a$10$...',  -- bcrypt hash of 'password'
  'admin@inarbit.work',
  true
);
```

### 6. 配置环境变量

创建 `.env` 文件：

```bash
# 服务器配置
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
ENV=production

# 数据库配置
DB_HOST=localhost
DB_PORT=5432
DB_USER=inarbit
DB_PASSWORD=inarbit_password
DB_NAME=inarbit

# JWT配置
JWT_SECRET=your-secret-key-change-in-production

# CORS配置
CORS_ALLOWED_ORIGINS=https://inarbit.work,https://www.inarbit.work

# 日志配置
LOG_LEVEL=info
```

### 7. 编译后端

```bash
cd /root/inarbit/backend
go build -o inarbit-backend main.go config.go
```

### 8. 运行后端

```bash
# 开发模式
go run main.go config.go

# 生产模式
./inarbit-backend
```

## API文档

### 认证API

#### 登录

```
POST /api/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "password"
}

响应:
{
  "status": "success",
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": 1,
      "username": "admin",
      "email": "admin@inarbit.work",
      "is_active": true,
      "created_at": "2024-01-01T00:00:00Z"
    }
  }
}
```

#### 验证认证

```
GET /api/auth/verify
Authorization: Bearer <token>

响应:
{
  "status": "success",
  "message": "认证有效",
  "data": {
    "id": 1,
    "username": "admin",
    "email": "admin@inarbit.work",
    "is_active": true
  }
}
```

#### 登出

```
POST /api/auth/logout
Authorization: Bearer <token>

响应:
{
  "status": "success",
  "message": "登出成功"
}
```

### 机器人API

#### 获取机器人列表

```
GET /api/bots
Authorization: Bearer <token>

响应:
{
  "status": "success",
  "message": "获取机器人列表成功",
  "data": {
    "bots": [
      {
        "id": 1,
        "user_id": 1,
        "name": "三角套利机器人-1",
        "strategy_type": "triangular",
        "exchange_id": 1,
        "is_running": true,
        "is_simulation": false,
        "update_frequency": 5,
        "created_at": "2024-01-01T00:00:00Z",
        "total_profit": 1250.50,
        "total_trades": 45
      }
    ]
  }
}
```

#### 创建机器人

```
POST /api/bots
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "三角套利机器人-1",
  "strategy_type": "triangular",
  "exchange_id": 1,
  "is_simulation": false,
  "update_frequency": 5
}

响应: (同上)
```

#### 获取单个机器人

```
GET /api/bots/{id}
Authorization: Bearer <token>
```

#### 更新机器人

```
PUT /api/bots/{id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "三角套利机器人-1",
  "strategy_type": "triangular",
  "exchange_id": 1,
  "is_simulation": false,
  "update_frequency": 5
}
```

#### 删除机器人

```
DELETE /api/bots/{id}
Authorization: Bearer <token>
```

#### 启动机器人

```
POST /api/bots/{id}/start
Authorization: Bearer <token>
```

#### 停止机器人

```
POST /api/bots/{id}/stop
Authorization: Bearer <token>
```

#### 切换模式（虚拟/实盘）

```
POST /api/bots/{id}/switch-mode
Authorization: Bearer <token>
Content-Type: application/json

{
  "is_simulation": true
}
```

### 仪表板API

#### 获取统计数据

```
GET /api/dashboard/stats
Authorization: Bearer <token>

响应:
{
  "status": "success",
  "message": "获取统计数据成功",
  "data": {
    "total_profit": 5000.75,
    "total_trades": 150,
    "active_bots": 3,
    "win_rate": 72.5
  }
}
```

#### 获取图表数据

```
GET /api/dashboard/chart-data
Authorization: Bearer <token>

响应:
{
  "status": "success",
  "message": "获取图表数据成功",
  "data": {
    "data": [
      {
        "date": "2024-01-01",
        "profit": 120.50,
        "trades": 8
      },
      {
        "date": "2024-01-02",
        "profit": 290.75,
        "trades": 12
      }
    ]
  }
}
```

### 交易所API

#### 获取交易所列表

```
GET /api/exchanges
Authorization: Bearer <token>

响应:
{
  "status": "success",
  "message": "获取交易所列表成功",
  "data": {
    "exchanges": [
      {
        "id": 1,
        "user_id": 1,
        "name": "Binance",
        "is_testnet": false,
        "is_active": true,
        "created_at": "2024-01-01T00:00:00Z"
      }
    ]
  }
}
```

#### 创建交易所

```
POST /api/exchanges
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Binance",
  "api_key": "your-api-key",
  "api_secret": "your-api-secret",
  "is_testnet": false
}
```

#### 删除交易所

```
DELETE /api/exchanges/{id}
Authorization: Bearer <token>
```

### WebSocket API

#### 连接WebSocket

```
ws://localhost:8080/ws

初始化消息:
{
  "type": "auth",
  "token": "<jwt_token>"
}

连接成功响应:
{
  "type": "connected",
  "payload": {
    "message": "连接成功",
    "user_id": 1
  }
}
```

#### WebSocket消息类型

##### 统计更新

```json
{
  "type": "stats_update",
  "payload": {
    "total_profit": 5000.75,
    "total_trades": 150,
    "active_bots": 3,
    "win_rate": 72.5,
    "timestamp": "2024-01-01T12:00:00Z"
  }
}
```

##### 交易更新

```json
{
  "type": "trade_update",
  "payload": {
    "bot_id": 1,
    "trade_id": 123,
    "status": "completed",
    "profit": 125.50,
    "timestamp": "2024-01-01T12:00:00Z"
  }
}
```

##### 机器人状态

```json
{
  "type": "bot_status",
  "payload": {
    "bot_id": 1,
    "is_running": true,
    "timestamp": "2024-01-01T12:00:00Z"
  }
}
```

##### 日志

```json
{
  "type": "log",
  "payload": {
    "bot_id": 1,
    "level": "info",
    "message": "机器人已启动",
    "time": "2024-01-01T12:00:00Z"
  }
}
```

## 性能优化

### 数据库优化

1. **连接池**：已配置最大25个连接，最小5个空闲连接
2. **索引**：确保以下字段有索引：
   - users.username
   - bots.user_id
   - bots.is_running
   - trades.bot_id
   - trades.created_at

3. **查询优化**：使用预编译语句防止SQL注入

### WebSocket优化

1. **消息缓冲**：每个客户端有256字节的消息缓冲
2. **心跳检测**：30秒发送一次ping消息
3. **超时处理**：60秒无响应自动断开连接

### API优化

1. **CORS**：已配置跨域资源共享
2. **超时**：读取超时15秒，写入超时15秒
3. **并发**：支持多个并发连接

## 监控和日志

### 日志级别

- **info**：信息级别日志（默认）
- **debug**：调试级别日志
- **warn**：警告级别日志
- **error**：错误级别日志

### 关键日志

- ✓ 用户登录成功
- ✓ 数据库连接成功
- ✓ WebSocket客户端连接
- ✗ 数据库连接失败
- ✗ 认证失败
- ✗ API错误

## 故障排除

### 数据库连接失败

```bash
# 检查PostgreSQL是否运行
sudo systemctl status postgresql

# 检查数据库凭证
psql -U inarbit -d inarbit -h localhost

# 检查防火墙
sudo ufw allow 5432
```

### WebSocket连接失败

```bash
# 检查后端是否运行
curl http://localhost:8080/api/health

# 检查WebSocket升级
# 确保浏览器支持WebSocket
# 检查代理配置（如Nginx）
```

### 性能问题

1. **增加数据库连接数**：修改 `db.SetMaxOpenConns()`
2. **增加WebSocket缓冲**：修改 `upgrader.ReadBufferSize` 和 `WriteBufferSize`
3. **使用缓存**：考虑添加Redis缓存

## 部署到生产环境

### 1. 使用Systemd服务

创建 `/etc/systemd/system/inarbit-backend.service`：

```ini
[Unit]
Description=iNarbit Backend Service
After=network.target postgresql.service

[Service]
Type=simple
User=root
WorkingDirectory=/root/inarbit/backend
EnvironmentFile=/root/inarbit/backend/.env
ExecStart=/root/inarbit/backend/inarbit-backend
Restart=on-failure
RestartSec=10

[Install]
WantedBy=multi-user.target
```

启动服务：

```bash
sudo systemctl daemon-reload
sudo systemctl enable inarbit-backend
sudo systemctl start inarbit-backend
sudo systemctl status inarbit-backend
```

### 2. 使用Nginx反向代理

```nginx
upstream inarbit_backend {
    server 127.0.0.1:8080;
}

server {
    listen 80;
    server_name inarbit.work www.inarbit.work;

    location / {
        proxy_pass http://inarbit_backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    location /ws {
        proxy_pass http://inarbit_backend;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 3. SSL证书配置

```bash
# 使用Let's Encrypt
sudo certbot certonly --standalone -d inarbit.work -d www.inarbit.work

# 在Nginx中启用HTTPS
listen 443 ssl;
ssl_certificate /etc/letsencrypt/live/inarbit.work/fullchain.pem;
ssl_certificate_key /etc/letsencrypt/live/inarbit.work/privkey.pem;
```

## 技术栈

- **语言**：Go 1.21+
- **数据库**：PostgreSQL 12+
- **认证**：JWT (golang-jwt)
- **WebSocket**：Gorilla WebSocket
- **路由**：Gorilla Mux
- **CORS**：rs/cors
- **密码加密**：bcrypt

## 常见问题

### Q1：如何修改JWT过期时间？

编辑 `services/auth.go` 中的 `GenerateToken` 函数：

```go
expiresAt := now.Add(24 * time.Hour) // 修改这里
```

### Q2：如何添加新的API端点？

1. 在 `handlers/api.go` 中添加处理函数
2. 在 `RegisterRoutes` 中注册路由
3. 如果需要认证，使用 `AuthMiddleware` 包装

### Q3：如何实现实时数据推送？

使用WebSocket管理器的广播函数：

```go
wsManager.BroadcastStatsUpdate(userID, stats)
wsManager.BroadcastTradeUpdate(userID, trade)
wsManager.BroadcastBotStatus(userID, botID, isRunning)
```

### Q4：如何处理数据库迁移？

1. 修改 `scripts/init_db.sql`
2. 运行迁移脚本
3. 重新启动后端服务

## 下一步

1. **实现交易引擎**：完成三角/四角/五角套利算法
2. **集成Binance API**：实现实时价格获取和下单
3. **添加策略管理**：完整的策略CRUD操作
4. **性能监控**：添加Prometheus指标
5. **日志系统**：集成ELK日志系统

---

**最后更新**：2024年1月
**版本**：1.0.0

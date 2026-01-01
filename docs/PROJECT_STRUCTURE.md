# iNarbit 三角套利机器人 - Web控制系统

## 项目结构

```
/root/inarbit/
├── backend/                    # Go后端应用
│   ├── main.go                # 主程序入口
│   ├── go.mod                 # Go模块定义
│   ├── go.sum                 # Go依赖锁定
│   ├── config/                # 配置文件
│   │   ├── config.go          # 配置结构体
│   │   └── config.yaml        # 配置文件
│   ├── models/                # 数据模型
│   │   ├── user.go
│   │   ├── bot.go
│   │   ├── strategy.go
│   │   ├── trade.go
│   │   └── exchange.go
│   ├── handlers/              # HTTP处理器
│   │   ├── auth.go            # 认证相关
│   │   ├── bot.go             # 机器人管理
│   │   ├── strategy.go        # 策略管理
│   │   ├── trade.go           # 交易查询
│   │   └── exchange.go        # 交易所管理
│   ├── services/              # 业务逻辑
│   │   ├── auth_service.go
│   │   ├── bot_service.go
│   │   ├── trade_service.go
│   │   ├── arbitrage_engine.go # 套利引擎
│   │   └── exchange_service.go
│   ├── websocket/             # WebSocket实时推送
│   │   ├── hub.go
│   │   └── client.go
│   ├── database/              # 数据库操作
│   │   ├── db.go
│   │   └── migrations/
│   │       └── init.sql
│   └── utils/                 # 工具函数
│       ├── logger.go
│       ├── crypto.go
│       └── helpers.go
├── frontend/                  # React前端应用
│   ├── public/
│   │   └── index.html
│   ├── src/
│   │   ├── components/
│   │   ├── pages/
│   │   ├── services/
│   │   ├── styles/
│   │   ├── App.jsx
│   │   └── main.jsx
│   ├── package.json
│   └── vite.config.js
├── nginx/                     # Nginx配置
│   ├── nginx.conf
│   └── ssl/                   # SSL证书存放
├── scripts/                   # 部署脚本
│   ├── 01_init_system.sh      # 系统初始化
│   ├── 02_install_deps.sh     # 安装依赖
│   ├── 03_setup_database.sh   # 数据库初始化
│   ├── 04_build_backend.sh    # 编译后端
│   ├── 05_build_frontend.sh   # 编译前端
│   ├── 06_setup_ssl.sh        # SSL证书配置
│   ├── 07_start_services.sh   # 启动服务
│   └── deploy.sh              # 一键部署脚本
├── docker-compose.yml         # Docker编排（可选）
├── .env.example               # 环境变量示例
└── README.md                  # 项目文档
```

## 部署步骤

1. 在服务器上创建项目目录：`mkdir -p /root/inarbit && cd /root/inarbit`
2. 将所有生成的文件上传到服务器
3. 按顺序执行部署脚本：
   ```bash
   chmod +x scripts/*.sh
   bash scripts/deploy.sh
   ```

## 技术栈

- **后端**：Go 1.21+
- **前端**：React 18 + Vite
- **数据库**：PostgreSQL 14+
- **Web服务器**：Nginx
- **实时通信**：WebSocket
- **SSL**：Let's Encrypt (Certbot)

## 功能特性

- ✅ 用户登录认证（JWT）
- ✅ 机器人创建、编辑、删除、启停
- ✅ 三角/四角/五角套利策略配置
- ✅ 虚拟盘和实盘切换
- ✅ 实时交易数据推送（WebSocket）
- ✅ 参数动态修改（秒级到小时级更新频率）
- ✅ Binance交易所集成（支持添加其他交易所）
- ✅ 交易历史记录查询
- ✅ 收益统计分析
- ✅ SSL安全连接


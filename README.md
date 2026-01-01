# iNarbit - 三角套利交易机器人

完整的加密货币三角套利交易机器人系统，支持3/4/5角套利策略。

## 项目结构

```
inarbit/
├── backend/          # Go后端服务
│   ├── main.go       # 主入口
│   ├── go.mod        # Go模块定义
│   └── go.sum        # Go依赖锁定
├── frontend/         # React前端应用
│   ├── src/
│   │   ├── main.jsx  # 入口文件
│   │   ├── App.jsx   # 主应用组件
│   │   └── index.css # 全局样式
│   ├── package.json  # npm依赖
│   ├── vite.config.js # Vite配置
│   └── tailwind.config.js # Tailwind配置
├── database/         # 数据库脚本
├── scripts/          # 部署脚本
└── docs/             # 文档
```

## 环境要求

- Go 1.21+
- Node.js 18+
- npm 9+
- PostgreSQL 12+

## 本地开发

### 后端

```bash
cd backend
go mod download
go run main.go
```

### 前端

```bash
cd frontend
npm install
npm run dev
```

## 编译

### 后端

```bash
cd backend
go build -o inarbit-server main.go
```

### 前端

```bash
cd frontend
npm install
npm run build
```

## 部署

详见部署文档。

## 许可证

MIT

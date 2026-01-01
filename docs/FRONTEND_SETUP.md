# iNarbit 前端应用设置指南

## 文件结构

将以下生成的文件放入相应的目录：

```
/root/inarbit/frontend/
├── public/
│   └── index.html                    # HTML入口（inarbit_frontend_index.html）
├── src/
│   ├── main.jsx                      # 应用入口（inarbit_frontend_main.jsx）
│   ├── App.jsx                       # 应用主组件（inarbit_frontend_App.jsx）
│   ├── index.css                     # 全局样式（inarbit_frontend_index.css）
│   ├── pages/
│   │   ├── LoginPage.jsx             # 登录页面（inarbit_LoginPage.jsx）
│   │   ├── LoginPage.css             # 登录页面样式（inarbit_LoginPage.css）
│   │   ├── DashboardPage.jsx         # 仪表板页面（inarbit_DashboardPage.jsx）
│   │   ├── DashboardPage.css         # 仪表板样式（inarbit_DashboardPage.css）
│   │   ├── BotsPage.jsx              # 机器人管理页面（待开发）
│   │   ├── StrategiesPage.jsx        # 策略管理页面（待开发）
│   │   ├── TradesPage.jsx            # 交易记录页面（待开发）
│   │   ├── SettingsPage.jsx          # 设置页面（待开发）
│   │   └── NotFoundPage.jsx          # 404页面（待开发）
│   ├── components/
│   │   ├── Layout.jsx                # 布局组件（inarbit_Layout.jsx）
│   │   ├── Layout.css                # 布局样式（inarbit_Layout.css）
│   │   ├── ProtectedRoute.jsx        # 受保护路由（inarbit_ProtectedRoute.jsx）
│   │   ├── StatCard.jsx              # 统计卡片（inarbit_StatCard.jsx）
│   │   ├── StatCard.css              # 统计卡片样式（inarbit_StatCard.css）
│   │   ├── BotCard.jsx               # 机器人卡片（inarbit_BotCard.jsx）
│   │   ├── BotCard.css               # 机器人卡片样式（inarbit_BotCard.css）
│   │   └── ErrorBoundary.jsx         # 错误边界（待开发）
│   ├── hooks/
│   │   └── useWebSocket.js           # WebSocket Hook（inarbit_useWebSocket.js）
│   ├── store/
│   │   └── authStore.js              # 认证状态管理（inarbit_authStore.js）
│   └── services/
│       ├── api.js                    # API服务（待开发）
│       └── websocket.js              # WebSocket服务（待开发）
├── package.json                      # 依赖定义（inarbit_frontend_package.json）
├── vite.config.js                    # Vite配置（inarbit_frontend_vite.config.js）
├── tailwind.config.js                # Tailwind配置（inarbit_frontend_tailwind.config.js）
└── postcss.config.js                 # PostCSS配置（inarbit_frontend_postcss.config.js）
```

## 文件对应关系

| 生成的文件 | 目标位置 |
|-----------|--------|
| inarbit_frontend_index.html | frontend/public/index.html |
| inarbit_frontend_main.jsx | frontend/src/main.jsx |
| inarbit_frontend_App.jsx | frontend/src/App.jsx |
| inarbit_frontend_index.css | frontend/src/index.css |
| inarbit_LoginPage.jsx | frontend/src/pages/LoginPage.jsx |
| inarbit_LoginPage.css | frontend/src/pages/LoginPage.css |
| inarbit_DashboardPage.jsx | frontend/src/pages/DashboardPage.jsx |
| inarbit_DashboardPage.css | frontend/src/pages/DashboardPage.css |
| inarbit_Layout.jsx | frontend/src/components/Layout.jsx |
| inarbit_Layout.css | frontend/src/components/Layout.css |
| inarbit_ProtectedRoute.jsx | frontend/src/components/ProtectedRoute.jsx |
| inarbit_StatCard.jsx | frontend/src/components/StatCard.jsx |
| inarbit_StatCard.css | frontend/src/components/StatCard.css |
| inarbit_BotCard.jsx | frontend/src/components/BotCard.jsx |
| inarbit_BotCard.css | frontend/src/components/BotCard.css |
| inarbit_useWebSocket.js | frontend/src/hooks/useWebSocket.js |
| inarbit_authStore.js | frontend/src/store/authStore.js |
| inarbit_frontend_package.json | frontend/package.json |
| inarbit_frontend_vite.config.js | frontend/vite.config.js |
| inarbit_frontend_tailwind.config.js | frontend/tailwind.config.js |
| inarbit_frontend_postcss.config.js | frontend/postcss.config.js |

## 安装步骤

### 1. 创建目录结构

```bash
cd /root/inarbit/frontend

# 创建必要的目录
mkdir -p src/{pages,components,hooks,store,services}
mkdir -p public
```

### 2. 复制文件

将所有生成的文件复制到相应的目录。

### 3. 安装依赖

```bash
cd /root/inarbit/frontend
npm install
```

### 4. 开发模式运行

```bash
npm run dev
```

访问 http://localhost:3000

### 5. 生产构建

```bash
npm run build
```

## 设计风格说明

### 配色方案

- **主色**：#9d8b7e（棕色）
- **辅色**：#b5a293（浅棕色）
- **背景**：#f5f1e8（米色）
- **文本**：#5f5550（深棕色）
- **边框**：#e5ddd0（浅灰色）

### 字体

- **主字体**：Inter
- **备用字体**：-apple-system, BlinkMacSystemFont, Segoe UI, sans-serif

### 组件风格

- **圆角**：8px - 14px
- **阴影**：0 4px 12px rgba(0, 0, 0, 0.06)
- **过渡**：0.3s ease
- **边框**：1px solid rgba(157, 139, 126, 0.1)

## 已实现的页面

### 1. 登录页面 (LoginPage)

**功能**：
- 用户名和密码输入
- 记住我功能
- 密码显示/隐藏
- 错误提示
- 加载状态

**特点**：
- 居中对称布局
- 米色/棕色配色
- 平滑动画效果
- 响应式设计

### 2. 仪表板页面 (DashboardPage)

**功能**：
- 统计卡片（总收益、总交易数、活跃机器人、胜率）
- 收益趋势图表
- 交易统计图表
- 交易成功率饼图
- 最近活动列表
- 活跃机器人卡片

**特点**：
- 实时数据更新（WebSocket）
- 交互式图表
- 卡片式布局
- 响应式网格

### 3. 布局组件 (Layout)

**功能**：
- 侧边栏导航
- 顶部用户菜单
- 响应式设计
- 移动端折叠菜单

**特点**：
- 固定侧边栏
- 可折叠导航
- 用户信息展示
- 登出功能

## 待开发的页面

### 1. 机器人管理页面 (BotsPage)

需要实现：
- 机器人列表
- 创建/编辑/删除机器人
- 启动/停止机器人
- 机器人详情页

### 2. 策略管理页面 (StrategiesPage)

需要实现：
- 策略列表
- 创建/编辑/删除策略
- 策略参数配置
- 策略测试

### 3. 交易记录页面 (TradesPage)

需要实现：
- 交易列表
- 交易详情
- 交易筛选和搜索
- 交易导出

### 4. 设置页面 (SettingsPage)

需要实现：
- 用户设置
- API密钥管理
- 通知设置
- 系统设置

### 5. 404页面 (NotFoundPage)

需要实现：
- 404错误提示
- 返回首页链接

## API集成

### 认证API

```javascript
// 登录
POST /api/auth/login
{
  "username": "admin",
  "password": "password"
}

// 验证token
GET /api/auth/verify
Headers: { "Authorization": "Bearer <token>" }

// 登出
POST /api/auth/logout
```

### 机器人API

```javascript
// 获取机器人列表
GET /api/bots

// 获取单个机器人
GET /api/bots/:id

// 创建机器人
POST /api/bots

// 更新机器人
PUT /api/bots/:id

// 删除机器人
DELETE /api/bots/:id

// 启动机器人
POST /api/bots/:id/start

// 停止机器人
POST /api/bots/:id/stop

// 切换模式（虚拟/实盘）
POST /api/bots/:id/switch-mode
```

### 仪表板API

```javascript
// 获取统计数据
GET /api/dashboard/stats

// 获取图表数据
GET /api/dashboard/chart-data
```

## WebSocket事件

### 认证

```javascript
{
  "type": "auth",
  "token": "<jwt_token>"
}
```

### 统计更新

```javascript
{
  "type": "stats_update",
  "payload": {
    "totalProfit": 1250.50,
    "totalTrades": 45,
    "activeBots": 3,
    "winRate": 72
  }
}
```

### 交易更新

```javascript
{
  "type": "trade_update",
  "payload": {
    "botId": 1,
    "tradeId": 123,
    "status": "completed",
    "profit": 125.50
  }
}
```

## 开发建议

1. **组件复用**：充分利用StatCard和BotCard等可复用组件
2. **状态管理**：使用Zustand管理全局状态
3. **API调用**：创建统一的API服务层
4. **错误处理**：添加错误边界和错误提示
5. **加载状态**：为所有异步操作添加加载状态
6. **响应式设计**：确保所有页面在移动设备上正常显示
7. **性能优化**：使用React.memo和useMemo优化性能
8. **可访问性**：添加ARIA标签和键盘导航支持

## 常见问题

### Q1：如何修改配色？

编辑 `frontend/src/index.css` 中的CSS变量：

```css
:root {
  --color-primary: #9d8b7e;
  --color-secondary: #b5a293;
  /* ... */
}
```

### Q2：如何添加新页面？

1. 在 `src/pages/` 创建新的 `.jsx` 文件
2. 在 `App.jsx` 中添加路由
3. 在 `Layout.jsx` 中添加导航项

### Q3：如何调试WebSocket？

在浏览器开发者工具中：
1. 打开 Network 标签
2. 筛选 WS 类型
3. 查看消息内容

### Q4：如何处理API错误？

所有API调用都应该有try-catch：

```javascript
try {
  const response = await fetch('/api/...')
  if (!response.ok) throw new Error('API错误')
  const data = await response.json()
  // 处理数据
} catch (error) {
  setError(error.message)
}
```

## 技术栈

- **框架**：React 18
- **路由**：React Router 6
- **状态管理**：Zustand
- **UI库**：Recharts（图表）
- **图标**：Lucide React
- **样式**：Tailwind CSS
- **构建工具**：Vite
- **包管理**：npm

## 性能指标

目标性能指标：
- **首屏加载**：< 2s
- **交互响应**：< 100ms
- **Lighthouse得分**：> 90

## 部署

### 开发环境

```bash
npm run dev
```

### 生产构建

```bash
npm run build
```

### 预览生产构建

```bash
npm run preview
```

生成的文件在 `dist/` 目录中，由后端Nginx提供服务。

---

**最后更新**：2024年1月
**版本**：1.0.0

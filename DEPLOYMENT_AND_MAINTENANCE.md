# www.inarbit.work 部署和维护指南

**项目**: 开源三角套利机器人研究报告  
**网址**: http://8.211.158.208/ (www.inarbit.work)  
**服务器**: 阿里云 (8.211.158.208)  
**更新日期**: 2026-01-01

---

## 📋 项目概览

| 项目 | 详情 |
|------|------|
| **项目名称** | 开源三角套利机器人研究报告 |
| **前端技术** | React 19 + TypeScript + Vite + Tailwind CSS 4 |
| **后端技术** | Go 1.21 |
| **路由库** | Wouter |
| **设计主题** | Solarized Light |
| **部署位置** | 阿里云 (8.211.158.208) |
| **GitHub仓库** | https://github.com/zillafan80-Maxzilla/inarbit |
| **主分支** | master |

---

## 🏗️ 项目结构

### 本地开发环境 (`/home/ubuntu/www.inarbit.work`)

```
www.inarbit.work/
├── client/                      # 前端代码
│   ├── public/                  # 静态资源
│   │   ├── images/              # 图片资源
│   │   └── favicon.ico
│   ├── src/
│   │   ├── pages/               # 页面组件
│   │   │   ├── Home.tsx
│   │   │   ├── HowItWorks.tsx
│   │   │   ├── Technical.tsx
│   │   │   ├── Dashboard.tsx
│   │   │   └── NotFound.tsx
│   │   ├── components/          # 可复用组件
│   │   ├── contexts/            # React上下文
│   │   ├── hooks/               # 自定义hooks
│   │   ├── lib/                 # 工具函数
│   │   ├── App.tsx              # 主应用组件（路由配置）
│   │   ├── main.tsx             # React入口
│   │   └── index.css            # 全局样式（Solarized Light主题）
│   └── index.html               # HTML模板
├── server/                      # 后端代码（Go）
│   └── index.ts                 # 后端入口
├── dist/                        # 生产构建输出
│   ├── public/                  # 前端生产构建
│   │   ├── index.html
│   │   └── assets/
│   └── index.js                 # 后端生产构建
├── vite.config.ts               # Vite配置
├── tsconfig.json                # TypeScript配置
├── package.json                 # 项目依赖
└── BROWSER_COMPATIBILITY_REPORT.md  # 浏览器兼容性报告
```

### 生产服务器 (`/var/www/inarbit`)

```
/var/www/inarbit/
├── frontend/
│   └── dist/                    # 前端生产构建（由Nginx提供）
│       ├── index.html
│       ├── assets/
│       └── images/
└── backend/                     # 后端服务（Go）
    └── inarbit                  # 可执行文件
```

---

## 🚀 部署流程

### 1. 本地开发

```bash
# 进入项目目录
cd /home/ubuntu/www.inarbit.work

# 安装依赖
pnpm install

# 启动开发服务器（端口3001）
pnpm run dev

# 访问 http://localhost:3001
```

### 2. 生产构建

```bash
# 构建前端和后端
pnpm run build

# 输出位置：
# - 前端: dist/public/
# - 后端: dist/index.js
```

### 3. 上传到服务器

```bash
# 上传前端构建
scp -r dist/public/* root@8.211.158.208:/var/www/inarbit/frontend/dist/

# 上传后端构建
scp dist/index.js root@8.211.158.208:/var/www/inarbit/backend/
```

### 4. 服务器配置

#### 4.1 Nginx配置

```bash
# 编辑Nginx配置
nano /etc/nginx/sites-available/inarbit

# 关键配置：
# location / {
#     try_files $uri $uri/ /index.html;  # 支持客户端路由
# }

# 测试配置
nginx -t

# 重新加载
systemctl reload nginx
```

#### 4.2 后端服务

```bash
# 启动后端服务（Go）
cd /var/www/inarbit/backend
./inarbit

# 或使用systemd管理
systemctl start inarbit
systemctl status inarbit
```

---

## 🔧 常见操作

### 查看日志

```bash
# Nginx访问日志
tail -f /var/log/nginx/inarbit_access.log

# Nginx错误日志
tail -f /var/log/nginx/inarbit_error.log

# 后端服务日志
journalctl -u inarbit -f
```

### 重启服务

```bash
# 重启Nginx
systemctl restart nginx

# 重启后端服务
systemctl restart inarbit

# 重新加载Nginx（不中断连接）
systemctl reload nginx
```

### 检查服务状态

```bash
# 检查Nginx
systemctl status nginx

# 检查后端服务
systemctl status inarbit

# 检查端口占用
netstat -tlnp | grep -E "80|443|5000"
```

### 清除缓存

```bash
# 清除浏览器缓存
# 在浏览器中按 Ctrl+Shift+Delete

# 清除Nginx缓存
rm -rf /var/cache/nginx/*
systemctl reload nginx
```

---

## 🔄 代码更新流程

### 步骤1: 更新本地代码

```bash
cd /home/ubuntu/www.inarbit.work

# 进行代码修改
# ...

# 提交到GitHub
git add .
git commit -m "描述你的更改"
git push origin master
```

### 步骤2: 本地测试

```bash
# 启动开发服务器
pnpm run dev

# 在浏览器中测试所有页面
# - http://localhost:3001/
# - http://localhost:3001/how-it-works
# - http://localhost:3001/technical
# - http://localhost:3001/dashboard
```

### 步骤3: 生产构建

```bash
# 构建
pnpm run build

# 验证构建输出
ls -la dist/public/
```

### 步骤4: 部署到服务器

```bash
# 备份当前版本
ssh root@8.211.158.208 "cp -r /var/www/inarbit/frontend/dist /var/www/inarbit/frontend/dist.backup.$(date +%Y%m%d_%H%M%S)"

# 上传新版本
scp -r dist/public/* root@8.211.158.208:/var/www/inarbit/frontend/dist/

# 重新加载Nginx
ssh root@8.211.158.208 "systemctl reload nginx"
```

### 步骤5: 验证部署

```bash
# 测试首页
curl -I http://8.211.158.208/

# 测试路由
curl -I http://8.211.158.208/how-it-works

# 在浏览器中验证
# http://8.211.158.208/
# http://8.211.158.208/how-it-works
# http://8.211.158.208/technical
# http://8.211.158.208/dashboard
```

---

## 🎨 设计主题: Solarized Light

### 颜色方案

| 用途 | 颜色 | 十六进制 |
|------|------|--------|
| **背景** | 米色 | #FDF6E3 |
| **文本** | 深灰色 | #657B83 |
| **主色** | 绿色 | #859900 |
| **辅色** | 棕色 | #B58900 |
| **强调色** | 蓝色 | #268BD2 |
| **错误** | 红色 | #DC322F |
| **警告** | 橙色 | #CB4B16 |

### 字体

- **标题**: Georgia, serif
- **正文**: Segoe UI, sans-serif
- **代码**: Courier New, monospace

### 修改主题

编辑 `client/src/index.css`：

```css
@layer base {
  :root {
    --background: #FDF6E3;
    --foreground: #657B83;
    --primary: #859900;
    --secondary: #B58900;
    --accent: #268BD2;
    /* ... 其他颜色 ... */
  }
}
```

---

## 📊 性能优化

### 1. Gzip压缩

Nginx配置已启用Gzip压缩。验证：

```bash
curl -I -H "Accept-Encoding: gzip" http://8.211.158.208/
```

预期看到：
```
Content-Encoding: gzip
```

### 2. 缓存策略

- **HTML**: 不缓存（Cache-Control: no-cache）
- **JS/CSS**: 1年缓存（Cache-Control: public, immutable）
- **图片**: 1年缓存（Cache-Control: public, immutable）

### 3. 代码分割

当前JavaScript包大小：645KB (gzip: 175KB)

优化建议：
- 使用动态导入分割路由
- 使用`React.lazy()`延迟加载组件
- 移除未使用的依赖

### 4. 图片优化

- 使用WebP格式（带PNG回退）
- 使用`<picture>`标签进行响应式图片
- 压缩图片大小

---

## 🔒 安全配置

### Nginx安全头

已配置以下安全头：

```nginx
add_header X-Frame-Options "SAMEORIGIN" always;
add_header X-Content-Type-Options "nosniff" always;
add_header X-XSS-Protection "1; mode=block" always;
add_header Referrer-Policy "no-referrer-when-downgrade" always;
```

### HTTPS配置

推荐启用HTTPS。使用Let's Encrypt免费证书：

```bash
# 安装certbot
apt-get install certbot python3-certbot-nginx

# 获取证书
certbot certonly --nginx -d www.inarbit.work -d inarbit.work

# 更新Nginx配置以使用HTTPS
nano /etc/nginx/sites-available/inarbit

# 重新加载
systemctl reload nginx
```

---

## 📈 监控和维护

### 1. 定期检查

```bash
# 检查磁盘空间
df -h

# 检查内存使用
free -h

# 检查进程
ps aux | grep -E "nginx|inarbit"

# 检查端口
netstat -tlnp | grep -E "80|443|5000"
```

### 2. 日志轮转

配置日志轮转以防止日志文件过大：

```bash
# 编辑logrotate配置
nano /etc/logrotate.d/nginx

# 手动运行logrotate
logrotate -f /etc/logrotate.d/nginx
```

### 3. 备份

定期备份重要文件：

```bash
# 备份前端
tar -czf inarbit_frontend_$(date +%Y%m%d).tar.gz /var/www/inarbit/frontend/

# 备份配置
tar -czf inarbit_config_$(date +%Y%m%d).tar.gz /etc/nginx/sites-available/inarbit

# 上传到本地
scp root@8.211.158.208:/root/inarbit_*.tar.gz ./backups/
```

---

## 🐛 故障排查

### 问题1: 访问路由显示首页

**原因**: Nginx没有正确配置`try_files`  
**解决方案**: 参考 `ROUTING_FIX_GUIDE.md`

### 问题2: 404错误

**原因**: 文件不存在或权限问题  
**解决方案**:
```bash
# 检查文件权限
ls -la /var/www/inarbit/frontend/dist/

# 修复权限
chmod -R 755 /var/www/inarbit/frontend/dist/
chown -R www-data:www-data /var/www/inarbit/frontend/dist/
```

### 问题3: Nginx无法启动

**原因**: 配置文件语法错误  
**解决方案**:
```bash
# 测试配置
nginx -t

# 查看详细错误
nginx -T

# 查看日志
tail -f /var/log/nginx/error.log
```

### 问题4: 后端服务无响应

**原因**: 服务未运行或端口被占用  
**解决方案**:
```bash
# 检查服务状态
systemctl status inarbit

# 重启服务
systemctl restart inarbit

# 检查端口
netstat -tlnp | grep 5000

# 查看日志
journalctl -u inarbit -n 50
```

---

## 📚 相关文档

- [Wouter文档](https://github.com/molefrog/wouter)
- [Vite文档](https://vitejs.dev/)
- [Tailwind CSS文档](https://tailwindcss.com/)
- [Nginx文档](https://nginx.org/en/docs/)
- [React文档](https://react.dev/)

---

## 👥 团队信息

| 角色 | 信息 |
|------|------|
| **GitHub用户** | zillafan80-Maxzilla |
| **服务器IP** | 8.211.158.208 |
| **服务器用户** | root |
| **本地开发目录** | /home/ubuntu/www.inarbit.work |

---

## 📞 支持和反馈

如有问题或建议，请：

1. 查看本文档中的故障排查部分
2. 检查相关日志文件
3. 参考GitHub仓库的Issue
4. 联系项目维护者

---

**最后更新**: 2026-01-01  
**维护者**: Manus AI Agent  
**版本**: 1.0.0

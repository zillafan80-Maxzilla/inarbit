# www.inarbit.work 路由问题修复指南

**问题**: 前端路由在生产环境中失效，所有路由都显示首页内容  
**原因**: Nginx的`try_files`指令需要正确配置以支持客户端路由  
**状态**: 需要在生产服务器上执行修复

---

## 🔍 问题诊断

### 症状
- URL: `http://8.211.158.208/how-it-works` → 显示首页内容
- URL: `http://8.211.158.208/technical` → 显示首页内容
- URL: `http://8.211.158.208/dashboard` → 显示首页内容
- URL: `http://8.211.158.208/` → 正确显示首页

### 根本原因
Wouter是一个客户端路由库，所有路由都应该由浏览器中的JavaScript处理。但Nginx在找不到对应的物理文件时返回404，而不是回退到`index.html`。

---

## ✅ 解决方案

### 步骤1: 检查当前Nginx配置

连接到服务器：
```bash
ssh root@8.211.158.208
```

查看Nginx配置：
```bash
cat /etc/nginx/sites-available/inarbit
```

### 步骤2: 更新Nginx配置

编辑Nginx配置文件：
```bash
nano /etc/nginx/sites-available/inarbit
```

**关键配置**: 在`location /`块中添加或修改`try_files`指令：

```nginx
server {
    listen 80;
    listen [::]:80;
    server_name www.inarbit.work inarbit.work 8.211.158.208;

    # 重定向HTTP到HTTPS（如果需要）
    # return 301 https://$server_name$request_uri;

    root /var/www/inarbit/frontend/dist;
    index index.html;

    # 关键：为客户端路由配置try_files
    location / {
        # 尝试文件 → 尝试目录 → 回退到index.html
        try_files $uri $uri/ /index.html;
    }

    # 静态资源缓存
    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }

    # 禁用缓存HTML文件
    location ~* \.html$ {
        expires -1;
        add_header Cache-Control "no-cache, no-store, must-revalidate";
    }

    # 反向代理到后端API（如果需要）
    # location /api/ {
    #     proxy_pass http://localhost:5000;
    #     proxy_set_header Host $host;
    #     proxy_set_header X-Real-IP $remote_addr;
    #     proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    #     proxy_set_header X-Forwarded-Proto $scheme;
    # }
}
```

### 步骤3: 测试Nginx配置

```bash
nginx -t
```

预期输出：
```
nginx: the configuration file /etc/nginx/nginx.conf syntax is ok
nginx: configuration file /etc/nginx/nginx.conf test is successful
```

### 步骤4: 重新加载Nginx

```bash
systemctl reload nginx
```

或者：
```bash
service nginx reload
```

### 步骤5: 验证修复

在浏览器中测试以下URL：

| URL | 预期结果 |
|-----|--------|
| `http://8.211.158.208/` | 显示首页 |
| `http://8.211.158.208/how-it-works` | 显示"原理解析"页面 |
| `http://8.211.158.208/technical` | 显示"技术实现"页面 |
| `http://8.211.158.208/dashboard` | 显示"仪表板"页面 |
| `http://8.211.158.208/nonexistent` | 显示404页面 |

---

## 🔧 完整的Nginx配置示例

以下是完整的推荐配置：

```nginx
# /etc/nginx/sites-available/inarbit

upstream backend {
    server localhost:5000;
}

server {
    listen 80;
    listen [::]:80;
    server_name www.inarbit.work inarbit.work 8.211.158.208;

    # 日志
    access_log /var/log/nginx/inarbit_access.log;
    error_log /var/log/nginx/inarbit_error.log;

    # 前端静态文件
    root /var/www/inarbit/frontend/dist;
    index index.html;

    # Gzip压缩
    gzip on;
    gzip_types text/plain text/css text/javascript application/javascript application/json;
    gzip_min_length 1000;
    gzip_vary on;

    # 安全头
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "no-referrer-when-downgrade" always;

    # 客户端路由配置
    location / {
        try_files $uri $uri/ /index.html;
        
        # 禁用缓存HTML
        add_header Cache-Control "no-cache, no-store, must-revalidate";
    }

    # 静态资源缓存
    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }

    # 后端API反向代理
    location /api/ {
        proxy_pass http://backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_redirect off;
    }

    # 拒绝访问隐藏文件
    location ~ /\. {
        deny all;
    }
}

# 可选：HTTPS重定向
# server {
#     listen 443 ssl http2;
#     listen [::]:443 ssl http2;
#     server_name www.inarbit.work inarbit.work;
#
#     ssl_certificate /etc/letsencrypt/live/inarbit.work/fullchain.pem;
#     ssl_certificate_key /etc/letsencrypt/live/inarbit.work/privkey.pem;
#
#     # ... 其他配置同上 ...
# }
#
# server {
#     listen 80;
#     listen [::]:80;
#     server_name www.inarbit.work inarbit.work;
#     return 301 https://$server_name$request_uri;
# }
```

---

## 📋 检查清单

- [ ] SSH连接到服务器：`ssh root@8.211.158.208`
- [ ] 查看当前Nginx配置：`cat /etc/nginx/sites-available/inarbit`
- [ ] 编辑Nginx配置：`nano /etc/nginx/sites-available/inarbit`
- [ ] 添加`try_files $uri $uri/ /index.html;`到`location /`块
- [ ] 测试配置：`nginx -t`
- [ ] 重新加载Nginx：`systemctl reload nginx`
- [ ] 在浏览器中测试所有路由
- [ ] 验证404页面正常工作
- [ ] 检查Nginx日志：`tail -f /var/log/nginx/inarbit_error.log`

---

## 🚨 常见问题

### Q: 修改后仍然显示首页？
**A**: 
1. 确保Nginx配置已正确保存
2. 运行`nginx -t`验证语法
3. 运行`systemctl reload nginx`重新加载
4. 清除浏览器缓存（Ctrl+Shift+Delete）
5. 检查Nginx错误日志：`tail -f /var/log/nginx/inarbit_error.log`

### Q: 404页面不显示？
**A**: 这是正常的。当访问不存在的路由时，Nginx会回退到`index.html`，由React的Wouter库显示404页面。

### Q: 如何验证修复成功？
**A**: 
```bash
# 测试首页
curl -I http://8.211.158.208/

# 测试子路由
curl -I http://8.211.158.208/how-it-works

# 检查返回的HTML是否相同（都应该是index.html）
curl http://8.211.158.208/ | head -20
curl http://8.211.158.208/how-it-works | head -20
```

---

## 📚 相关文档

- [Wouter文档](https://github.com/molefrog/wouter)
- [Nginx try_files指令](https://nginx.org/en/docs/http/ngx_http_rewrite_module.html#try_files)
- [React Router配置](https://reactrouter.com/en/main/guides/ssr)

---

## 🔄 后续步骤

1. **应用修复**: 按照上述步骤修改Nginx配置
2. **测试验证**: 在浏览器中测试所有路由
3. **监控日志**: 观察Nginx错误日志以发现任何问题
4. **性能优化**: 考虑启用HTTPS和HTTP/2
5. **备份配置**: 将修改后的配置保存到GitHub

---

**最后更新**: 2026-01-01  
**作者**: Manus AI Agent  
**状态**: 待执行

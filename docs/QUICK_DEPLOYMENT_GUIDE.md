# iNarbit GitHub Actions 快速部署指南

## 🚀 5分钟快速开始

### 第一步：准备工作（1分钟）

1. **生成SSH密钥**（在您的服务器上）
```bash
ssh-keygen -t rsa -b 4096 -f /root/.ssh/github_deploy -N ""
cat /root/.ssh/github_deploy  # 复制私钥
cat /root/.ssh/github_deploy.pub >> /root/.ssh/authorized_keys
```

2. **生成Binance API密钥**（在Binance账户中）
   - 登录 https://www.binance.com
   - 进入 Account → API Management
   - 创建新的API密钥
   - **重要**：设置以下限制：
     - ✅ 仅允许现货交易
     - ✅ 仅允许特定IP（8.211.158.208）
     - ✅ 禁用提现权限

### 第二步：配置GitHub Secrets（2分钟）

1. 进入 https://github.com/zillafan80-Maxzilla/inarbit
2. 点击 **Settings → Secrets and variables → Actions**
3. 点击 **New repository secret**，添加以下4个secrets：

| Secret名称 | 值 |
|-----------|-----|
| `DEPLOY_SSH_KEY` | 从第一步复制的私钥 |
| `DEPLOY_HOST` | `8.211.158.208` |
| `BINANCE_API_KEY` | 您的Binance API密钥 |
| `BINANCE_API_SECRET` | 您的Binance API Secret |

### 第三步：上传工作流文件（1分钟）

```bash
cd /root/inarbit
mkdir -p .github/workflows
cp /home/ubuntu/github_actions_deploy.yml .github/workflows/deploy.yml
git add .github/workflows/deploy.yml
git commit -m "添加GitHub Actions自动化部署"
git push origin main
```

### 第四步：监控部署（1分钟）

1. 进入 https://github.com/zillafan80-Maxzilla/inarbit/actions
2. 查看 **自动化部署 iNarbit** 工作流
3. 等待部署完成（约5-10分钟）

---

## ✅ 部署完成后的验证

### 检查项目状态

```bash
# SSH登录服务器
ssh root@8.211.158.208

# 检查服务状态
supervisorctl status

# 查看日志
tail -f /root/inarbit/logs/inarbit.log

# 测试API
curl http://localhost:8080/api/health

# 测试Web界面
curl https://inarbit.work
```

### 访问应用

- **Web界面**：https://inarbit.work
- **API文档**：https://inarbit.work/api/docs
- **仪表板**：https://inarbit.work/dashboard

### 默认登录凭证

```
用户名: admin
密码: password
```

⚠️ **立即修改默认密码！**

---

## 🔄 自动化工作流说明

### 工作流触发

每次您推送代码到 `main` 分支时，GitHub Actions会自动：

1. ✅ 检查代码质量
2. ✅ 运行单元测试
3. ✅ 编译后端和前端
4. ✅ 部署到服务器
5. ✅ 运行集成测试
6. ✅ 运行性能测试
7. ✅ 生成部署报告

### 工作流状态

在 **Actions** 标签中查看：
- 🟢 **Success**：部署成功
- 🔴 **Failed**：部署失败
- 🟡 **In Progress**：正在部署

---

## 🛠️ 常见操作

### 查看部署日志

```
Actions → 选择工作流 → 点击运行 → 查看日志
```

### 手动触发部署

```
Actions → 自动化部署 iNarbit → Run workflow → 选择分支 → Run workflow
```

### 查看部署报告

```
Actions → 最近的运行 → Artifacts → 下载 deployment-report
```

---

## 🚨 故障排除

### 部署失败？

1. **查看GitHub Actions日志**
   - 进入 Actions → 查看失败的运行
   - 展开失败的步骤查看错误信息

2. **常见错误**：

| 错误 | 原因 | 解决方案 |
|------|------|--------|
| SSH连接失败 | 密钥配置错误 | 检查 DEPLOY_SSH_KEY secret |
| 数据库错误 | 数据库未初始化 | 在服务器上手动运行初始化脚本 |
| 编译失败 | 依赖缺失 | 检查Go和Node.js版本 |
| 部署超时 | 网络问题 | 增加超时时间或检查网络 |

### 服务无法启动？

```bash
# SSH登录服务器
ssh root@8.211.158.208

# 查看Supervisor日志
supervisorctl tail inarbit-backend

# 查看Nginx日志
tail -f /var/log/nginx/error.log

# 重启服务
supervisorctl restart inarbit-backend
supervisorctl restart inarbit-frontend
systemctl restart nginx
```

---

## 📊 部署后的建议

### 立即执行（高优先级）

- [ ] 修改默认密码
- [ ] 测试Binance连接
- [ ] 运行虚拟盘测试
- [ ] 配置数据库备份
- [ ] 配置监控告警

### 本周执行（中优先级）

- [ ] 完善机器人管理功能
- [ ] 实现交易记录查询
- [ ] 性能优化
- [ ] 安全加固

### 本月执行（低优先级）

- [ ] 四角/五角套利
- [ ] 机器学习集成
- [ ] 多交易所支持
- [ ] 期货交易支持

---

## 📞 需要帮助？

### 查看详细文档

- **GitHub Actions配置**：`GITHUB_ACTIONS_SETUP.md`
- **部署后建议**：`NEXT_STEPS.md`
- **Git同步指南**：`GIT_SYNC_GUIDE.md`
- **测试指南**：`TESTING_GUIDE.md`

### 常用命令

```bash
# 查看工作流列表
gh workflow list

# 查看最近的运行
gh run list --workflow=deploy.yml

# 查看运行日志
gh run view <run-id> --log

# 手动触发工作流
gh workflow run deploy.yml
```

---

## ✨ 总结

通过GitHub Actions自动化部署，您可以：

1. ✅ **完全自动化**：推送代码 → 自动部署
2. ✅ **快速反馈**：5-10分钟内完成部署
3. ✅ **质量保证**：自动化测试确保代码质量
4. ✅ **零停机部署**：无需手动操作
5. ✅ **完整报告**：自动生成部署报告

**下一步**：

1. 按照上面的步骤配置GitHub Actions
2. 推送代码到GitHub
3. 监控部署过程
4. 验证应用是否正常运行

祝您的自动化部署顺利！🎉

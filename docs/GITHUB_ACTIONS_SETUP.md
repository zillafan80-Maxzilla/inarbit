# GitHub Actions 自动化部署配置指南

## 目录

1. [概述](#概述)
2. [前置条件](#前置条件)
3. [配置步骤](#配置步骤)
4. [工作流说明](#工作流说明)
5. [故障排除](#故障排除)
6. [监控和维护](#监控和维护)

---

## 概述

本指南将帮助您配置GitHub Actions，实现完全自动化的部署、测试和同步流程。

**工作流包括**：
- ✅ 代码质量检查
- ✅ 后端单元测试
- ✅ 前端构建测试
- ✅ 自动部署到服务器
- ✅ 集成测试
- ✅ 性能测试
- ✅ 部署报告生成

---

## 前置条件

### 1. GitHub仓库

您需要有一个GitHub仓库：
```
https://github.com/zillafan80-Maxzilla/inarbit
```

### 2. 服务器SSH密钥

首先，在您的服务器上生成SSH密钥对：

```bash
# 在服务器上生成密钥对
ssh-keygen -t rsa -b 4096 -f /root/.ssh/github_deploy -N ""

# 显示私钥（用于GitHub）
cat /root/.ssh/github_deploy

# 显示公钥（用于服务器）
cat /root/.ssh/github_deploy.pub
```

### 3. GitHub Personal Access Token

创建GitHub Personal Access Token用于GitHub Actions：

1. 登录GitHub账户
2. 进入 **Settings → Developer settings → Personal access tokens**
3. 点击 **Generate new token (classic)**
4. 选择以下权限：
   - `repo` (完整的仓库访问)
   - `workflow` (GitHub Actions)
   - `admin:repo_hook` (仓库钩子)
5. 生成token并保存

---

## 配置步骤

### 第一步：将工作流文件添加到仓库

```bash
# 克隆仓库
git clone https://github.com/zillafan80-Maxzilla/inarbit.git
cd inarbit

# 创建GitHub Actions目录
mkdir -p .github/workflows

# 复制工作流文件
cp /home/ubuntu/github_actions_deploy.yml .github/workflows/deploy.yml

# 提交并推送
git add .github/workflows/deploy.yml
git commit -m "添加GitHub Actions自动化部署工作流"
git push origin main
```

### 第二步：配置GitHub Secrets

在GitHub仓库中配置以下secrets：

#### 2.1 添加SSH部署密钥

1. 进入仓库 **Settings → Secrets and variables → Actions**
2. 点击 **New repository secret**
3. 创建以下secrets：

**Secret 1: DEPLOY_SSH_KEY**
```
Name: DEPLOY_SSH_KEY
Value: (粘贴您的私钥内容，从 ssh-keygen 生成的)
```

**Secret 2: DEPLOY_HOST**
```
Name: DEPLOY_HOST
Value: 8.211.158.208
```

**Secret 3: BINANCE_API_KEY**
```
Name: BINANCE_API_KEY
Value: (您的新的受限Binance API密钥)
```

**Secret 4: BINANCE_API_SECRET**
```
Name: BINANCE_API_SECRET
Value: (您的Binance API密钥对应的Secret)
```

### 第三步：配置服务器端

在服务器上配置SSH密钥：

```bash
# 添加GitHub的公钥到authorized_keys
cat /root/.ssh/github_deploy.pub >> /root/.ssh/authorized_keys

# 设置正确的权限
chmod 700 /root/.ssh
chmod 600 /root/.ssh/authorized_keys

# 验证SSH连接
ssh -i /root/.ssh/github_deploy -o StrictHostKeyChecking=no root@8.211.158.208 "echo 'SSH连接成功'"
```

### 第四步：配置GitHub仓库设置

1. 进入仓库 **Settings**
2. 在 **Actions → General** 中：
   - 选择 **Allow all actions and reusable workflows**
   - 选择 **Read and write permissions**
   - 选择 **Allow GitHub Actions to create and approve pull requests**

### 第五步：验证配置

```bash
# 在本地测试SSH连接
ssh -i ~/.ssh/github_deploy -o StrictHostKeyChecking=no root@8.211.158.208 "whoami"

# 应该输出：root
```

---

## 工作流说明

### 工作流触发条件

工作流在以下情况下自动触发：

1. **推送到main分支**
   ```bash
   git push origin main
   ```

2. **推送到develop分支**（仅代码检查和测试）
   ```bash
   git push origin develop
   ```

3. **手动触发**
   - 进入仓库 **Actions** 标签
   - 选择 **自动化部署 iNarbit** 工作流
   - 点击 **Run workflow**

### 工作流步骤

#### 1. 代码质量检查 (code-quality)
- Go代码格式检查
- Go代码静态分析
- 前端代码检查

#### 2. 后端单元测试 (backend-test)
- 启动PostgreSQL服务
- 初始化数据库
- 运行后端单元测试
- 生成覆盖率报告

#### 3. 前端构建测试 (frontend-test)
- 安装npm依赖
- 构建React应用
- 上传构建产物

#### 4. 部署到服务器 (deploy)
- 设置SSH连接
- 拉取最新代码
- 初始化数据库
- 编译后端
- 编译前端
- 重启服务
- 执行健康检查

#### 5. 集成测试 (integration-test)
- 测试Binance连接
- 测试数据库连接
- 测试API端点
- 测试Web界面
- 测试WebSocket

#### 6. 性能测试 (performance-test)
- API响应时间测试
- 数据库查询性能测试

#### 7. 生成报告 (generate-report)
- 生成部署报告
- 上传到GitHub Pages

### 工作流依赖关系

```
代码检查
    ↓
┌───────────────────────┐
│ 后端测试  前端测试     │
└───────────────────────┘
    ↓
部署到服务器
    ↓
┌───────────────────────┐
│ 集成测试  性能测试     │
└───────────────────────┘
    ↓
生成报告
```

---

## 故障排除

### 问题1：SSH连接失败

**错误信息**：`Permission denied (publickey)`

**解决方案**：
```bash
# 检查服务器上的公钥
cat /root/.ssh/authorized_keys

# 确保包含GitHub的公钥
cat /root/.ssh/github_deploy.pub >> /root/.ssh/authorized_keys

# 检查权限
chmod 700 /root/.ssh
chmod 600 /root/.ssh/authorized_keys
```

### 问题2：部署失败

**错误信息**：`git: command not found`

**解决方案**：
```bash
# 在服务器上安装Git
apt-get update
apt-get install -y git

# 配置Git
git config --global user.name "GitHub Actions"
git config --global user.email "actions@github.com"
```

### 问题3：数据库连接失败

**错误信息**：`FATAL: role "inarbit" does not exist`

**解决方案**：
```bash
# 在服务器上初始化数据库
sudo -u postgres psql << EOF
CREATE USER inarbit WITH PASSWORD 'inarbit_password';
CREATE DATABASE inarbit OWNER inarbit;
GRANT ALL PRIVILEGES ON DATABASE inarbit TO inarbit;
EOF
```

### 问题4：Go编译失败

**错误信息**：`go: command not found`

**解决方案**：
```bash
# 在服务器上安装Go
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```

### 问题5：前端编译失败

**错误信息**：`npm: command not found`

**解决方案**：
```bash
# 在服务器上安装Node.js
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
apt-get install -y nodejs
```

---

## 监控和维护

### 1. 查看工作流运行

```
进入仓库 → Actions → 选择工作流 → 查看运行历史
```

### 2. 查看工作流日志

```
点击具体的运行 → 查看各个步骤的日志
```

### 3. 配置工作流通知

1. 进入GitHub **Settings → Notifications**
2. 启用 **Actions** 通知
3. 选择通知方式（邮件、Web等）

### 4. 定期检查

```bash
# 检查服务器部署状态
ssh root@8.211.158.208 "supervisorctl status"

# 检查服务日志
ssh root@8.211.158.208 "tail -f /root/inarbit/logs/inarbit.log"

# 检查数据库
ssh root@8.211.158.208 "PGPASSWORD=inarbit_password psql -h localhost -U inarbit -d inarbit -c 'SELECT COUNT(*) FROM trades;'"
```

### 5. 性能优化

```yaml
# 在 .github/workflows/deploy.yml 中调整
- 增加缓存以加快构建速度
- 并行运行独立的测试
- 使用自托管运行器以获得更好的性能
```

---

## 高级配置

### 1. 条件部署

```yaml
# 仅在特定标签上部署
on:
  push:
    tags:
      - 'v*'
```

### 2. 环境特定配置

```yaml
# 不同环境的不同配置
env:
  STAGING_DOMAIN: staging.inarbit.work
  PRODUCTION_DOMAIN: inarbit.work
```

### 3. 自定义通知

```yaml
# 使用Slack通知
- name: Slack通知
  uses: slackapi/slack-github-action@v1
  with:
    webhook-url: ${{ secrets.SLACK_WEBHOOK }}
```

### 4. 自动化发布

```yaml
# 自动创建GitHub Release
- name: 创建Release
  uses: actions/create-release@v1
  env:
    GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

---

## 安全最佳实践

### 1. 密钥管理

- ✅ 使用GitHub Secrets存储敏感信息
- ✅ 定期轮换SSH密钥
- ✅ 不要在日志中打印密钥
- ✅ 使用受限的API密钥

### 2. 访问控制

- ✅ 限制谁可以触发工作流
- ✅ 使用分支保护规则
- ✅ 要求代码审查
- ✅ 启用两因素认证

### 3. 审计日志

```bash
# 查看GitHub Actions审计日志
# 进入仓库 → Settings → Audit log
```

---

## 完整的部署流程

```
1. 开发者推送代码到GitHub
   ↓
2. GitHub Actions自动触发工作流
   ↓
3. 代码质量检查和测试
   ↓
4. 如果测试通过，自动部署到服务器
   ↓
5. 执行集成测试和性能测试
   ↓
6. 生成部署报告
   ↓
7. 部署完成，系统自动更新
```

---

## 常用命令

```bash
# 查看工作流状态
gh workflow list

# 查看最近的运行
gh run list --workflow=deploy.yml

# 查看特定运行的日志
gh run view <run-id> --log

# 手动触发工作流
gh workflow run deploy.yml -f deploy_type=full

# 取消正在运行的工作流
gh run cancel <run-id>
```

---

## 总结

通过GitHub Actions自动化部署，您可以：

1. ✅ 完全自动化的部署流程
2. ✅ 自动化的测试和质量检查
3. ✅ 快速发现和修复问题
4. ✅ 减少人工操作错误
5. ✅ 提高开发效率

**下一步**：

1. 按照配置步骤设置GitHub Actions
2. 推送代码到GitHub
3. 监控工作流运行
4. 根据需要调整工作流

祝您的自动化部署顺利！

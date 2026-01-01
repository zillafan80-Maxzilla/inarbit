# iNarbit Git同步和部署指南

## 目录

1. [Git配置](#git配置)
2. [代码同步](#代码同步)
3. [部署流程](#部署流程)
4. [故障排除](#故障排除)

---

## Git配置

### 1. 初始化Git仓库

```bash
cd /root/inarbit

# 初始化Git
git init

# 添加远程仓库
git remote add origin https://github.com/zillafan80-Maxzilla/inarbit.git

# 配置用户信息
git config user.name "zillafan80-Maxzilla"
git config user.email "your-email@example.com"

# 配置凭证（使用Personal Access Token）
git config credential.helper store
```

### 2. 生成GitHub Personal Access Token

1. 登录GitHub账户
2. 进入 Settings → Developer settings → Personal access tokens
3. 点击 "Generate new token"
4. 选择以下权限：
   - `repo` (完整的仓库访问)
   - `workflow` (GitHub Actions)
5. 生成token并保存

### 3. 配置SSH密钥（推荐）

```bash
# 生成SSH密钥
ssh-keygen -t rsa -b 4096 -f ~/.ssh/id_rsa -N ""

# 显示公钥
cat ~/.ssh/id_rsa.pub

# 将公钥添加到GitHub账户
# Settings → SSH and GPG keys → New SSH key
```

### 4. 配置Git使用SSH

```bash
# 修改远程仓库URL为SSH格式
git remote set-url origin git@github.com:zillafan80-Maxzilla/inarbit.git

# 验证配置
git remote -v
```

---

## 代码同步

### 1. 从GitHub拉取代码

```bash
cd /root/inarbit

# 拉取最新代码
git pull origin main

# 或者强制拉取（覆盖本地更改）
git fetch origin
git reset --hard origin/main
```

### 2. 提交本地更改到GitHub

```bash
cd /root/inarbit

# 查看修改的文件
git status

# 添加所有更改
git add .

# 提交更改
git commit -m "描述你的更改"

# 推送到GitHub
git push origin main
```

### 3. 完整的同步流程

```bash
#!/bin/bash
# 同步脚本：sync_to_github.sh

cd /root/inarbit

# 1. 拉取最新代码
echo "拉取最新代码..."
git pull origin main

# 2. 添加所有更改
echo "添加更改..."
git add .

# 3. 提交更改
echo "提交更改..."
git commit -m "自动同步: $(date '+%Y-%m-%d %H:%M:%S')"

# 4. 推送到GitHub
echo "推送到GitHub..."
git push origin main

echo "同步完成！"
```

### 4. 自动同步脚本

```bash
#!/bin/bash
# 自动同步脚本：auto_sync.sh

PROJECT_DIR="/root/inarbit"
SYNC_INTERVAL=3600  # 1小时同步一次

while true; do
    cd ${PROJECT_DIR}
    
    # 检查是否有更改
    if [ -n "$(git status --porcelain)" ]; then
        echo "[$(date '+%Y-%m-%d %H:%M:%S')] 检测到更改，开始同步..."
        
        git add .
        git commit -m "自动同步: $(date '+%Y-%m-%d %H:%M:%S')"
        git push origin main
        
        echo "同步完成"
    fi
    
    sleep ${SYNC_INTERVAL}
done
```

### 5. 使用Cron定期同步

```bash
# 编辑crontab
crontab -e

# 添加以下行（每小时同步一次）
0 * * * * cd /root/inarbit && git add . && git commit -m "自动同步: $(date '+\%Y-\%m-\%d \%H:\%M:\%S')" && git push origin main

# 或者每天同步一次
0 2 * * * cd /root/inarbit && git add . && git commit -m "自动同步: $(date '+\%Y-\%m-\%d \%H:\%M:\%S')" && git push origin main
```

---

## 部署流程

### 完整的部署和同步流程

```bash
#!/bin/bash
# 完整部署脚本：deploy_and_sync.sh

set -e

PROJECT_DIR="/root/inarbit"
GITHUB_REPO="https://github.com/zillafan80-Maxzilla/inarbit.git"

echo "=========================================="
echo "iNarbit 完整部署和同步流程"
echo "=========================================="

# 第一步：从GitHub拉取最新代码
echo ""
echo "[步骤1] 从GitHub拉取最新代码..."
cd ${PROJECT_DIR}
git pull origin main

# 第二步：编译后端
echo ""
echo "[步骤2] 编译后端..."
cd ${PROJECT_DIR}/backend
go mod tidy
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o inarbit-server main.go

# 第三步：编译前端
echo ""
echo "[步骤3] 编译前端..."
cd ${PROJECT_DIR}/frontend
npm install
npm run build

# 第四步：重启服务
echo ""
echo "[步骤4] 重启服务..."
supervisorctl restart inarbit-backend
supervisorctl restart inarbit-frontend
systemctl restart nginx

# 第五步：执行测试
echo ""
echo "[步骤5] 执行测试..."
cd ${PROJECT_DIR}
bash COMPLETE_TEST_SUITE.sh all

# 第六步：提交更改到GitHub
echo ""
echo "[步骤6] 提交更改到GitHub..."
cd ${PROJECT_DIR}
git add .
git commit -m "部署: $(date '+%Y-%m-%d %H:%M:%S')" || true
git push origin main

echo ""
echo "=========================================="
echo "部署和同步完成！"
echo "=========================================="
```

---

## 故障排除

### 1. Git认证失败

**问题**：`fatal: Authentication failed`

**解决方案**：
```bash
# 检查SSH配置
ssh -T git@github.com

# 如果失败，重新配置SSH密钥
ssh-keygen -t rsa -b 4096 -f ~/.ssh/id_rsa -N ""
cat ~/.ssh/id_rsa.pub  # 复制到GitHub

# 或者使用HTTPS和Personal Access Token
git remote set-url origin https://github.com/zillafan80-Maxzilla/inarbit.git
git config credential.helper store
```

### 2. 合并冲突

**问题**：`CONFLICT (content merge)`

**解决方案**：
```bash
# 查看冲突文件
git status

# 编辑冲突文件，解决冲突
vim <conflict-file>

# 标记为已解决
git add <conflict-file>

# 完成合并
git commit -m "解决合并冲突"
```

### 3. 推送被拒绝

**问题**：`rejected ... (non-fast-forward)`

**解决方案**：
```bash
# 拉取最新代码
git pull origin main

# 解决冲突（如果有）
# 然后重新推送
git push origin main

# 或者强制推送（谨慎使用）
git push -f origin main
```

### 4. 大文件限制

**问题**：`File is too large`

**解决方案**：
```bash
# 使用Git LFS处理大文件
git lfs install

# 跟踪大文件
git lfs track "*.bin"
git lfs track "*.zip"

# 添加和提交
git add .gitattributes
git commit -m "配置Git LFS"
```

---

## 最佳实践

### 1. 提交信息规范

```bash
# 好的提交信息
git commit -m "feat: 添加三角套利引擎"
git commit -m "fix: 修复WebSocket连接问题"
git commit -m "docs: 更新部署文档"
git commit -m "refactor: 优化数据库查询"
git commit -m "test: 添加单元测试"

# 不好的提交信息
git commit -m "update"
git commit -m "fix bug"
git commit -m "changes"
```

### 2. 分支管理

```bash
# 创建功能分支
git checkout -b feature/triangular-arbitrage

# 开发完成后合并到main
git checkout main
git merge feature/triangular-arbitrage

# 删除功能分支
git branch -d feature/triangular-arbitrage

# 推送到GitHub
git push origin main
```

### 3. 定期备份

```bash
# 创建备份
tar -czf /root/backups/inarbit-$(date +%Y%m%d).tar.gz /root/inarbit

# 或使用Git标签
git tag -a v1.0.0 -m "版本1.0.0"
git push origin v1.0.0
```

### 4. 忽略文件

创建 `.gitignore` 文件：

```
# 依赖
node_modules/
vendor/

# 编译输出
dist/
build/
*.o
*.a

# 环境变量
.env
.env.local

# 日志
logs/
*.log

# 临时文件
*.tmp
*.swp
.DS_Store

# IDE
.vscode/
.idea/
*.iml

# 数据库
*.db
*.sqlite

# 密钥
*.key
*.pem
```

---

## 监控和维护

### 1. 检查仓库状态

```bash
cd /root/inarbit

# 查看当前分支
git branch -v

# 查看提交历史
git log --oneline -10

# 查看远程仓库状态
git remote -v

# 查看未推送的提交
git log origin/main..main
```

### 2. 清理本地仓库

```bash
# 删除未跟踪的文件
git clean -fd

# 删除本地已删除的分支
git remote prune origin

# 压缩仓库
git gc --aggressive
```

### 3. 恢复历史版本

```bash
# 查看特定文件的历史
git log --oneline <file>

# 恢复特定版本
git checkout <commit-hash> -- <file>

# 恢复整个项目到特定版本
git reset --hard <commit-hash>
```

---

## 自动化部署

### 使用GitHub Actions

创建 `.github/workflows/deploy.yml`：

```yaml
name: Deploy

on:
  push:
    branches: [ main ]

jobs:
  deploy:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v2
    
    - name: Deploy to Server
      env:
        DEPLOY_KEY: ${{ secrets.DEPLOY_KEY }}
        DEPLOY_HOST: ${{ secrets.DEPLOY_HOST }}
      run: |
        mkdir -p ~/.ssh
        echo "$DEPLOY_KEY" > ~/.ssh/deploy_key
        chmod 600 ~/.ssh/deploy_key
        ssh -i ~/.ssh/deploy_key root@$DEPLOY_HOST 'cd /root/inarbit && git pull && bash COMPLETE_DEPLOYMENT.sh deploy'
```

---

## 总结

通过正确的Git配置和同步流程，您可以：

1. ✅ 安全地管理代码
2. ✅ 自动化部署流程
3. ✅ 跟踪所有更改
4. ✅ 快速恢复错误
5. ✅ 团队协作开发

**关键命令速查**：

```bash
# 拉取最新代码
git pull origin main

# 提交本地更改
git add . && git commit -m "描述" && git push origin main

# 查看状态
git status

# 查看历史
git log --oneline

# 恢复文件
git checkout -- <file>

# 撤销提交
git revert <commit-hash>
```

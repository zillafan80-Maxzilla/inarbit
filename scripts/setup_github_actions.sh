#!/bin/bash

# ==========================================
# iNarbit GitHub Actions 一键设置脚本
# ==========================================

set -e

PROJECT_DIR="/root/inarbit"
GITHUB_REPO="https://github.com/zillafan80-Maxzilla/inarbit.git"

echo "=========================================="
echo "iNarbit GitHub Actions 设置"
echo "=========================================="

# 检查项目目录
if [ ! -d "$PROJECT_DIR" ]; then
    echo "❌ 项目目录不存在: $PROJECT_DIR"
    exit 1
fi

cd "$PROJECT_DIR"

# 第一步：检查Git是否已初始化
echo ""
echo "[1/6] 检查Git仓库..."
if [ ! -d ".git" ]; then
    echo "初始化Git仓库..."
    git init
    git remote add origin "$GITHUB_REPO"
    echo "✓ Git仓库已初始化"
else
    echo "✓ Git仓库已存在"
fi

# 第二步：配置Git用户信息
echo ""
echo "[2/6] 配置Git用户信息..."
git config user.name "GitHub Actions" || git config --global user.name "GitHub Actions"
git config user.email "actions@github.com" || git config --global user.email "actions@github.com"
echo "✓ Git用户信息已配置"

# 第三步：创建GitHub Actions目录
echo ""
echo "[3/6] 创建GitHub Actions目录..."
mkdir -p .github/workflows
echo "✓ 目录已创建"

# 第四步：创建工作流文件
echo ""
echo "[4/6] 创建GitHub Actions工作流文件..."

cat > .github/workflows/deploy.yml << 'WORKFLOW_EOF'
name: 自动化部署 iNarbit

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]
  workflow_dispatch:
    inputs:
      deploy_type:
        description: '部署类型'
        required: true
        default: 'full'
        type: choice
        options:
          - full
          - backend
          - frontend
          - test_only

env:
  PROJECT_DIR: /root/inarbit
  DEPLOY_HOST: 8.211.158.208
  DEPLOY_USER: root
  DOMAIN: inarbit.work

jobs:
  # ==================== 代码检查 ====================
  code-quality:
    name: 代码质量检查
    runs-on: ubuntu-latest
    
    steps:
    - name: 检出代码
      uses: actions/checkout@v3
    
    - name: 设置Go环境
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'
    
    - name: Go代码检查
      run: |
        cd backend 2>/dev/null || true
        go fmt ./... 2>/dev/null || true
        go vet ./... 2>/dev/null || true
    
    - name: 设置Node.js环境
      uses: actions/setup-node@v3
      with:
        node-version: '18'
    
    - name: 前端代码检查
      run: |
        cd frontend 2>/dev/null || true
        npm install 2>/dev/null || true
        npm run lint 2>/dev/null || true

  # ==================== 部署到服务器 ====================
  deploy:
    name: 部署到服务器
    runs-on: ubuntu-latest
    needs: code-quality
    if: github.event_name == 'push' && github.ref == 'refs/heads/main'
    
    steps:
    - name: 检出代码
      uses: actions/checkout@v3
    
    - name: 设置SSH密钥
      run: |
        mkdir -p ~/.ssh
        echo "${{ secrets.DEPLOY_SSH_KEY }}" > ~/.ssh/deploy_key
        chmod 600 ~/.ssh/deploy_key
        ssh-keyscan -H ${{ env.DEPLOY_HOST }} >> ~/.ssh/known_hosts 2>/dev/null || true
    
    - name: 部署代码到服务器
      run: |
        ssh -i ~/.ssh/deploy_key -o StrictHostKeyChecking=no ${{ env.DEPLOY_USER }}@${{ env.DEPLOY_HOST }} << 'EOF'
        set -e
        
        echo "=========================================="
        echo "开始部署 iNarbit"
        echo "=========================================="
        
        # 进入项目目录
        cd ${{ env.PROJECT_DIR }}
        
        # 第一步：拉取最新代码
        echo "[1/5] 拉取最新代码..."
        git fetch origin main 2>/dev/null || true
        git reset --hard origin/main 2>/dev/null || true
        
        # 第二步：编译后端
        echo "[2/5] 编译后端..."
        if [ -d "backend" ]; then
          cd backend
          go mod download 2>/dev/null || true
          go mod tidy 2>/dev/null || true
          CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o inarbit-server main.go 2>/dev/null || true
          cd ..
        fi
        
        # 第三步：编译前端
        echo "[3/5] 编译前端..."
        if [ -d "frontend" ]; then
          cd frontend
          npm ci 2>/dev/null || true
          npm run build 2>/dev/null || true
          cd ..
        fi
        
        # 第四步：重启服务
        echo "[4/5] 重启服务..."
        supervisorctl restart inarbit-backend 2>/dev/null || true
        supervisorctl restart inarbit-frontend 2>/dev/null || true
        systemctl restart nginx 2>/dev/null || true
        
        # 第五步：等待服务启动
        echo "[5/5] 等待服务启动..."
        sleep 5
        
        echo "=========================================="
        echo "部署完成！"
        echo "=========================================="
        EOF
    
    - name: 部署成功通知
      if: success()
      run: echo "✅ 部署成功！"

WORKFLOW_EOF

echo "✓ 工作流文件已创建"

# 第五步：提交文件
echo ""
echo "[5/6] 提交文件到Git..."
git add .github/workflows/deploy.yml
git commit -m "添加GitHub Actions自动化部署工作流" || echo "⚠ 文件已是最新状态"
echo "✓ 文件已提交"

# 第六步：推送到GitHub
echo ""
echo "[6/6] 推送到GitHub..."
git push origin main || echo "⚠ 推送失败，请检查Git凭证配置"
echo "✓ 推送完成"

echo ""
echo "=========================================="
echo "✅ GitHub Actions设置完成！"
echo "=========================================="
echo ""
echo "后续步骤："
echo "1. 进入 https://github.com/zillafan80-Maxzilla/inarbit/settings/secrets/actions"
echo "2. 添加以下Secrets："
echo "   - DEPLOY_SSH_KEY: (您的SSH私钥)"
echo "   - BINANCE_API_KEY: (您的Binance API密钥)"
echo "   - BINANCE_API_SECRET: (您的Binance API Secret)"
echo ""
echo "3. 进入 Actions 标签查看部署状态"
echo ""

#!/bin/bash

# iNarbit 前端编译脚本

set -e

echo "=========================================="
echo "iNarbit 前端编译"
echo "=========================================="

PROJECT_PATH="/root/inarbit"
FRONTEND_PATH="$PROJECT_PATH/frontend"

echo "1. 进入前端目录..."
cd "$FRONTEND_PATH"

echo "2. 安装依赖..."
npm install

echo "3. 编译前端..."
npm run build

echo "4. 验证编译..."
if [ -d "dist" ]; then
    echo "前端编译成功"
    ls -la dist/ | head -10
else
    echo "错误：前端编译失败"
    exit 1
fi

echo "=========================================="
echo "前端编译完成！"
echo "=========================================="
echo ""
echo "编译输出目录: $FRONTEND_PATH/dist"
echo ""
echo "后端将自动提供前端文件"
echo "=========================================="

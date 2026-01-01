#!/bin/bash

# iNarbit 后端编译脚本

set -e

echo "=========================================="
echo "iNarbit 后端编译"
echo "=========================================="

PROJECT_PATH="/root/inarbit"
BACKEND_PATH="$PROJECT_PATH/backend"

echo "1. 进入后端目录..."
cd "$BACKEND_PATH"

echo "2. 初始化Go模块..."
go mod init github.com/zillafan80-Maxzilla/inarbit 2>/dev/null || true

echo "3. 下载依赖..."
go mod download
go mod tidy

echo "4. 编译后端..."
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o inarbit-server main.go

echo "5. 设置执行权限..."
chmod +x inarbit-server

echo "6. 验证编译..."
./inarbit-server --version 2>/dev/null || echo "后端编译完成"

echo "=========================================="
echo "后端编译完成！"
echo "=========================================="
echo ""
echo "编译输出: $BACKEND_PATH/inarbit-server"
echo ""
echo "运行后端服务器："
echo "  cd $BACKEND_PATH"
echo "  ./inarbit-server"
echo "=========================================="

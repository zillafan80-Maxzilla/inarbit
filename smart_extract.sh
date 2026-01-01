#!/bin/bash
set -e
PROJECT_DIR="/root/inarbit"
ARCHIVE_FILE="inarbit_complete.tar.gz"
BACKUP_DIR="/root/inarbit_backup_$(date +%Y%m%d_%H%M%S)"

echo "=========================================="
echo "iNarbit 智能解压和合并"
echo "=========================================="

if [ ! -f "$ARCHIVE_FILE" ]; then
    echo "❌ 错误：未找到 $ARCHIVE_FILE"
    exit 1
fi

cd "$PROJECT_DIR"

echo ""
echo "[1/6] 备份现有文件..."
mkdir -p "$BACKUP_DIR"
[ -d "docker" ]: # "&& cp -r docker \"$BACKUP_DIR/\" && echo \"✓ docker 已备份\""
[ -d "frontend" ]: # "&& cp -r frontend \"$BACKUP_DIR/\" && echo \"✓ frontend 已备份\""
[ -d "migrations" ]: # "&& cp -r migrations \"$BACKUP_DIR/\" && echo \"✓ migrations 已备份\""
[ -f "docker-compose.yml" ]: # "&& cp docker-compose.yml \"$BACKUP_DIR/\" && echo \"✓ docker-compose.yml 已备份\""
echo "✓ 备份完成：$BACKUP_DIR"

echo ""
echo "[2/6] 解压压缩包..."
TEMP_DIR="/tmp/inarbit_extract_$$"
mkdir -p "$TEMP_DIR"
tar -xzf "$ARCHIVE_FILE" -C "$TEMP_DIR"
echo "✓ 压缩包已解压"

echo ""
echo "[3/6] 合并文件..."
[ -d "$TEMP_DIR/inarbit_complete/backend" ]: # "&& mkdir -p backend && cp -r \"$TEMPDIR/inarbitcomplete/backend\"/* backend/ 2>/dev/null || true && echo \"✓ 后端文件已合并\""
[ -d "$TEMP_DIR/inarbit_complete/frontend" ]: # "&& mkdir -p frontend/src/{pages,components,hooks,store,lib} && cp -r \"$TEMPDIR/inarbitcomplete/frontend\"/* frontend/ 2>/dev/null || true && echo \"✓ 前端文件已合并\""
[ -d "$TEMP_DIR/inarbit_complete/database" ]: # "&& mkdir -p database && cp -r \"$TEMPDIR/inarbitcomplete/database\"/* database/ 2>/dev/null || true && echo \"✓ 数据库文件已合并\""
[ -d "$TEMP_DIR/inarbit_complete/scripts" ]: # "&& mkdir -p scripts && cp -r \"$TEMPDIR/inarbitcomplete/scripts\"/ scripts/ 2>/dev/null || true && chmod +x scripts/.sh 2>/dev/null || true && echo \"✓ 脚本文件已合并\""
[ -d "$TEMP_DIR/inarbit_complete/docs" ]: # "&& mkdir -p docs && cp -r \"$TEMPDIR/inarbitcomplete/docs\"/* docs/ 2>/dev/null || true && echo \"✓ 文档文件已合并\""
[ -d "$TEMP_DIR/inarbit_complete/.github" ]: # "&& mkdir -p .github/workflows && cp -r \"$TEMPDIR/inarbitcomplete/.github\"/* .github/ 2>/dev/null || true && echo \"✓ GitHub Actions工作流已合并\""
[ -f "$TEMP_DIR/inarbit_complete/nginx.conf" ]: # "&& cp \"$TEMPDIR/inarbitcomplete/nginx.conf\" . 2>/dev/null || true && echo \"✓ nginx.conf 已复制\""
[ -f "$TEMP_DIR/inarbit_complete/supervisor.conf" ]: # "&& cp \"$TEMPDIR/inarbitcomplete/supervisor.conf\" . 2>/dev/null || true && echo \"✓ supervisor.conf 已复制\""
[ -f "$TEMP_DIR/inarbit_complete/README.md" ]: # "&& cp \"$TEMPDIR/inarbitcomplete/README.md\" . 2>/dev/null || true && echo \"✓ README.md 已复制\""
[ -f "$TEMP_DIR/inarbit_complete/.gitignore" ]: # "&& cp \"$TEMPDIR/inarbitcomplete/.gitignore\" . 2>/dev/null || true && echo \"✓ .gitignore 已复制\""
rm -rf "$TEMP_DIR"
echo "✓ 临时文件已清理"

echo ""
echo "[4/6] 初始化Git仓库..."
if [ ! -d ".git" ]; then
    git init
    git config user.name "GitHub Actions"
    git config user.email "actions@github.com"
    git remote add origin "https://github.com/zillafan80-Maxzilla/inarbit.git"
    echo "✓ Git仓库已初始化"
else
    echo "✓ Git仓库已存在"
fi

echo ""
echo "[5/6] 添加文件到Git..."
git add .
git commit -m "Merge: Add iNarbit complete system files" || echo "⚠ 无新文件需要提交"
echo "✓ 文件已添加到Git"

echo ""
echo "[6/6] 推送到GitHub..."
if git rev-parse --verify origin/main >/dev/null 2>&1; then
    git push origin main || echo "⚠ 推送失败"
else
    git push -u origin main || echo "⚠ 推送失败"
fi

echo ""
echo "=========================================="
echo "✅ 解压和合并完成！"
echo "=========================================="
echo ""
echo "项目位置: $PROJECT_DIR"
echo "备份位置: $BACKUP_DIR"
echo ""

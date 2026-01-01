#!/bin/bash

################################################################################
# iNarbit Go项目编译脚本
# 功能：编译Go后端和前端，生成可执行文件
# 使用：./build.sh [backend|frontend|all]
################################################################################

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 项目配置
PROJECT_ROOT="/root/inarbit"
BACKEND_DIR="${PROJECT_ROOT}/backend"
FRONTEND_DIR="${PROJECT_ROOT}/frontend"
DIST_DIR="${PROJECT_ROOT}/dist"
BUILD_TIME=$(date -u '+%Y-%m-%d %H:%M:%S')
GIT_COMMIT=$(cd ${PROJECT_ROOT} && git rev-parse --short HEAD 2>/dev/null || echo "unknown")
VERSION="1.0.0"

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查环境
check_environment() {
    log_info "检查构建环境..."

    # 检查Go
    if ! command -v go &> /dev/null; then
        log_error "Go未安装，请先安装Go"
        exit 1
    fi
    GO_VERSION=$(go version | awk '{print $3}')
    log_success "Go版本: $GO_VERSION"

    # 检查Node.js（前端需要）
    if ! command -v node &> /dev/null; then
        log_warn "Node.js未安装，跳过前端编译"
        SKIP_FRONTEND=true
    else
        NODE_VERSION=$(node --version)
        log_success "Node.js版本: $NODE_VERSION"
    fi

    # 检查npm
    if ! command -v npm &> /dev/null; then
        log_warn "npm未安装，跳过前端编译"
        SKIP_FRONTEND=true
    else
        NPM_VERSION=$(npm --version)
        log_success "npm版本: $NPM_VERSION"
    fi
}

# 创建输出目录
create_dist_dir() {
    log_info "创建输出目录..."
    mkdir -p ${DIST_DIR}/backend
    mkdir -p ${DIST_DIR}/frontend
    mkdir -p ${DIST_DIR}/config
    mkdir -p ${DIST_DIR}/scripts
    log_success "输出目录已创建"
}

# 编译后端
build_backend() {
    log_info "开始编译后端..."

    cd ${BACKEND_DIR}

    # 检查go.mod
    if [ ! -f "go.mod" ]; then
        log_error "go.mod不存在，请先运行 go mod init"
        exit 1
    fi

    # 下载依赖
    log_info "下载依赖..."
    go mod download
    go mod tidy

    # 编译
    log_info "编译Go代码..."
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
        -ldflags="-X main.Version=${VERSION} -X 'main.BuildTime=${BUILD_TIME}' -X main.GitCommit=${GIT_COMMIT}" \
        -o ${DIST_DIR}/backend/inarbit-backend \
        .

    # 检查编译结果
    if [ -f "${DIST_DIR}/backend/inarbit-backend" ]; then
        log_success "后端编译成功: ${DIST_DIR}/backend/inarbit-backend"
        
        # 显示文件大小
        SIZE=$(du -h ${DIST_DIR}/backend/inarbit-backend | cut -f1)
        log_info "文件大小: $SIZE"
    else
        log_error "后端编译失败"
        exit 1
    fi
}

# 编译前端
build_frontend() {
    if [ "$SKIP_FRONTEND" = true ]; then
        log_warn "跳过前端编译"
        return
    fi

    log_info "开始编译前端..."

    cd ${FRONTEND_DIR}

    # 检查package.json
    if [ ! -f "package.json" ]; then
        log_error "package.json不存在"
        exit 1
    fi

    # 安装依赖
    log_info "安装npm依赖..."
    npm ci --prefer-offline --no-audit

    # 编译
    log_info "编译前端代码..."
    npm run build

    # 检查编译结果
    if [ -d "dist" ]; then
        log_success "前端编译成功"
        
        # 复制到输出目录
        cp -r dist/* ${DIST_DIR}/frontend/
        log_success "前端文件已复制到输出目录"
    else
        log_error "前端编译失败"
        exit 1
    fi
}

# 复制配置文件
copy_config() {
    log_info "复制配置文件..."

    # 复制示例配置
    if [ -f "${PROJECT_ROOT}/.env.example" ]; then
        cp ${PROJECT_ROOT}/.env.example ${DIST_DIR}/config/.env.example
    fi

    # 复制Supervisor配置
    if [ -f "${PROJECT_ROOT}/supervisor/inarbit.conf" ]; then
        cp ${PROJECT_ROOT}/supervisor/inarbit.conf ${DIST_DIR}/config/
    fi

    # 复制Nginx配置
    if [ -f "${PROJECT_ROOT}/nginx/inarbit.conf" ]; then
        cp ${PROJECT_ROOT}/nginx/inarbit.conf ${DIST_DIR}/config/
    fi

    log_success "配置文件已复制"
}

# 复制脚本
copy_scripts() {
    log_info "复制脚本文件..."

    cp ${PROJECT_ROOT}/scripts/*.sh ${DIST_DIR}/scripts/ 2>/dev/null || true

    log_success "脚本文件已复制"
}

# 生成构建信息
generate_build_info() {
    log_info "生成构建信息..."

    cat > ${DIST_DIR}/BUILD_INFO.txt << EOF
===========================================
iNarbit 构建信息
===========================================
版本: ${VERSION}
构建时间: ${BUILD_TIME}
Git提交: ${GIT_COMMIT}
Go版本: ${GO_VERSION}
操作系统: Linux
架构: amd64
===========================================

后端可执行文件: backend/inarbit-backend
前端文件目录: frontend/

配置文件:
- config/.env.example
- config/inarbit.conf (Supervisor)
- config/inarbit.conf (Nginx)

脚本文件:
- scripts/deploy.sh
- scripts/start.sh
- scripts/stop.sh
- scripts/restart.sh
- scripts/health-check.sh

使用说明:
1. 配置环境变量: cp config/.env.example .env
2. 启动服务: ./scripts/start.sh
3. 查看日志: tail -f logs/inarbit.log

===========================================
EOF

    log_success "构建信息已生成"
}

# 验证编译结果
verify_build() {
    log_info "验证编译结果..."

    # 检查后端可执行文件
    if [ ! -f "${DIST_DIR}/backend/inarbit-backend" ]; then
        log_error "后端可执行文件不存在"
        exit 1
    fi

    # 检查可执行权限
    if [ ! -x "${DIST_DIR}/backend/inarbit-backend" ]; then
        chmod +x ${DIST_DIR}/backend/inarbit-backend
    fi

    # 测试后端
    log_info "测试后端可执行文件..."
    ${DIST_DIR}/backend/inarbit-backend -version 2>/dev/null || true

    log_success "编译结果验证完成"
}

# 生成打包文件
create_archive() {
    log_info "创建发布包..."

    cd ${PROJECT_ROOT}
    
    ARCHIVE_NAME="inarbit-${VERSION}-$(date +%Y%m%d-%H%M%S).tar.gz"
    tar -czf ${ARCHIVE_NAME} dist/ BUILD_INFO.txt

    if [ -f "${ARCHIVE_NAME}" ]; then
        SIZE=$(du -h ${ARCHIVE_NAME} | cut -f1)
        log_success "发布包已创建: ${ARCHIVE_NAME} ($SIZE)"
    else
        log_error "创建发布包失败"
        exit 1
    fi
}

# 显示总结
show_summary() {
    echo ""
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}编译完成${NC}"
    echo -e "${GREEN}========================================${NC}"
    echo ""
    echo "输出目录: ${DIST_DIR}"
    echo ""
    echo "文件列表:"
    ls -lh ${DIST_DIR}/backend/inarbit-backend
    if [ -d "${DIST_DIR}/frontend" ] && [ "$(ls -A ${DIST_DIR}/frontend)" ]; then
        echo "前端文件: $(ls ${DIST_DIR}/frontend | wc -l) 个"
    fi
    echo ""
    echo "下一步:"
    echo "1. 配置环境变量: cp ${DIST_DIR}/config/.env.example .env"
    echo "2. 启动服务: ./scripts/start.sh"
    echo "3. 查看日志: tail -f logs/inarbit.log"
    echo ""
}

# 主函数
main() {
    local BUILD_TARGET=${1:-all}

    log_info "=========================================="
    log_info "iNarbit Go项目编译脚本"
    log_info "=========================================="
    log_info "项目根目录: ${PROJECT_ROOT}"
    log_info "输出目录: ${DIST_DIR}"
    log_info "编译目标: ${BUILD_TARGET}"
    log_info ""

    # 检查环境
    check_environment

    # 创建输出目录
    create_dist_dir

    # 编译
    case ${BUILD_TARGET} in
        backend)
            build_backend
            ;;
        frontend)
            build_frontend
            ;;
        all)
            build_backend
            build_frontend
            ;;
        *)
            log_error "未知的编译目标: ${BUILD_TARGET}"
            echo "用法: $0 [backend|frontend|all]"
            exit 1
            ;;
    esac

    # 复制配置和脚本
    copy_config
    copy_scripts

    # 生成构建信息
    generate_build_info

    # 验证编译结果
    verify_build

    # 创建发布包
    create_archive

    # 显示总结
    show_summary
}

# 执行主函数
main "$@"

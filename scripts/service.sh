#!/bin/bash

################################################################################
# iNarbit 服务管理脚本集合
# 包含：启动、停止、重启、状态检查、日志查看等功能
################################################################################

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# 项目配置
PROJECT_ROOT="/root/inarbit"
LOG_DIR="${PROJECT_ROOT}/logs"
APP_NAME="inarbit-backend"
APP_PORT=8080

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

# ===== 启动脚本 (start.sh) =====

start_service() {
    log_info "启动iNarbit服务..."

    # 检查Supervisor是否运行
    if ! systemctl is-active --quiet supervisor; then
        log_warn "Supervisor未运行，正在启动..."
        systemctl start supervisor
        sleep 2
    fi

    # 启动后端服务
    supervisorctl start ${APP_NAME}

    # 等待服务启动
    sleep 3

    # 检查服务状态
    if supervisorctl status ${APP_NAME} | grep -q "RUNNING"; then
        log_success "后端服务已启动"
    else
        log_error "后端服务启动失败"
        supervisorctl tail ${APP_NAME}
        return 1
    fi

    # 启动Nginx
    systemctl start nginx
    log_success "Nginx已启动"

    # 等待服务完全启动
    sleep 2

    # 健康检查
    if curl -s http://localhost:${APP_PORT}/api/health > /dev/null 2>&1; then
        log_success "服务健康检查通过"
        echo ""
        echo "访问地址:"
        echo "  http://inarbit.work"
        echo "  https://inarbit.work"
        return 0
    else
        log_warn "服务健康检查失败，请检查日志"
        return 1
    fi
}

# ===== 停止脚本 (stop.sh) =====

stop_service() {
    log_info "停止iNarbit服务..."

    # 停止后端服务
    supervisorctl stop ${APP_NAME}
    log_success "后端服务已停止"

    # 停止Nginx
    systemctl stop nginx
    log_success "Nginx已停止"
}

# ===== 重启脚本 (restart.sh) =====

restart_service() {
    log_info "重启iNarbit服务..."

    stop_service
    sleep 2
    start_service
}

# ===== 状态检查脚本 (status.sh) =====

check_status() {
    log_info "检查iNarbit服务状态..."
    echo ""

    # 检查后端服务
    echo -n "后端服务: "
    if supervisorctl status ${APP_NAME} | grep -q "RUNNING"; then
        echo -e "${GREEN}运行中${NC}"
    else
        echo -e "${RED}已停止${NC}"
    fi

    # 检查Nginx
    echo -n "Nginx: "
    if systemctl is-active --quiet nginx; then
        echo -e "${GREEN}运行中${NC}"
    else
        echo -e "${RED}已停止${NC}"
    fi

    # 检查PostgreSQL
    echo -n "PostgreSQL: "
    if systemctl is-active --quiet postgresql; then
        echo -e "${GREEN}运行中${NC}"
    else
        echo -e "${RED}已停止${NC}"
    fi

    # 检查Supervisor
    echo -n "Supervisor: "
    if systemctl is-active --quiet supervisor; then
        echo -e "${GREEN}运行中${NC}"
    else
        echo -e "${RED}已停止${NC}"
    fi

    echo ""

    # 检查端口
    echo "端口占用情况:"
    echo -n "  8080 (后端): "
    if netstat -tuln 2>/dev/null | grep -q ":8080 "; then
        echo -e "${GREEN}已占用${NC}"
    else
        echo -e "${RED}未占用${NC}"
    fi

    echo -n "  80 (HTTP): "
    if netstat -tuln 2>/dev/null | grep -q ":80 "; then
        echo -e "${GREEN}已占用${NC}"
    else
        echo -e "${RED}未占用${NC}"
    fi

    echo -n "  443 (HTTPS): "
    if netstat -tuln 2>/dev/null | grep -q ":443 "; then
        echo -e "${GREEN}已占用${NC}"
    else
        echo -e "${RED}未占用${NC}"
    fi

    echo ""

    # 检查磁盘空间
    echo "磁盘空间:"
    df -h ${PROJECT_ROOT} | tail -1 | awk '{printf "  使用率: %s/%s (%s)\n", $3, $2, $5}'

    echo ""

    # 检查进程
    echo "进程信息:"
    ps aux | grep inarbit-backend | grep -v grep || echo "  后端进程未运行"
}

# ===== 日志查看脚本 (logs.sh) =====

view_logs() {
    local log_type=${1:-backend}

    case ${log_type} in
        backend)
            log_info "查看后端日志..."
            tail -f ${LOG_DIR}/inarbit.log
            ;;
        nginx-access)
            log_info "查看Nginx访问日志..."
            tail -f ${LOG_DIR}/nginx_access.log
            ;;
        nginx-error)
            log_info "查看Nginx错误日志..."
            tail -f ${LOG_DIR}/nginx_error.log
            ;;
        supervisor)
            log_info "查看Supervisor日志..."
            tail -f /var/log/supervisor/supervisord.log
            ;;
        all)
            log_info "查看所有日志..."
            echo "后端日志:"
            tail -20 ${LOG_DIR}/inarbit.log
            echo ""
            echo "Nginx访问日志:"
            tail -20 ${LOG_DIR}/nginx_access.log
            echo ""
            echo "Nginx错误日志:"
            tail -20 ${LOG_DIR}/nginx_error.log
            ;;
        *)
            log_error "未知的日志类型: ${log_type}"
            echo "用法: logs.sh [backend|nginx-access|nginx-error|supervisor|all]"
            return 1
            ;;
    esac
}

# ===== 健康检查脚本 (health-check.sh) =====

health_check() {
    log_info "执行服务健康检查..."
    echo ""

    local failed=0

    # 检查后端API
    echo -n "检查后端API (http://localhost:${APP_PORT}/api/health): "
    if curl -s http://localhost:${APP_PORT}/api/health > /dev/null 2>&1; then
        echo -e "${GREEN}✓${NC}"
    else
        echo -e "${RED}✗${NC}"
        ((failed++))
    fi

    # 检查前端
    echo -n "检查前端 (http://localhost/): "
    if curl -s http://localhost/ > /dev/null 2>&1; then
        echo -e "${GREEN}✓${NC}"
    else
        echo -e "${RED}✗${NC}"
        ((failed++))
    fi

    # 检查WebSocket
    echo -n "检查WebSocket (ws://localhost:${APP_PORT}/ws): "
    # 简单的WebSocket检查（实际应该使用WebSocket客户端）
    if curl -s http://localhost:${APP_PORT}/ws 2>&1 | grep -q "400\|upgrade"; then
        echo -e "${GREEN}✓${NC}"
    else
        echo -e "${YELLOW}?${NC}"
    fi

    # 检查数据库连接
    echo -n "检查数据库连接: "
    if sudo -u postgres psql -d inarbit -c "SELECT 1" > /dev/null 2>&1; then
        echo -e "${GREEN}✓${NC}"
    else
        echo -e "${RED}✗${NC}"
        ((failed++))
    fi

    echo ""

    if [ ${failed} -eq 0 ]; then
        log_success "所有检查通过"
        return 0
    else
        log_error "有${failed}项检查失败"
        return 1
    fi
}

# ===== 性能监控脚本 (monitor.sh) =====

monitor_performance() {
    log_info "监控服务性能..."
    echo ""

    # 监听间隔
    local interval=${1:-5}

    while true; do
        clear
        echo "iNarbit 性能监控 (更新间隔: ${interval}秒)"
        echo "=========================================="
        echo ""

        # CPU和内存使用
        echo "CPU和内存使用:"
        ps aux | grep inarbit-backend | grep -v grep | awk '{printf "  CPU: %.1f%%, 内存: %.1f%%\n", $3, $4}'

        echo ""

        # 网络连接数
        echo "网络连接数:"
        echo "  HTTP连接: $(netstat -an 2>/dev/null | grep -c ':80 ')"
        echo "  HTTPS连接: $(netstat -an 2>/dev/null | grep -c ':443 ')"
        echo "  后端连接: $(netstat -an 2>/dev/null | grep -c ':8080 ')"

        echo ""

        # 磁盘I/O
        echo "磁盘使用:"
        df -h ${PROJECT_ROOT} | tail -1 | awk '{printf "  使用率: %s/%s (%s)\n", $3, $2, $5}'

        echo ""

        # 日志大小
        echo "日志文件大小:"
        ls -lh ${LOG_DIR}/*.log 2>/dev/null | awk '{printf "  %s: %s\n", $9, $5}'

        echo ""
        echo "按 Ctrl+C 退出监控"
        echo ""

        sleep ${interval}
    done
}

# ===== 数据库备份脚本 (backup.sh) =====

backup_database() {
    log_info "备份数据库..."

    local backup_dir="${PROJECT_ROOT}/backups"
    local backup_file="${backup_dir}/inarbit_$(date +%Y%m%d_%H%M%S).sql"

    mkdir -p ${backup_dir}

    # 执行备份
    sudo -u postgres pg_dump inarbit > ${backup_file}

    if [ -f "${backup_file}" ]; then
        # 压缩备份文件
        gzip ${backup_file}
        local size=$(du -h ${backup_file}.gz | cut -f1)
        log_success "数据库备份完成: ${backup_file}.gz ($size)"

        # 保留最近7天的备份
        find ${backup_dir} -name "inarbit_*.sql.gz" -mtime +7 -delete
        log_info "已清理7天前的备份文件"
    else
        log_error "数据库备份失败"
        return 1
    fi
}

# ===== 清理日志脚本 (cleanup.sh) =====

cleanup_logs() {
    log_info "清理日志文件..."

    # 清理7天前的日志
    find ${LOG_DIR} -name "*.log" -mtime +7 -delete
    log_success "已清理7天前的日志文件"

    # 显示当前日志大小
    local total_size=$(du -sh ${LOG_DIR} | cut -f1)
    log_info "日志目录总大小: $total_size"
}

# ===== 主函数 =====

show_help() {
    cat << EOF
iNarbit 服务管理脚本

用法: $0 <command> [options]

命令:
  start              启动服务
  stop               停止服务
  restart            重启服务
  status             查看服务状态
  logs [type]        查看日志 (backend|nginx-access|nginx-error|supervisor|all)
  health-check       执行健康检查
  monitor [interval] 监控性能 (默认间隔5秒)
  backup             备份数据库
  cleanup            清理日志文件
  help               显示帮助信息

示例:
  $0 start           # 启动服务
  $0 stop            # 停止服务
  $0 restart         # 重启服务
  $0 status          # 查看状态
  $0 logs backend    # 查看后端日志
  $0 logs all        # 查看所有日志
  $0 health-check    # 执行健康检查
  $0 monitor 10      # 每10秒监控一次性能
  $0 backup          # 备份数据库
  $0 cleanup         # 清理日志

EOF
}

# 主程序
main() {
    local command=${1:-help}

    case ${command} in
        start)
            start_service
            ;;
        stop)
            stop_service
            ;;
        restart)
            restart_service
            ;;
        status)
            check_status
            ;;
        logs)
            view_logs "$2"
            ;;
        health-check)
            health_check
            ;;
        monitor)
            monitor_performance "$2"
            ;;
        backup)
            backup_database
            ;;
        cleanup)
            cleanup_logs
            ;;
        help)
            show_help
            ;;
        *)
            log_error "未知的命令: ${command}"
            show_help
            exit 1
            ;;
    esac
}

# 执行主程序
main "$@"

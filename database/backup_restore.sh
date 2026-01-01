#!/bin/bash

################################################################################
# iNarbit PostgreSQL 数据库备份和恢复脚本
# 功能：备份、恢复、导出、导入数据库
# 使用：./db_backup_restore.sh [backup|restore|export|import|cleanup]
################################################################################

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# 数据库配置
DB_HOST="localhost"
DB_PORT="5432"
DB_USER="inarbit"
DB_PASSWORD="inarbit_password"
DB_NAME="inarbit"
BACKUP_DIR="/root/inarbit/backups"
EXPORT_DIR="/root/inarbit/exports"

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
    # 检查PostgreSQL客户端
    if ! command -v psql &> /dev/null; then
        log_error "psql未安装"
        exit 1
    fi

    # 检查pg_dump
    if ! command -v pg_dump &> /dev/null; then
        log_error "pg_dump未安装"
        exit 1
    fi

    # 检查pg_restore
    if ! command -v pg_restore &> /dev/null; then
        log_error "pg_restore未安装"
        exit 1
    fi

    # 创建备份目录
    mkdir -p ${BACKUP_DIR}
    mkdir -p ${EXPORT_DIR}

    log_success "环境检查通过"
}

# 备份数据库（SQL格式）
backup_sql() {
    log_info "备份数据库（SQL格式）..."

    local backup_file="${BACKUP_DIR}/inarbit_$(date +%Y%m%d_%H%M%S).sql"

    PGPASSWORD=${DB_PASSWORD} pg_dump \
        -h ${DB_HOST} \
        -p ${DB_PORT} \
        -U ${DB_USER} \
        -d ${DB_NAME} \
        --no-password \
        > ${backup_file}

    if [ -f "${backup_file}" ]; then
        # 压缩备份文件
        gzip ${backup_file}
        local size=$(du -h ${backup_file}.gz | cut -f1)
        log_success "数据库备份完成: ${backup_file}.gz ($size)"
        return 0
    else
        log_error "数据库备份失败"
        return 1
    fi
}

# 备份数据库（自定义格式）
backup_custom() {
    log_info "备份数据库（自定义格式）..."

    local backup_file="${BACKUP_DIR}/inarbit_$(date +%Y%m%d_%H%M%S).dump"

    PGPASSWORD=${DB_PASSWORD} pg_dump \
        -h ${DB_HOST} \
        -p ${DB_PORT} \
        -U ${DB_USER} \
        -d ${DB_NAME} \
        --no-password \
        -Fc \
        -f ${backup_file}

    if [ -f "${backup_file}" ]; then
        local size=$(du -h ${backup_file} | cut -f1)
        log_success "数据库备份完成: ${backup_file} ($size)"
        return 0
    else
        log_error "数据库备份失败"
        return 1
    fi
}

# 恢复数据库（SQL格式）
restore_sql() {
    local backup_file=$1

    if [ -z "${backup_file}" ]; then
        log_error "请指定备份文件"
        return 1
    fi

    if [ ! -f "${backup_file}" ]; then
        log_error "备份文件不存在: ${backup_file}"
        return 1
    fi

    log_warn "即将恢复数据库，现有数据将被覆盖"
    read -p "是否继续？(y/N): " confirm
    if [ "${confirm}" != "y" ]; then
        log_info "已取消恢复"
        return 0
    fi

    log_info "恢复数据库..."

    # 如果是gzip压缩的，先解压
    if [[ ${backup_file} == *.gz ]]; then
        local temp_file="${backup_file%.gz}"
        gunzip -c ${backup_file} > ${temp_file}
        backup_file=${temp_file}
    fi

    PGPASSWORD=${DB_PASSWORD} psql \
        -h ${DB_HOST} \
        -p ${DB_PORT} \
        -U ${DB_USER} \
        -d ${DB_NAME} \
        --no-password \
        < ${backup_file}

    if [ $? -eq 0 ]; then
        log_success "数据库恢复完成"
        return 0
    else
        log_error "数据库恢复失败"
        return 1
    fi
}

# 恢复数据库（自定义格式）
restore_custom() {
    local backup_file=$1

    if [ -z "${backup_file}" ]; then
        log_error "请指定备份文件"
        return 1
    fi

    if [ ! -f "${backup_file}" ]; then
        log_error "备份文件不存在: ${backup_file}"
        return 1
    fi

    log_warn "即将恢复数据库，现有数据将被覆盖"
    read -p "是否继续？(y/N): " confirm
    if [ "${confirm}" != "y" ]; then
        log_info "已取消恢复"
        return 0
    fi

    log_info "恢复数据库..."

    PGPASSWORD=${DB_PASSWORD} pg_restore \
        -h ${DB_HOST} \
        -p ${DB_PORT} \
        -U ${DB_USER} \
        -d ${DB_NAME} \
        --no-password \
        -c \
        ${backup_file}

    if [ $? -eq 0 ]; then
        log_success "数据库恢复完成"
        return 0
    else
        log_error "数据库恢复失败"
        return 1
    fi
}

# 导出表数据为CSV
export_csv() {
    local table_name=$1

    if [ -z "${table_name}" ]; then
        log_error "请指定表名"
        return 1
    fi

    log_info "导出表 ${table_name} 为CSV..."

    local export_file="${EXPORT_DIR}/${table_name}_$(date +%Y%m%d_%H%M%S).csv"

    PGPASSWORD=${DB_PASSWORD} psql \
        -h ${DB_HOST} \
        -p ${DB_PORT} \
        -U ${DB_USER} \
        -d ${DB_NAME} \
        --no-password \
        -c "COPY ${table_name} TO STDOUT WITH CSV HEADER" \
        > ${export_file}

    if [ -f "${export_file}" ]; then
        local size=$(du -h ${export_file} | cut -f1)
        log_success "表数据导出完成: ${export_file} ($size)"
        return 0
    else
        log_error "表数据导出失败"
        return 1
    fi
}

# 导入CSV数据到表
import_csv() {
    local table_name=$1
    local csv_file=$2

    if [ -z "${table_name}" ] || [ -z "${csv_file}" ]; then
        log_error "请指定表名和CSV文件"
        return 1
    fi

    if [ ! -f "${csv_file}" ]; then
        log_error "CSV文件不存在: ${csv_file}"
        return 1
    fi

    log_warn "即将导入数据到表 ${table_name}"
    read -p "是否继续？(y/N): " confirm
    if [ "${confirm}" != "y" ]; then
        log_info "已取消导入"
        return 0
    fi

    log_info "导入CSV数据..."

    PGPASSWORD=${DB_PASSWORD} psql \
        -h ${DB_HOST} \
        -p ${DB_PORT} \
        -U ${DB_USER} \
        -d ${DB_NAME} \
        --no-password \
        -c "COPY ${table_name} FROM STDIN WITH CSV HEADER" \
        < ${csv_file}

    if [ $? -eq 0 ]; then
        log_success "数据导入完成"
        return 0
    else
        log_error "数据导入失败"
        return 1
    fi
}

# 列出备份文件
list_backups() {
    log_info "备份文件列表:"
    echo ""
    
    if [ ! -d "${BACKUP_DIR}" ] || [ -z "$(ls -A ${BACKUP_DIR})" ]; then
        log_warn "没有备份文件"
        return
    fi

    ls -lh ${BACKUP_DIR} | tail -n +2 | awk '{printf "  %-40s %10s  %s %s %s\n", $9, $5, $6, $7, $8}'
}

# 清理旧备份
cleanup_backups() {
    local days=${1:-30}

    log_info "清理${days}天前的备份文件..."

    local count=$(find ${BACKUP_DIR} -name "*.sql.gz" -o -name "*.dump" | wc -l)

    find ${BACKUP_DIR} \( -name "*.sql.gz" -o -name "*.dump" \) -mtime +${days} -delete

    local new_count=$(find ${BACKUP_DIR} -name "*.sql.gz" -o -name "*.dump" | wc -l)
    local deleted=$((count - new_count))

    log_success "已清理${deleted}个备份文件"
}

# 验证数据库连接
verify_connection() {
    log_info "验证数据库连接..."

    PGPASSWORD=${DB_PASSWORD} psql \
        -h ${DB_HOST} \
        -p ${DB_PORT} \
        -U ${DB_USER} \
        -d ${DB_NAME} \
        --no-password \
        -c "SELECT 1" > /dev/null 2>&1

    if [ $? -eq 0 ]; then
        log_success "数据库连接成功"
        return 0
    else
        log_error "数据库连接失败"
        return 1
    fi
}

# 显示数据库统计信息
show_statistics() {
    log_info "数据库统计信息:"
    echo ""

    PGPASSWORD=${DB_PASSWORD} psql \
        -h ${DB_HOST} \
        -p ${DB_PORT} \
        -U ${DB_USER} \
        -d ${DB_NAME} \
        --no-password \
        << EOF
-- 表统计
SELECT 
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) as size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;

-- 数据库大小
SELECT pg_size_pretty(pg_database_size('${DB_NAME}')) as database_size;

-- 行数统计
SELECT 
    tablename,
    n_live_tup as row_count
FROM pg_stat_user_tables
ORDER BY n_live_tup DESC;
EOF
}

# 显示帮助信息
show_help() {
    cat << EOF
iNarbit PostgreSQL 数据库备份和恢复脚本

用法: $0 <command> [options]

命令:
  backup              备份数据库（SQL格式）
  backup-custom       备份数据库（自定义格式）
  restore <file>      恢复数据库（SQL格式）
  restore-custom <file> 恢复数据库（自定义格式）
  export <table>      导出表数据为CSV
  import <table> <file> 导入CSV数据到表
  list                列出备份文件
  cleanup [days]      清理旧备份（默认30天）
  verify              验证数据库连接
  stats               显示数据库统计信息
  help                显示帮助信息

示例:
  $0 backup                           # 备份数据库
  $0 backup-custom                    # 备份数据库（自定义格式）
  $0 restore /path/to/backup.sql.gz   # 恢复数据库
  $0 export users                     # 导出users表
  $0 import users /path/to/users.csv  # 导入users表
  $0 list                             # 列出备份文件
  $0 cleanup 30                       # 清理30天前的备份
  $0 verify                           # 验证数据库连接
  $0 stats                            # 显示统计信息

EOF
}

# 主函数
main() {
    local command=${1:-help}

    check_environment

    case ${command} in
        backup)
            backup_sql
            ;;
        backup-custom)
            backup_custom
            ;;
        restore)
            restore_sql "$2"
            ;;
        restore-custom)
            restore_custom "$2"
            ;;
        export)
            export_csv "$2"
            ;;
        import)
            import_csv "$2" "$3"
            ;;
        list)
            list_backups
            ;;
        cleanup)
            cleanup_backups "$2"
            ;;
        verify)
            verify_connection
            ;;
        stats)
            show_statistics
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

# 执行主函数
main "$@"

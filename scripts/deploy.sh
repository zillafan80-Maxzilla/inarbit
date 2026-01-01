#!/bin/bash

# iNarbit 一键部署脚本
# 用法: bash deploy.sh

set -e

echo "=========================================="
echo "iNarbit 完整部署脚本"
echo "=========================================="
echo ""

PROJECT_PATH="/root/inarbit"
SCRIPTS_PATH="$PROJECT_PATH/scripts"

# 检查是否为root用户
if [[ $EUID -ne 0 ]]; then
   echo "错误：此脚本必须以root用户身份运行"
   exit 1
fi

# 检查脚本是否存在
check_script() {
    if [ ! -f "$1" ]; then
        echo "错误：找不到脚本 $1"
        exit 1
    fi
}

# 执行脚本
run_script() {
    local script=$1
    local name=$2
    
    echo ""
    echo "=========================================="
    echo "执行: $name"
    echo "=========================================="
    
    check_script "$script"
    bash "$script"
    
    if [ $? -eq 0 ]; then
        echo "✓ $name 完成"
    else
        echo "✗ $name 失败"
        exit 1
    fi
}

# 显示菜单
show_menu() {
    echo ""
    echo "选择要执行的步骤："
    echo "1) 完整部署（执行所有步骤）"
    echo "2) 仅系统初始化"
    echo "3) 仅数据库初始化"
    echo "4) 仅编译后端"
    echo "5) 仅编译前端"
    echo "6) 仅配置SSL"
    echo "7) 启动服务"
    echo "0) 退出"
    echo ""
    read -p "请选择 [0-7]: " choice
}

# 完整部署
full_deploy() {
    run_script "$SCRIPTS_PATH/01_init_system.sh" "系统初始化"
    run_script "$SCRIPTS_PATH/02_install_deps.sh" "依赖安装"
    run_script "$SCRIPTS_PATH/03_setup_database.sh" "数据库初始化"
    run_script "$SCRIPTS_PATH/04_build_backend.sh" "后端编译"
    run_script "$SCRIPTS_PATH/05_build_frontend.sh" "前端编译"
    run_script "$SCRIPTS_PATH/06_setup_ssl.sh" "SSL证书配置"
    run_script "$SCRIPTS_PATH/07_start_services.sh" "启动服务"
}

# 主菜单循环
while true; do
    show_menu
    
    case $choice in
        1)
            full_deploy
            echo ""
            echo "=========================================="
            echo "完整部署完成！"
            echo "=========================================="
            echo ""
            echo "访问应用:"
            echo "  https://inarbit.work"
            echo "  https://www.inarbit.work"
            echo ""
            echo "默认用户:"
            echo "  用户名: admin"
            echo "  密码: password (请立即修改)"
            echo ""
            break
            ;;
        2)
            run_script "$SCRIPTS_PATH/01_init_system.sh" "系统初始化"
            ;;
        3)
            run_script "$SCRIPTS_PATH/03_setup_database.sh" "数据库初始化"
            ;;
        4)
            run_script "$SCRIPTS_PATH/04_build_backend.sh" "后端编译"
            ;;
        5)
            run_script "$SCRIPTS_PATH/05_build_frontend.sh" "前端编译"
            ;;
        6)
            run_script "$SCRIPTS_PATH/06_setup_ssl.sh" "SSL证书配置"
            ;;
        7)
            run_script "$SCRIPTS_PATH/07_start_services.sh" "启动服务"
            ;;
        0)
            echo "退出"
            exit 0
            ;;
        *)
            echo "无效选择，请重试"
            ;;
    esac
done

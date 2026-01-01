#!/bin/bash

# iNarbit 启动服务脚本

set -e

echo "=========================================="
echo "iNarbit 启动服务"
echo "=========================================="

PROJECT_PATH="/root/inarbit"
BACKEND_PATH="$PROJECT_PATH/backend"
NGINX_PATH="$PROJECT_PATH/nginx"

echo "1. 配置Nginx..."
cp "$NGINX_PATH/nginx.conf" /etc/nginx/sites-available/inarbit
ln -sf /etc/nginx/sites-available/inarbit /etc/nginx/sites-enabled/inarbit
rm -f /etc/nginx/sites-enabled/default

echo "2. 测试Nginx配置..."
nginx -t

echo "3. 启动Nginx..."
systemctl start nginx
systemctl enable nginx

echo "4. 创建Systemd服务文件..."
cat > /etc/systemd/system/inarbit.service << 'EOF'
[Unit]
Description=iNarbit Arbitrage Bot Service
After=network.target postgresql.service

[Service]
Type=simple
User=inarbit
WorkingDirectory=/root/inarbit/backend
ExecStart=/root/inarbit/backend/inarbit-server
Restart=always
RestartSec=10
StandardOutput=journal
StandardError=journal
Environment="PATH=/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"

[Install]
WantedBy=multi-user.target
EOF

echo "5. 重新加载Systemd..."
systemctl daemon-reload

echo "6. 启动iNarbit服务..."
systemctl start inarbit
systemctl enable inarbit

echo "7. 检查服务状态..."
sleep 2
systemctl status inarbit --no-pager

echo "8. 检查日志..."
journalctl -u inarbit -n 20 --no-pager

echo "=========================================="
echo "服务启动完成！"
echo "=========================================="
echo ""
echo "服务管理命令:"
echo "  启动:   systemctl start inarbit"
echo "  停止:   systemctl stop inarbit"
echo "  重启:   systemctl restart inarbit"
echo "  状态:   systemctl status inarbit"
echo "  日志:   journalctl -u inarbit -f"
echo ""
echo "Nginx管理命令:"
echo "  启动:   systemctl start nginx"
echo "  停止:   systemctl stop nginx"
echo "  重启:   systemctl restart nginx"
echo "  测试:   nginx -t"
echo ""
echo "访问应用:"
echo "  https://inarbit.work"
echo "  https://www.inarbit.work"
echo ""
echo "API文档:"
echo "  https://inarbit.work/api/docs"
echo ""
echo "WebSocket连接:"
echo "  wss://inarbit.work/ws"
echo "=========================================="

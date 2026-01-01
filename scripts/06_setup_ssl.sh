#!/bin/bash

# iNarbit SSL证书配置脚本

set -e

echo "=========================================="
echo "iNarbit SSL证书配置"
echo "=========================================="

PROJECT_PATH="/root/inarbit"
NGINX_PATH="$PROJECT_PATH/nginx"
DOMAINS="inarbit.work www.inarbit.work"
EMAIL="admin@inarbit.work"

echo "1. 检查Certbot..."
if ! command -v certbot &> /dev/null; then
    echo "安装Certbot..."
    apt-get install -y certbot python3-certbot-nginx
fi

echo "2. 停止Nginx..."
systemctl stop nginx || true

echo "3. 申请SSL证书..."
certbot certonly --standalone \
    --non-interactive \
    --agree-tos \
    --email "$EMAIL" \
    -d inarbit.work \
    -d www.inarbit.work \
    --cert-name inarbit-cert \
    --force-renewal

echo "4. 复制证书到项目目录..."
CERT_PATH="/etc/letsencrypt/live/inarbit-cert"
if [ -d "$CERT_PATH" ]; then
    cp "$CERT_PATH/fullchain.pem" "$NGINX_PATH/ssl/cert.pem"
    cp "$CERT_PATH/privkey.pem" "$NGINX_PATH/ssl/key.pem"
    chmod 644 "$NGINX_PATH/ssl/cert.pem"
    chmod 600 "$NGINX_PATH/ssl/key.pem"
    echo "证书复制成功"
else
    echo "警告：证书目录不存在，使用自签名证书"
    # 生成自签名证书（用于测试）
    openssl req -x509 -newkey rsa:4096 -keyout "$NGINX_PATH/ssl/key.pem" -out "$NGINX_PATH/ssl/cert.pem" -days 365 -nodes \
        -subj "/C=JP/ST=Tokyo/L=Tokyo/O=iNarbit/CN=inarbit.work"
fi

echo "5. 设置证书权限..."
chown -R root:root "$NGINX_PATH/ssl"
chmod 755 "$NGINX_PATH/ssl"

echo "6. 配置自动续期..."
# 添加到crontab
(crontab -l 2>/dev/null | grep -v "certbot renew" || true; echo "0 3 * * * certbot renew --quiet && systemctl reload nginx") | crontab -

echo "=========================================="
echo "SSL证书配置完成！"
echo "=========================================="
echo ""
echo "证书信息："
echo "  证书路径: $NGINX_PATH/ssl/cert.pem"
echo "  密钥路径: $NGINX_PATH/ssl/key.pem"
echo ""
echo "自动续期："
echo "  每天凌晨3点自动检查并续期证书"
echo ""
echo "验证证书："
echo "  openssl x509 -in $NGINX_PATH/ssl/cert.pem -text -noout"
echo "=========================================="

# iNarbit PostgreSQL 数据库完整指南

## 目录

1. [系统要求](#系统要求)
2. [安装和初始化](#安装和初始化)
3. [数据库结构](#数据库结构)
4. [备份和恢复](#备份和恢复)
5. [性能优化](#性能优化)
6. [故障排除](#故障排除)
7. [安全加固](#安全加固)
8. [常见问题](#常见问题)

---

## 系统要求

### 硬件要求

- **CPU**: 2核及以上
- **内存**: 4GB及以上
- **磁盘**: 50GB SSD（ESSD）
- **网络**: 1Mbps及以上

### 软件要求

- **操作系统**: Ubuntu 22.04 LTS
- **PostgreSQL**: 14.0 或更高版本
- **Python**: 3.8 或更高版本（用于管理脚本）

### 依赖包

```bash
# 安装PostgreSQL
sudo apt-get update
sudo apt-get install -y postgresql postgresql-contrib

# 安装PostgreSQL开发工具
sudo apt-get install -y postgresql-client postgresql-client-common

# 安装备份工具
sudo apt-get install -y pg-dump-restore
```

---

## 安装和初始化

### 1. 安装PostgreSQL

```bash
# 更新包列表
sudo apt-get update

# 安装PostgreSQL
sudo apt-get install -y postgresql postgresql-contrib

# 启动PostgreSQL服务
sudo systemctl start postgresql
sudo systemctl enable postgresql

# 验证安装
psql --version
```

### 2. 初始化数据库

```bash
# 切换到postgres用户
sudo -u postgres psql

# 在psql中执行以下命令
-- 创建应用用户
CREATE USER inarbit WITH PASSWORD 'inarbit_password';

-- 创建数据库
CREATE DATABASE inarbit OWNER inarbit;

-- 授予权限
GRANT ALL PRIVILEGES ON DATABASE inarbit TO inarbit;

-- 退出
\q
```

### 3. 导入SQL脚本

```bash
# 使用SQL脚本初始化数据库
psql -U postgres -f init_db_complete.sql

# 或者使用以下命令
sudo -u postgres psql -f init_db_complete.sql
```

### 4. 验证初始化

```bash
# 连接到数据库
psql -h localhost -U inarbit -d inarbit

# 在psql中查看表
\dt

# 查看视图
\dv

# 查看函数
\df

# 退出
\q
```

---

## 数据库结构

### 核心表

#### 1. users（用户表）

存储应用用户信息。

```sql
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255),
    is_active BOOLEAN DEFAULT true,
    is_admin BOOLEAN DEFAULT false,
    last_login TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
```

**字段说明**：
- `id`: 用户ID（自增主键）
- `username`: 用户名（唯一）
- `email`: 邮箱（唯一）
- `password_hash`: 密码哈希值
- `full_name`: 全名
- `is_active`: 是否激活
- `is_admin`: 是否为管理员
- `last_login`: 最后登录时间
- `created_at`: 创建时间
- `updated_at`: 更新时间
- `deleted_at`: 删除时间（软删除）

#### 2. exchanges（交易所配置表）

存储用户配置的交易所信息。

```sql
CREATE TABLE exchanges (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    exchange_type VARCHAR(50) NOT NULL,
    api_key VARCHAR(500) NOT NULL,
    api_secret VARCHAR(500) NOT NULL,
    passphrase VARCHAR(500),
    is_testnet BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    UNIQUE(user_id, exchange_type)
);
```

**字段说明**：
- `id`: 交易所配置ID
- `user_id`: 用户ID（外键）
- `name`: 交易所名称
- `exchange_type`: 交易所类型（binance, okex, huobi等）
- `api_key`: API密钥
- `api_secret`: API密钥
- `passphrase`: 密码短语（某些交易所需要）
- `is_testnet`: 是否为测试网
- `is_active`: 是否激活
- `description`: 描述

#### 3. bots（机器人表）

存储交易机器人配置。

```sql
CREATE TABLE bots (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    exchange_id BIGINT NOT NULL REFERENCES exchanges(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    strategy_type VARCHAR(50) NOT NULL,
    is_active BOOLEAN DEFAULT false,
    is_simulation BOOLEAN DEFAULT true,
    min_profit_percent FLOAT DEFAULT 0.1,
    max_concurrent_trades INT DEFAULT 5,
    update_frequency INT DEFAULT 5,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
```

**字段说明**：
- `id`: 机器人ID
- `user_id`: 用户ID（外键）
- `exchange_id`: 交易所ID（外键）
- `name`: 机器人名称
- `description`: 描述
- `strategy_type`: 策略类型（triangular, quadrangular, pentagonal）
- `is_active`: 是否激活
- `is_simulation`: 是否为模拟交易
- `min_profit_percent`: 最小利润百分比
- `max_concurrent_trades`: 最大并发交易数
- `update_frequency`: 更新频率（秒）

#### 4. strategies（策略配置表）

存储交易策略配置。

```sql
CREATE TABLE strategies (
    id BIGSERIAL PRIMARY KEY,
    bot_id BIGINT NOT NULL REFERENCES bots(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    strategy_type VARCHAR(50) NOT NULL,
    trading_pairs TEXT NOT NULL,
    base_currency VARCHAR(20) NOT NULL,
    initial_amount DECIMAL(20, 8) NOT NULL DEFAULT 100,
    max_loss_percent FLOAT DEFAULT 5.0,
    max_loss_amount DECIMAL(20, 8),
    take_profit_percent FLOAT DEFAULT 2.0,
    stop_loss_percent FLOAT DEFAULT 1.0,
    taker_fee_percent FLOAT DEFAULT 0.1,
    maker_fee_percent FLOAT DEFAULT 0.1,
    slippage_percent FLOAT DEFAULT 0.05,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
```

#### 5. trades（交易记录表）

存储所有交易记录。

```sql
CREATE TABLE trades (
    id BIGSERIAL PRIMARY KEY,
    bot_id BIGINT NOT NULL REFERENCES bots(id) ON DELETE CASCADE,
    strategy_id BIGINT REFERENCES strategies(id) ON DELETE SET NULL,
    execution_id VARCHAR(100) UNIQUE,
    status VARCHAR(50) NOT NULL,
    strategy_type VARCHAR(50) NOT NULL,
    trading_path TEXT NOT NULL,
    initial_amount DECIMAL(20, 8) NOT NULL,
    final_amount DECIMAL(20, 8),
    gross_profit DECIMAL(20, 8),
    net_profit DECIMAL(20, 8),
    profit_percent FLOAT,
    total_fees DECIMAL(20, 8) DEFAULT 0,
    execution_time_ms INT,
    orders_count INT DEFAULT 0,
    is_simulation BOOLEAN DEFAULT true,
    error_message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
```

#### 6. orders（订单表）

存储订单详细信息。

```sql
CREATE TABLE orders (
    id BIGSERIAL PRIMARY KEY,
    trade_id BIGINT NOT NULL REFERENCES trades(id) ON DELETE CASCADE,
    exchange_order_id VARCHAR(100),
    symbol VARCHAR(50) NOT NULL,
    side VARCHAR(10) NOT NULL,
    order_type VARCHAR(20) NOT NULL,
    quantity DECIMAL(20, 8) NOT NULL,
    price DECIMAL(20, 8),
    executed_quantity DECIMAL(20, 8) DEFAULT 0,
    executed_price DECIMAL(20, 8),
    commission DECIMAL(20, 8) DEFAULT 0,
    commission_asset VARCHAR(20),
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
```

### 统计和日志表

#### 7. bot_statistics（机器人统计表）

存储机器人的统计数据。

#### 8. daily_statistics（每日统计表）

存储每日的交易统计。

#### 9. system_logs（系统日志表）

存储系统运行日志。

#### 10. market_data（市场数据缓存表）

缓存实时行情数据。

#### 11. arbitrage_opportunities（套利机会表）

存储发现的套利机会。

### 视图

#### v_user_trade_summary

用户交易汇总视图，包含用户的总交易数、成功率、总利润等。

```sql
SELECT * FROM v_user_trade_summary WHERE username = 'admin';
```

#### v_bot_activity

机器人活动视图，包含机器人的交易数、利润、最后交易时间等。

```sql
SELECT * FROM v_bot_activity WHERE user_id = 1;
```

---

## 备份和恢复

### 使用备份脚本

#### 1. 备份数据库

```bash
# SQL格式备份（推荐）
./db_backup_restore.sh backup

# 自定义格式备份
./db_backup_restore.sh backup-custom
```

#### 2. 恢复数据库

```bash
# 恢复SQL格式备份
./db_backup_restore.sh restore /path/to/backup.sql.gz

# 恢复自定义格式备份
./db_backup_restore.sh restore-custom /path/to/backup.dump
```

#### 3. 导出表数据

```bash
# 导出users表
./db_backup_restore.sh export users

# 导出trades表
./db_backup_restore.sh export trades
```

#### 4. 导入表数据

```bash
# 导入users表
./db_backup_restore.sh import users /path/to/users.csv
```

#### 5. 列出备份

```bash
./db_backup_restore.sh list
```

#### 6. 清理旧备份

```bash
# 清理30天前的备份
./db_backup_restore.sh cleanup 30

# 清理60天前的备份
./db_backup_restore.sh cleanup 60
```

### 手动备份

#### 完整备份

```bash
# SQL格式
pg_dump -h localhost -U inarbit -d inarbit > backup.sql

# 压缩
gzip backup.sql

# 自定义格式
pg_dump -h localhost -U inarbit -d inarbit -Fc -f backup.dump
```

#### 表级备份

```bash
# 备份特定表
pg_dump -h localhost -U inarbit -d inarbit -t users > users.sql

# 备份多个表
pg_dump -h localhost -U inarbit -d inarbit -t users -t bots > users_bots.sql
```

### 恢复

```bash
# 从SQL文件恢复
psql -h localhost -U inarbit -d inarbit < backup.sql

# 从压缩文件恢复
gunzip -c backup.sql.gz | psql -h localhost -U inarbit -d inarbit

# 从自定义格式恢复
pg_restore -h localhost -U inarbit -d inarbit backup.dump
```

---

## 性能优化

### 1. 连接池配置

编辑 `/etc/postgresql/14/main/postgresql.conf`：

```conf
# 连接参数
max_connections = 200
superuser_reserved_connections = 3

# 内存参数
shared_buffers = 256MB
effective_cache_size = 1GB
work_mem = 4MB
maintenance_work_mem = 64MB

# 日志参数
log_min_duration_statement = 1000
log_statement = 'all'
```

### 2. 索引优化

查看未使用的索引：

```sql
SELECT schemaname, tablename, indexname, idx_scan
FROM pg_stat_user_indexes
WHERE idx_scan = 0
ORDER BY pg_relation_size(indexrelid) DESC;
```

删除未使用的索引：

```sql
DROP INDEX IF EXISTS index_name;
```

### 3. 查询优化

使用EXPLAIN分析查询：

```sql
EXPLAIN ANALYZE
SELECT * FROM trades WHERE bot_id = 1 AND created_at > NOW() - INTERVAL '7 days';
```

### 4. 表维护

```sql
-- 分析表
ANALYZE trades;

-- 清理表（删除死元组）
VACUUM trades;

-- 重建索引
REINDEX TABLE trades;
```

### 5. 自动清理配置

编辑 `/etc/postgresql/14/main/postgresql.conf`：

```conf
autovacuum = on
autovacuum_naptime = 10s
autovacuum_vacuum_threshold = 50
autovacuum_analyze_threshold = 50
```

---

## 故障排除

### 1. 连接失败

**症状**: `psql: could not connect to server`

**解决方案**:

```bash
# 检查PostgreSQL服务状态
sudo systemctl status postgresql

# 启动PostgreSQL
sudo systemctl start postgresql

# 检查监听端口
sudo netstat -tlnp | grep postgres

# 检查PostgreSQL日志
sudo tail -f /var/log/postgresql/postgresql-14-main.log
```

### 2. 权限错误

**症状**: `permission denied`

**解决方案**:

```bash
# 重新授予权限
sudo -u postgres psql -d inarbit -c "GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO inarbit;"
sudo -u postgres psql -d inarbit -c "GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO inarbit;"
```

### 3. 磁盘空间不足

**症状**: `disk full`

**解决方案**:

```bash
# 检查磁盘使用情况
df -h

# 清理日志
sudo truncate -s 0 /var/log/postgresql/postgresql-14-main.log

# 清理备份
./db_backup_restore.sh cleanup 30

# 执行VACUUM
psql -h localhost -U inarbit -d inarbit -c "VACUUM FULL;"
```

### 4. 性能下降

**症状**: 查询变慢

**解决方案**:

```bash
# 分析慢查询
sudo -u postgres psql -d inarbit << EOF
SELECT query, calls, total_time, mean_time
FROM pg_stat_statements
ORDER BY mean_time DESC
LIMIT 10;
EOF

# 重建索引
psql -h localhost -U inarbit -d inarbit -c "REINDEX DATABASE inarbit;"

# 执行ANALYZE
psql -h localhost -U inarbit -d inarbit -c "ANALYZE;"
```

### 5. 数据损坏

**症状**: 数据库无法启动或查询错误

**解决方案**:

```bash
# 停止PostgreSQL
sudo systemctl stop postgresql

# 检查数据库
sudo -u postgres pg_verify_checksums -D /var/lib/postgresql/14/main

# 重建数据库（最后手段）
sudo -u postgres pg_dump -d inarbit > backup.sql
sudo -u postgres dropdb inarbit
sudo -u postgres createdb -O inarbit inarbit
psql -h localhost -U inarbit -d inarbit < backup.sql
```

---

## 安全加固

### 1. 密码安全

```bash
# 修改用户密码
psql -h localhost -U inarbit -d inarbit -c "ALTER USER inarbit WITH PASSWORD 'new_password';"
```

### 2. 连接安全

编辑 `/etc/postgresql/14/main/pg_hba.conf`：

```conf
# 只允许本地连接
local   all             all                                     trust
host    all             all             127.0.0.1/32            md5
host    all             all             ::1/128                 md5

# 允许特定IP连接
host    inarbit         inarbit         192.168.1.0/24          md5
```

### 3. 防火墙配置

```bash
# 只允许本地连接
sudo ufw allow from 127.0.0.1 to any port 5432

# 允许特定IP
sudo ufw allow from 192.168.1.100 to any port 5432
```

### 4. SSL连接

```bash
# 生成自签名证书
sudo -u postgres openssl req -new -x509 -days 365 -nodes \
  -out /etc/postgresql/14/main/server.crt \
  -keyout /etc/postgresql/14/main/server.key

# 设置权限
sudo chmod 600 /etc/postgresql/14/main/server.key
sudo chown postgres:postgres /etc/postgresql/14/main/server.*

# 启用SSL
sudo -u postgres psql -c "ALTER SYSTEM SET ssl = on;"
sudo systemctl restart postgresql
```

### 5. 审计日志

编辑 `/etc/postgresql/14/main/postgresql.conf`：

```conf
log_connections = on
log_disconnections = on
log_statement = 'all'
log_min_duration_statement = 0
```

---

## 常见问题

### Q1: 如何修改数据库用户密码？

```bash
psql -h localhost -U inarbit -d inarbit -c "ALTER USER inarbit WITH PASSWORD 'new_password';"
```

### Q2: 如何导出所有数据？

```bash
pg_dump -h localhost -U inarbit -d inarbit > full_backup.sql
gzip full_backup.sql
```

### Q3: 如何查看数据库大小？

```bash
psql -h localhost -U inarbit -d inarbit -c "SELECT pg_size_pretty(pg_database_size('inarbit'));"
```

### Q4: 如何查看表的行数？

```bash
psql -h localhost -U inarbit -d inarbit -c "SELECT tablename, n_live_tup FROM pg_stat_user_tables ORDER BY n_live_tup DESC;"
```

### Q5: 如何删除旧数据？

```bash
-- 删除30天前的交易记录
DELETE FROM trades WHERE created_at < NOW() - INTERVAL '30 days';

-- 删除30天前的日志
DELETE FROM system_logs WHERE created_at < NOW() - INTERVAL '30 days';

-- 执行VACUUM清理
VACUUM ANALYZE;
```

### Q6: 如何监控数据库性能？

```bash
# 查看活跃连接
psql -h localhost -U inarbit -d inarbit -c "SELECT * FROM pg_stat_activity;"

# 查看缓存命中率
psql -h localhost -U inarbit -d inarbit -c "SELECT sum(heap_blks_read) as heap_read, sum(heap_blks_hit) as heap_hit, sum(heap_blks_hit) / (sum(heap_blks_hit) + sum(heap_blks_read)) as ratio FROM pg_statio_user_tables;"

# 查看慢查询
psql -h localhost -U inarbit -d inarbit -c "SELECT query, calls, total_time, mean_time FROM pg_stat_statements ORDER BY mean_time DESC LIMIT 10;"
```

### Q7: 如何设置自动备份？

```bash
# 创建cron任务
crontab -e

# 添加以下行（每天凌晨2点备份）
0 2 * * * /root/inarbit/db_backup_restore.sh backup

# 添加以下行（每周清理备份）
0 3 * * 0 /root/inarbit/db_backup_restore.sh cleanup 30
```

### Q8: 如何处理事务锁定？

```bash
-- 查看锁定
SELECT * FROM pg_locks WHERE NOT granted;

-- 查看阻塞的查询
SELECT pid, usename, application_name, state, query
FROM pg_stat_activity
WHERE state != 'idle';

-- 杀死进程
SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE pid <> pg_backend_pid();
```

---

## 总结

本指南涵盖了iNarbit PostgreSQL数据库的完整部署、维护和优化。定期备份、监控性能和及时处理问题是确保数据库稳定运行的关键。

如有问题，请参考PostgreSQL官方文档：https://www.postgresql.org/docs/

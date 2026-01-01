-- ============================================================================
-- iNarbit PostgreSQL 数据库初始化脚本（完整版）
-- 功能：创建用户、数据库、表结构、索引、函数和初始数据
-- 使用：psql -U postgres -f init_db_complete.sql
-- ============================================================================

-- ============================================================================
-- 第一部分：创建用户和数据库
-- ============================================================================

-- 创建应用用户
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_user WHERE usename = 'inarbit') THEN
        CREATE USER inarbit WITH PASSWORD 'inarbit_password';
    END IF;
END
$$;

-- 授予用户权限
ALTER USER inarbit CREATEDB;

-- 创建数据库
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'inarbit') THEN
        CREATE DATABASE inarbit OWNER inarbit;
    END IF;
END
$$;

-- 连接到数据库
\c inarbit

-- 授予数据库权限
GRANT ALL PRIVILEGES ON DATABASE inarbit TO inarbit;

-- ============================================================================
-- 第二部分：创建扩展
-- ============================================================================

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ============================================================================
-- 第三部分：创建表结构
-- ============================================================================

-- 1. 用户表
CREATE TABLE IF NOT EXISTS users (
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

-- 2. 交易所配置表
CREATE TABLE IF NOT EXISTS exchanges (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    exchange_type VARCHAR(50) NOT NULL, -- binance, okex, huobi, etc.
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

-- 3. 机器人表
CREATE TABLE IF NOT EXISTS bots (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    exchange_id BIGINT NOT NULL REFERENCES exchanges(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    strategy_type VARCHAR(50) NOT NULL, -- triangular, quadrangular, pentagonal
    is_active BOOLEAN DEFAULT false,
    is_simulation BOOLEAN DEFAULT true,
    min_profit_percent FLOAT DEFAULT 0.1,
    max_concurrent_trades INT DEFAULT 5,
    update_frequency INT DEFAULT 5, -- 秒
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 4. 策略配置表
CREATE TABLE IF NOT EXISTS strategies (
    id BIGSERIAL PRIMARY KEY,
    bot_id BIGINT NOT NULL REFERENCES bots(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    strategy_type VARCHAR(50) NOT NULL,
    trading_pairs TEXT NOT NULL, -- JSON数组
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

-- 5. 交易记录表
CREATE TABLE IF NOT EXISTS trades (
    id BIGSERIAL PRIMARY KEY,
    bot_id BIGINT NOT NULL REFERENCES bots(id) ON DELETE CASCADE,
    strategy_id BIGINT REFERENCES strategies(id) ON DELETE SET NULL,
    execution_id VARCHAR(100) UNIQUE,
    status VARCHAR(50) NOT NULL, -- pending, executing, completed, failed, cancelled
    strategy_type VARCHAR(50) NOT NULL,
    trading_path TEXT NOT NULL, -- JSON数组
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

-- 6. 订单表
CREATE TABLE IF NOT EXISTS orders (
    id BIGSERIAL PRIMARY KEY,
    trade_id BIGINT NOT NULL REFERENCES trades(id) ON DELETE CASCADE,
    exchange_order_id VARCHAR(100),
    symbol VARCHAR(50) NOT NULL,
    side VARCHAR(10) NOT NULL, -- BUY, SELL
    order_type VARCHAR(20) NOT NULL, -- LIMIT, MARKET
    quantity DECIMAL(20, 8) NOT NULL,
    price DECIMAL(20, 8),
    executed_quantity DECIMAL(20, 8) DEFAULT 0,
    executed_price DECIMAL(20, 8),
    commission DECIMAL(20, 8) DEFAULT 0,
    commission_asset VARCHAR(20),
    status VARCHAR(50) NOT NULL, -- PENDING, PARTIALLY_FILLED, FILLED, CANCELLED, REJECTED
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 7. 机器人统计表
CREATE TABLE IF NOT EXISTS bot_statistics (
    id BIGSERIAL PRIMARY KEY,
    bot_id BIGINT NOT NULL REFERENCES bots(id) ON DELETE CASCADE,
    total_trades BIGINT DEFAULT 0,
    successful_trades BIGINT DEFAULT 0,
    failed_trades BIGINT DEFAULT 0,
    total_profit DECIMAL(20, 8) DEFAULT 0,
    total_loss DECIMAL(20, 8) DEFAULT 0,
    total_fees DECIMAL(20, 8) DEFAULT 0,
    win_rate FLOAT DEFAULT 0,
    average_profit DECIMAL(20, 8) DEFAULT 0,
    best_trade DECIMAL(20, 8) DEFAULT 0,
    worst_trade DECIMAL(20, 8) DEFAULT 0,
    uptime_percent FLOAT DEFAULT 0,
    last_trade_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(bot_id)
);

-- 8. 每日统计表
CREATE TABLE IF NOT EXISTS daily_statistics (
    id BIGSERIAL PRIMARY KEY,
    bot_id BIGINT NOT NULL REFERENCES bots(id) ON DELETE CASCADE,
    trade_date DATE NOT NULL,
    total_trades INT DEFAULT 0,
    successful_trades INT DEFAULT 0,
    failed_trades INT DEFAULT 0,
    daily_profit DECIMAL(20, 8) DEFAULT 0,
    daily_loss DECIMAL(20, 8) DEFAULT 0,
    daily_fees DECIMAL(20, 8) DEFAULT 0,
    win_rate FLOAT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(bot_id, trade_date)
);

-- 9. 系统日志表
CREATE TABLE IF NOT EXISTS system_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    bot_id BIGINT REFERENCES bots(id) ON DELETE SET NULL,
    log_level VARCHAR(20) NOT NULL, -- INFO, WARN, ERROR, DEBUG
    message TEXT NOT NULL,
    details JSONB,
    ip_address VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 10. 市场数据缓存表
CREATE TABLE IF NOT EXISTS market_data (
    id BIGSERIAL PRIMARY KEY,
    exchange_id BIGINT NOT NULL REFERENCES exchanges(id) ON DELETE CASCADE,
    symbol VARCHAR(50) NOT NULL,
    bid_price DECIMAL(20, 8) NOT NULL,
    ask_price DECIMAL(20, 8) NOT NULL,
    last_price DECIMAL(20, 8) NOT NULL,
    bid_quantity DECIMAL(20, 8),
    ask_quantity DECIMAL(20, 8),
    volume_24h DECIMAL(20, 8),
    price_change_percent FLOAT,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(exchange_id, symbol)
);

-- 11. 套利机会表
CREATE TABLE IF NOT EXISTS arbitrage_opportunities (
    id BIGSERIAL PRIMARY KEY,
    bot_id BIGINT REFERENCES bots(id) ON DELETE SET NULL,
    strategy_type VARCHAR(50) NOT NULL,
    trading_path TEXT NOT NULL, -- JSON数组
    initial_amount DECIMAL(20, 8) NOT NULL,
    final_amount DECIMAL(20, 8) NOT NULL,
    gross_profit DECIMAL(20, 8) NOT NULL,
    net_profit DECIMAL(20, 8) NOT NULL,
    profit_percent FLOAT NOT NULL,
    confidence_score FLOAT DEFAULT 0,
    risk_level VARCHAR(20), -- LOW, MEDIUM, HIGH
    is_executed BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expired_at TIMESTAMP
);

-- 12. 配置表
CREATE TABLE IF NOT EXISTS configurations (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    config_key VARCHAR(255) NOT NULL,
    config_value TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, config_key)
);

-- 13. API密钥表
CREATE TABLE IF NOT EXISTS api_keys (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    key_name VARCHAR(255) NOT NULL,
    api_key VARCHAR(500) NOT NULL UNIQUE,
    api_secret VARCHAR(500) NOT NULL,
    permissions TEXT[],
    is_active BOOLEAN DEFAULT true,
    last_used_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 14. 审计日志表
CREATE TABLE IF NOT EXISTS audit_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(100) NOT NULL,
    resource_type VARCHAR(100) NOT NULL,
    resource_id BIGINT,
    old_values JSONB,
    new_values JSONB,
    ip_address VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================================
-- 第四部分：创建索引
-- ============================================================================

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_is_active ON users(is_active);
CREATE INDEX idx_users_created_at ON users(created_at);

CREATE INDEX idx_exchanges_user_id ON exchanges(user_id);
CREATE INDEX idx_exchanges_exchange_type ON exchanges(exchange_type);
CREATE INDEX idx_exchanges_is_active ON exchanges(is_active);

CREATE INDEX idx_bots_user_id ON bots(user_id);
CREATE INDEX idx_bots_exchange_id ON bots(exchange_id);
CREATE INDEX idx_bots_strategy_type ON bots(strategy_type);
CREATE INDEX idx_bots_is_active ON bots(is_active);
CREATE INDEX idx_bots_created_at ON bots(created_at);

CREATE INDEX idx_strategies_bot_id ON strategies(bot_id);
CREATE INDEX idx_strategies_strategy_type ON strategies(strategy_type);

CREATE INDEX idx_trades_bot_id ON trades(bot_id);
CREATE INDEX idx_trades_strategy_id ON trades(strategy_id);
CREATE INDEX idx_trades_status ON trades(status);
CREATE INDEX idx_trades_created_at ON trades(created_at);
CREATE INDEX idx_trades_is_simulation ON trades(is_simulation);

CREATE INDEX idx_orders_trade_id ON orders(trade_id);
CREATE INDEX idx_orders_symbol ON orders(symbol);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_created_at ON orders(created_at);

CREATE INDEX idx_bot_statistics_bot_id ON bot_statistics(bot_id);

CREATE INDEX idx_daily_statistics_bot_id ON daily_statistics(bot_id);
CREATE INDEX idx_daily_statistics_trade_date ON daily_statistics(trade_date);

CREATE INDEX idx_system_logs_user_id ON system_logs(user_id);
CREATE INDEX idx_system_logs_bot_id ON system_logs(bot_id);
CREATE INDEX idx_system_logs_log_level ON system_logs(log_level);
CREATE INDEX idx_system_logs_created_at ON system_logs(created_at);

CREATE INDEX idx_market_data_exchange_id ON market_data(exchange_id);
CREATE INDEX idx_market_data_symbol ON market_data(symbol);

CREATE INDEX idx_arbitrage_opportunities_bot_id ON arbitrage_opportunities(bot_id);
CREATE INDEX idx_arbitrage_opportunities_created_at ON arbitrage_opportunities(created_at);

CREATE INDEX idx_configurations_user_id ON configurations(user_id);
CREATE INDEX idx_configurations_config_key ON configurations(config_key);

CREATE INDEX idx_api_keys_user_id ON api_keys(user_id);
CREATE INDEX idx_api_keys_api_key ON api_keys(api_key);

CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_resource_type ON audit_logs(resource_type);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);

-- ============================================================================
-- 第五部分：创建触发器和函数
-- ============================================================================

-- 自动更新updated_at时间戳的函数
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- 为所有需要的表创建触发器
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_exchanges_updated_at BEFORE UPDATE ON exchanges
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_bots_updated_at BEFORE UPDATE ON bots
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_strategies_updated_at BEFORE UPDATE ON strategies
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_trades_updated_at BEFORE UPDATE ON trades
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_orders_updated_at BEFORE UPDATE ON orders
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_bot_statistics_updated_at BEFORE UPDATE ON bot_statistics
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_daily_statistics_updated_at BEFORE UPDATE ON daily_statistics
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_configurations_updated_at BEFORE UPDATE ON configurations
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_api_keys_updated_at BEFORE UPDATE ON api_keys
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- 计算机器人统计信息的函数
CREATE OR REPLACE FUNCTION calculate_bot_statistics(bot_id_param BIGINT)
RETURNS TABLE (
    total_trades BIGINT,
    successful_trades BIGINT,
    failed_trades BIGINT,
    total_profit DECIMAL,
    total_loss DECIMAL,
    win_rate FLOAT
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        COUNT(*) as total_trades,
        COUNT(*) FILTER (WHERE net_profit > 0) as successful_trades,
        COUNT(*) FILTER (WHERE net_profit <= 0) as failed_trades,
        COALESCE(SUM(net_profit) FILTER (WHERE net_profit > 0), 0) as total_profit,
        COALESCE(SUM(net_profit) FILTER (WHERE net_profit <= 0), 0) as total_loss,
        CASE 
            WHEN COUNT(*) = 0 THEN 0
            ELSE COUNT(*) FILTER (WHERE net_profit > 0)::FLOAT / COUNT(*)::FLOAT * 100
        END as win_rate
    FROM trades
    WHERE bot_id = bot_id_param AND deleted_at IS NULL;
END;
$$ LANGUAGE plpgsql;

-- 获取用户的总利润函数
CREATE OR REPLACE FUNCTION get_user_total_profit(user_id_param BIGINT)
RETURNS DECIMAL AS $$
DECLARE
    total_profit DECIMAL;
BEGIN
    SELECT COALESCE(SUM(net_profit), 0)
    INTO total_profit
    FROM trades t
    JOIN bots b ON t.bot_id = b.id
    WHERE b.user_id = user_id_param AND t.deleted_at IS NULL;
    
    RETURN total_profit;
END;
$$ LANGUAGE plpgsql;

-- ============================================================================
-- 第六部分：创建视图
-- ============================================================================

-- 用户交易汇总视图
CREATE OR REPLACE VIEW v_user_trade_summary AS
SELECT
    u.id as user_id,
    u.username,
    COUNT(DISTINCT t.id) as total_trades,
    COUNT(DISTINCT t.id) FILTER (WHERE t.net_profit > 0) as successful_trades,
    COUNT(DISTINCT t.id) FILTER (WHERE t.net_profit <= 0) as failed_trades,
    COALESCE(SUM(t.net_profit), 0) as total_profit,
    COALESCE(SUM(t.total_fees), 0) as total_fees,
    CASE 
        WHEN COUNT(DISTINCT t.id) = 0 THEN 0
        ELSE COUNT(DISTINCT t.id) FILTER (WHERE t.net_profit > 0)::FLOAT / COUNT(DISTINCT t.id)::FLOAT * 100
    END as win_rate,
    MAX(t.created_at) as last_trade_at
FROM users u
LEFT JOIN bots b ON u.id = b.user_id
LEFT JOIN trades t ON b.id = t.bot_id AND t.deleted_at IS NULL
WHERE u.deleted_at IS NULL
GROUP BY u.id, u.username;

-- 机器人活动视图
CREATE OR REPLACE VIEW v_bot_activity AS
SELECT
    b.id as bot_id,
    b.name,
    b.user_id,
    b.is_active,
    b.is_simulation,
    COUNT(DISTINCT t.id) as total_trades,
    COALESCE(SUM(t.net_profit), 0) as total_profit,
    MAX(t.created_at) as last_trade_at,
    COUNT(DISTINCT CASE WHEN t.created_at >= CURRENT_TIMESTAMP - INTERVAL '24 hours' THEN t.id END) as trades_24h,
    COALESCE(SUM(t.net_profit) FILTER (WHERE t.created_at >= CURRENT_TIMESTAMP - INTERVAL '24 hours'), 0) as profit_24h
FROM bots b
LEFT JOIN trades t ON b.id = t.bot_id AND t.deleted_at IS NULL
WHERE b.deleted_at IS NULL
GROUP BY b.id, b.name, b.user_id, b.is_active, b.is_simulation;

-- ============================================================================
-- 第七部分：插入初始数据
-- ============================================================================

-- 创建默认管理员用户
INSERT INTO users (username, email, password_hash, full_name, is_active, is_admin)
VALUES (
    'admin',
    'admin@inarbit.work',
    crypt('password', gen_salt('bf')),
    'Administrator',
    true,
    true
) ON CONFLICT (username) DO NOTHING;

-- 创建默认配置
INSERT INTO configurations (user_id, config_key, config_value, description)
SELECT 
    u.id,
    'min_profit_threshold',
    '0.1',
    '最小利润阈值（百分比）'
FROM users u WHERE u.username = 'admin'
ON CONFLICT (user_id, config_key) DO NOTHING;

INSERT INTO configurations (user_id, config_key, config_value, description)
SELECT 
    u.id,
    'max_concurrent_trades',
    '5',
    '最大并发交易数'
FROM users u WHERE u.username = 'admin'
ON CONFLICT (user_id, config_key) DO NOTHING;

INSERT INTO configurations (user_id, config_key, config_value, description)
SELECT 
    u.id,
    'market_update_frequency',
    '5',
    '市场数据更新频率（秒）'
FROM users u WHERE u.username = 'admin'
ON CONFLICT (user_id, config_key) DO NOTHING;

-- ============================================================================
-- 第八部分：授予权限
-- ============================================================================

GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO inarbit;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO inarbit;
GRANT EXECUTE ON ALL FUNCTIONS IN SCHEMA public TO inarbit;

-- ============================================================================
-- 初始化完成
-- ============================================================================

\echo '=========================================='
\echo 'iNarbit 数据库初始化完成'
\echo '=========================================='
\echo ''
\echo '已创建的用户:'
\echo '  用户名: admin'
\echo '  密码: password'
\echo ''
\echo '数据库连接信息:'
\echo '  主机: localhost'
\echo '  端口: 5432'
\echo '  用户: inarbit'
\echo '  密码: inarbit_password'
\echo '  数据库: inarbit'
\echo ''
\echo '=========================================='

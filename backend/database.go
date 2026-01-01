package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// Database 数据库连接包装
type Database struct {
	DB *sql.DB
}

// InitDatabase 初始化数据库连接
func InitDatabase(config *Config) (*Database, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.DBHost,
		config.DBPort,
		config.DBUser,
		config.DBPassword,
		config.DBName,
		"disable", // 开发环境禁用SSL，生产环境改为require
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("数据库连接失败: %w", err)
	}

	// 测试连接
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("数据库ping失败: %w", err)
	}

	// 设置连接池参数
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	log.Println("✓ 数据库连接成功")

	return &Database{DB: db}, nil
}

// Close 关闭数据库连接
func (d *Database) Close() error {
	return d.DB.Close()
}

// GetUser 获取用户信息
func (d *Database) GetUser(username string) (*User, error) {
	user := &User{}
	err := d.DB.QueryRow(
		"SELECT id, username, password_hash, email, is_active FROM users WHERE username = $1",
		username,
	).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Email, &user.IsActive)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, err
	}

	return user, nil
}

// GetUserByID 根据ID获取用户
func (d *Database) GetUserByID(id int64) (*User, error) {
	user := &User{}
	err := d.DB.QueryRow(
		"SELECT id, username, email, is_active FROM users WHERE id = $1",
		id,
	).Scan(&user.ID, &user.Username, &user.Email, &user.IsActive)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, err
	}

	return user, nil
}

// GetBots 获取用户的所有机器人
func (d *Database) GetBots(userID int64) ([]*Bot, error) {
	rows, err := d.DB.Query(
		`SELECT id, user_id, name, strategy_type, exchange_id, is_running, is_simulation, 
		        update_frequency, created_at, total_profit, total_trades
		 FROM bots WHERE user_id = $1 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bots []*Bot
	for rows.Next() {
		bot := &Bot{}
		err := rows.Scan(
			&bot.ID, &bot.UserID, &bot.Name, &bot.StrategyType, &bot.ExchangeID,
			&bot.IsRunning, &bot.IsSimulation, &bot.UpdateFrequency,
			&bot.CreatedAt, &bot.TotalProfit, &bot.TotalTrades,
		)
		if err != nil {
			return nil, err
		}
		bots = append(bots, bot)
	}

	return bots, rows.Err()
}

// GetBotByID 获取单个机器人
func (d *Database) GetBotByID(id int64, userID int64) (*Bot, error) {
	bot := &Bot{}
	err := d.DB.QueryRow(
		`SELECT id, user_id, name, strategy_type, exchange_id, is_running, is_simulation, 
		        update_frequency, created_at, total_profit, total_trades
		 FROM bots WHERE id = $1 AND user_id = $2`,
		id, userID,
	).Scan(
		&bot.ID, &bot.UserID, &bot.Name, &bot.StrategyType, &bot.ExchangeID,
		&bot.IsRunning, &bot.IsSimulation, &bot.UpdateFrequency,
		&bot.CreatedAt, &bot.TotalProfit, &bot.TotalTrades,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("机器人不存在")
		}
		return nil, err
	}

	return bot, nil
}

// CreateBot 创建机器人
func (d *Database) CreateBot(bot *Bot) error {
	err := d.DB.QueryRow(
		`INSERT INTO bots (user_id, name, strategy_type, exchange_id, is_running, is_simulation, update_frequency)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id, created_at`,
		bot.UserID, bot.Name, bot.StrategyType, bot.ExchangeID,
		bot.IsRunning, bot.IsSimulation, bot.UpdateFrequency,
	).Scan(&bot.ID, &bot.CreatedAt)

	return err
}

// UpdateBot 更新机器人
func (d *Database) UpdateBot(bot *Bot) error {
	result, err := d.DB.Exec(
		`UPDATE bots SET name = $1, strategy_type = $2, exchange_id = $3, 
		        is_running = $4, is_simulation = $5, update_frequency = $6, updated_at = NOW()
		 WHERE id = $7 AND user_id = $8`,
		bot.Name, bot.StrategyType, bot.ExchangeID,
		bot.IsRunning, bot.IsSimulation, bot.UpdateFrequency,
		bot.ID, bot.UserID,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("机器人不存在或无权限修改")
	}

	return nil
}

// DeleteBot 删除机器人
func (d *Database) DeleteBot(id int64, userID int64) error {
	result, err := d.DB.Exec(
		"DELETE FROM bots WHERE id = $1 AND user_id = $2",
		id, userID,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("机器人不存在或无权限删除")
	}

	return nil
}

// UpdateBotStatus 更新机器人运行状态
func (d *Database) UpdateBotStatus(id int64, userID int64, isRunning bool) error {
	result, err := d.DB.Exec(
		`UPDATE bots SET is_running = $1, updated_at = NOW() WHERE id = $2 AND user_id = $3`,
		isRunning, id, userID,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("机器人不存在")
	}

	return nil
}

// GetDashboardStats 获取仪表板统计数据
func (d *Database) GetDashboardStats(userID int64) (*DashboardStats, error) {
	stats := &DashboardStats{}

	// 获取总收益
	err := d.DB.QueryRow(
		"SELECT COALESCE(SUM(profit), 0) FROM trades WHERE bot_id IN (SELECT id FROM bots WHERE user_id = $1) AND status = 'completed'",
		userID,
	).Scan(&stats.TotalProfit)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// 获取总交易数
	err = d.DB.QueryRow(
		"SELECT COUNT(*) FROM trades WHERE bot_id IN (SELECT id FROM bots WHERE user_id = $1) AND status = 'completed'",
		userID,
	).Scan(&stats.TotalTrades)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// 获取活跃机器人数
	err = d.DB.QueryRow(
		"SELECT COUNT(*) FROM bots WHERE user_id = $1 AND is_running = true",
		userID,
	).Scan(&stats.ActiveBots)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// 计算胜率
	var successCount int64
	err = d.DB.QueryRow(
		`SELECT COUNT(*) FROM trades 
		 WHERE bot_id IN (SELECT id FROM bots WHERE user_id = $1) 
		 AND status = 'completed' AND profit > 0`,
		userID,
	).Scan(&successCount)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if stats.TotalTrades > 0 {
		stats.WinRate = float64(successCount) / float64(stats.TotalTrades) * 100
	}

	return stats, nil
}

// GetChartData 获取图表数据（最近7天）
func (d *Database) GetChartData(userID int64) ([]*ChartDataPoint, error) {
	rows, err := d.DB.Query(
		`SELECT DATE(created_at) as date, 
		        COALESCE(SUM(profit), 0) as profit, 
		        COUNT(*) as trades
		 FROM trades 
		 WHERE bot_id IN (SELECT id FROM bots WHERE user_id = $1) 
		 AND created_at >= NOW() - INTERVAL '7 days'
		 AND status = 'completed'
		 GROUP BY DATE(created_at)
		 ORDER BY date ASC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []*ChartDataPoint
	for rows.Next() {
		point := &ChartDataPoint{}
		err := rows.Scan(&point.Date, &point.Profit, &point.Trades)
		if err != nil {
			return nil, err
		}
		data = append(data, point)
	}

	return data, rows.Err()
}

// RecordTrade 记录交易
func (d *Database) RecordTrade(trade *Trade) error {
	err := d.DB.QueryRow(
		`INSERT INTO trades (bot_id, strategy_id, trade_type, status, pair1, pair2, pair3, 
		                      initial_amount, final_amount, profit, profit_percentage, total_fees, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, NOW())
		 RETURNING id, created_at`,
		trade.BotID, trade.StrategyID, trade.TradeType, trade.Status,
		trade.Pair1, trade.Pair2, trade.Pair3,
		trade.InitialAmount, trade.FinalAmount, trade.Profit, trade.ProfitPercentage,
		trade.TotalFees,
	).Scan(&trade.ID, &trade.CreatedAt)

	return err
}

// GetExchanges 获取用户的交易所配置
func (d *Database) GetExchanges(userID int64) ([]*Exchange, error) {
	rows, err := d.DB.Query(
		"SELECT id, user_id, name, is_testnet, is_active FROM exchanges WHERE user_id = $1 ORDER BY created_at DESC",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exchanges []*Exchange
	for rows.Next() {
		exchange := &Exchange{}
		err := rows.Scan(&exchange.ID, &exchange.UserID, &exchange.Name, &exchange.IsTestnet, &exchange.IsActive)
		if err != nil {
			return nil, err
		}
		exchanges = append(exchanges, exchange)
	}

	return exchanges, rows.Err()
}

// GetExchangeByID 获取单个交易所配置
func (d *Database) GetExchangeByID(id int64, userID int64) (*Exchange, error) {
	exchange := &Exchange{}
	err := d.DB.QueryRow(
		"SELECT id, user_id, name, api_key, api_secret, is_testnet, is_active FROM exchanges WHERE id = $1 AND user_id = $2",
		id, userID,
	).Scan(&exchange.ID, &exchange.UserID, &exchange.Name, &exchange.APIKey, &exchange.APISecret, &exchange.IsTestnet, &exchange.IsActive)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("交易所配置不存在")
		}
		return nil, err
	}

	return exchange, nil
}

// CreateExchange 创建交易所配置
func (d *Database) CreateExchange(exchange *Exchange) error {
	err := d.DB.QueryRow(
		`INSERT INTO exchanges (user_id, name, api_key, api_secret, is_testnet, is_active)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, created_at`,
		exchange.UserID, exchange.Name, exchange.APIKey, exchange.APISecret, exchange.IsTestnet, exchange.IsActive,
	).Scan(&exchange.ID, &exchange.CreatedAt)

	return err
}

// DeleteExchange 删除交易所配置
func (d *Database) DeleteExchange(id int64, userID int64) error {
	result, err := d.DB.Exec(
		"DELETE FROM exchanges WHERE id = $1 AND user_id = $2",
		id, userID,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("交易所配置不存在或无权限删除")
	}

	return nil
}

// LogSystemEvent 记录系统日志
func (d *Database) LogSystemEvent(botID *int64, logLevel string, message string, details interface{}) error {
	_, err := d.DB.Exec(
		`INSERT INTO system_logs (bot_id, log_level, message, created_at)
		 VALUES ($1, $2, $3, NOW())`,
		botID, logLevel, message,
	)
	return err
}

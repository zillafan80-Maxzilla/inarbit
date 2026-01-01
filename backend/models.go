package main

import "time"

// User 用户模型
type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"` // 不序列化密码哈希
	Email        string    `json:"email"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	LastLogin    *time.Time `json:"last_login"`
}

// Bot 机器人模型
type Bot struct {
	ID              int64     `json:"id"`
	UserID          int64     `json:"user_id"`
	Name            string    `json:"name"`
	StrategyType    string    `json:"strategy_type"` // triangular, quadrangular, pentagonal
	ExchangeID      int64     `json:"exchange_id"`
	IsRunning       bool      `json:"is_running"`
	IsSimulation    bool      `json:"is_simulation"`
	UpdateFrequency int       `json:"update_frequency"` // 秒
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	StartedAt       *time.Time `json:"started_at"`
	StoppedAt       *time.Time `json:"stopped_at"`
	TotalProfit     float64   `json:"total_profit"`
	TotalTrades     int64     `json:"total_trades"`
}

// Strategy 策略模型
type Strategy struct {
	ID                   int64     `json:"id"`
	BotID                int64     `json:"bot_id"`
	Name                 string    `json:"name"`
	StrategyType         string    `json:"strategy_type"`
	BasePair             string    `json:"base_pair"`
	QuoteCurrency        string    `json:"quote_currency"`
	Pair1                string    `json:"pair1"`
	Pair2                string    `json:"pair2"`
	Pair3                string    `json:"pair3"`
	Pair4                *string   `json:"pair4"`
	Pair5                *string   `json:"pair5"`
	MinProfitPercentage  float64   `json:"min_profit_percentage"`
	MaxTradeAmount       float64   `json:"max_trade_amount"`
	MinTradeAmount       float64   `json:"min_trade_amount"`
	MaxLossPercentage    float64   `json:"max_loss_percentage"`
	MaxConcurrentTrades  int       `json:"max_concurrent_trades"`
	UseMargin            bool      `json:"use_margin"`
	Leverage             float64   `json:"leverage"`
	IsActive             bool      `json:"is_active"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// Trade 交易模型
type Trade struct {
	ID                int64      `json:"id"`
	BotID             int64      `json:"bot_id"`
	StrategyID        int64      `json:"strategy_id"`
	TradeType         string     `json:"trade_type"` // triangular, quadrangular, pentagonal
	Status            string     `json:"status"`     // pending, executing, completed, failed, cancelled
	Pair1             string     `json:"pair1"`
	Pair2             string     `json:"pair2"`
	Pair3             string     `json:"pair3"`
	Pair4             *string    `json:"pair4"`
	Pair5             *string    `json:"pair5"`
	InitialAmount     float64    `json:"initial_amount"`
	FinalAmount       *float64   `json:"final_amount"`
	Profit            *float64   `json:"profit"`
	ProfitPercentage  *float64   `json:"profit_percentage"`
	TotalFees         float64    `json:"total_fees"`
	CreatedAt         time.Time  `json:"created_at"`
	CompletedAt       *time.Time `json:"completed_at"`
	ExecutionTimeMs   *int       `json:"execution_time_ms"`
	Details           interface{} `json:"details"`
	ErrorMessage      *string    `json:"error_message"`
}

// Exchange 交易所模型
type Exchange struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Name      string    `json:"name"`
	APIKey    string    `json:"-"` // 不序列化API密钥
	APISecret string    `json:"-"` // 不序列化API密钥
	IsTestnet bool      `json:"is_testnet"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// DashboardStats 仪表板统计数据
type DashboardStats struct {
	TotalProfit float64 `json:"total_profit"`
	TotalTrades int64   `json:"total_trades"`
	ActiveBots  int64   `json:"active_bots"`
	WinRate     float64 `json:"win_rate"`
}

// ChartDataPoint 图表数据点
type ChartDataPoint struct {
	Date   string  `json:"date"`
	Profit float64 `json:"profit"`
	Trades int64   `json:"trades"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

// CreateBotRequest 创建机器人请求
type CreateBotRequest struct {
	Name            string `json:"name" binding:"required"`
	StrategyType    string `json:"strategy_type" binding:"required"`
	ExchangeID      int64  `json:"exchange_id" binding:"required"`
	IsSimulation    bool   `json:"is_simulation"`
	UpdateFrequency int    `json:"update_frequency"`
}

// UpdateBotRequest 更新机器人请求
type UpdateBotRequest struct {
	Name            string `json:"name"`
	StrategyType    string `json:"strategy_type"`
	ExchangeID      int64  `json:"exchange_id"`
	IsSimulation    bool   `json:"is_simulation"`
	UpdateFrequency int    `json:"update_frequency"`
}

// SwitchModeRequest 切换模式请求
type SwitchModeRequest struct {
	IsSimulation bool `json:"is_simulation"`
}

// CreateStrategyRequest 创建策略请求
type CreateStrategyRequest struct {
	Name                string   `json:"name" binding:"required"`
	StrategyType        string   `json:"strategy_type" binding:"required"`
	BasePair            string   `json:"base_pair" binding:"required"`
	QuoteCurrency       string   `json:"quote_currency" binding:"required"`
	Pair1               string   `json:"pair1" binding:"required"`
	Pair2               string   `json:"pair2" binding:"required"`
	Pair3               string   `json:"pair3" binding:"required"`
	Pair4               *string  `json:"pair4"`
	Pair5               *string  `json:"pair5"`
	MinProfitPercentage float64  `json:"min_profit_percentage"`
	MaxTradeAmount      float64  `json:"max_trade_amount"`
	MinTradeAmount      float64  `json:"min_trade_amount"`
	MaxLossPercentage   float64  `json:"max_loss_percentage"`
	MaxConcurrentTrades int      `json:"max_concurrent_trades"`
	UseMargin           bool     `json:"use_margin"`
	Leverage            float64  `json:"leverage"`
}

// CreateExchangeRequest 创建交易所请求
type CreateExchangeRequest struct {
	Name      string `json:"name" binding:"required"`
	APIKey    string `json:"api_key" binding:"required"`
	APISecret string `json:"api_secret" binding:"required"`
	IsTestnet bool   `json:"is_testnet"`
}

// APIResponse 通用API响应
type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// WebSocketMessage WebSocket消息
type WebSocketMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload,omitempty"`
	Token   string      `json:"token,omitempty"`
}

// StatsUpdatePayload 统计更新负载
type StatsUpdatePayload struct {
	TotalProfit float64 `json:"total_profit"`
	TotalTrades int64   `json:"total_trades"`
	ActiveBots  int64   `json:"active_bots"`
	WinRate     float64 `json:"win_rate"`
	Timestamp   time.Time `json:"timestamp"`
}

// TradeUpdatePayload 交易更新负载
type TradeUpdatePayload struct {
	BotID     int64     `json:"bot_id"`
	TradeID   int64     `json:"trade_id"`
	Status    string    `json:"status"`
	Profit    float64   `json:"profit"`
	Timestamp time.Time `json:"timestamp"`
}

// BotStatusPayload 机器人状态负载
type BotStatusPayload struct {
	BotID     int64     `json:"bot_id"`
	IsRunning bool      `json:"is_running"`
	Timestamp time.Time `json:"timestamp"`
}

// LogPayload 日志负载
type LogPayload struct {
	BotID   *int64    `json:"bot_id"`
	Level   string    `json:"level"`
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}

// Claims JWT声明
type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	// 标准声明
	ExpiresAt int64 `json:"exp"`
	IssuedAt  int64 `json:"iat"`
	NotBefore int64 `json:"nbf"`
}

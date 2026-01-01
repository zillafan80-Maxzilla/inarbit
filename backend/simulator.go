package simulator

import (
	"fmt"
	"sync"
	"time"
)

// SimulatedTrade 模拟交易
type SimulatedTrade struct {
	ID              string
	Symbol          string
	Side            string
	Quantity        float64
	Price           float64
	ExecutedQty     float64
	ExecutedPrice   float64
	Status          string
	Commission      float64
	CommissionAsset string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// SimulatedAccount 模拟账户
type SimulatedAccount struct {
	Balances map[string]float64
	mu       sync.RWMutex
}

// SimulatedExchange 模拟交易所
type SimulatedExchange struct {
	account *SimulatedAccount
	trades  map[string]*SimulatedTrade
	prices  map[string]float64
	mu      sync.RWMutex
}

// NewSimulatedExchange 创建新的模拟交易所
func NewSimulatedExchange(initialBalance map[string]float64) *SimulatedExchange {
	balances := make(map[string]float64)
	for asset, amount := range initialBalance {
		balances[asset] = amount
	}

	return &SimulatedExchange{
		account: &SimulatedAccount{
			Balances: balances,
		},
		trades: make(map[string]*SimulatedTrade),
		prices: make(map[string]float64),
	}
}

// SetPrice 设置交易对价格
func (se *SimulatedExchange) SetPrice(symbol string, price float64) {
	se.mu.Lock()
	defer se.mu.Unlock()
	se.prices[symbol] = price
}

// GetPrice 获取交易对价格
func (se *SimulatedExchange) GetPrice(symbol string) float64 {
	se.mu.RLock()
	defer se.mu.RUnlock()
	return se.prices[symbol]
}

// PlaceOrder 下单（模拟）
func (se *SimulatedExchange) PlaceOrder(symbol, side string, quantity, price float64) (*SimulatedTrade, error) {
	se.mu.Lock()
	defer se.mu.Unlock()

	// 解析交易对
	baseAsset, quoteAsset := parseSymbol(symbol)

	// 检查余额
	if side == "BUY" {
		requiredBalance := quantity * price
		if se.account.Balances[quoteAsset] < requiredBalance {
			return nil, fmt.Errorf("余额不足: 需要 %f %s，实际 %f", requiredBalance, quoteAsset, se.account.Balances[quoteAsset])
		}
	} else if side == "SELL" {
		if se.account.Balances[baseAsset] < quantity {
			return nil, fmt.Errorf("余额不足: 需要 %f %s，实际 %f", quantity, baseAsset, se.account.Balances[baseAsset])
		}
	}

	// 创建交易
	trade := &SimulatedTrade{
		ID:        fmt.Sprintf("SIM_%d", time.Now().UnixNano()),
		Symbol:    symbol,
		Side:      side,
		Quantity:  quantity,
		Price:     price,
		Status:    "PENDING",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	se.trades[trade.ID] = trade

	// 立即成交（模拟）
	se.executeOrder(trade, baseAsset, quoteAsset)

	return trade, nil
}

// executeOrder 执行订单
func (se *SimulatedExchange) executeOrder(trade *SimulatedTrade, baseAsset, quoteAsset string) {
	// 计算手续费（0.1%）
	commission := trade.Quantity * trade.Price * 0.001

	if trade.Side == "BUY" {
		// 买入
		se.account.Balances[quoteAsset] -= (trade.Quantity * trade.Price)
		se.account.Balances[baseAsset] += trade.Quantity
		trade.CommissionAsset = baseAsset
		trade.Commission = trade.Quantity * 0.001
	} else if trade.Side == "SELL" {
		// 卖出
		se.account.Balances[baseAsset] -= trade.Quantity
		se.account.Balances[quoteAsset] += (trade.Quantity * trade.Price) - commission
		trade.CommissionAsset = quoteAsset
		trade.Commission = commission
	}

	trade.ExecutedQty = trade.Quantity
	trade.ExecutedPrice = trade.Price
	trade.Status = "FILLED"
	trade.UpdatedAt = time.Now()
}

// GetBalance 获取余额
func (se *SimulatedExchange) GetBalance(asset string) float64 {
	se.mu.RLock()
	defer se.mu.RUnlock()
	return se.account.Balances[asset]
}

// GetAllBalances 获取所有余额
func (se *SimulatedExchange) GetAllBalances() map[string]float64 {
	se.mu.RLock()
	defer se.mu.RUnlock()

	balances := make(map[string]float64)
	for asset, amount := range se.account.Balances {
		balances[asset] = amount
	}
	return balances
}

// GetTrade 获取交易信息
func (se *SimulatedExchange) GetTrade(tradeID string) (*SimulatedTrade, error) {
	se.mu.RLock()
	defer se.mu.RUnlock()

	trade, ok := se.trades[tradeID]
	if !ok {
		return nil, fmt.Errorf("交易不存在: %s", tradeID)
	}

	return trade, nil
}

// GetAllTrades 获取所有交易
func (se *SimulatedExchange) GetAllTrades() []*SimulatedTrade {
	se.mu.RLock()
	defer se.mu.RUnlock()

	trades := make([]*SimulatedTrade, 0, len(se.trades))
	for _, trade := range se.trades {
		trades = append(trades, trade)
	}
	return trades
}

// CalculateProfit 计算利润
func (se *SimulatedExchange) CalculateProfit(initialBalance map[string]float64) map[string]float64 {
	se.mu.RLock()
	defer se.mu.RUnlock()

	profit := make(map[string]float64)
	for asset, initialAmount := range initialBalance {
		currentAmount := se.account.Balances[asset]
		profit[asset] = currentAmount - initialAmount
	}
	return profit
}

// CalculateProfitPercent 计算利润百分比
func (se *SimulatedExchange) CalculateProfitPercent(initialBalance map[string]float64) map[string]float64 {
	se.mu.RLock()
	defer se.mu.RUnlock()

	profitPercent := make(map[string]float64)
	for asset, initialAmount := range initialBalance {
		if initialAmount == 0 {
			continue
		}
		currentAmount := se.account.Balances[asset]
		profitPercent[asset] = ((currentAmount - initialAmount) / initialAmount) * 100
	}
	return profitPercent
}

// Reset 重置模拟账户
func (se *SimulatedExchange) Reset(initialBalance map[string]float64) {
	se.mu.Lock()
	defer se.mu.Unlock()

	se.account.Balances = make(map[string]float64)
	for asset, amount := range initialBalance {
		se.account.Balances[asset] = amount
	}
	se.trades = make(map[string]*SimulatedTrade)
}

// parseSymbol 解析交易对
func parseSymbol(symbol string) (string, string) {
	// 简单的解析逻辑，假设格式为 BASEUSDT 或 BASEQUOTE
	if len(symbol) >= 4 && symbol[len(symbol)-4:] == "USDT" {
		return symbol[:len(symbol)-4], "USDT"
	}
	if len(symbol) >= 3 && symbol[len(symbol)-3:] == "BNB" {
		return symbol[:len(symbol)-3], "BNB"
	}
	// 默认假设最后两个字符是引用资产
	if len(symbol) >= 2 {
		return symbol[:len(symbol)-2], symbol[len(symbol)-2:]
	}
	return symbol, ""
}

// SimulationResult 模拟结果
type SimulationResult struct {
	InitialBalance map[string]float64
	FinalBalance   map[string]float64
	Profit         map[string]float64
	ProfitPercent  map[string]float64
	TotalTrades    int
	SuccessfulTrades int
	FailedTrades   int
	TotalCommission float64
	ExecutionTime  time.Duration
	StartTime      time.Time
	EndTime        time.Time
}

// RunSimulation 运行模拟
func RunSimulation(exchange *SimulatedExchange, initialBalance map[string]float64, trades []struct {
	Symbol   string
	Side     string
	Quantity float64
	Price    float64
}) *SimulationResult {
	startTime := time.Now()

	result := &SimulationResult{
		InitialBalance: initialBalance,
		StartTime:      startTime,
	}

	successCount := 0
	failCount := 0
	totalCommission := 0.0

	for _, t := range trades {
		trade, err := exchange.PlaceOrder(t.Symbol, t.Side, t.Quantity, t.Price)
		if err != nil {
			failCount++
			continue
		}

		successCount++
		totalCommission += trade.Commission
	}

	endTime := time.Now()

	result.FinalBalance = exchange.GetAllBalances()
	result.Profit = exchange.CalculateProfit(initialBalance)
	result.ProfitPercent = exchange.CalculateProfitPercent(initialBalance)
	result.TotalTrades = len(trades)
	result.SuccessfulTrades = successCount
	result.FailedTrades = failCount
	result.TotalCommission = totalCommission
	result.ExecutionTime = endTime.Sub(startTime)
	result.EndTime = endTime

	return result
}

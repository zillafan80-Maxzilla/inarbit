package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// MarketManager 行情管理器
type MarketManager struct {
	client        *BinanceClient
	tickers       map[string]*Ticker     // 交易对行情缓存
	symbolInfo    map[string]*SymbolInfo // 交易对信息缓存
	mu            sync.RWMutex
	updateTicker  time.Duration
	lastUpdate    time.Time
	stopChan      chan struct{}
	updateChan    chan *Ticker
}

// NewMarketManager 创建行情管理器
func NewMarketManager(client *BinanceClient, updateInterval time.Duration) *MarketManager {
	return &MarketManager{
		client:       client,
		tickers:      make(map[string]*Ticker),
		symbolInfo:   make(map[string]*SymbolInfo),
		updateTicker: updateInterval,
		stopChan:     make(chan struct{}),
		updateChan:   make(chan *Ticker, 100),
	}
}

// Start 启动行情管理器
func (m *MarketManager) Start() error {
	// 初始化交易对信息
	if err := m.initSymbolInfo(); err != nil {
		return fmt.Errorf("初始化交易对信息失败: %w", err)
	}

	// 启动定期更新goroutine
	go m.updateLoop()

	log.Println("✓ 行情管理器已启动")
	return nil
}

// Stop 停止行情管理器
func (m *MarketManager) Stop() {
	close(m.stopChan)
	log.Println("✓ 行情管理器已停止")
}

// initSymbolInfo 初始化交易对信息
func (m *MarketManager) initSymbolInfo() error {
	info, err := m.client.GetExchangeInfo()
	if err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	for i, symbol := range info.Symbols {
		if symbol.Status == "TRADING" {
			m.symbolInfo[symbol.Symbol] = &info.Symbols[i]
		}
	}

	log.Printf("✓ 已加载 %d 个交易对信息", len(m.symbolInfo))
	return nil
}

// updateLoop 定期更新行情
func (m *MarketManager) updateLoop() {
	ticker := time.NewTicker(m.updateTicker)
	defer ticker.Stop()

	for {
		select {
		case <-m.stopChan:
			return

		case <-ticker.C:
			if err := m.updateAllTickers(); err != nil {
				log.Printf("更新行情失败: %v", err)
			}
		}
	}
}

// updateAllTickers 更新所有交易对行情
func (m *MarketManager) updateAllTickers() error {
	tickers, err := m.client.GetAllTickers()
	if err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	for _, ticker := range tickers {
		m.tickers[ticker.Symbol] = ticker

		// 发送更新事件
		select {
		case m.updateChan <- ticker:
		default:
			// 通道满，跳过
		}
	}

	m.lastUpdate = time.Now()
	return nil
}

// GetTicker 获取交易对行情
func (m *MarketManager) GetTicker(symbol string) *Ticker {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.tickers[symbol]
}

// GetTickers 获取多个交易对行情
func (m *MarketManager) GetTickers(symbols []string) map[string]*Ticker {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[string]*Ticker)
	for _, symbol := range symbols {
		if ticker, ok := m.tickers[symbol]; ok {
			result[symbol] = ticker
		}
	}
	return result
}

// GetSymbolInfo 获取交易对信息
func (m *MarketManager) GetSymbolInfo(symbol string) *SymbolInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.symbolInfo[symbol]
}

// GetAllSymbols 获取所有交易对
func (m *MarketManager) GetAllSymbols() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	symbols := make([]string, 0, len(m.symbolInfo))
	for symbol := range m.symbolInfo {
		symbols = append(symbols, symbol)
	}
	return symbols
}

// ValidateSymbol 验证交易对是否存在
func (m *MarketManager) ValidateSymbol(symbol string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.symbolInfo[symbol]
	return ok
}

// GetLastUpdate 获取最后更新时间
func (m *MarketManager) GetLastUpdate() time.Time {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.lastUpdate
}

// IsDataFresh 检查数据是否新鲜
func (m *MarketManager) IsDataFresh(maxAge time.Duration) bool {
	return time.Since(m.GetLastUpdate()) < maxAge
}

// ===== 价格查询辅助方法 =====

// GetBidPrice 获取买价
func (m *MarketManager) GetBidPrice(symbol string) float64 {
	ticker := m.GetTicker(symbol)
	if ticker == nil {
		return 0
	}
	return ticker.BidPrice
}

// GetAskPrice 获取卖价
func (m *MarketManager) GetAskPrice(symbol string) float64 {
	ticker := m.GetTicker(symbol)
	if ticker == nil {
		return 0
	}
	return ticker.AskPrice
}

// GetLastPrice 获取最新价格
func (m *MarketManager) GetLastPrice(symbol string) float64 {
	ticker := m.GetTicker(symbol)
	if ticker == nil {
		return 0
	}
	return ticker.LastPrice
}

// GetMidPrice 获取中间价格 (买价+卖价)/2
func (m *MarketManager) GetMidPrice(symbol string) float64 {
	ticker := m.GetTicker(symbol)
	if ticker == nil {
		return 0
	}
	return (ticker.BidPrice + ticker.AskPrice) / 2
}

// GetSpread 获取买卖价差
func (m *MarketManager) GetSpread(symbol string) float64 {
	ticker := m.GetTicker(symbol)
	if ticker == nil {
		return 0
	}
	return ticker.AskPrice - ticker.BidPrice
}

// GetSpreadPercentage 获取买卖价差百分比
func (m *MarketManager) GetSpreadPercentage(symbol string) float64 {
	ticker := m.GetTicker(symbol)
	if ticker == nil {
		return 0
	}
	if ticker.BidPrice == 0 {
		return 0
	}
	return (ticker.AskPrice - ticker.BidPrice) / ticker.BidPrice * 100
}

// ===== 交易对精度处理 =====

// RoundQuantity 四舍五入数量到合法精度
func (m *MarketManager) RoundQuantity(symbol string, quantity float64) (float64, error) {
	symbolInfo := m.GetSymbolInfo(symbol)
	if symbolInfo == nil {
		return 0, fmt.Errorf("交易对 %s 不存在", symbol)
	}

	// 获取步长
	stepSize := 0.00000001
	for _, filter := range symbolInfo.Filters {
		if filter.FilterType == "LOT_SIZE" {
			var err error
			stepSize, err = parseFloat(filter.StepSize)
			if err != nil {
				return 0, err
			}
			break
		}
	}

	// 四舍五入
	return roundToStep(quantity, stepSize), nil
}

// RoundPrice 四舍五入价格到合法精度
func (m *MarketManager) RoundPrice(symbol string, price float64) (float64, error) {
	symbolInfo := m.GetSymbolInfo(symbol)
	if symbolInfo == nil {
		return 0, fmt.Errorf("交易对 %s 不存在", symbol)
	}

	// 获取步长
	tickSize := 0.00000001
	for _, filter := range symbolInfo.Filters {
		if filter.FilterType == "PRICE_FILTER" {
			var err error
			tickSize, err = parseFloat(filter.TickSize)
			if err != nil {
				return 0, err
			}
			break
		}
	}

	// 四舍五入
	return roundToStep(price, tickSize), nil
}

// ===== 辅助函数 =====

// parseFloat 安全地解析float
func parseFloat(s string) (float64, error) {
	if s == "" {
		return 0, fmt.Errorf("空字符串")
	}

	// 使用简单的字符串转换
	var result float64
	_, err := fmt.Sscanf(s, "%f", &result)
	if err != nil {
		return 0, err
	}
	return result, nil
}

// roundToStep 四舍五入到指定步长
func roundToStep(value float64, step float64) float64 {
	if step == 0 {
		return value
	}
	return float64(int64(value/step)) * step
}

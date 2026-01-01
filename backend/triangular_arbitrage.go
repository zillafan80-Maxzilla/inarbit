package arbitrage

import (
	"fmt"
	"log"
	"math"
	"sort"
	"sync"
	"time"

	"inarbit/exchange"
)

// ArbitrageOpportunity 套利机会
type ArbitrageOpportunity struct {
	ID              string
	Path            []string      // 交易路径，例如 ["BTCUSDT", "ETHBTC", "ETHUSDT"]
	StartAsset      string        // 起始资产，例如 "USDT"
	InitialAmount   float64       // 初始金额
	FinalAmount     float64       // 最终金额
	GrossProfit     float64       // 毛利润
	NetProfit       float64       // 净利润（扣除手续费）
	ProfitPercent   float64       // 利润百分比
	ConfidenceScore float64       // 信心度评分 (0-100)
	ExecutionTime   int64         // 预计执行时间（毫秒）
	Timestamp       time.Time     // 发现时间
	IsValid         bool          // 是否有效
	ErrorMessage    string        // 错误信息
}

// TriangularArbitrageEngine 三角套利引擎
type TriangularArbitrageEngine struct {
	client              *exchange.BinanceClient
	tickers             map[string]*exchange.Ticker
	opportunities       []*ArbitrageOpportunity
	minProfitPercent    float64
	takerFeePercent     float64
	slippagePercent     float64
	maxConcurrentTrades int
	mu                  sync.RWMutex
	stopChan            chan bool
	isRunning           bool
}

// NewTriangularArbitrageEngine 创建新的三角套利引擎
func NewTriangularArbitrageEngine(client *exchange.BinanceClient, minProfitPercent float64) *TriangularArbitrageEngine {
	return &TriangularArbitrageEngine{
		client:              client,
		tickers:             make(map[string]*exchange.Ticker),
		opportunities:       make([]*ArbitrageOpportunity, 0),
		minProfitPercent:    minProfitPercent,
		takerFeePercent:     0.1,  // Binance taker费用 0.1%
		slippagePercent:     0.05, // 滑点 0.05%
		maxConcurrentTrades: 5,
		stopChan:            make(chan bool),
		isRunning:           false,
	}
}

// Start 启动套利引擎
func (tae *TriangularArbitrageEngine) Start(updateInterval time.Duration) error {
	tae.mu.Lock()
	if tae.isRunning {
		tae.mu.Unlock()
		return fmt.Errorf("引擎已在运行")
	}
	tae.isRunning = true
	tae.mu.Unlock()

	log.Println("三角套利引擎启动...")

	// 定期更新行情
	go func() {
		ticker := time.NewTicker(updateInterval)
		defer ticker.Stop()

		for {
			select {
			case <-tae.stopChan:
				return
			case <-ticker.C:
				if err := tae.updateTickers(); err != nil {
					log.Printf("更新行情失败: %v", err)
				}

				// 扫描套利机会
				tae.scanOpportunities()
			}
		}
	}()

	return nil
}

// Stop 停止套利引擎
func (tae *TriangularArbitrageEngine) Stop() {
	tae.mu.Lock()
	defer tae.mu.Unlock()

	if tae.isRunning {
		tae.isRunning = false
		close(tae.stopChan)
		log.Println("三角套利引擎已停止")
	}
}

// updateTickers 更新行情数据
func (tae *TriangularArbitrageEngine) updateTickers() error {
	tickers, err := tae.client.GetAllTickers()
	if err != nil {
		return fmt.Errorf("获取行情失败: %v", err)
	}

	tae.mu.Lock()
	defer tae.mu.Unlock()

	for _, ticker := range tickers {
		tae.tickers[ticker.Symbol] = &ticker
	}

	return nil
}

// scanOpportunities 扫描套利机会
func (tae *TriangularArbitrageEngine) scanOpportunities() {
	tae.mu.Lock()
	defer tae.mu.Unlock()

	// 清空旧的机会
	tae.opportunities = make([]*ArbitrageOpportunity, 0)

	// 获取所有交易对
	symbols := make([]string, 0, len(tae.tickers))
	for symbol := range tae.tickers {
		symbols = append(symbols, symbol)
	}

	// 扫描所有可能的三角套利路径
	for i := 0; i < len(symbols); i++ {
		for j := 0; j < len(symbols); j++ {
			for k := 0; k < len(symbols); k++ {
				if i != j && j != k && i != k {
					path := []string{symbols[i], symbols[j], symbols[k]}
					opp := tae.calculateArbitrage(path, 1.0)
					if opp != nil && opp.IsValid && opp.ProfitPercent >= tae.minProfitPercent {
						tae.opportunities = append(tae.opportunities, opp)
					}
				}
			}
		}
	}

	// 按利润排序
	sort.Slice(tae.opportunities, func(i, j int) bool {
		return tae.opportunities[i].NetProfit > tae.opportunities[j].NetProfit
	})
}

// calculateArbitrage 计算套利机会
func (tae *TriangularArbitrageEngine) calculateArbitrage(path []string, initialAmount float64) *ArbitrageOpportunity {
	if len(path) != 3 {
		return nil
	}

	opp := &ArbitrageOpportunity{
		ID:            fmt.Sprintf("ARB_%d", time.Now().UnixNano()),
		Path:          path,
		InitialAmount: initialAmount,
		Timestamp:     time.Now(),
	}

	// 获取交易对信息
	ticker1, ok1 := tae.tickers[path[0]]
	ticker2, ok2 := tae.tickers[path[1]]
	ticker3, ok3 := tae.tickers[path[2]]

	if !ok1 || !ok2 || !ok3 {
		opp.IsValid = false
		opp.ErrorMessage = "缺少行情数据"
		return opp
	}

	if ticker1 == nil || ticker2 == nil || ticker3 == nil {
		opp.IsValid = false
		opp.ErrorMessage = "行情数据为空"
		return opp
	}

	// 解析价格
	price1, err1 := parsePrice(ticker1.AskPrice)
	price2, err2 := parsePrice(ticker2.AskPrice)
	price3, err3 := parsePrice(ticker3.BidPrice)

	if err1 != nil || err2 != nil || err3 != nil {
		opp.IsValid = false
		opp.ErrorMessage = "价格解析失败"
		return opp
	}

	// 计算交易路径
	// 路径1: 用初始资产买入第一个交易对的基础资产
	amount1 := initialAmount / price1 * (1 - tae.takerFeePercent/100)

	// 路径2: 用第一个交易对的基础资产买入第二个交易对的基础资产
	amount2 := amount1 / price2 * (1 - tae.takerFeePercent/100)

	// 路径3: 用第二个交易对的基础资产卖出获得初始资产
	finalAmount := amount2 * price3 * (1 - tae.takerFeePercent/100)

	// 计算利润
	opp.FinalAmount = finalAmount
	opp.GrossProfit = finalAmount - initialAmount
	opp.NetProfit = opp.GrossProfit - (initialAmount * tae.slippagePercent / 100)
	opp.ProfitPercent = (opp.NetProfit / initialAmount) * 100

	// 计算信心度评分
	opp.ConfidenceScore = tae.calculateConfidenceScore(opp)

	// 计算预计执行时间
	opp.ExecutionTime = 3000 // 3秒（保守估计）

	// 验证套利机会
	opp.IsValid = opp.NetProfit > 0 && opp.ProfitPercent >= tae.minProfitPercent

	if !opp.IsValid && opp.ErrorMessage == "" {
		opp.ErrorMessage = fmt.Sprintf("利润不足: %.4f%%", opp.ProfitPercent)
	}

	return opp
}

// calculateConfidenceScore 计算信心度评分
func (tae *TriangularArbitrageEngine) calculateConfidenceScore(opp *ArbitrageOpportunity) float64 {
	score := 50.0 // 基础分数

	// 利润越高，信心度越高
	if opp.ProfitPercent > 0.5 {
		score += 30
	} else if opp.ProfitPercent > 0.2 {
		score += 20
	} else if opp.ProfitPercent > 0.1 {
		score += 10
	}

	// 执行时间越短，信心度越高
	if opp.ExecutionTime < 2000 {
		score += 15
	} else if opp.ExecutionTime < 5000 {
		score += 10
	}

	// 路径越简单，信心度越高
	score += 5

	// 确保分数在0-100之间
	if score > 100 {
		score = 100
	}

	return score
}

// GetTopOpportunities 获取前N个最佳套利机会
func (tae *TriangularArbitrageEngine) GetTopOpportunities(limit int) []*ArbitrageOpportunity {
	tae.mu.RLock()
	defer tae.mu.RUnlock()

	if limit > len(tae.opportunities) {
		limit = len(tae.opportunities)
	}

	return tae.opportunities[:limit]
}

// ExecuteArbitrage 执行套利交易
func (tae *TriangularArbitrageEngine) ExecuteArbitrage(opp *ArbitrageOpportunity, initialAmount float64) (*ArbitrageExecutionResult, error) {
	if !opp.IsValid {
		return nil, fmt.Errorf("套利机会无效: %s", opp.ErrorMessage)
	}

	result := &ArbitrageExecutionResult{
		OpportunityID: opp.ID,
		Path:          opp.Path,
		StartTime:     time.Now(),
		Status:        "PENDING",
	}

	// 第一步：买入第一个交易对的基础资产
	ticker1, ok := tae.tickers[opp.Path[0]]
	if !ok {
		result.Status = "FAILED"
		result.ErrorMessage = "缺少第一个交易对的行情"
		return result, fmt.Errorf("缺少行情数据")
	}

	price1, _ := parsePrice(ticker1.AskPrice)
	quantity1 := initialAmount / price1

	log.Printf("第一步: 买入 %s, 数量: %.8f, 价格: %.2f", opp.Path[0], quantity1, price1)

	order1, err := tae.client.PlaceOrder(opp.Path[0], "BUY", quantity1, price1)
	if err != nil {
		result.Status = "FAILED"
		result.ErrorMessage = fmt.Sprintf("第一步下单失败: %v", err)
		return result, err
	}

	result.Orders = append(result.Orders, order1)

	// 第二步：买入第二个交易对的基础资产
	ticker2, ok := tae.tickers[opp.Path[1]]
	if !ok {
		result.Status = "FAILED"
		result.ErrorMessage = "缺少第二个交易对的行情"
		return result, fmt.Errorf("缺少行情数据")
	}

	price2, _ := parsePrice(ticker2.AskPrice)
	quantity2 := quantity1 / price2

	log.Printf("第二步: 买入 %s, 数量: %.8f, 价格: %.2f", opp.Path[1], quantity2, price2)

	order2, err := tae.client.PlaceOrder(opp.Path[1], "BUY", quantity2, price2)
	if err != nil {
		// 撤销第一个订单
		tae.client.CancelOrder(opp.Path[0], order1.OrderID)
		result.Status = "FAILED"
		result.ErrorMessage = fmt.Sprintf("第二步下单失败: %v", err)
		return result, err
	}

	result.Orders = append(result.Orders, order2)

	// 第三步：卖出获得初始资产
	ticker3, ok := tae.tickers[opp.Path[2]]
	if !ok {
		result.Status = "FAILED"
		result.ErrorMessage = "缺少第三个交易对的行情"
		return result, fmt.Errorf("缺少行情数据")
	}

	price3, _ := parsePrice(ticker3.BidPrice)

	log.Printf("第三步: 卖出 %s, 数量: %.8f, 价格: %.2f", opp.Path[2], quantity2, price3)

	order3, err := tae.client.PlaceOrder(opp.Path[2], "SELL", quantity2, price3)
	if err != nil {
		// 撤销前两个订单
		tae.client.CancelOrder(opp.Path[0], order1.OrderID)
		tae.client.CancelOrder(opp.Path[1], order2.OrderID)
		result.Status = "FAILED"
		result.ErrorMessage = fmt.Sprintf("第三步下单失败: %v", err)
		return result, err
	}

	result.Orders = append(result.Orders, order3)

	// 计算实际利润
	finalAmount, _ := parsePrice(order3.CummulativeQuoteQty)
	result.FinalAmount = finalAmount
	result.Profit = finalAmount - initialAmount
	result.ProfitPercent = (result.Profit / initialAmount) * 100
	result.Status = "SUCCESS"
	result.EndTime = time.Now()

	log.Printf("套利执行完成: 初始 %.2f, 最终 %.2f, 利润 %.4f%%", initialAmount, finalAmount, result.ProfitPercent)

	return result, nil
}

// ArbitrageExecutionResult 套利执行结果
type ArbitrageExecutionResult struct {
	OpportunityID  string
	Path           []string
	Orders         []*exchange.Order
	InitialAmount  float64
	FinalAmount    float64
	Profit         float64
	ProfitPercent  float64
	Status         string
	ErrorMessage   string
	StartTime      time.Time
	EndTime        time.Time
	ExecutionTime  time.Duration
}

// GetExecutionTime 获取执行时间
func (aer *ArbitrageExecutionResult) GetExecutionTime() time.Duration {
	return aer.EndTime.Sub(aer.StartTime)
}

// TriangleValidator 三角形验证器
type TriangleValidator struct {
	baseAssets  []string
	quoteAssets []string
}

// NewTriangleValidator 创建新的三角形验证器
func NewTriangleValidator() *TriangleValidator {
	return &TriangleValidator{
		baseAssets:  []string{"BTC", "ETH", "BNB", "ADA", "XRP"},
		quoteAssets: []string{"USDT", "BUSD", "USDC"},
	}
}

// IsValidTriangle 验证是否为有效的三角形
func (tv *TriangleValidator) IsValidTriangle(path []string) bool {
	if len(path) != 3 {
		return false
	}

	// 简单的验证逻辑
	// 可以扩展为更复杂的验证

	return true
}

// QuadrangularArbitrageEngine 四角套利引擎
type QuadrangularArbitrageEngine struct {
	client           *exchange.BinanceClient
	tickers          map[string]*exchange.Ticker
	opportunities    []*QuadrangularOpportunity
	minProfitPercent float64
	mu               sync.RWMutex
}

// QuadrangularOpportunity 四角套利机会
type QuadrangularOpportunity struct {
	ID            string
	Path          []string  // 4个交易对
	InitialAmount float64
	FinalAmount   float64
	NetProfit     float64
	ProfitPercent float64
	Timestamp     time.Time
	IsValid       bool
}

// NewQuadrangularArbitrageEngine 创建新的四角套利引擎
func NewQuadrangularArbitrageEngine(client *exchange.BinanceClient, minProfitPercent float64) *QuadrangularArbitrageEngine {
	return &QuadrangularArbitrageEngine{
		client:           client,
		tickers:          make(map[string]*exchange.Ticker),
		opportunities:    make([]*QuadrangularOpportunity, 0),
		minProfitPercent: minProfitPercent,
	}
}

// PentagonalArbitrageEngine 五角套利引擎
type PentagonalArbitrageEngine struct {
	client           *exchange.BinanceClient
	tickers          map[string]*exchange.Ticker
	opportunities    []*PentagonalOpportunity
	minProfitPercent float64
	mu               sync.RWMutex
}

// PentagonalOpportunity 五角套利机会
type PentagonalOpportunity struct {
	ID            string
	Path          []string  // 5个交易对
	InitialAmount float64
	FinalAmount   float64
	NetProfit     float64
	ProfitPercent float64
	Timestamp     time.Time
	IsValid       bool
}

// NewPentagonalArbitrageEngine 创建新的五角套利引擎
func NewPentagonalArbitrageEngine(client *exchange.BinanceClient, minProfitPercent float64) *PentagonalArbitrageEngine {
	return &PentagonalArbitrageEngine{
		client:           client,
		tickers:          make(map[string]*exchange.Ticker),
		opportunities:    make([]*PentagonalOpportunity, 0),
		minProfitPercent: minProfitPercent,
	}
}

// parsePrice 解析价格字符串
func parsePrice(priceStr string) (float64, error) {
	var price float64
	_, err := fmt.Sscanf(priceStr, "%f", &price)
	return price, err
}

// RiskAssessment 风险评估
type RiskAssessment struct {
	ExecutionRisk    float64 // 执行风险 (0-100)
	LiquidityRisk    float64 // 流动性风险 (0-100)
	SlippageRisk     float64 // 滑点风险 (0-100)
	OverallRisk      float64 // 总体风险 (0-100)
	Recommendation   string  // 建议
}

// AssessRisk 评估风险
func (tae *TriangularArbitrageEngine) AssessRisk(opp *ArbitrageOpportunity) *RiskAssessment {
	assessment := &RiskAssessment{}

	// 执行风险：基于利润百分比
	if opp.ProfitPercent > 0.5 {
		assessment.ExecutionRisk = 20
	} else if opp.ProfitPercent > 0.2 {
		assessment.ExecutionRisk = 40
	} else {
		assessment.ExecutionRisk = 60
	}

	// 流动性风险：基于交易对
	assessment.LiquidityRisk = 30 // 假设Binance主流交易对流动性良好

	// 滑点风险
	assessment.SlippageRisk = tae.slippagePercent * 10

	// 总体风险
	assessment.OverallRisk = (assessment.ExecutionRisk + assessment.LiquidityRisk + assessment.SlippageRisk) / 3

	// 建议
	if assessment.OverallRisk < 30 {
		assessment.Recommendation = "低风险，可以执行"
	} else if assessment.OverallRisk < 60 {
		assessment.Recommendation = "中等风险，谨慎执行"
	} else {
		assessment.Recommendation = "高风险，不建议执行"
	}

	return assessment
}

// Statistics 统计信息
type Statistics struct {
	TotalOpportunities   int
	ValidOpportunities   int
	AverageProfitPercent float64
	MaxProfitPercent     float64
	MinProfitPercent     float64
	LastUpdateTime       time.Time
}

// GetStatistics 获取统计信息
func (tae *TriangularArbitrageEngine) GetStatistics() *Statistics {
	tae.mu.RLock()
	defer tae.mu.RUnlock()

	stats := &Statistics{
		TotalOpportunities: len(tae.opportunities),
		LastUpdateTime:     time.Now(),
	}

	if len(tae.opportunities) == 0 {
		return stats
	}

	totalProfit := 0.0
	maxProfit := tae.opportunities[0].ProfitPercent
	minProfit := tae.opportunities[0].ProfitPercent

	for _, opp := range tae.opportunities {
		if opp.IsValid {
			stats.ValidOpportunities++
		}
		totalProfit += opp.ProfitPercent
		if opp.ProfitPercent > maxProfit {
			maxProfit = opp.ProfitPercent
		}
		if opp.ProfitPercent < minProfit {
			minProfit = opp.ProfitPercent
		}
	}

	stats.AverageProfitPercent = totalProfit / float64(len(tae.opportunities))
	stats.MaxProfitPercent = maxProfit
	stats.MinProfitPercent = minProfit

	return stats
}

// OptimizePath 优化交易路径
func (tae *TriangularArbitrageEngine) OptimizePath(opp *ArbitrageOpportunity) *ArbitrageOpportunity {
	// 这里可以实现路径优化算法
	// 例如：选择流动性最好的交易对组合

	return opp
}

// CalculateSlippage 计算滑点
func (tae *TriangularArbitrageEngine) CalculateSlippage(quantity float64, price float64) float64 {
	// 根据订单大小计算滑点
	// 大订单会有更大的滑点

	baseSlippage := tae.slippagePercent
	orderSizeMultiplier := math.Log(quantity*price/1000 + 1)

	return baseSlippage * orderSizeMultiplier
}

// MonitorOpportunity 监控套利机会
func (tae *TriangularArbitrageEngine) MonitorOpportunity(opp *ArbitrageOpportunity, duration time.Duration) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	endTime := time.Now().Add(duration)

	for {
		select {
		case <-ticker.C:
			if time.Now().After(endTime) {
				return
			}

			// 重新计算套利机会
			updatedOpp := tae.calculateArbitrage(opp.Path, opp.InitialAmount)
			if updatedOpp != nil && updatedOpp.ProfitPercent > opp.ProfitPercent {
				log.Printf("套利机会改善: %.4f%% -> %.4f%%", opp.ProfitPercent, updatedOpp.ProfitPercent)
			}
		}
	}
}

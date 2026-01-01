package main

import (
	"fmt"
	"log"
	"math"
	"sync"
	"time"
)

// ArbitrageOpportunity 套利机会
type ArbitrageOpportunity struct {
	ID                 string
	Type               string    // triangular, quadrangular, pentagonal
	Pair1              string
	Pair2              string
	Pair3              string
	Pair4              *string
	Pair5              *string
	Path               []string  // 交易路径
	InitialAmount      float64
	FinalAmount        float64
	GrossProfit        float64   // 毛利润
	NetProfit          float64   // 净利润
	ProfitPercentage   float64   // 利润百分比
	ExecutionTime      int       // 预计执行时间（毫秒）
	Confidence         float64   // 信心度 (0-100)
	Timestamp          time.Time
	Details            *ArbitrageDetails
}

// ArbitrageDetails 套利详情
type ArbitrageDetails struct {
	Step1 *TradeStep
	Step2 *TradeStep
	Step3 *TradeStep
	Step4 *TradeStep
	Step5 *TradeStep
	TotalFees float64
	Slippage  float64
}

// TradeStep 交易步骤
type TradeStep struct {
	Symbol       string
	Side         string  // BUY or SELL
	Price        float64
	Quantity     float64
	Amount       float64
	Fee          float64
	FeePercentage float64
}

// ArbitrageEngine 套利引擎
type ArbitrageEngine struct {
	marketManager    *MarketManager
	minProfitPercent float64
	takerFeePercent  float64
	makerFeePercent  float64
	mu               sync.RWMutex
	opportunities    []*ArbitrageOpportunity
	stopChan         chan struct{}
}

// NewArbitrageEngine 创建套利引擎
func NewArbitrageEngine(marketManager *MarketManager, minProfitPercent float64) *ArbitrageEngine {
	return &ArbitrageEngine{
		marketManager:    marketManager,
		minProfitPercent: minProfitPercent,
		takerFeePercent:  0.001, // 0.1%
		makerFeePercent:  0.001, // 0.1%
		opportunities:    make([]*ArbitrageOpportunity, 0),
		stopChan:         make(chan struct{}),
	}
}

// Start 启动套利引擎
func (e *ArbitrageEngine) Start() {
	go e.scanLoop()
	log.Println("✓ 套利引擎已启动")
}

// Stop 停止套利引擎
func (e *ArbitrageEngine) Stop() {
	close(e.stopChan)
	log.Println("✓ 套利引擎已停止")
}

// scanLoop 扫描套利机会
func (e *ArbitrageEngine) scanLoop() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-e.stopChan:
			return

		case <-ticker.C:
			e.scanTriangularArbitrages()
		}
	}
}

// ===== 三角套利扫描 =====

// scanTriangularArbitrages 扫描三角套利机会
func (e *ArbitrageEngine) scanTriangularArbitrages() {
	symbols := e.marketManager.GetAllSymbols()
	
	// 按基础货币分组
	baseCurrencies := make(map[string][]string)
	for _, symbol := range symbols {
		if len(symbol) >= 6 {
			// 假设交易对格式为 XXXYYY（6个字符）
			// 这是一个简化的实现，实际应该从交易所信息获取
			baseCurrency := symbol[3:6]
			baseCurrencies[baseCurrency] = append(baseCurrencies[baseCurrency], symbol)
		}
	}

	// 对每个基础货币扫描三角套利
	for baseCurrency, pairs := range baseCurrencies {
		if len(pairs) >= 3 {
			e.findTriangularOpportunities(baseCurrency, pairs)
		}
	}
}

// findTriangularOpportunities 查找三角套利机会
func (e *ArbitrageEngine) findTriangularOpportunities(baseCurrency string, pairs []string) {
	// 这是一个简化的实现
	// 实际应该根据交易对的基础货币和报价货币来构建交易路径
	
	if len(pairs) < 3 {
		return
	}

	// 示例：对于USDT基础货币，查找如下路径：
	// 1. USDT -> BTC (买入BTC)
	// 2. BTC -> ETH (卖出BTC买入ETH)
	// 3. ETH -> USDT (卖出ETH买入USDT)

	// 这需要更复杂的图论算法来找到所有可能的路径
	// 这里只是演示基本逻辑
}

// CalculateTriangularArbitrage 计算三角套利
func (e *ArbitrageEngine) CalculateTriangularArbitrage(
	pair1, pair2, pair3 string,
	initialAmount float64,
) *ArbitrageOpportunity {
	
	// 获取行情
	ticker1 := e.marketManager.GetTicker(pair1)
	ticker2 := e.marketManager.GetTicker(pair2)
	ticker3 := e.marketManager.GetTicker(pair3)

	if ticker1 == nil || ticker2 == nil || ticker3 == nil {
		return nil
	}

	// 验证交易对组合
	if !e.validatePairCombination(pair1, pair2, pair3) {
		return nil
	}

	// 计算交易步骤
	step1 := &TradeStep{
		Symbol:        pair1,
		Side:          "BUY",
		Price:         ticker1.AskPrice,
		Quantity:      initialAmount / ticker1.AskPrice,
		Amount:        initialAmount,
		FeePercentage: e.takerFeePercent,
	}
	step1.Fee = step1.Amount * step1.FeePercentage
	step1.Amount -= step1.Fee

	step2 := &TradeStep{
		Symbol:        pair2,
		Side:          "SELL",
		Price:         ticker2.BidPrice,
		Quantity:      step1.Quantity,
		Amount:        step1.Quantity * ticker2.BidPrice,
		FeePercentage: e.takerFeePercent,
	}
	step2.Fee = step2.Amount * step2.FeePercentage
	step2.Amount -= step2.Fee

	step3 := &TradeStep{
		Symbol:        pair3,
		Side:          "SELL",
		Price:         ticker3.BidPrice,
		Quantity:      step2.Amount / ticker3.BidPrice,
		Amount:        step2.Amount,
		FeePercentage: e.takerFeePercent,
	}
	step3.Fee = step3.Amount * step3.FeePercentage
	step3.Amount -= step3.Fee

	// 计算利润
	finalAmount := step3.Amount
	grossProfit := finalAmount - initialAmount
	totalFees := step1.Fee + step2.Fee + step3.Fee
	netProfit := grossProfit - totalFees
	profitPercentage := (netProfit / initialAmount) * 100

	// 检查是否值得执行
	if profitPercentage < e.minProfitPercent {
		return nil
	}

	opportunity := &ArbitrageOpportunity{
		ID:               generateOpportunityID(),
		Type:             "triangular",
		Pair1:            pair1,
		Pair2:            pair2,
		Pair3:            pair3,
		Path:             []string{pair1, pair2, pair3},
		InitialAmount:    initialAmount,
		FinalAmount:      finalAmount,
		GrossProfit:      grossProfit,
		NetProfit:        netProfit,
		ProfitPercentage: profitPercentage,
		ExecutionTime:    3000, // 预计3秒
		Confidence:       calculateConfidence(profitPercentage),
		Timestamp:        time.Now(),
		Details: &ArbitrageDetails{
			Step1:     step1,
			Step2:     step2,
			Step3:     step3,
			TotalFees: totalFees,
			Slippage:  0, // 实际执行时会有滑点
		},
	}

	return opportunity
}

// validatePairCombination 验证交易对组合
func (e *ArbitrageEngine) validatePairCombination(pair1, pair2, pair3 string) bool {
	// 验证交易对是否存在
	if !e.marketManager.ValidateSymbol(pair1) ||
		!e.marketManager.ValidateSymbol(pair2) ||
		!e.marketManager.ValidateSymbol(pair3) {
		return false
	}

	// 验证交易对是否能形成闭合回路
	// 这需要检查基础货币和报价货币的匹配
	// 简化实现：假设所有交易对都有效
	return true
}

// AddOpportunity 添加套利机会
func (e *ArbitrageEngine) AddOpportunity(opp *ArbitrageOpportunity) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.opportunities = append(e.opportunities, opp)

	// 只保留最近的1000个机会
	if len(e.opportunities) > 1000 {
		e.opportunities = e.opportunities[len(e.opportunities)-1000:]
	}
}

// GetOpportunities 获取套利机会
func (e *ArbitrageEngine) GetOpportunities() []*ArbitrageOpportunity {
	e.mu.RLock()
	defer e.mu.RUnlock()

	result := make([]*ArbitrageOpportunity, len(e.opportunities))
	copy(result, e.opportunities)
	return result
}

// GetBestOpportunity 获取最佳套利机会
func (e *ArbitrageEngine) GetBestOpportunity() *ArbitrageOpportunity {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if len(e.opportunities) == 0 {
		return nil
	}

	best := e.opportunities[0]
	for _, opp := range e.opportunities {
		if opp.ProfitPercentage > best.ProfitPercentage {
			best = opp
		}
	}

	return best
}

// ClearOpportunities 清除套利机会
func (e *ArbitrageEngine) ClearOpportunities() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.opportunities = make([]*ArbitrageOpportunity, 0)
}

// ===== 辅助函数 =====

// generateOpportunityID 生成机会ID
func generateOpportunityID() string {
	return fmt.Sprintf("opp_%d", time.Now().UnixNano())
}

// calculateConfidence 计算信心度
func calculateConfidence(profitPercentage float64) float64 {
	// 利润越高，信心度越高
	// 0.1% 利润 -> 50% 信心度
	// 0.5% 利润 -> 80% 信心度
	// 1.0% 利润 -> 95% 信心度
	
	confidence := math.Min(50+profitPercentage*100, 100)
	return math.Max(confidence, 0)
}

// ===== 风险评估 =====

// AssessRisk 评估风险
func (e *ArbitrageEngine) AssessRisk(opp *ArbitrageOpportunity) *RiskAssessment {
	assessment := &RiskAssessment{
		OpportunityID: opp.ID,
		Timestamp:     time.Now(),
	}

	// 评估滑点风险
	assessment.SlippageRisk = e.assessSlippageRisk(opp)

	// 评估流动性风险
	assessment.LiquidityRisk = e.assessLiquidityRisk(opp)

	// 评估执行风险
	assessment.ExecutionRisk = e.assessExecutionRisk(opp)

	// 计算总体风险评分
	assessment.OverallRisk = (assessment.SlippageRisk + assessment.LiquidityRisk + assessment.ExecutionRisk) / 3

	// 计算风险调整后的利润
	assessment.RiskAdjustedProfit = opp.NetProfit * (1 - assessment.OverallRisk/100)

	return assessment
}

// assessSlippageRisk 评估滑点风险
func (e *ArbitrageEngine) assessSlippageRisk(opp *ArbitrageOpportunity) float64 {
	// 基于买卖价差评估滑点风险
	// 价差越大，滑点风险越高
	
	risk := 0.0
	for _, pair := range opp.Path {
		spread := e.marketManager.GetSpreadPercentage(pair)
		risk += spread * 10 // 权重
	}

	return math.Min(risk/float64(len(opp.Path)), 100)
}

// assessLiquidityRisk 评估流动性风险
func (e *ArbitrageEngine) assessLiquidityRisk(opp *ArbitrageOpportunity) float64 {
	// 基于交易量评估流动性风险
	// 交易量越小，流动性风险越高
	
	risk := 0.0
	for _, pair := range opp.Path {
		ticker := e.marketManager.GetTicker(pair)
		if ticker != nil {
			// 如果交易量很小，风险高
			if ticker.Volume < 100 {
				risk += 50
			} else if ticker.Volume < 1000 {
				risk += 30
			} else if ticker.Volume < 10000 {
				risk += 10
			}
		}
	}

	return math.Min(risk/float64(len(opp.Path)), 100)
}

// assessExecutionRisk 评估执行风险
func (e *ArbitrageEngine) assessExecutionRisk(opp *ArbitrageOpportunity) float64 {
	// 基于执行时间评估执行风险
	// 执行时间越长，风险越高
	
	// 每1秒增加5%风险
	risk := float64(opp.ExecutionTime/1000) * 5
	return math.Min(risk, 100)
}

// RiskAssessment 风险评估结果
type RiskAssessment struct {
	OpportunityID      string
	SlippageRisk       float64
	LiquidityRisk      float64
	ExecutionRisk      float64
	OverallRisk        float64
	RiskAdjustedProfit float64
	Timestamp          time.Time
}

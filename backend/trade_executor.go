package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// TradeExecution 交易执行记录
type TradeExecution struct {
	ID                  string
	BotID               int64
	OpportunityID       string
	Status              string // pending, executing, completed, failed, cancelled
	Type                string // triangular, quadrangular, pentagonal
	Path                []string
	InitialAmount       float64
	FinalAmount         float64
	ActualProfit        float64
	ActualProfitPercent float64
	TotalFees           float64
	Slippage            float64
	Orders              []*ExecutedOrder
	StartTime           time.Time
	EndTime             time.Time
	ExecutionTime       int64 // 毫秒
	ErrorMessage        string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

// ExecutedOrder 已执行的订单
type ExecutedOrder struct {
	OrderID       int64
	Symbol        string
	Side          string
	Type          string
	Price         float64
	Quantity      float64
	ExecutedQty   float64
	CummulativeQty float64
	Status        string
	Fee           float64
	FeeAsset      string
	ExecutedAt    time.Time
}

// TradeExecutor 交易执行器
type TradeExecutor struct {
	client              *BinanceClient
	marketManager       *MarketManager
	db                  *Database
	maxConcurrentTrades int
	executingTrades     map[string]*TradeExecution
	mu                  sync.RWMutex
	stopChan            chan struct{}
}

// NewTradeExecutor 创建交易执行器
func NewTradeExecutor(client *BinanceClient, marketManager *MarketManager, db *Database) *TradeExecutor {
	return &TradeExecutor{
		client:              client,
		marketManager:       marketManager,
		db:                  db,
		maxConcurrentTrades: 5,
		executingTrades:     make(map[string]*TradeExecution),
		stopChan:            make(chan struct{}),
	}
}

// ExecuteArbitrage 执行套利交易
func (e *TradeExecutor) ExecuteArbitrage(botID int64, opp *ArbitrageOpportunity, isSimulation bool) (*TradeExecution, error) {
	// 检查并发限制
	e.mu.Lock()
	if len(e.executingTrades) >= e.maxConcurrentTrades {
		e.mu.Unlock()
		return nil, fmt.Errorf("并发交易数已达上限 (%d)", e.maxConcurrentTrades)
	}
	e.mu.Unlock()

	execution := &TradeExecution{
		ID:            generateTradeID(),
		BotID:         botID,
		OpportunityID: opp.ID,
		Status:        "pending",
		Type:          opp.Type,
		Path:          opp.Path,
		InitialAmount: opp.InitialAmount,
		Orders:        make([]*ExecutedOrder, 0),
		StartTime:     time.Now(),
		CreatedAt:     time.Now(),
	}

	// 添加到执行中的交易列表
	e.mu.Lock()
	e.executingTrades[execution.ID] = execution
	e.mu.Unlock()

	// 执行交易
	if isSimulation {
		go e.executeSimulation(execution, opp)
	} else {
		go e.executeReal(execution, opp)
	}

	return execution, nil
}

// executeSimulation 执行模拟交易
func (e *TradeExecutor) executeSimulation(execution *TradeExecution, opp *ArbitrageOpportunity) {
	execution.Status = "executing"

	// 模拟交易延迟
	time.Sleep(time.Duration(opp.ExecutionTime) * time.Millisecond)

	// 模拟滑点
	slippage := opp.NetProfit * 0.1 // 假设滑点为利润的10%

	execution.FinalAmount = opp.FinalAmount - slippage
	execution.ActualProfit = execution.FinalAmount - execution.InitialAmount
	execution.ActualProfitPercent = (execution.ActualProfit / execution.InitialAmount) * 100
	execution.TotalFees = opp.Details.TotalFees
	execution.Slippage = slippage
	execution.EndTime = time.Now()
	execution.ExecutionTime = execution.EndTime.Sub(execution.StartTime).Milliseconds()
	execution.Status = "completed"
	execution.UpdatedAt = time.Now()

	// 记录交易
	e.recordExecution(execution)

	log.Printf("✓ 模拟交易完成: %s, 利润: %.2f (%+.2f%%)", execution.ID, execution.ActualProfit, execution.ActualProfitPercent)

	// 清除执行记录
	e.mu.Lock()
	delete(e.executingTrades, execution.ID)
	e.mu.Unlock()
}

// executeReal 执行真实交易
func (e *TradeExecutor) executeReal(execution *TradeExecution, opp *ArbitrageOpportunity) {
	execution.Status = "executing"

	// 执行第一步：买入第一个交易对
	order1, err := e.executeStep(execution, opp.Details.Step1, 1)
	if err != nil {
		execution.Status = "failed"
		execution.ErrorMessage = fmt.Sprintf("第一步失败: %v", err)
		e.recordExecution(execution)
		e.mu.Lock()
		delete(e.executingTrades, execution.ID)
		e.mu.Unlock()
		log.Printf("✗ 交易失败: %s, 错误: %s", execution.ID, execution.ErrorMessage)
		return
	}
	execution.Orders = append(execution.Orders, order1)

	// 等待订单成交
	if !e.waitForOrder(opp.Path[0], order1.OrderID, 30*time.Second) {
		execution.Status = "failed"
		execution.ErrorMessage = "第一步订单超时"
		e.recordExecution(execution)
		e.mu.Lock()
		delete(e.executingTrades, execution.ID)
		e.mu.Unlock()
		return
	}

	// 执行第二步：卖出第一个交易对，买入第二个
	order2, err := e.executeStep(execution, opp.Details.Step2, 2)
	if err != nil {
		execution.Status = "failed"
		execution.ErrorMessage = fmt.Sprintf("第二步失败: %v", err)
		e.recordExecution(execution)
		e.mu.Lock()
		delete(e.executingTrades, execution.ID)
		e.mu.Unlock()
		return
	}
	execution.Orders = append(execution.Orders, order2)

	// 等待订单成交
	if !e.waitForOrder(opp.Path[1], order2.OrderID, 30*time.Second) {
		execution.Status = "failed"
		execution.ErrorMessage = "第二步订单超时"
		e.recordExecution(execution)
		e.mu.Lock()
		delete(e.executingTrades, execution.ID)
		e.mu.Unlock()
		return
	}

	// 执行第三步：卖出第二个交易对
	order3, err := e.executeStep(execution, opp.Details.Step3, 3)
	if err != nil {
		execution.Status = "failed"
		execution.ErrorMessage = fmt.Sprintf("第三步失败: %v", err)
		e.recordExecution(execution)
		e.mu.Lock()
		delete(e.executingTrades, execution.ID)
		e.mu.Unlock()
		return
	}
	execution.Orders = append(execution.Orders, order3)

	// 等待订单成交
	if !e.waitForOrder(opp.Path[2], order3.OrderID, 30*time.Second) {
		execution.Status = "failed"
		execution.ErrorMessage = "第三步订单超时"
		e.recordExecution(execution)
		e.mu.Lock()
		delete(e.executingTrades, execution.ID)
		e.mu.Unlock()
		return
	}

	// 计算实际结果
	execution.FinalAmount = order3.CummulativeQty
	execution.ActualProfit = execution.FinalAmount - execution.InitialAmount
	execution.ActualProfitPercent = (execution.ActualProfit / execution.InitialAmount) * 100
	execution.TotalFees = order1.Fee + order2.Fee + order3.Fee
	execution.EndTime = time.Now()
	execution.ExecutionTime = execution.EndTime.Sub(execution.StartTime).Milliseconds()
	execution.Status = "completed"
	execution.UpdatedAt = time.Now()

	// 记录交易
	e.recordExecution(execution)

	log.Printf("✓ 交易完成: %s, 利润: %.2f (%+.2f%%)", execution.ID, execution.ActualProfit, execution.ActualProfitPercent)

	// 清除执行记录
	e.mu.Lock()
	delete(e.executingTrades, execution.ID)
	e.mu.Unlock()
}

// executeStep 执行交易步骤
func (e *TradeExecutor) executeStep(execution *TradeExecution, step *TradeStep, stepNum int) (*ExecutedOrder, error) {
	log.Printf("执行第%d步: %s %s %.8f @ %.8f", stepNum, step.Side, step.Symbol, step.Quantity, step.Price)

	var order *Order
	var err error

	if step.Side == "BUY" {
		order, err = e.client.PlaceOrder(step.Symbol, "BUY", step.Quantity, step.Price)
	} else {
		order, err = e.client.PlaceOrder(step.Symbol, "SELL", step.Quantity, step.Price)
	}

	if err != nil {
		return nil, fmt.Errorf("下单失败: %w", err)
	}

	executedOrder := &ExecutedOrder{
		OrderID:       order.OrderID,
		Symbol:        order.Symbol,
		Side:          order.Side,
		Type:          order.Type,
		Price:         order.Price,
		Quantity:      order.OrigQty,
		ExecutedQty:   order.ExecutedQty,
		CummulativeQty: order.CummulativeQuoteQty,
		Status:        order.Status,
		ExecutedAt:    time.Now(),
	}

	return executedOrder, nil
}

// waitForOrder 等待订单成交
func (e *TradeExecutor) waitForOrder(symbol string, orderID int64, timeout time.Duration) bool {
	startTime := time.Now()

	for {
		if time.Since(startTime) > timeout {
			return false
		}

		order, err := e.client.GetOrder(symbol, orderID)
		if err != nil {
			log.Printf("查询订单失败: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}

		if order.Status == "FILLED" || order.Status == "PARTIALLY_FILLED" {
			return true
		}

		if order.Status == "CANCELED" || order.Status == "REJECTED" {
			return false
		}

		time.Sleep(500 * time.Millisecond)
	}
}

// recordExecution 记录交易执行
func (e *TradeExecutor) recordExecution(execution *TradeExecution) {
	// 保存到数据库
	// TODO: 实现数据库保存逻辑
	
	log.Printf("记录交易: %s, 状态: %s, 利润: %.2f", execution.ID, execution.Status, execution.ActualProfit)
}

// GetExecution 获取交易执行记录
func (e *TradeExecutor) GetExecution(executionID string) *TradeExecution {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.executingTrades[executionID]
}

// GetExecutingTrades 获取正在执行的交易
func (e *TradeExecutor) GetExecutingTrades() []*TradeExecution {
	e.mu.RLock()
	defer e.mu.RUnlock()

	result := make([]*TradeExecution, 0, len(e.executingTrades))
	for _, trade := range e.executingTrades {
		result = append(result, trade)
	}
	return result
}

// CancelExecution 取消交易执行
func (e *TradeExecutor) CancelExecution(executionID string) error {
	e.mu.Lock()
	execution, ok := e.executingTrades[executionID]
	e.mu.Unlock()

	if !ok {
		return fmt.Errorf("交易不存在")
	}

	if execution.Status != "executing" && execution.Status != "pending" {
		return fmt.Errorf("无法取消已完成的交易")
	}

	// 取消所有未成交的订单
	for _, order := range execution.Orders {
		if order.Status != "FILLED" {
			_, err := e.client.CancelOrder(order.Symbol, order.OrderID)
			if err != nil {
				log.Printf("取消订单失败: %v", err)
			}
		}
	}

	execution.Status = "cancelled"
	execution.UpdatedAt = time.Now()

	return nil
}

// ===== 辅助函数 =====

// generateTradeID 生成交易ID
func generateTradeID() string {
	return fmt.Sprintf("trade_%d", time.Now().UnixNano())
}

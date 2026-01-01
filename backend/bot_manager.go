package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// BotManager 机器人管理器
type BotManager struct {
	db                *Database
	binanceClient     *BinanceClient
	marketManager     *MarketManager
	arbitrageEngine   *ArbitrageEngine
	tradeExecutor     *TradeExecutor
	activeBots        map[int64]*BotInstance
	mu                sync.RWMutex
	stopChan          chan struct{}
	wsManager         *WebSocketManager
}

// BotInstance 机器人实例
type BotInstance struct {
	Bot                *Bot
	IsRunning          bool
	MarketManager      *MarketManager
	ArbitrageEngine    *ArbitrageEngine
	TradeExecutor      *TradeExecutor
	LastOpportunity    *ArbitrageOpportunity
	LastExecution      *TradeExecution
	Statistics         *BotStatistics
	UpdateFrequency    time.Duration
	stopChan           chan struct{}
	mu                 sync.RWMutex
}

// BotStatistics 机器人统计信息
type BotStatistics struct {
	TotalTrades        int64
	SuccessfulTrades   int64
	FailedTrades       int64
	TotalProfit        float64
	TotalLoss          float64
	WinRate            float64
	AverageProfitPerTrade float64
	BestTrade          float64
	WorstTrade         float64
	TotalFees          float64
	StartTime          time.Time
	UpdatedAt          time.Time
}

// NewBotManager 创建机器人管理器
func NewBotManager(
	db *Database,
	binanceClient *BinanceClient,
	marketManager *MarketManager,
	arbitrageEngine *ArbitrageEngine,
	tradeExecutor *TradeExecutor,
	wsManager *WebSocketManager,
) *BotManager {
	return &BotManager{
		db:              db,
		binanceClient:   binanceClient,
		marketManager:   marketManager,
		arbitrageEngine: arbitrageEngine,
		tradeExecutor:   tradeExecutor,
		activeBots:      make(map[int64]*BotInstance),
		stopChan:        make(chan struct{}),
		wsManager:       wsManager,
	}
}

// Start 启动机器人管理器
func (bm *BotManager) Start() error {
	log.Println("✓ 机器人管理器已启动")
	return nil
}

// Stop 停止机器人管理器
func (bm *BotManager) Stop() {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	// 停止所有活跃的机器人
	for _, botInstance := range bm.activeBots {
		botInstance.Stop()
	}

	close(bm.stopChan)
	log.Println("✓ 机器人管理器已停止")
}

// StartBot 启动机器人
func (bm *BotManager) StartBot(botID int64) error {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	// 检查机器人是否已在运行
	if _, ok := bm.activeBots[botID]; ok {
		return fmt.Errorf("机器人已在运行")
	}

	// 从数据库获取机器人信息
	bot, err := bm.db.GetBotByID(botID, 0) // TODO: 获取正确的userID
	if err != nil {
		return fmt.Errorf("获取机器人信息失败: %w", err)
	}

	// 创建机器人实例
	botInstance := &BotInstance{
		Bot:               bot,
		IsRunning:         true,
		MarketManager:     bm.marketManager,
		ArbitrageEngine:   bm.arbitrageEngine,
		TradeExecutor:     bm.tradeExecutor,
		UpdateFrequency:   time.Duration(bot.UpdateFrequency) * time.Second,
		stopChan:          make(chan struct{}),
		Statistics: &BotStatistics{
			StartTime: time.Now(),
		},
	}

	// 启动机器人
	go botInstance.Run()

	// 添加到活跃机器人列表
	bm.activeBots[botID] = botInstance

	log.Printf("✓ 机器人 %d 已启动", botID)
	return nil
}

// StopBot 停止机器人
func (bm *BotManager) StopBot(botID int64) error {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	botInstance, ok := bm.activeBots[botID]
	if !ok {
		return fmt.Errorf("机器人不存在或未运行")
	}

	botInstance.Stop()
	delete(bm.activeBots, botID)

	log.Printf("✓ 机器人 %d 已停止", botID)
	return nil
}

// GetBotInstance 获取机器人实例
func (bm *BotManager) GetBotInstance(botID int64) *BotInstance {
	bm.mu.RLock()
	defer bm.mu.RUnlock()
	return bm.activeBots[botID]
}

// GetActiveBots 获取所有活跃机器人
func (bm *BotManager) GetActiveBots() []*BotInstance {
	bm.mu.RLock()
	defer bm.mu.RUnlock()

	result := make([]*BotInstance, 0, len(bm.activeBots))
	for _, bot := range bm.activeBots {
		result = append(result, bot)
	}
	return result
}

// ===== BotInstance 方法 =====

// Run 运行机器人
func (bi *BotInstance) Run() {
	log.Printf("机器人 %d 开始运行", bi.Bot.ID)

	ticker := time.NewTicker(bi.UpdateFrequency)
	defer ticker.Stop()

	for {
		select {
		case <-bi.stopChan:
			return

		case <-ticker.C:
			bi.scan()
		}
	}
}

// Stop 停止机器人
func (bi *BotInstance) Stop() {
	bi.mu.Lock()
	defer bi.mu.Unlock()

	bi.IsRunning = false
	close(bi.stopChan)
	log.Printf("机器人 %d 已停止", bi.Bot.ID)
}

// scan 扫描套利机会
func (bi *BotInstance) scan() {
	bi.mu.Lock()
	defer bi.mu.Unlock()

	// 检查行情数据是否新鲜
	if !bi.MarketManager.IsDataFresh(10 * time.Second) {
		log.Printf("机器人 %d: 行情数据过旧，跳过本轮扫描", bi.Bot.ID)
		return
	}

	// 扫描套利机会
	// TODO: 实现具体的套利机会扫描逻辑
	// 这里是一个示例：扫描三角套利

	// 获取最佳套利机会
	bestOpp := bi.ArbitrageEngine.GetBestOpportunity()
	if bestOpp == nil {
		return
	}

	// 评估风险
	riskAssessment := bi.ArbitrageEngine.AssessRisk(bestOpp)
	if riskAssessment.OverallRisk > 50 {
		log.Printf("机器人 %d: 风险过高 (%.2f%%), 跳过", bi.Bot.ID, riskAssessment.OverallRisk)
		return
	}

	// 执行交易
	execution, err := bi.TradeExecutor.ExecuteArbitrage(bi.Bot.ID, bestOpp, bi.Bot.IsSimulation)
	if err != nil {
		log.Printf("机器人 %d: 执行交易失败: %v", bi.Bot.ID, err)
		return
	}

	bi.LastOpportunity = bestOpp
	bi.LastExecution = execution

	// 更新统计信息
	bi.updateStatistics(execution)

	log.Printf("机器人 %d: 执行交易 %s, 利润: %.2f", bi.Bot.ID, execution.ID, execution.ActualProfit)
}

// updateStatistics 更新统计信息
func (bi *BotInstance) updateStatistics(execution *TradeExecution) {
	stats := bi.Statistics

	stats.TotalTrades++
	stats.UpdatedAt = time.Now()

	if execution.Status == "completed" {
		if execution.ActualProfit > 0 {
			stats.SuccessfulTrades++
			stats.TotalProfit += execution.ActualProfit
		} else {
			stats.FailedTrades++
			stats.TotalLoss += -execution.ActualProfit
		}

		// 更新最好和最坏的交易
		if execution.ActualProfit > stats.BestTrade {
			stats.BestTrade = execution.ActualProfit
		}
		if execution.ActualProfit < stats.WorstTrade {
			stats.WorstTrade = execution.ActualProfit
		}

		// 计算胜率
		if stats.TotalTrades > 0 {
			stats.WinRate = float64(stats.SuccessfulTrades) / float64(stats.TotalTrades) * 100
		}

		// 计算平均利润
		if stats.SuccessfulTrades > 0 {
			stats.AverageProfitPerTrade = stats.TotalProfit / float64(stats.SuccessfulTrades)
		}
	}

	stats.TotalFees += execution.TotalFees
}

// GetStatistics 获取统计信息
func (bi *BotInstance) GetStatistics() *BotStatistics {
	bi.mu.RLock()
	defer bi.mu.RUnlock()

	return bi.Statistics
}

// SetUpdateFrequency 设置更新频率
func (bi *BotInstance) SetUpdateFrequency(frequency time.Duration) {
	bi.mu.Lock()
	defer bi.mu.Unlock()

	bi.UpdateFrequency = frequency
	log.Printf("机器人 %d: 更新频率已设置为 %v", bi.Bot.ID, frequency)
}

// SwitchMode 切换模式（虚拟/实盘）
func (bi *BotInstance) SwitchMode(isSimulation bool) {
	bi.mu.Lock()
	defer bi.mu.Unlock()

	bi.Bot.IsSimulation = isSimulation
	mode := "实盘"
	if isSimulation {
		mode = "虚拟盘"
	}
	log.Printf("机器人 %d: 已切换为%s", bi.Bot.ID, mode)
}

// GetStatus 获取机器人状态
func (bi *BotInstance) GetStatus() map[string]interface{} {
	bi.mu.RLock()
	defer bi.mu.RUnlock()

	return map[string]interface{}{
		"bot_id":              bi.Bot.ID,
		"is_running":          bi.IsRunning,
		"is_simulation":       bi.Bot.IsSimulation,
		"update_frequency":    bi.UpdateFrequency.String(),
		"last_opportunity":    bi.LastOpportunity,
		"last_execution":      bi.LastExecution,
		"statistics":          bi.Statistics,
	}
}

// ===== 机器人协调器 =====

// BotCoordinator 机器人协调器（用于多个机器人之间的协调）
type BotCoordinator struct {
	botManager *BotManager
	mu         sync.RWMutex
}

// NewBotCoordinator 创建机器人协调器
func NewBotCoordinator(botManager *BotManager) *BotCoordinator {
	return &BotCoordinator{
		botManager: botManager,
	}
}

// CheckConflicts 检查机器人之间的冲突
func (bc *BotCoordinator) CheckConflicts() []string {
	activeBots := bc.botManager.GetActiveBots()
	conflicts := make([]string, 0)

	// 检查是否有两个机器人在同时执行相同的交易对
	for i, bot1 := range activeBots {
		for j, bot2 := range activeBots {
			if i >= j {
				continue
			}

			if bot1.LastExecution != nil && bot2.LastExecution != nil {
				// 检查是否有重叠的交易对
				for _, pair1 := range bot1.LastExecution.Path {
					for _, pair2 := range bot2.LastExecution.Path {
						if pair1 == pair2 {
							conflicts = append(conflicts, fmt.Sprintf(
								"机器人 %d 和 %d 在交易对 %s 上有冲突",
								bot1.Bot.ID, bot2.Bot.ID, pair1,
							))
						}
					}
				}
			}
		}
	}

	return conflicts
}

// BalanceLoad 平衡负载
func (bc *BotCoordinator) BalanceLoad() {
	activeBots := bc.botManager.GetActiveBots()

	// 如果有机器人失败，重新启动它们
	for _, bot := range activeBots {
		if !bot.IsRunning {
			log.Printf("重新启动失败的机器人 %d", bot.Bot.ID)
			bc.botManager.StartBot(bot.Bot.ID)
		}
	}
}

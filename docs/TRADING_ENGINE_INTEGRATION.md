# iNarbit 交易引擎集成指南

## 概述

这份文档说明如何将所有交易引擎模块集成到主应用中。

## 模块关系图

```
┌─────────────────────────────────────────────────────────────┐
│                      BotManager (机器人管理)                  │
│  ┌──────────────────────────────────────────────────────┐   │
│  │ BotInstance (机器人实例)                              │   │
│  │  ├─ MarketManager (行情管理)                         │   │
│  │  ├─ ArbitrageEngine (套利引擎)                       │   │
│  │  ├─ TradeExecutor (交易执行)                         │   │
│  │  └─ BotStatistics (统计信息)                         │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
         ↓
┌─────────────────────────────────────────────────────────────┐
│                    MarketManager (行情管理)                   │
│  ├─ BinanceClient (Binance API客户端)                       │
│  ├─ Ticker Cache (行情缓存)                                 │
│  └─ SymbolInfo Cache (交易对信息缓存)                       │
└─────────────────────────────────────────────────────────────┘
         ↓
┌─────────────────────────────────────────────────────────────┐
│                 ArbitrageEngine (套利引擎)                    │
│  ├─ 三角套利计算                                             │
│  ├─ 四角套利计算                                             │
│  ├─ 五角套利计算                                             │
│  ├─ 风险评估                                                 │
│  └─ 机会管理                                                 │
└─────────────────────────────────────────────────────────────┘
         ↓
┌─────────────────────────────────────────────────────────────┐
│                  TradeExecutor (交易执行)                     │
│  ├─ 订单下单                                                 │
│  ├─ 订单监控                                                 │
│  ├─ 交易记录                                                 │
│  └─ 并发控制                                                 │
└─────────────────────────────────────────────────────────────┘
```

## 集成步骤

### 1. 在主程序中初始化交易引擎

```go
// 在 main.go 中添加以下代码

func main() {
    // ... 现有的初始化代码 ...

    // 初始化Binance客户端
    exchange, err := db.GetExchangeByID(1, 1) // 获取用户的第一个交易所
    if err != nil {
        log.Fatalf("获取交易所信息失败: %v", err)
    }

    binanceClient := NewBinanceClient(exchange.APIKey, exchange.APISecret, exchange.IsTestnet)

    // 初始化行情管理器
    marketManager := NewMarketManager(binanceClient, 5*time.Second)
    if err := marketManager.Start(); err != nil {
        log.Fatalf("启动行情管理器失败: %v", err)
    }
    defer marketManager.Stop()

    // 初始化套利引擎
    arbitrageEngine := NewArbitrageEngine(marketManager, 0.1) // 最小利润0.1%
    arbitrageEngine.Start()
    defer arbitrageEngine.Stop()

    // 初始化交易执行器
    tradeExecutor := NewTradeExecutor(binanceClient, marketManager, db)

    // 初始化机器人管理器
    botManager := NewBotManager(
        db,
        binanceClient,
        marketManager,
        arbitrageEngine,
        tradeExecutor,
        wsManager,
    )
    if err := botManager.Start(); err != nil {
        log.Fatalf("启动机器人管理器失败: %v", err)
    }
    defer botManager.Stop()

    // 在API处理器中注册机器人管理器
    apiHandler.botManager = botManager

    // ... 启动HTTP服务器 ...
}
```

### 2. 在API处理器中添加机器人控制端点

```go
// 在 handlers/api.go 中添加以下方法

// StartBot 启动机器人
func (h *APIHandler) StartBot(w http.ResponseWriter, r *http.Request) {
    userID, err := h.GetUserID(r)
    if err != nil {
        h.RespondError(w, http.StatusUnauthorized, "获取用户ID失败")
        return
    }

    vars := mux.Vars(r)
    botID, err := strconv.ParseInt(vars["id"], 10, 64)
    if err != nil {
        h.RespondError(w, http.StatusBadRequest, "无效的机器人ID")
        return
    }

    // 启动机器人
    if err := h.botManager.StartBot(botID); err != nil {
        h.RespondError(w, http.StatusInternalServerError, err.Error())
        return
    }

    h.RespondSuccess(w, http.StatusOK, "机器人已启动", nil)
}

// StopBot 停止机器人
func (h *APIHandler) StopBot(w http.ResponseWriter, r *http.Request) {
    userID, err := h.GetUserID(r)
    if err != nil {
        h.RespondError(w, http.StatusUnauthorized, "获取用户ID失败")
        return
    }

    vars := mux.Vars(r)
    botID, err := strconv.ParseInt(vars["id"], 10, 64)
    if err != nil {
        h.RespondError(w, http.StatusBadRequest, "无效的机器人ID")
        return
    }

    // 停止机器人
    if err := h.botManager.StopBot(botID); err != nil {
        h.RespondError(w, http.StatusInternalServerError, err.Error())
        return
    }

    h.RespondSuccess(w, http.StatusOK, "机器人已停止", nil)
}

// GetBotStatus 获取机器人状态
func (h *APIHandler) GetBotStatus(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    botID, err := strconv.ParseInt(vars["id"], 10, 64)
    if err != nil {
        h.RespondError(w, http.StatusBadRequest, "无效的机器人ID")
        return
    }

    botInstance := h.botManager.GetBotInstance(botID)
    if botInstance == nil {
        h.RespondError(w, http.StatusNotFound, "机器人未运行")
        return
    }

    h.RespondSuccess(w, http.StatusOK, "获取机器人状态成功", botInstance.GetStatus())
}
```

### 3. 在WebSocket中推送实时数据

```go
// 在 websocket/manager.go 中添加以下方法

// BroadcastBotStatus 广播机器人状态
func (m *WebSocketManager) BroadcastBotStatus(userID int64, botID int64, isRunning bool) {
    message := WebSocketMessage{
        Type: "bot_status",
        Payload: map[string]interface{}{
            "bot_id":    botID,
            "is_running": isRunning,
            "timestamp": time.Now(),
        },
    }

    m.broadcastToUser(userID, message)
}

// BroadcastTradeUpdate 广播交易更新
func (m *WebSocketManager) BroadcastTradeUpdate(userID int64, execution *TradeExecution) {
    message := WebSocketMessage{
        Type: "trade_update",
        Payload: map[string]interface{}{
            "execution_id":  execution.ID,
            "bot_id":        execution.BotID,
            "status":        execution.Status,
            "profit":        execution.ActualProfit,
            "profit_percent": execution.ActualProfitPercent,
            "timestamp":     time.Now(),
        },
    }

    m.broadcastToUser(userID, message)
}

// BroadcastBotStatistics 广播机器人统计信息
func (m *WebSocketManager) BroadcastBotStatistics(userID int64, botID int64, stats *BotStatistics) {
    message := WebSocketMessage{
        Type: "bot_statistics",
        Payload: map[string]interface{}{
            "bot_id":              botID,
            "total_trades":        stats.TotalTrades,
            "successful_trades":   stats.SuccessfulTrades,
            "failed_trades":       stats.FailedTrades,
            "total_profit":        stats.TotalProfit,
            "win_rate":            stats.WinRate,
            "average_profit":      stats.AverageProfitPerTrade,
            "best_trade":          stats.BestTrade,
            "worst_trade":         stats.WorstTrade,
            "total_fees":          stats.TotalFees,
            "timestamp":           time.Now(),
        },
    }

    m.broadcastToUser(userID, message)
}
```

## 配置参数

### MarketManager 配置

```go
// 更新频率：1秒、5秒、10秒等
marketManager := NewMarketManager(binanceClient, 5*time.Second)
```

### ArbitrageEngine 配置

```go
// 最小利润百分比：0.1%、0.5%、1.0%等
arbitrageEngine := NewArbitrageEngine(marketManager, 0.1)

// 手续费配置（在ArbitrageEngine中修改）
arbitrageEngine.takerFeePercent = 0.001  // 0.1%
arbitrageEngine.makerFeePercent = 0.001  // 0.1%
```

### BotManager 配置

```go
// 最大并发交易数
tradeExecutor.maxConcurrentTrades = 5

// 订单超时时间（在TradeExecutor中修改）
timeout := 30 * time.Second
```

## 性能优化建议

### 1. 行情更新频率

- **低频（10秒）**：降低API调用，适合长期持仓
- **中频（5秒）**：平衡性能和延迟，推荐使用
- **高频（1秒）**：高实时性，但增加API调用和CPU使用

### 2. 套利机会扫描

- **全量扫描**：扫描所有交易对，准确但耗时
- **增量扫描**：只扫描行情变化的交易对，更高效
- **分批扫描**：分批处理交易对，避免阻塞

### 3. 并发控制

```go
// 限制并发交易数
tradeExecutor.maxConcurrentTrades = 5

// 使用信号量控制并发
semaphore := make(chan struct{}, 5)
```

### 4. 缓存策略

```go
// 缓存行情数据
marketManager.tickers // 内存缓存

// 缓存交易对信息
marketManager.symbolInfo // 内存缓存

// 定期更新缓存
marketManager.updateAllTickers()
```

## 监控和日志

### 关键指标

```
- 行情更新延迟
- 套利机会发现数
- 交易执行成功率
- 平均利润
- 总手续费
- 机器人运行时间
```

### 日志级别

```
INFO: 机器人启动/停止、交易执行
WARN: 风险过高、行情延迟
ERROR: 交易失败、连接错误
DEBUG: 详细的计算过程
```

## 故障恢复

### 自动重启

```go
// 在BotCoordinator中实现
func (bc *BotCoordinator) BalanceLoad() {
    activeBots := bc.botManager.GetActiveBots()
    for _, bot := range activeBots {
        if !bot.IsRunning {
            bc.botManager.StartBot(bot.Bot.ID)
        }
    }
}
```

### 订单恢复

```go
// 在TradeExecutor中实现
func (e *TradeExecutor) RecoverPendingOrders() {
    // 查询所有未成交订单
    openOrders, err := e.client.GetOpenOrders("")
    if err != nil {
        log.Printf("查询未成交订单失败: %v", err)
        return
    }

    // 处理未成交订单
    for _, order := range openOrders {
        // 检查订单是否超时
        // 如果超时，取消订单
    }
}
```

## 测试

### 单元测试

```bash
go test ./... -v
```

### 集成测试

```bash
# 使用Binance测试网络
export BINANCE_TESTNET=true
go test ./... -v
```

### 性能测试

```bash
go test -bench=. -benchmem
```

## 部署检查清单

- [ ] 配置Binance API密钥
- [ ] 设置最小利润阈值
- [ ] 配置行情更新频率
- [ ] 设置最大并发交易数
- [ ] 配置WebSocket推送频率
- [ ] 测试虚拟盘交易
- [ ] 测试实盘交易（小额）
- [ ] 配置监控和告警
- [ ] 备份数据库
- [ ] 准备应急预案

## 常见问题

### Q1：如何添加新的套利策略？

在 `arbitrage_engine.go` 中添加新的计算方法：

```go
// 四角套利
func (e *ArbitrageEngine) CalculateQuadrangularArbitrage(...) *ArbitrageOpportunity {
    // 实现逻辑
}
```

### Q2：如何修改手续费？

在 `ArbitrageEngine` 中修改：

```go
arbitrageEngine.takerFeePercent = 0.002  // 0.2%
arbitrageEngine.makerFeePercent = 0.001  // 0.1%
```

### Q3：如何实现止损？

在 `TradeExecutor` 中添加：

```go
if execution.ActualProfit < -maxLoss {
    // 立即平仓
    e.CancelExecution(execution.ID)
}
```

### Q4：如何处理API限流？

实现限流器：

```go
type RateLimiter struct {
    tokens int
    limit  int
    ticker *time.Ticker
}
```

## 下一步

1. **完善套利算法**：支持更多交易对组合
2. **添加机器学习**：预测最佳交易时机
3. **实现风险管理**：动态调整头寸
4. **添加通知系统**：邮件/短信告警
5. **性能优化**：使用缓存和异步处理
6. **可视化仪表板**：实时展示交易数据

---

**最后更新**：2024年1月
**版本**：1.0.0

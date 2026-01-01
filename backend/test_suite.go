package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"inarbit/exchange"
	"inarbit/simulator"
)

// TestResult 测试结果
type TestResult struct {
	TestName  string
	Status    string
	Message   string
	Duration  time.Duration
	Timestamp time.Time
}

// TestSuite 测试套件
type TestSuite struct {
	results []TestResult
}

// AddResult 添加测试结果
func (ts *TestSuite) AddResult(testName, status, message string, duration time.Duration) {
	ts.results = append(ts.results, TestResult{
		TestName:  testName,
		Status:    status,
		Message:   message,
		Duration:  duration,
		Timestamp: time.Now(),
	})
}

// PrintResults 打印测试结果
func (ts *TestSuite) PrintResults() {
	fmt.Println("\n" + "="*80)
	fmt.Println("测试结果总结")
	fmt.Println("="*80)

	passCount := 0
	failCount := 0

	for _, result := range ts.results {
		status := "✓"
		if result.Status == "FAIL" {
			status = "✗"
			failCount++
		} else {
			passCount++
		}

		fmt.Printf("%s [%-30s] %s (%v)\n", status, result.TestName, result.Message, result.Duration)
	}

	fmt.Println("="*80)
	fmt.Printf("总计: %d 通过, %d 失败\n", passCount, failCount)
	fmt.Println("="*80 + "\n")
}

// Test1_BinanceConnection 测试1: Binance连接
func Test1_BinanceConnection(ts *TestSuite, apiKey, apiSecret string) {
	start := time.Now()
	testName := "Binance连接测试"

	client := exchange.NewBinanceClient(apiKey, apiSecret, false)
	
	err := client.TestConnection()
	duration := time.Since(start)

	if err != nil {
		ts.AddResult(testName, "FAIL", fmt.Sprintf("连接失败: %v", err), duration)
		return
	}

	ts.AddResult(testName, "PASS", "成功连接到Binance API", duration)
}

// Test2_GetServerTime 测试2: 获取服务器时间
func Test2_GetServerTime(ts *TestSuite, apiKey, apiSecret string) {
	start := time.Now()
	testName := "获取服务器时间"

	client := exchange.NewBinanceClient(apiKey, apiSecret, false)
	
	serverTime, err := client.GetServerTime()
	duration := time.Since(start)

	if err != nil {
		ts.AddResult(testName, "FAIL", fmt.Sprintf("获取失败: %v", err), duration)
		return
	}

	ts.AddResult(testName, "PASS", fmt.Sprintf("服务器时间: %d", serverTime), duration)
}

// Test3_GetExchangeInfo 测试3: 获取交易所信息
func Test3_GetExchangeInfo(ts *TestSuite, apiKey, apiSecret string) {
	start := time.Now()
	testName := "获取交易所信息"

	client := exchange.NewBinanceClient(apiKey, apiSecret, false)
	
	info, err := client.GetExchangeInfo()
	duration := time.Since(start)

	if err != nil {
		ts.AddResult(testName, "FAIL", fmt.Sprintf("获取失败: %v", err), duration)
		return
	}

	ts.AddResult(testName, "PASS", fmt.Sprintf("获取到 %d 个交易对", len(info.Symbols)), duration)
}

// Test4_GetTickers 测试4: 获取行情数据
func Test4_GetTickers(ts *TestSuite, apiKey, apiSecret string) {
	start := time.Now()
	testName := "获取行情数据"

	client := exchange.NewBinanceClient(apiKey, apiSecret, false)
	
	tickers, err := client.GetAllTickers()
	duration := time.Since(start)

	if err != nil {
		ts.AddResult(testName, "FAIL", fmt.Sprintf("获取失败: %v", err), duration)
		return
	}

	ts.AddResult(testName, "PASS", fmt.Sprintf("获取到 %d 个行情数据", len(tickers)), duration)
}

// Test5_GetSpecificTicker 测试5: 获取特定交易对行情
func Test5_GetSpecificTicker(ts *TestSuite, apiKey, apiSecret string) {
	start := time.Now()
	testName := "获取特定交易对行情"

	client := exchange.NewBinanceClient(apiKey, apiSecret, false)
	
	ticker, err := client.GetTicker("BTCUSDT")
	duration := time.Since(start)

	if err != nil {
		ts.AddResult(testName, "FAIL", fmt.Sprintf("获取失败: %v", err), duration)
		return
	}

	ts.AddResult(testName, "PASS", fmt.Sprintf("BTCUSDT 价格: %s", ticker.LastPrice), duration)
}

// Test6_GetAccount 测试6: 获取账户信息
func Test6_GetAccount(ts *TestSuite, apiKey, apiSecret string) {
	start := time.Now()
	testName := "获取账户信息"

	client := exchange.NewBinanceClient(apiKey, apiSecret, false)
	
	account, err := client.GetAccount()
	duration := time.Since(start)

	if err != nil {
		ts.AddResult(testName, "FAIL", fmt.Sprintf("获取失败: %v", err), duration)
		return
	}

	// 计算总余额
	totalAssets := 0
	for _, balance := range account.Balances {
		free, _ := parseFloat(balance.Free)
		locked, _ := parseFloat(balance.Locked)
		if free > 0 || locked > 0 {
			totalAssets++
		}
	}

	ts.AddResult(testName, "PASS", fmt.Sprintf("账户有 %d 个资产", totalAssets), duration)
}

// Test7_VirtualTradingSimulation 测试7: 虚拟盘交易模拟
func Test7_VirtualTradingSimulation(ts *TestSuite) {
	start := time.Now()
	testName := "虚拟盘交易模拟"

	// 创建模拟交易所
	initialBalance := map[string]float64{
		"USDT": 1000.0,
		"BTC":  0.0,
		"ETH":  0.0,
	}

	exchange := simulator.NewSimulatedExchange(initialBalance)

	// 设置价格
	exchange.SetPrice("BTCUSDT", 45000.0)
	exchange.SetPrice("ETHUSDT", 3000.0)

	// 执行模拟交易
	trades := []struct {
		Symbol   string
		Side     string
		Quantity float64
		Price    float64
	}{
		{"BTCUSDT", "BUY", 0.01, 45000.0},   // 买入0.01 BTC
		{"ETHUSDT", "BUY", 0.1, 3000.0},     // 买入0.1 ETH
		{"BTCUSDT", "SELL", 0.01, 45500.0},  // 卖出0.01 BTC
	}

	result := simulator.RunSimulation(exchange, initialBalance, trades)
	duration := time.Since(start)

	if result.FailedTrades > 0 {
		ts.AddResult(testName, "FAIL", fmt.Sprintf("失败交易数: %d", result.FailedTrades), duration)
		return
	}

	profitPercent := result.ProfitPercent["USDT"]
	ts.AddResult(testName, "PASS", fmt.Sprintf("完成 %d 笔交易，利润: %.2f%%", result.TotalTrades, profitPercent), duration)
}

// Test8_TriangularArbitrageSimulation 测试8: 三角套利模拟
func Test8_TriangularArbitrageSimulation(ts *TestSuite) {
	start := time.Now()
	testName := "三角套利模拟"

	// 创建模拟交易所
	initialBalance := map[string]float64{
		"USDT": 1000.0,
		"BTC":  0.0,
		"ETH":  0.0,
	}

	exchange := simulator.NewSimulatedExchange(initialBalance)

	// 设置价格（制造套利机会）
	exchange.SetPrice("BTCUSDT", 45000.0)
	exchange.SetPrice("ETHUSDT", 3000.0)
	exchange.SetPrice("ETHBTC", 0.0667)  // 正常应该是 3000/45000 = 0.0667

	// 三角套利交易路径: USDT -> BTC -> ETH -> USDT
	trades := []struct {
		Symbol   string
		Side     string
		Quantity float64
		Price    float64
	}{
		{"BTCUSDT", "BUY", 1.0/45000, 45000.0},   // 用1 USDT买入BTC
		{"ETHBTC", "BUY", 1.0/45000/0.0667, 0.0667},  // 用BTC买入ETH
		{"ETHUSDT", "SELL", 1.0/45000/0.0667, 3000.0},  // 用ETH卖出USDT
	}

	result := simulator.RunSimulation(exchange, initialBalance, trades)
	duration := time.Since(start)

	if result.FailedTrades > 0 {
		ts.AddResult(testName, "FAIL", fmt.Sprintf("失败交易数: %d", result.FailedTrades), duration)
		return
	}

	profitPercent := result.ProfitPercent["USDT"]
	ts.AddResult(testName, "PASS", fmt.Sprintf("三角套利完成，利润: %.4f%%", profitPercent), duration)
}

// Test9_PerformanceTest 测试9: 性能测试
func Test9_PerformanceTest(ts *TestSuite) {
	start := time.Now()
	testName := "性能测试（1000笔交易）"

	initialBalance := map[string]float64{
		"USDT": 100000.0,
		"BTC":  0.0,
	}

	exchange := simulator.NewSimulatedExchange(initialBalance)
	exchange.SetPrice("BTCUSDT", 45000.0)

	// 生成1000笔交易
	trades := make([]struct {
		Symbol   string
		Side     string
		Quantity float64
		Price    float64
	}, 1000)

	for i := 0; i < 1000; i++ {
		if i%2 == 0 {
			trades[i] = struct {
				Symbol   string
				Side     string
				Quantity float64
				Price    float64
			}{"BTCUSDT", "BUY", 0.001, 45000.0}
		} else {
			trades[i] = struct {
				Symbol   string
				Side     string
				Quantity float64
				Price    float64
			}{"BTCUSDT", "SELL", 0.001, 45100.0}
		}
	}

	result := simulator.RunSimulation(exchange, initialBalance, trades)
	duration := time.Since(start)

	ts.AddResult(testName, "PASS", fmt.Sprintf("执行1000笔交易耗时 %v", result.ExecutionTime), duration)
}

// Test10_DataValidation 测试10: 数据验证
func Test10_DataValidation(ts *TestSuite, apiKey, apiSecret string) {
	start := time.Now()
	testName := "数据验证"

	client := exchange.NewBinanceClient(apiKey, apiSecret, false)

	// 获取多个行情数据
	symbols := []string{"BTCUSDT", "ETHUSDT", "BNBUSDT"}
	validCount := 0

	for _, symbol := range symbols {
		ticker, err := client.GetTicker(symbol)
		if err == nil && ticker.LastPrice != "" {
			validCount++
		}
	}

	duration := time.Since(start)

	if validCount != len(symbols) {
		ts.AddResult(testName, "FAIL", fmt.Sprintf("验证失败: %d/%d", validCount, len(symbols)), duration)
		return
	}

	ts.AddResult(testName, "PASS", fmt.Sprintf("验证成功: %d/%d", validCount, len(symbols)), duration)
}

// parseFloat 解析浮点数
func parseFloat(s string) (float64, error) {
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	return f, err
}

// main 主函数
func main() {
	// 从环境变量读取API密钥
	apiKey := os.Getenv("BINANCE_API_KEY")
	apiSecret := os.Getenv("BINANCE_API_SECRET")

	if apiKey == "" || apiSecret == "" {
		log.Fatal("请设置 BINANCE_API_KEY 和 BINANCE_API_SECRET 环境变量")
	}

	fmt.Println("="*80)
	fmt.Println("iNarbit 完整测试套件")
	fmt.Println("="*80)

	ts := &TestSuite{}

	// 执行测试
	fmt.Println("\n[虚拟盘测试]")
	Test7_VirtualTradingSimulation(ts)
	Test8_TriangularArbitrageSimulation(ts)
	Test9_PerformanceTest(ts)

	fmt.Println("\n[实盘连接测试]")
	Test1_BinanceConnection(ts, apiKey, apiSecret)
	Test2_GetServerTime(ts, apiKey, apiSecret)
	Test3_GetExchangeInfo(ts, apiKey, apiSecret)

	fmt.Println("\n[行情数据测试]")
	Test4_GetTickers(ts, apiKey, apiSecret)
	Test5_GetSpecificTicker(ts, apiKey, apiSecret)
	Test10_DataValidation(ts, apiKey, apiSecret)

	fmt.Println("\n[账户测试]")
	Test6_GetAccount(ts, apiKey, apiSecret)

	// 打印结果
	ts.PrintResults()

	// 生成测试报告
	generateTestReport(ts)
}

// generateTestReport 生成测试报告
func generateTestReport(ts *TestSuite) {
	reportFile := "test_report.txt"
	file, err := os.Create(reportFile)
	if err != nil {
		log.Printf("无法创建报告文件: %v", err)
		return
	}
	defer file.Close()

	fmt.Fprintf(file, "iNarbit 测试报告\n")
	fmt.Fprintf(file, "生成时间: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Fprintf(file, "="*80+"\n\n")

	passCount := 0
	failCount := 0

	for _, result := range ts.results {
		if result.Status == "PASS" {
			passCount++
		} else {
			failCount++
		}

		fmt.Fprintf(file, "测试: %s\n", result.TestName)
		fmt.Fprintf(file, "状态: %s\n", result.Status)
		fmt.Fprintf(file, "消息: %s\n", result.Message)
		fmt.Fprintf(file, "耗时: %v\n", result.Duration)
		fmt.Fprintf(file, "时间: %s\n", result.Timestamp.Format("2006-01-02 15:04:05"))
		fmt.Fprintf(file, "-"*80+"\n\n")
	}

	fmt.Fprintf(file, "="*80+"\n")
	fmt.Fprintf(file, "总计: %d 通过, %d 失败\n", passCount, failCount)
	fmt.Fprintf(file, "="*80+"\n")

	fmt.Printf("测试报告已生成: %s\n", reportFile)
}

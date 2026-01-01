# iNarbit Binance集成和测试完整指南

## 目录

1. [快速开始](#快速开始)
2. [虚拟盘测试](#虚拟盘测试)
3. [实盘连接](#实盘连接)
4. [测试套件](#测试套件)
5. [故障排除](#故障排除)
6. [安全最佳实践](#安全最佳实践)

---

## 快速开始

### 1. 项目结构

```
inarbit/
├── exchange/
│   ├── binance_client.go      # Binance API客户端
│   └── types.go               # 数据类型定义
├── simulator/
│   └── simulator.go           # 虚拟盘模拟器
├── crypto/
│   └── crypto_utils.go        # 加密和密钥管理
├── test/
│   └── test_suite.go          # 完整测试套件
├── main.go                    # 主程序
├── go.mod                     # Go模块定义
└── README.md                  # 项目文档
```

### 2. 安装依赖

```bash
# 进入项目目录
cd /root/inarbit

# 初始化Go模块
go mod init inarbit

# 下载依赖
go mod download
```

### 3. 环境配置

创建 `.env` 文件：

```bash
# Binance API配置
BINANCE_API_KEY=your_api_key_here
BINANCE_API_SECRET=your_api_secret_here
BINANCE_TESTNET=false

# 加密密钥
MASTER_KEY=inarbit-master-key-32-bytes-long!

# 数据库配置
DB_HOST=localhost
DB_PORT=5432
DB_USER=inarbit
DB_PASSWORD=inarbit_password
DB_NAME=inarbit
```

---

## 虚拟盘测试

### 1. 虚拟盘模拟器功能

虚拟盘模拟器（Simulator）提供以下功能：

- **账户模拟**: 模拟真实账户余额和交易
- **订单执行**: 模拟订单下单和成交
- **手续费计算**: 自动计算交易手续费
- **利润计算**: 计算交易利润和利润百分比
- **性能测试**: 支持大规模交易模拟

### 2. 基本使用

```go
package main

import (
	"fmt"
	"inarbit/simulator"
)

func main() {
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

	// 下单
	trade, err := exchange.PlaceOrder("BTCUSDT", "BUY", 0.01, 45000.0)
	if err != nil {
		fmt.Printf("下单失败: %v\n", err)
		return
	}

	fmt.Printf("订单成功: %+v\n", trade)

	// 查看余额
	balances := exchange.GetAllBalances()
	fmt.Printf("当前余额: %+v\n", balances)

	// 计算利润
	profit := exchange.CalculateProfit(initialBalance)
	fmt.Printf("利润: %+v\n", profit)
}
```

### 3. 三角套利模拟

```go
// 三角套利交易路径: USDT -> BTC -> ETH -> USDT
trades := []struct {
	Symbol   string
	Side     string
	Quantity float64
	Price    float64
}{
	{"BTCUSDT", "BUY", 1.0/45000, 45000.0},      // 用1 USDT买入BTC
	{"ETHBTC", "BUY", 1.0/45000/0.0667, 0.0667}, // 用BTC买入ETH
	{"ETHUSDT", "SELL", 1.0/45000/0.0667, 3000.0}, // 用ETH卖出USDT
}

result := simulator.RunSimulation(exchange, initialBalance, trades)
fmt.Printf("利润: %.4f%%\n", result.ProfitPercent["USDT"])
```

### 4. 性能测试

```bash
# 运行性能测试（1000笔交易）
go test -bench=BenchmarkSimulation -benchtime=10s

# 输出示例:
# BenchmarkSimulation-4   1000000   1234 ns/op
```

---

## 实盘连接

### 1. Binance API密钥配置

**重要**: 在使用实盘前，请确保：

1. ✅ 生成新的API密钥（不要使用已暴露的密钥）
2. ✅ 设置IP白名单（仅允许您的服务器IP）
3. ✅ 禁用提现权限
4. ✅ 限制交易金额
5. ✅ 启用2FA认证

### 2. 设置环境变量

```bash
# 导出API密钥
export BINANCE_API_KEY="your_new_api_key"
export BINANCE_API_SECRET="your_new_api_secret"
export BINANCE_TESTNET="false"

# 验证设置
echo $BINANCE_API_KEY
echo $BINANCE_API_SECRET
```

### 3. 测试连接

```go
package main

import (
	"fmt"
	"os"
	"inarbit/exchange"
)

func main() {
	apiKey := os.Getenv("BINANCE_API_KEY")
	apiSecret := os.Getenv("BINANCE_API_SECRET")

	client := exchange.NewBinanceClient(apiKey, apiSecret, false)

	// 测试连接
	err := client.TestConnection()
	if err != nil {
		fmt.Printf("连接失败: %v\n", err)
		return
	}

	fmt.Println("连接成功！")

	// 获取账户信息
	account, err := client.GetAccount()
	if err != nil {
		fmt.Printf("获取账户失败: %v\n", err)
		return
	}

	fmt.Printf("账户信息: %+v\n", account)
}
```

### 4. 实盘交易

```go
// 下单（实盘）
order, err := client.PlaceOrder("BTCUSDT", "BUY", 0.001, 45000.0)
if err != nil {
	fmt.Printf("下单失败: %v\n", err)
	return
}

fmt.Printf("订单成功: %+v\n", order)

// 查询订单
order, err = client.GetOrder("BTCUSDT", order.OrderID)
if err != nil {
	fmt.Printf("查询失败: %v\n", err)
	return
}

fmt.Printf("订单状态: %s\n", order.Status)

// 撤销订单
order, err = client.CancelOrder("BTCUSDT", order.OrderID)
if err != nil {
	fmt.Printf("撤销失败: %v\n", err)
	return
}

fmt.Printf("订单已撤销\n")
```

---

## 测试套件

### 1. 运行完整测试

```bash
# 设置环境变量
export BINANCE_API_KEY="your_api_key"
export BINANCE_API_SECRET="your_api_secret"

# 运行测试
go run test/test_suite.go

# 输出示例:
# ================================================================================
# iNarbit 完整测试套件
# ================================================================================
#
# [虚拟盘测试]
# ✓ [虚拟盘交易模拟          ] 完成 3 笔交易，利润: 0.50% (1.234ms)
# ✓ [三角套利模拟            ] 三角套利完成，利润: 0.01% (2.456ms)
# ✓ [性能测试（1000笔交易）  ] 执行1000笔交易耗时 123.456ms (456.789ms)
#
# [实盘连接测试]
# ✓ [Binance连接测试         ] 成功连接到Binance API (234.567ms)
# ✓ [获取服务器时间          ] 服务器时间: 1234567890 (123.456ms)
# ✓ [获取交易所信息          ] 获取到 1234 个交易对 (567.890ms)
#
# [行情数据测试]
# ✓ [获取行情数据            ] 获取到 1234 个行情数据 (890.123ms)
# ✓ [获取特定交易对行情      ] BTCUSDT 价格: 45123.45 (234.567ms)
# ✓ [数据验证                ] 验证成功: 3/3 (345.678ms)
#
# [账户测试]
# ✓ [获取账户信息            ] 账户有 5 个资产 (456.789ms)
#
# ================================================================================
# 总计: 10 通过, 0 失败
# ================================================================================
```

### 2. 测试项目详解

| 测试项 | 描述 | 类型 |
|--------|------|------|
| 虚拟盘交易模拟 | 模拟3笔交易并计算利润 | 虚拟盘 |
| 三角套利模拟 | 模拟三角套利交易 | 虚拟盘 |
| 性能测试 | 执行1000笔交易性能 | 虚拟盘 |
| Binance连接测试 | 测试API连接 | 实盘 |
| 获取服务器时间 | 获取Binance服务器时间 | 实盘 |
| 获取交易所信息 | 获取所有交易对信息 | 实盘 |
| 获取行情数据 | 获取所有交易对行情 | 实盘 |
| 获取特定交易对行情 | 获取BTCUSDT行情 | 实盘 |
| 数据验证 | 验证多个交易对数据 | 实盘 |
| 获取账户信息 | 获取账户余额和资产 | 实盘 |

### 3. 单个测试运行

```bash
# 只运行虚拟盘测试
go test -run TestVirtual ./test

# 只运行实盘测试
go test -run TestLive ./test

# 运行特定测试
go test -run TestBinanceConnection ./test
```

---

## 故障排除

### 问题1: API连接失败

**症状**: `连接失败: connection refused`

**解决方案**:

```bash
# 检查网络连接
ping api.binance.com

# 检查防火墙
sudo ufw allow 443

# 检查API密钥
echo $BINANCE_API_KEY
echo $BINANCE_API_SECRET
```

### 问题2: 无效的API密钥

**症状**: `API错误 (HTTP 401): Invalid API-key`

**解决方案**:

```bash
# 确保API密钥正确
# 1. 检查Binance账户设置
# 2. 确认密钥未过期
# 3. 重新生成新的API密钥
# 4. 更新环境变量
```

### 问题3: IP白名单限制

**症状**: `API错误 (HTTP 403): Forbidden`

**解决方案**:

```bash
# 在Binance账户设置中添加服务器IP
# 1. 获取服务器IP: curl ifconfig.me
# 2. 登录Binance账户
# 3. API管理 -> IP白名单 -> 添加IP
# 4. 重试连接
```

### 问题4: 时间同步问题

**症状**: `API错误 (HTTP 400): Timestamp for this request was 1000ms ahead of the server's time`

**解决方案**:

```bash
# 同步系统时间
sudo timedatectl set-ntp true

# 检查时间
date

# 手动同步
sudo ntpdate -s time.nist.gov
```

### 问题5: 余额不足

**症状**: `余额不足: 需要 X，实际 Y`

**解决方案**:

```bash
# 检查账户余额
# 1. 登录Binance账户
# 2. 查看钱包 -> 现货账户
# 3. 确保有足够的资金
# 4. 检查锁定的资金（未成交订单）
```

---

## 安全最佳实践

### 1. API密钥安全

✅ **必须做**:
- 使用环境变量存储API密钥
- 定期轮换API密钥
- 为不同的机器人使用不同的密钥
- 启用IP白名单
- 禁用提现权限
- 限制交易金额

❌ **不要做**:
- 在代码中硬编码API密钥
- 在Git中提交API密钥
- 在聊天或邮件中分享API密钥
- 使用过期的API密钥
- 在公网上暴露API密钥

### 2. 密钥加密存储

```go
// 使用加密工具存储敏感信息
km, err := crypto.NewKeyManager("")
if err != nil {
	log.Fatal(err)
}

// 加密API密钥
encryptedKey, err := km.EncryptAPIKey("your_api_key")
if err != nil {
	log.Fatal(err)
}

// 存储加密后的密钥到数据库
// ...

// 使用时解密
decryptedKey, err := km.DecryptAPIKey(encryptedKey)
if err != nil {
	log.Fatal(err)
}
```

### 3. 账户安全

```bash
# 1. 启用2FA认证
# 2. 设置IP白名单
# 3. 限制API权限
# 4. 定期检查账户活动
# 5. 设置提现地址白名单
```

### 4. 交易安全

```go
// 1. 始终在虚拟盘测试后再进行实盘交易
// 2. 从小额开始测试
// 3. 设置止损和止盈
// 4. 监控交易执行
// 5. 记录所有交易日志
```

---

## 监控和日志

### 1. 启用日志

```go
import "log"

// 创建日志文件
logFile, err := os.OpenFile("trading.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
if err != nil {
	log.Fatal(err)
}
defer logFile.Close()

// 设置日志输出
log.SetOutput(logFile)
log.SetFlags(log.LstdFlags | log.Lshortfile)

// 记录交易
log.Printf("下单: %s %s %.8f @ %.2f", symbol, side, quantity, price)
```

### 2. 性能监控

```bash
# 监控CPU使用率
top -p $(pgrep -f inarbit)

# 监控内存使用
ps aux | grep inarbit

# 监控网络连接
netstat -an | grep :443
```

### 3. 错误处理

```go
// 重试机制
func retryRequest(fn func() error, maxRetries int) error {
	var lastErr error
	for i := 0; i < maxRetries; i++ {
		if err := fn(); err == nil {
			return nil
		} else {
			lastErr = err
			time.Sleep(time.Duration(2^i) * time.Second) // 指数退避
		}
	}
	return lastErr
}
```

---

## 总结

本指南涵盖了iNarbit与Binance集成的完整过程，包括虚拟盘测试、实盘连接和安全最佳实践。

**建议的部署流程**:

1. ✅ 虚拟盘测试（确保逻辑正确）
2. ✅ 小额实盘测试（测试API连接）
3. ✅ 监控和日志（记录所有交易）
4. ✅ 逐步增加交易金额（风险管理）
5. ✅ 定期审计和优化（持续改进）

如有问题，请参考Binance官方文档：https://binance-docs.github.io/apidocs/

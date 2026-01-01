package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// BinanceClient Binance API客户端
type BinanceClient struct {
	APIKey    string
	APISecret string
	BaseURL   string
	IsTestnet bool
	HTTPClient *http.Client
}

// NewBinanceClient 创建Binance客户端
func NewBinanceClient(apiKey, apiSecret string, isTestnet bool) *BinanceClient {
	baseURL := "https://api.binance.com"
	if isTestnet {
		baseURL = "https://testnet.binance.vision"
	}

	return &BinanceClient{
		APIKey:     apiKey,
		APISecret:  apiSecret,
		BaseURL:    baseURL,
		IsTestnet:  isTestnet,
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// ===== 行情数据结构 =====

// Ticker 交易对信息
type Ticker struct {
	Symbol             string  `json:"symbol"`
	BidPrice           float64 `json:"bidPrice,string"`
	BidQty             float64 `json:"bidQty,string"`
	AskPrice           float64 `json:"askPrice,string"`
	AskQty             float64 `json:"askQty,string"`
	LastPrice          float64 `json:"lastPrice,string"`
	LastQty            float64 `json:"lastQty,string"`
	OpenPrice          float64 `json:"openPrice,string"`
	HighPrice          float64 `json:"highPrice,string"`
	LowPrice           float64 `json:"lowPrice,string"`
	Volume             float64 `json:"volume,string"`
	QuoteAssetVolume   float64 `json:"quoteAssetVolume,string"`
	OpenTime           int64   `json:"openTime"`
	CloseTime          int64   `json:"closeTime"`
	FirstID            int64   `json:"firstId"`
	LastID             int64   `json:"lastId"`
	Count              int64   `json:"count"`
	PriceChange        float64 `json:"priceChange,string"`
	PriceChangePercent float64 `json:"priceChangePercent,string"`
	WeightedAvgPrice   float64 `json:"weightedAvgPrice,string"`
}

// ExchangeInfo 交易所信息
type ExchangeInfo struct {
	Symbols []SymbolInfo `json:"symbols"`
}

// SymbolInfo 交易对信息
type SymbolInfo struct {
	Symbol              string        `json:"symbol"`
	Status              string        `json:"status"`
	BaseAsset           string        `json:"baseAsset"`
	BaseAssetPrecision  int           `json:"baseAssetPrecision"`
	QuoteAsset          string        `json:"quoteAsset"`
	QuoteAssetPrecision int           `json:"quoteAssetPrecision"`
	OrderTypes          []string      `json:"orderTypes"`
	IcebergAllowed      bool          `json:"icebergAllowed"`
	Filters             []FilterInfo  `json:"filters"`
	Permissions         []string      `json:"permissions"`
}

// FilterInfo 过滤器信息
type FilterInfo struct {
	FilterType string `json:"filterType"`
	MinPrice   string `json:"minPrice,omitempty"`
	MaxPrice   string `json:"maxPrice,omitempty"`
	TickSize   string `json:"tickSize,omitempty"`
	MinQty     string `json:"minQty,omitempty"`
	MaxQty     string `json:"maxQty,omitempty"`
	StepSize   string `json:"stepSize,omitempty"`
	MinNotional string `json:"minNotional,omitempty"`
	ApplyToMarket bool `json:"applyToMarket,omitempty"`
}

// Account 账户信息
type Account struct {
	MakerCommission  int      `json:"makerCommission"`
	TakerCommission  int      `json:"takerCommission"`
	BuyerCommission  int      `json:"buyerCommission"`
	SellerCommission int      `json:"sellerCommission"`
	CanTrade         bool     `json:"canTrade"`
	CanDeposit       bool     `json:"canDeposit"`
	CanWithdraw      bool     `json:"canWithdraw"`
	UpdateTime       int64    `json:"updateTime"`
	Balances         []Balance `json:"balances"`
}

// Balance 余额信息
type Balance struct {
	Asset  string  `json:"asset"`
	Free   float64 `json:"free,string"`
	Locked float64 `json:"locked,string"`
}

// Order 订单信息
type Order struct {
	Symbol            string  `json:"symbol"`
	OrderID           int64   `json:"orderId"`
	OrderListID       int64   `json:"orderListId"`
	ClientOrderID     string  `json:"clientOrderId"`
	Price             float64 `json:"price,string"`
	OrigQty           float64 `json:"origQty,string"`
	ExecutedQty       float64 `json:"executedQty,string"`
	CummulativeQuoteQty float64 `json:"cummulativeQuoteQty,string"`
	Status            string  `json:"status"`
	TimeInForce       string  `json:"timeInForce"`
	Type              string  `json:"type"`
	Side              string  `json:"side"`
	StopPrice         float64 `json:"stopPrice,string"`
	IcebergQty        float64 `json:"icebergQty,string"`
	Time              int64   `json:"time"`
	UpdateTime        int64   `json:"updateTime"`
	IsWorking         bool    `json:"isWorking"`
	OrigQuoteOrderQty float64 `json:"origQuoteOrderQty,string"`
}

// ===== 公开API方法 =====

// GetTicker 获取交易对行情
func (c *BinanceClient) GetTicker(symbol string) (*Ticker, error) {
	params := url.Values{}
	params.Add("symbol", symbol)

	body, err := c.doRequest("GET", "/api/v3/ticker/24hr", params, false)
	if err != nil {
		return nil, err
	}

	var ticker Ticker
	if err := json.Unmarshal(body, &ticker); err != nil {
		return nil, fmt.Errorf("解析行情数据失败: %w", err)
	}

	return &ticker, nil
}

// GetAllTickers 获取所有交易对行情
func (c *BinanceClient) GetAllTickers() ([]*Ticker, error) {
	body, err := c.doRequest("GET", "/api/v3/ticker/24hr", url.Values{}, false)
	if err != nil {
		return nil, err
	}

	var tickers []*Ticker
	if err := json.Unmarshal(body, &tickers); err != nil {
		return nil, fmt.Errorf("解析行情数据失败: %w", err)
	}

	return tickers, nil
}

// GetExchangeInfo 获取交易所信息
func (c *BinanceClient) GetExchangeInfo() (*ExchangeInfo, error) {
	body, err := c.doRequest("GET", "/api/v3/exchangeInfo", url.Values{}, false)
	if err != nil {
		return nil, err
	}

	var info ExchangeInfo
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, fmt.Errorf("解析交易所信息失败: %w", err)
	}

	return &info, nil
}

// ===== 账户API方法 =====

// GetAccount 获取账户信息
func (c *BinanceClient) GetAccount() (*Account, error) {
	params := url.Values{}
	params.Add("timestamp", strconv.FormatInt(time.Now().UnixMilli(), 10))

	body, err := c.doRequest("GET", "/api/v3/account", params, true)
	if err != nil {
		return nil, err
	}

	var account Account
	if err := json.Unmarshal(body, &account); err != nil {
		return nil, fmt.Errorf("解析账户信息失败: %w", err)
	}

	return &account, nil
}

// GetBalance 获取特定资产余额
func (c *BinanceClient) GetBalance(asset string) (*Balance, error) {
	account, err := c.GetAccount()
	if err != nil {
		return nil, err
	}

	for _, balance := range account.Balances {
		if balance.Asset == asset {
			return &balance, nil
		}
	}

	return nil, fmt.Errorf("资产 %s 不存在", asset)
}

// ===== 交易API方法 =====

// PlaceOrder 下单
func (c *BinanceClient) PlaceOrder(symbol string, side string, quantity float64, price float64) (*Order, error) {
	params := url.Values{}
	params.Add("symbol", symbol)
	params.Add("side", side)
	params.Add("type", "LIMIT")
	params.Add("timeInForce", "GTC")
	params.Add("quantity", fmt.Sprintf("%.8f", quantity))
	params.Add("price", fmt.Sprintf("%.8f", price))
	params.Add("timestamp", strconv.FormatInt(time.Now().UnixMilli(), 10))

	body, err := c.doRequest("POST", "/api/v3/order", params, true)
	if err != nil {
		return nil, err
	}

	var order Order
	if err := json.Unmarshal(body, &order); err != nil {
		return nil, fmt.Errorf("解析订单信息失败: %w", err)
	}

	return &order, nil
}

// PlaceMarketOrder 市价下单
func (c *BinanceClient) PlaceMarketOrder(symbol string, side string, quantity float64) (*Order, error) {
	params := url.Values{}
	params.Add("symbol", symbol)
	params.Add("side", side)
	params.Add("type", "MARKET")
	params.Add("quantity", fmt.Sprintf("%.8f", quantity))
	params.Add("timestamp", strconv.FormatInt(time.Now().UnixMilli(), 10))

	body, err := c.doRequest("POST", "/api/v3/order", params, true)
	if err != nil {
		return nil, err
	}

	var order Order
	if err := json.Unmarshal(body, &order); err != nil {
		return nil, fmt.Errorf("解析订单信息失败: %w", err)
	}

	return &order, nil
}

// CancelOrder 撤销订单
func (c *BinanceClient) CancelOrder(symbol string, orderID int64) (*Order, error) {
	params := url.Values{}
	params.Add("symbol", symbol)
	params.Add("orderId", strconv.FormatInt(orderID, 10))
	params.Add("timestamp", strconv.FormatInt(time.Now().UnixMilli(), 10))

	body, err := c.doRequest("DELETE", "/api/v3/order", params, true)
	if err != nil {
		return nil, err
	}

	var order Order
	if err := json.Unmarshal(body, &order); err != nil {
		return nil, fmt.Errorf("解析订单信息失败: %w", err)
	}

	return &order, nil
}

// GetOrder 查询订单
func (c *BinanceClient) GetOrder(symbol string, orderID int64) (*Order, error) {
	params := url.Values{}
	params.Add("symbol", symbol)
	params.Add("orderId", strconv.FormatInt(orderID, 10))
	params.Add("timestamp", strconv.FormatInt(time.Now().UnixMilli(), 10))

	body, err := c.doRequest("GET", "/api/v3/order", params, true)
	if err != nil {
		return nil, err
	}

	var order Order
	if err := json.Unmarshal(body, &order); err != nil {
		return nil, fmt.Errorf("解析订单信息失败: %w", err)
	}

	return &order, nil
}

// GetOpenOrders 获取未成交订单
func (c *BinanceClient) GetOpenOrders(symbol string) ([]*Order, error) {
	params := url.Values{}
	if symbol != "" {
		params.Add("symbol", symbol)
	}
	params.Add("timestamp", strconv.FormatInt(time.Now().UnixMilli(), 10))

	body, err := c.doRequest("GET", "/api/v3/openOrders", params, true)
	if err != nil {
		return nil, err
	}

	var orders []*Order
	if err := json.Unmarshal(body, &orders); err != nil {
		return nil, fmt.Errorf("解析订单列表失败: %w", err)
	}

	return orders, nil
}

// ===== 私有请求辅助方法 =====

// doRequest 执行HTTP请求
func (c *BinanceClient) doRequest(method string, endpoint string, params url.Values, signed bool) ([]byte, error) {
	if signed {
		// 添加签名
		queryString := params.Encode()
		signature := c.sign(queryString)
		params.Add("signature", signature)
	}

	fullURL := c.BaseURL + endpoint
	if len(params) > 0 {
		fullURL += "?" + params.Encode()
	}

	req, err := http.NewRequest(method, fullURL, nil)
	if err != nil {
		return nil, err
	}

	// 添加API密钥
	req.Header.Add("X-MBX-APIKEY", c.APIKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API错误 (HTTP %d): %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// sign 对请求进行签名
func (c *BinanceClient) sign(message string) string {
	mac := hmac.New(sha256.New, []byte(c.APISecret))
	mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}

// ===== WebSocket行情推送 =====

// TickerStream WebSocket行情流
type TickerStream struct {
	Symbol string
	Price  float64
	Time   time.Time
}

// SubscribeTickerStream 订阅行情流
// 注意：这是一个简化的实现，实际应该使用gorilla/websocket
func (c *BinanceClient) SubscribeTickerStream(symbols []string, callback func(*TickerStream)) error {
	// 构建WebSocket URL
	wsURL := "wss://stream.binance.com:9443/ws"
	
	// 添加订阅流
	streams := make([]string, len(symbols))
	for i, symbol := range symbols {
		streams[i] = strings.ToLower(symbol) + "@ticker"
	}
	
	wsURL += "/" + strings.Join(streams, "/")
	
	log.Printf("订阅WebSocket行情流: %s", wsURL)
	
	// TODO: 实现WebSocket连接和消息处理
	// 这需要使用gorilla/websocket库
	
	return nil
}

// ===== 辅助方法 =====

// ValidateSymbol 验证交易对是否存在
func (c *BinanceClient) ValidateSymbol(symbol string) (bool, error) {
	info, err := c.GetExchangeInfo()
	if err != nil {
		return false, err
	}

	for _, s := range info.Symbols {
		if s.Symbol == symbol && s.Status == "TRADING" {
			return true, nil
		}
	}

	return false, nil
}

// GetSymbolInfo 获取交易对详细信息
func (c *BinanceClient) GetSymbolInfo(symbol string) (*SymbolInfo, error) {
	info, err := c.GetExchangeInfo()
	if err != nil {
		return nil, err
	}

	for _, s := range info.Symbols {
		if s.Symbol == symbol {
			return &s, nil
		}
	}

	return nil, fmt.Errorf("交易对 %s 不存在", symbol)
}

// GetMinNotional 获取最小成交额
func (c *BinanceClient) GetMinNotional(symbol string) (float64, error) {
	symbolInfo, err := c.GetSymbolInfo(symbol)
	if err != nil {
		return 0, err
	}

	for _, filter := range symbolInfo.Filters {
		if filter.FilterType == "MIN_NOTIONAL" {
			minNotional, err := strconv.ParseFloat(filter.MinNotional, 64)
			if err != nil {
				return 0, err
			}
			return minNotional, nil
		}
	}

	return 10, nil // 默认最小成交额
}

// GetStepSize 获取数量精度
func (c *BinanceClient) GetStepSize(symbol string) (float64, error) {
	symbolInfo, err := c.GetSymbolInfo(symbol)
	if err != nil {
		return 0, err
	}

	for _, filter := range symbolInfo.Filters {
		if filter.FilterType == "LOT_SIZE" {
			stepSize, err := strconv.ParseFloat(filter.StepSize, 64)
			if err != nil {
				return 0, err
			}
			return stepSize, nil
		}
	}

	return 0.00000001, nil // 默认精度
}

// GetTickSize 获取价格精度
func (c *BinanceClient) GetTickSize(symbol string) (float64, error) {
	symbolInfo, err := c.GetSymbolInfo(symbol)
	if err != nil {
		return 0, err
	}

	for _, filter := range symbolInfo.Filters {
		if filter.FilterType == "PRICE_FILTER" {
			tickSize, err := strconv.ParseFloat(filter.TickSize, 64)
			if err != nil {
				return 0, err
			}
			return tickSize, nil
		}
	}

	return 0.00000001, nil // 默认精度
}

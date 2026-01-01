package exchange

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// BinanceClient Binance API客户端
type BinanceClient struct {
	apiKey    string
	apiSecret string
	baseURL   string
	httpClient *http.Client
	isTestnet bool
}

// NewBinanceClient 创建新的Binance客户端
func NewBinanceClient(apiKey, apiSecret string, isTestnet bool) *BinanceClient {
	baseURL := "https://api.binance.com"
	if isTestnet {
		baseURL = "https://testnet.binance.vision"
	}

	return &BinanceClient{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		baseURL:   baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		isTestnet: isTestnet,
	}
}

// ExchangeInfo 交易所信息
type ExchangeInfo struct {
	Timezone        string      `json:"timezone"`
	ServerTime      int64       `json:"serverTime"`
	Symbols         []SymbolInfo `json:"symbols"`
}

// SymbolInfo 交易对信息
type SymbolInfo struct {
	Symbol             string `json:"symbol"`
	Status             string `json:"status"`
	BaseAsset          string `json:"baseAsset"`
	QuoteAsset         string `json:"quoteAsset"`
	BaseAssetPrecision int    `json:"baseAssetPrecision"`
	QuotePrecision     int    `json:"quotePrecision"`
	OrderTypes         []string `json:"orderTypes"`
	Filters            []map[string]interface{} `json:"filters"`
}

// Ticker 行情数据
type Ticker struct {
	Symbol   string `json:"symbol"`
	BidPrice string `json:"bidPrice"`
	BidQty   string `json:"bidQty"`
	AskPrice string `json:"askPrice"`
	AskQty   string `json:"askQty"`
	LastPrice string `json:"lastPrice"`
	Volume   string `json:"volume"`
}

// Account 账户信息
type Account struct {
	MakerCommission  int       `json:"makerCommission"`
	TakerCommission  int       `json:"takerCommission"`
	BuyerCommission  int       `json:"buyerCommission"`
	SellerCommission int       `json:"sellerCommission"`
	CanTrade         bool      `json:"canTrade"`
	CanDeposit       bool      `json:"canDeposit"`
	CanWithdraw      bool      `json:"canWithdraw"`
	UpdateTime       int64     `json:"updateTime"`
	Balances         []Balance `json:"balances"`
}

// Balance 余额信息
type Balance struct {
	Asset  string `json:"asset"`
	Free   string `json:"free"`
	Locked string `json:"locked"`
}

// Order 订单信息
type Order struct {
	Symbol            string `json:"symbol"`
	OrderID           int64  `json:"orderId"`
	ClientOrderID     string `json:"clientOrderId"`
	Price             string `json:"price"`
	OrigQty           string `json:"origQty"`
	ExecutedQty       string `json:"executedQty"`
	CummulativeQuoteQty string `json:"cummulativeQuoteQty"`
	Status            string `json:"status"`
	TimeInForce       string `json:"timeInForce"`
	Type              string `json:"type"`
	Side              string `json:"side"`
	StopPrice         string `json:"stopPrice"`
	IcebergQty        string `json:"icebergQty"`
	Time              int64  `json:"time"`
	UpdateTime        int64  `json:"updateTime"`
	IsWorking         bool   `json:"isWorking"`
	Fills             []OrderFill `json:"fills"`
}

// OrderFill 订单成交信息
type OrderFill struct {
	Price           string `json:"price"`
	Qty             string `json:"qty"`
	Commission      string `json:"commission"`
	CommissionAsset string `json:"commissionAsset"`
	TradeID         int64  `json:"tradeId"`
}

// GetExchangeInfo 获取交易所信息
func (bc *BinanceClient) GetExchangeInfo() (*ExchangeInfo, error) {
	path := "/api/v3/exchangeInfo"
	resp, err := bc.request("GET", path, nil, false)
	if err != nil {
		return nil, err
	}

	var info ExchangeInfo
	if err := json.Unmarshal(resp, &info); err != nil {
		return nil, fmt.Errorf("解析交易所信息失败: %v", err)
	}

	return &info, nil
}

// GetTicker 获取行情数据
func (bc *BinanceClient) GetTicker(symbol string) (*Ticker, error) {
	path := "/api/v3/ticker/bookTicker"
	params := url.Values{}
	params.Add("symbol", symbol)

	resp, err := bc.request("GET", path, params, false)
	if err != nil {
		return nil, err
	}

	var ticker Ticker
	if err := json.Unmarshal(resp, &ticker); err != nil {
		return nil, fmt.Errorf("解析行情数据失败: %v", err)
	}

	return &ticker, nil
}

// GetAllTickers 获取所有行情数据
func (bc *BinanceClient) GetAllTickers() ([]Ticker, error) {
	path := "/api/v3/ticker/bookTicker"
	resp, err := bc.request("GET", path, nil, false)
	if err != nil {
		return nil, err
	}

	var tickers []Ticker
	if err := json.Unmarshal(resp, &tickers); err != nil {
		return nil, fmt.Errorf("解析行情数据失败: %v", err)
	}

	return tickers, nil
}

// GetAccount 获取账户信息
func (bc *BinanceClient) GetAccount() (*Account, error) {
	path := "/api/v3/account"
	resp, err := bc.request("GET", path, nil, true)
	if err != nil {
		return nil, err
	}

	var account Account
	if err := json.Unmarshal(resp, &account); err != nil {
		return nil, fmt.Errorf("解析账户信息失败: %v", err)
	}

	return &account, nil
}

// PlaceOrder 下单
func (bc *BinanceClient) PlaceOrder(symbol, side, orderType string, quantity, price float64) (*Order, error) {
	path := "/api/v3/order"
	params := url.Values{}
	params.Add("symbol", symbol)
	params.Add("side", side)
	params.Add("type", orderType)
	params.Add("timeInForce", "GTC")
	params.Add("quantity", fmt.Sprintf("%.8f", quantity))
	
	if price > 0 {
		params.Add("price", fmt.Sprintf("%.8f", price))
	}

	resp, err := bc.request("POST", path, params, true)
	if err != nil {
		return nil, err
	}

	var order Order
	if err := json.Unmarshal(resp, &order); err != nil {
		return nil, fmt.Errorf("解析订单信息失败: %v", err)
	}

	return &order, nil
}

// CancelOrder 撤销订单
func (bc *BinanceClient) CancelOrder(symbol string, orderID int64) (*Order, error) {
	path := "/api/v3/order"
	params := url.Values{}
	params.Add("symbol", symbol)
	params.Add("orderId", fmt.Sprintf("%d", orderID))

	resp, err := bc.request("DELETE", path, params, true)
	if err != nil {
		return nil, err
	}

	var order Order
	if err := json.Unmarshal(resp, &order); err != nil {
		return nil, fmt.Errorf("解析订单信息失败: %v", err)
	}

	return &order, nil
}

// GetOrder 查询订单
func (bc *BinanceClient) GetOrder(symbol string, orderID int64) (*Order, error) {
	path := "/api/v3/order"
	params := url.Values{}
	params.Add("symbol", symbol)
	params.Add("orderId", fmt.Sprintf("%d", orderID))

	resp, err := bc.request("GET", path, params, true)
	if err != nil {
		return nil, err
	}

	var order Order
	if err := json.Unmarshal(resp, &order); err != nil {
		return nil, fmt.Errorf("解析订单信息失败: %v", err)
	}

	return &order, nil
}

// GetOpenOrders 获取未成交订单
func (bc *BinanceClient) GetOpenOrders(symbol string) ([]Order, error) {
	path := "/api/v3/openOrders"
	params := url.Values{}
	if symbol != "" {
		params.Add("symbol", symbol)
	}

	resp, err := bc.request("GET", path, params, true)
	if err != nil {
		return nil, err
	}

	var orders []Order
	if err := json.Unmarshal(resp, &orders); err != nil {
		return nil, fmt.Errorf("解析订单列表失败: %v", err)
	}

	return orders, nil
}

// request 发送HTTP请求
func (bc *BinanceClient) request(method, path string, params url.Values, signed bool) ([]byte, error) {
	if params == nil {
		params = url.Values{}
	}

	// 添加时间戳
	params.Add("timestamp", fmt.Sprintf("%d", time.Now().UnixMilli()))

	// 签名
	if signed {
		signature := bc.sign(params.Encode())
		params.Add("signature", signature)
	}

	// 构建完整URL
	fullURL := bc.baseURL + path
	if method == "GET" || method == "DELETE" {
		fullURL += "?" + params.Encode()
	}

	// 创建请求
	req, err := http.NewRequest(method, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 添加请求头
	req.Header.Add("X-MBX-APIKEY", bc.apiKey)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// 如果是POST/PUT，将参数放在body中
	if method == "POST" || method == "PUT" {
		req.Body = io.NopCloser(strings.NewReader(params.Encode()))
	}

	// 发送请求
	resp, err := bc.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API错误 (HTTP %d): %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// sign 签名请求
func (bc *BinanceClient) sign(message string) string {
	h := hmac.New(sha256.New, []byte(bc.apiSecret))
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

// TestConnection 测试连接
func (bc *BinanceClient) TestConnection() error {
	path := "/api/v3/ping"
	_, err := bc.request("GET", path, nil, false)
	return err
}

// GetServerTime 获取服务器时间
func (bc *BinanceClient) GetServerTime() (int64, error) {
	path := "/api/v3/time"
	resp, err := bc.request("GET", path, nil, false)
	if err != nil {
		return 0, err
	}

	var result map[string]int64
	if err := json.Unmarshal(resp, &result); err != nil {
		return 0, fmt.Errorf("解析时间失败: %v", err)
	}

	return result["serverTime"], nil
}

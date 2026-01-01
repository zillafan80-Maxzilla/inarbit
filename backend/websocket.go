package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// WebSocketManager WebSocket连接管理器
type WebSocketManager struct {
	clients    map[int64]*WebSocketClient // userID -> client
	broadcast  chan interface{}            // 广播消息通道
	register   chan *WebSocketClient       // 注册通道
	unregister chan *WebSocketClient       // 注销通道
	mu         sync.RWMutex
	authService *AuthService
	db         *Database
}

// WebSocketClient WebSocket客户端
type WebSocketClient struct {
	UserID   int64
	Conn     *websocket.Conn
	Send     chan interface{}
	Manager  *WebSocketManager
	Done     chan struct{}
	LastPing time.Time
}

// NewWebSocketManager 创建WebSocket管理器
func NewWebSocketManager(authService *AuthService, db *Database) *WebSocketManager {
	return &WebSocketManager{
		clients:     make(map[int64]*WebSocketClient),
		broadcast:   make(chan interface{}, 256),
		register:    make(chan *WebSocketClient),
		unregister:  make(chan *WebSocketClient),
		authService: authService,
		db:          db,
	}
}

// Run 运行WebSocket管理器
func (m *WebSocketManager) Run() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case client := <-m.register:
			m.mu.Lock()
			m.clients[client.UserID] = client
			m.mu.Unlock()
			log.Printf("✓ WebSocket客户端已连接 (UserID: %d)", client.UserID)

		case client := <-m.unregister:
			m.mu.Lock()
			if c, ok := m.clients[client.UserID]; ok && c == client {
				delete(m.clients, client.UserID)
				close(client.Send)
			}
			m.mu.Unlock()
			log.Printf("✗ WebSocket客户端已断开 (UserID: %d)", client.UserID)

		case message := <-m.broadcast:
			m.mu.RLock()
			for _, client := range m.clients {
				select {
				case client.Send <- message:
				default:
					// 如果发送通道满，跳过此客户端
					log.Printf("警告：无法发送消息给客户端 %d", client.UserID)
				}
			}
			m.mu.RUnlock()

		case <-ticker.C:
			// 定期检查客户端连接
			m.mu.RLock()
			for _, client := range m.clients {
				if time.Since(client.LastPing) > 2*time.Minute {
					m.mu.RUnlock()
					m.unregister <- client
					m.mu.RLock()
				}
			}
			m.mu.RUnlock()
		}
	}
}

// HandleConnection 处理WebSocket连接
func (m *WebSocketManager) HandleConnection(conn *websocket.Conn) {
	client := &WebSocketClient{
		Conn:     conn,
		Send:     make(chan interface{}, 256),
		Manager:  m,
		Done:     make(chan struct{}),
		LastPing: time.Now(),
	}

	// 设置连接参数
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		client.LastPing = time.Now()
		return nil
	})

	// 等待认证消息
	var msg WebSocketMessage
	err := conn.ReadJSON(&msg)
	if err != nil {
		log.Printf("WebSocket读取错误: %v", err)
		conn.Close()
		return
	}

	// 验证token
	if msg.Type != "auth" || msg.Token == "" {
		log.Printf("WebSocket认证失败：缺少token")
		conn.WriteJSON(map[string]string{"error": "认证失败"})
		conn.Close()
		return
	}

	claims, err := m.authService.VerifyToken(msg.Token)
	if err != nil {
		log.Printf("WebSocket token验证失败: %v", err)
		conn.WriteJSON(map[string]string{"error": "token无效"})
		conn.Close()
		return
	}

	client.UserID = claims.UserID

	// 注册客户端
	m.register <- client

	// 发送连接成功消息
	client.Send <- WebSocketMessage{
		Type: "connected",
		Payload: map[string]interface{}{
			"message": "连接成功",
			"user_id": client.UserID,
		},
	}

	// 启动读取和写入goroutine
	go client.ReadPump()
	go client.WritePump()
}

// ReadPump 读取消息
func (c *WebSocketClient) ReadPump() {
	defer func() {
		c.Manager.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		c.LastPing = time.Now()
		return nil
	})

	for {
		var msg WebSocketMessage
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket错误: %v", err)
			}
			return
		}

		c.LastPing = time.Now()

		// 处理不同类型的消息
		switch msg.Type {
		case "ping":
			c.Send <- WebSocketMessage{Type: "pong"}

		case "subscribe":
			// 订阅特定事件
			c.handleSubscribe(msg)

		case "unsubscribe":
			// 取消订阅
			c.handleUnsubscribe(msg)

		default:
			log.Printf("未知的WebSocket消息类型: %s", msg.Type)
		}
	}
}

// WritePump 写入消息
func (c *WebSocketClient) WritePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// 序列化消息
			data, err := json.Marshal(message)
			if err != nil {
				log.Printf("消息序列化失败: %v", err)
				continue
			}

			err = c.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}

		case <-c.Done:
			return
		}
	}
}

// handleSubscribe 处理订阅请求
func (c *WebSocketClient) handleSubscribe(msg WebSocketMessage) {
	payload, ok := msg.Payload.(map[string]interface{})
	if !ok {
		return
	}

	eventType, ok := payload["event"].(string)
	if !ok {
		return
	}

	log.Printf("用户 %d 订阅事件: %s", c.UserID, eventType)

	// 根据事件类型发送初始数据
	switch eventType {
	case "stats":
		c.sendStatsUpdate()
	case "trades":
		// 发送最近的交易
		c.sendRecentTrades()
	}
}

// handleUnsubscribe 处理取消订阅请求
func (c *WebSocketClient) handleUnsubscribe(msg WebSocketMessage) {
	payload, ok := msg.Payload.(map[string]interface{})
	if !ok {
		return
	}

	eventType, ok := payload["event"].(string)
	if !ok {
		return
	}

	log.Printf("用户 %d 取消订阅事件: %s", c.UserID, eventType)
}

// sendStatsUpdate 发送统计更新
func (c *WebSocketClient) sendStatsUpdate() {
	stats, err := c.Manager.db.GetDashboardStats(c.UserID)
	if err != nil {
		log.Printf("获取统计数据失败: %v", err)
		return
	}

	c.Send <- WebSocketMessage{
		Type: "stats_update",
		Payload: StatsUpdatePayload{
			TotalProfit: stats.TotalProfit,
			TotalTrades: stats.TotalTrades,
			ActiveBots:  stats.ActiveBots,
			WinRate:     stats.WinRate,
			Timestamp:   time.Now(),
		},
	}
}

// sendRecentTrades 发送最近的交易
func (c *WebSocketClient) sendRecentTrades() {
	// TODO: 实现获取最近交易的逻辑
	log.Printf("用户 %d 请求最近交易", c.UserID)
}

// BroadcastStatsUpdate 广播统计更新
func (m *WebSocketManager) BroadcastStatsUpdate(userID int64, stats *DashboardStats) {
	message := WebSocketMessage{
		Type: "stats_update",
		Payload: StatsUpdatePayload{
			TotalProfit: stats.TotalProfit,
			TotalTrades: stats.TotalTrades,
			ActiveBots:  stats.ActiveBots,
			WinRate:     stats.WinRate,
			Timestamp:   time.Now(),
		},
	}

	m.mu.RLock()
	if client, ok := m.clients[userID]; ok {
		select {
		case client.Send <- message:
		default:
			log.Printf("警告：无法发送统计更新给用户 %d", userID)
		}
	}
	m.mu.RUnlock()
}

// BroadcastTradeUpdate 广播交易更新
func (m *WebSocketManager) BroadcastTradeUpdate(userID int64, trade *Trade) {
	message := WebSocketMessage{
		Type: "trade_update",
		Payload: TradeUpdatePayload{
			BotID:     trade.BotID,
			TradeID:   trade.ID,
			Status:    trade.Status,
			Profit:    *trade.Profit,
			Timestamp: time.Now(),
		},
	}

	m.mu.RLock()
	if client, ok := m.clients[userID]; ok {
		select {
		case client.Send <- message:
		default:
			log.Printf("警告：无法发送交易更新给用户 %d", userID)
		}
	}
	m.mu.RUnlock()
}

// BroadcastBotStatus 广播机器人状态
func (m *WebSocketManager) BroadcastBotStatus(userID int64, botID int64, isRunning bool) {
	message := WebSocketMessage{
		Type: "bot_status",
		Payload: BotStatusPayload{
			BotID:     botID,
			IsRunning: isRunning,
			Timestamp: time.Now(),
		},
	}

	m.mu.RLock()
	if client, ok := m.clients[userID]; ok {
		select {
		case client.Send <- message:
		default:
			log.Printf("警告：无法发送机器人状态给用户 %d", userID)
		}
	}
	m.mu.RUnlock()
}

// BroadcastLog 广播日志
func (m *WebSocketManager) BroadcastLog(userID int64, botID *int64, level string, message string) {
	wsMsg := WebSocketMessage{
		Type: "log",
		Payload: LogPayload{
			BotID:   botID,
			Level:   level,
			Message: message,
			Time:    time.Now(),
		},
	}

	m.mu.RLock()
	if client, ok := m.clients[userID]; ok {
		select {
		case client.Send <- wsMsg:
		default:
			log.Printf("警告：无法发送日志给用户 %d", userID)
		}
	}
	m.mu.RUnlock()
}

// GetConnectedUsers 获取已连接的用户列表
func (m *WebSocketManager) GetConnectedUsers() []int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	users := make([]int64, 0, len(m.clients))
	for userID := range m.clients {
		users = append(users, userID)
	}
	return users
}

// IsUserConnected 检查用户是否已连接
func (m *WebSocketManager) IsUserConnected(userID int64) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.clients[userID]
	return ok
}

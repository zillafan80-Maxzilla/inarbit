package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// APIHandler API处理器
type APIHandler struct {
	db          *Database
	authService *AuthService
	wsManager   *WebSocketManager
}

// NewAPIHandler 创建API处理器
func NewAPIHandler(db *Database, authService *AuthService, wsManager *WebSocketManager) *APIHandler {
	return &APIHandler{
		db:          db,
		authService: authService,
		wsManager:   wsManager,
	}
}

// RegisterRoutes 注册API路由
func (h *APIHandler) RegisterRoutes(router *mux.Router) {
	// 认证路由
	router.HandleFunc("/api/auth/login", h.Login).Methods("POST")
	router.HandleFunc("/api/auth/verify", h.VerifyAuth).Methods("GET")
	router.HandleFunc("/api/auth/logout", h.Logout).Methods("POST")

	// 机器人路由
	router.HandleFunc("/api/bots", h.AuthMiddleware(h.GetBots)).Methods("GET")
	router.HandleFunc("/api/bots", h.AuthMiddleware(h.CreateBot)).Methods("POST")
	router.HandleFunc("/api/bots/{id}", h.AuthMiddleware(h.GetBot)).Methods("GET")
	router.HandleFunc("/api/bots/{id}", h.AuthMiddleware(h.UpdateBot)).Methods("PUT")
	router.HandleFunc("/api/bots/{id}", h.AuthMiddleware(h.DeleteBot)).Methods("DELETE")
	router.HandleFunc("/api/bots/{id}/start", h.AuthMiddleware(h.StartBot)).Methods("POST")
	router.HandleFunc("/api/bots/{id}/stop", h.AuthMiddleware(h.StopBot)).Methods("POST")
	router.HandleFunc("/api/bots/{id}/switch-mode", h.AuthMiddleware(h.SwitchMode)).Methods("POST")

	// 仪表板路由
	router.HandleFunc("/api/dashboard/stats", h.AuthMiddleware(h.GetDashboardStats)).Methods("GET")
	router.HandleFunc("/api/dashboard/chart-data", h.AuthMiddleware(h.GetChartData)).Methods("GET")

	// 交易所路由
	router.HandleFunc("/api/exchanges", h.AuthMiddleware(h.GetExchanges)).Methods("GET")
	router.HandleFunc("/api/exchanges", h.AuthMiddleware(h.CreateExchange)).Methods("POST")
	router.HandleFunc("/api/exchanges/{id}", h.AuthMiddleware(h.DeleteExchange)).Methods("DELETE")

	// WebSocket路由
	router.HandleFunc("/ws", h.HandleWebSocket).Methods("GET")

	// 健康检查
	router.HandleFunc("/api/health", h.Health).Methods("GET")
}

// AuthMiddleware 认证中间件
func (h *APIHandler) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			h.RespondError(w, http.StatusUnauthorized, "缺少认证令牌")
			return
		}

		// 解析Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			h.RespondError(w, http.StatusUnauthorized, "无效的认证令牌格式")
			return
		}

		// 验证token
		claims, err := h.authService.VerifyToken(parts[1])
		if err != nil {
			h.RespondError(w, http.StatusUnauthorized, "认证令牌无效或已过期")
			return
		}

		// 将用户ID存储在请求上下文中
		r.Header.Set("X-User-ID", strconv.FormatInt(claims.UserID, 10))

		next(w, r)
	}
}

// GetUserID 从请求中获取用户ID
func (h *APIHandler) GetUserID(r *http.Request) (int64, error) {
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		return 0, fmt.Errorf("缺少用户ID")
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("无效的用户ID")
	}

	return userID, nil
}

// ===== 认证处理器 =====

// Login 登录
func (h *APIHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.RespondError(w, http.StatusBadRequest, "请求体格式错误")
		return
	}

	// 验证输入
	if req.Username == "" || req.Password == "" {
		h.RespondError(w, http.StatusBadRequest, "用户名和密码不能为空")
		return
	}

	// 执行登录
	token, user, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		h.RespondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	h.RespondSuccess(w, http.StatusOK, "登录成功", LoginResponse{
		Token: token,
		User:  user,
	})
}

// VerifyAuth 验证认证
func (h *APIHandler) VerifyAuth(w http.ResponseWriter, r *http.Request) {
	userID, err := h.GetUserID(r)
	if err != nil {
		h.RespondError(w, http.StatusUnauthorized, "认证验证失败")
		return
	}

	user, err := h.db.GetUserByID(userID)
	if err != nil {
		h.RespondError(w, http.StatusUnauthorized, "用户不存在")
		return
	}

	h.RespondSuccess(w, http.StatusOK, "认证有效", user)
}

// Logout 登出
func (h *APIHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// 前端会删除本地token，这里只需返回成功
	h.RespondSuccess(w, http.StatusOK, "登出成功", nil)
}

// ===== 机器人处理器 =====

// GetBots 获取机器人列表
func (h *APIHandler) GetBots(w http.ResponseWriter, r *http.Request) {
	userID, err := h.GetUserID(r)
	if err != nil {
		h.RespondError(w, http.StatusUnauthorized, "获取用户ID失败")
		return
	}

	bots, err := h.db.GetBots(userID)
	if err != nil {
		log.Printf("获取机器人列表失败: %v", err)
		h.RespondError(w, http.StatusInternalServerError, "获取机器人列表失败")
		return
	}

	if bots == nil {
		bots = make([]*Bot, 0)
	}

	h.RespondSuccess(w, http.StatusOK, "获取机器人列表成功", map[string]interface{}{
		"bots": bots,
	})
}

// GetBot 获取单个机器人
func (h *APIHandler) GetBot(w http.ResponseWriter, r *http.Request) {
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

	bot, err := h.db.GetBotByID(botID, userID)
	if err != nil {
		h.RespondError(w, http.StatusNotFound, "机器人不存在")
		return
	}

	h.RespondSuccess(w, http.StatusOK, "获取机器人成功", bot)
}

// CreateBot 创建机器人
func (h *APIHandler) CreateBot(w http.ResponseWriter, r *http.Request) {
	userID, err := h.GetUserID(r)
	if err != nil {
		h.RespondError(w, http.StatusUnauthorized, "获取用户ID失败")
		return
	}

	var req CreateBotRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.RespondError(w, http.StatusBadRequest, "请求体格式错误")
		return
	}

	// 验证输入
	if req.Name == "" || req.StrategyType == "" || req.ExchangeID == 0 {
		h.RespondError(w, http.StatusBadRequest, "缺少必要参数")
		return
	}

	// 设置默认值
	if req.UpdateFrequency == 0 {
		req.UpdateFrequency = 5 // 默认5秒
	}

	bot := &Bot{
		UserID:          userID,
		Name:            req.Name,
		StrategyType:    req.StrategyType,
		ExchangeID:      req.ExchangeID,
		IsRunning:       false,
		IsSimulation:    req.IsSimulation,
		UpdateFrequency: req.UpdateFrequency,
	}

	err = h.db.CreateBot(bot)
	if err != nil {
		log.Printf("创建机器人失败: %v", err)
		h.RespondError(w, http.StatusInternalServerError, "创建机器人失败")
		return
	}

	log.Printf("✓ 用户 %d 创建机器人 %d", userID, bot.ID)

	h.RespondSuccess(w, http.StatusCreated, "创建机器人成功", bot)
}

// UpdateBot 更新机器人
func (h *APIHandler) UpdateBot(w http.ResponseWriter, r *http.Request) {
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

	var req UpdateBotRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.RespondError(w, http.StatusBadRequest, "请求体格式错误")
		return
	}

	// 获取现有机器人
	bot, err := h.db.GetBotByID(botID, userID)
	if err != nil {
		h.RespondError(w, http.StatusNotFound, "机器人不存在")
		return
	}

	// 更新字段
	if req.Name != "" {
		bot.Name = req.Name
	}
	if req.StrategyType != "" {
		bot.StrategyType = req.StrategyType
	}
	if req.ExchangeID != 0 {
		bot.ExchangeID = req.ExchangeID
	}
	if req.UpdateFrequency > 0 {
		bot.UpdateFrequency = req.UpdateFrequency
	}
	bot.IsSimulation = req.IsSimulation

	err = h.db.UpdateBot(bot)
	if err != nil {
		log.Printf("更新机器人失败: %v", err)
		h.RespondError(w, http.StatusInternalServerError, "更新机器人失败")
		return
	}

	log.Printf("✓ 用户 %d 更新机器人 %d", userID, botID)

	h.RespondSuccess(w, http.StatusOK, "更新机器人成功", bot)
}

// DeleteBot 删除机器人
func (h *APIHandler) DeleteBot(w http.ResponseWriter, r *http.Request) {
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

	err = h.db.DeleteBot(botID, userID)
	if err != nil {
		h.RespondError(w, http.StatusInternalServerError, "删除机器人失败")
		return
	}

	log.Printf("✓ 用户 %d 删除机器人 %d", userID, botID)

	h.RespondSuccess(w, http.StatusOK, "删除机器人成功", nil)
}

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

	err = h.db.UpdateBotStatus(botID, userID, true)
	if err != nil {
		h.RespondError(w, http.StatusInternalServerError, "启动机器人失败")
		return
	}

	log.Printf("✓ 用户 %d 启动机器人 %d", userID, botID)

	// 通过WebSocket通知客户端
	h.wsManager.BroadcastBotStatus(userID, botID, true)

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

	err = h.db.UpdateBotStatus(botID, userID, false)
	if err != nil {
		h.RespondError(w, http.StatusInternalServerError, "停止机器人失败")
		return
	}

	log.Printf("✓ 用户 %d 停止机器人 %d", userID, botID)

	// 通过WebSocket通知客户端
	h.wsManager.BroadcastBotStatus(userID, botID, false)

	h.RespondSuccess(w, http.StatusOK, "机器人已停止", nil)
}

// SwitchMode 切换模式（虚拟/实盘）
func (h *APIHandler) SwitchMode(w http.ResponseWriter, r *http.Request) {
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

	var req SwitchModeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.RespondError(w, http.StatusBadRequest, "请求体格式错误")
		return
	}

	// 获取机器人
	bot, err := h.db.GetBotByID(botID, userID)
	if err != nil {
		h.RespondError(w, http.StatusNotFound, "机器人不存在")
		return
	}

	// 更新模式
	bot.IsSimulation = req.IsSimulation
	err = h.db.UpdateBot(bot)
	if err != nil {
		h.RespondError(w, http.StatusInternalServerError, "切换模式失败")
		return
	}

	mode := "实盘"
	if req.IsSimulation {
		mode = "虚拟盘"
	}

	log.Printf("✓ 用户 %d 将机器人 %d 切换为%s", userID, botID, mode)

	h.RespondSuccess(w, http.StatusOK, fmt.Sprintf("已切换为%s", mode), bot)
}

// ===== 仪表板处理器 =====

// GetDashboardStats 获取仪表板统计数据
func (h *APIHandler) GetDashboardStats(w http.ResponseWriter, r *http.Request) {
	userID, err := h.GetUserID(r)
	if err != nil {
		h.RespondError(w, http.StatusUnauthorized, "获取用户ID失败")
		return
	}

	stats, err := h.db.GetDashboardStats(userID)
	if err != nil {
		log.Printf("获取统计数据失败: %v", err)
		h.RespondError(w, http.StatusInternalServerError, "获取统计数据失败")
		return
	}

	h.RespondSuccess(w, http.StatusOK, "获取统计数据成功", stats)
}

// GetChartData 获取图表数据
func (h *APIHandler) GetChartData(w http.ResponseWriter, r *http.Request) {
	userID, err := h.GetUserID(r)
	if err != nil {
		h.RespondError(w, http.StatusUnauthorized, "获取用户ID失败")
		return
	}

	data, err := h.db.GetChartData(userID)
	if err != nil {
		log.Printf("获取图表数据失败: %v", err)
		h.RespondError(w, http.StatusInternalServerError, "获取图表数据失败")
		return
	}

	if data == nil {
		data = make([]*ChartDataPoint, 0)
	}

	h.RespondSuccess(w, http.StatusOK, "获取图表数据成功", map[string]interface{}{
		"data": data,
	})
}

// ===== 交易所处理器 =====

// GetExchanges 获取交易所列表
func (h *APIHandler) GetExchanges(w http.ResponseWriter, r *http.Request) {
	userID, err := h.GetUserID(r)
	if err != nil {
		h.RespondError(w, http.StatusUnauthorized, "获取用户ID失败")
		return
	}

	exchanges, err := h.db.GetExchanges(userID)
	if err != nil {
		log.Printf("获取交易所列表失败: %v", err)
		h.RespondError(w, http.StatusInternalServerError, "获取交易所列表失败")
		return
	}

	if exchanges == nil {
		exchanges = make([]*Exchange, 0)
	}

	h.RespondSuccess(w, http.StatusOK, "获取交易所列表成功", map[string]interface{}{
		"exchanges": exchanges,
	})
}

// CreateExchange 创建交易所
func (h *APIHandler) CreateExchange(w http.ResponseWriter, r *http.Request) {
	userID, err := h.GetUserID(r)
	if err != nil {
		h.RespondError(w, http.StatusUnauthorized, "获取用户ID失败")
		return
	}

	var req CreateExchangeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.RespondError(w, http.StatusBadRequest, "请求体格式错误")
		return
	}

	if req.Name == "" || req.APIKey == "" || req.APISecret == "" {
		h.RespondError(w, http.StatusBadRequest, "缺少必要参数")
		return
	}

	exchange := &Exchange{
		UserID:    userID,
		Name:      req.Name,
		APIKey:    req.APIKey,
		APISecret: req.APISecret,
		IsTestnet: req.IsTestnet,
		IsActive:  true,
	}

	err = h.db.CreateExchange(exchange)
	if err != nil {
		log.Printf("创建交易所失败: %v", err)
		h.RespondError(w, http.StatusInternalServerError, "创建交易所失败")
		return
	}

	log.Printf("✓ 用户 %d 创建交易所 %d", userID, exchange.ID)

	h.RespondSuccess(w, http.StatusCreated, "创建交易所成功", exchange)
}

// DeleteExchange 删除交易所
func (h *APIHandler) DeleteExchange(w http.ResponseWriter, r *http.Request) {
	userID, err := h.GetUserID(r)
	if err != nil {
		h.RespondError(w, http.StatusUnauthorized, "获取用户ID失败")
		return
	}

	vars := mux.Vars(r)
	exchangeID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		h.RespondError(w, http.StatusBadRequest, "无效的交易所ID")
		return
	}

	err = h.db.DeleteExchange(exchangeID, userID)
	if err != nil {
		h.RespondError(w, http.StatusInternalServerError, "删除交易所失败")
		return
	}

	log.Printf("✓ 用户 %d 删除交易所 %d", userID, exchangeID)

	h.RespondSuccess(w, http.StatusOK, "删除交易所成功", nil)
}

// ===== WebSocket处理器 =====

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// 在生产环境中，应该检查Origin header
		return true
	},
}

// HandleWebSocket 处理WebSocket连接
func (h *APIHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket升级失败: %v", err)
		return
	}

	h.wsManager.HandleConnection(conn)
}

// ===== 健康检查 =====

// Health 健康检查
func (h *APIHandler) Health(w http.ResponseWriter, r *http.Request) {
	h.RespondSuccess(w, http.StatusOK, "服务正常", map[string]interface{}{
		"status": "ok",
	})
}

// ===== 响应辅助函数 =====

// RespondSuccess 返回成功响应
func (h *APIHandler) RespondSuccess(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := APIResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	}

	json.NewEncoder(w).Encode(response)
}

// RespondError 返回错误响应
func (h *APIHandler) RespondError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := APIResponse{
		Status: "error",
		Error:  message,
	}

	json.NewEncoder(w).Encode(response)
}

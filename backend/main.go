package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// 加载环境变量
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 创建路由
	router := mux.NewRouter()

	// API路由
	router.HandleFunc("/api/health", healthHandler).Methods("GET")
	router.HandleFunc("/api/login", loginHandler).Methods("POST")
	router.HandleFunc("/api/bots", getBotsHandler).Methods("GET")
	router.HandleFunc("/api/bots", createBotHandler).Methods("POST")
	router.HandleFunc("/api/bots/{id}", updateBotHandler).Methods("PUT")
	router.HandleFunc("/api/bots/{id}", deleteBotHandler).Methods("DELETE")

	// 静态文件服务
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("../frontend/dist")))

	// 启动服务器
	log.Printf("iNarbit服务器启动在 http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

// 健康检查
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status":"ok"}`)
}

// 登录
func loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"token":"test-token"}`)
}

// 获取机器人列表
func getBotsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"bots":[]}`)
}

// 创建机器人
func createBotHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"id":"1","status":"created"}`)
}

// 更新机器人
func updateBotHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"id":"1","status":"updated"}`)
}

// 删除机器人
func deleteBotHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status":"deleted"}`)
}

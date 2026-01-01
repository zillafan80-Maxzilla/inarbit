package main

import (
	"fmt"
	"os"
	"strconv"
)

// Config 应用配置
type Config struct {
	// 服务器配置
	ServerHost string
	ServerPort string
	Env        string

	// 数据库配置
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// JWT配置
	JWTSecret string

	// CORS配置
	CORSAllowedOrigins string

	// 日志配置
	LogLevel string
}

// LoadConfig 加载配置
func LoadConfig() *Config {
	config := &Config{
		// 服务器配置
		ServerHost: getEnv("SERVER_HOST", "0.0.0.0"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
		Env:        getEnv("ENV", "development"),

		// 数据库配置
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "inarbit"),
		DBPassword: getEnv("DB_PASSWORD", "inarbit_password"),
		DBName:     getEnv("DB_NAME", "inarbit"),

		// JWT配置
		JWTSecret: getEnv("JWT_SECRET", "your-secret-key-change-in-production"),

		// CORS配置
		CORSAllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:80"),

		// 日志配置
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}

	return config
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvInt 获取整数环境变量
func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

// getEnvBool 获取布尔环境变量
func getEnvBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		return value == "true" || value == "1" || value == "yes"
	}
	return defaultValue
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.DBHost == "" {
		return fmt.Errorf("数据库主机未配置")
	}
	if c.DBUser == "" {
		return fmt.Errorf("数据库用户未配置")
	}
	if c.DBName == "" {
		return fmt.Errorf("数据库名称未配置")
	}
	if c.JWTSecret == "" {
		return fmt.Errorf("JWT密钥未配置")
	}
	return nil
}

// String 返回配置的字符串表示
func (c *Config) String() string {
	return fmt.Sprintf(`
配置信息:
  环境: %s
  服务器: %s:%s
  数据库: %s:%s/%s
  日志级别: %s
	`, c.Env, c.ServerHost, c.ServerPort, c.DBHost, c.DBPort, c.DBName, c.LogLevel)
}

package config

import (
	"os"
)

// Config 全局配置
type Config struct {
	// 腾讯云配置（内容安全审核 TMS/IMS 使用）
	TencentCloudSecretID  string
	TencentCloudSecretKey string
}

var globalConfig *Config

// InitConfig 初始化配置
func InitConfig() *Config {
	globalConfig = &Config{
		TencentCloudSecretID:  getEnv("TENCENTCLOUD_SECRET_ID", ""),
		TencentCloudSecretKey: getEnv("TENCENTCLOUD_SECRET_KEY", ""),
	}
	return globalConfig
}

// GetConfig 获取全局配置
func GetConfig() *Config {
	if globalConfig == nil {
		return InitConfig()
	}
	return globalConfig
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Code scaffolded by goctl. Safe to edit.
// goctl {{.version}}

package config

import (
	gobasecfg "github.com/qkja/gobase/config"
	"github.com/qkja/go-zero/rest"
)

// AppCfg 当前配置快照持有者（由入口调用 gobase Init 填充，支持热加载）。
// 加载/热加载机制全在 gobase（config.Init），这里只声明全局变量供各模块读取。
var AppCfg *gobasecfg.Config[Config]

// Config 统一配置结构体。
// 所有字段显式 toml tag。
type Config struct {
	Application ApplicationConfig `toml:"application"`
	Grpc        GrpcConfig        `toml:"grpc"`
	Logger      LoggerConfig      `toml:"logger"`
	Rest        rest.RestConf     `toml:"rest"`
	RateLimit   RateLimitConfig   `toml:"rateLimit"`
	{{.auth}}
	{{.jwtTrans}}
}

// LoggerConfig 日志配置（仅用于 toml 映射，实际初始化由 logger.InitLog 完成）。
type LoggerConfig struct {
	Level string `toml:"level"`
	Dir   string `toml:"dir"`
}

// GrpcConfig gRPC 客户端配置。
type GrpcConfig struct {
	// Timeout 单次 gRPC 调用超时（毫秒），默认 5000。
	Timeout int64 `toml:"timeout"`
}

type ApplicationConfig struct {
	Name string `toml:"name"`
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

// RateLimitConfig 全局限流配置。接口级覆盖路由配置见 config/ratelimit.json。
type RateLimitConfig struct {
	Enable bool `toml:"enable"`
	Rate   int  `toml:"rate"`
	Burst  int  `toml:"burst"`
}

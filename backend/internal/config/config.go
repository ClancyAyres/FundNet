package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config 应用配置
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	App      AppConfig      `yaml:"app"`
	Scraper  ScraperConfig  `yaml:"scraper"`
	CORS     CORSConfig     `yaml:"cors"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Path string `yaml:"path"`
}

// AppConfig 应用配置
type AppConfig struct {
	RefreshInterval int    `yaml:"refresh_interval"`
	LogLevel        string `yaml:"log_level"`
}

// ScraperConfig 爬虫配置
type ScraperConfig struct {
	Timeout    int `yaml:"timeout"`
	RetryCount int `yaml:"retry_count"`
}

// CORSConfig CORS配置
type CORSConfig struct {
	AllowedOrigins []string `yaml:"allowed_origins"`
	AllowedMethods []string `yaml:"allowed_methods"`
}

var cfg *Config

// Load 加载配置文件
func Load() (*Config, error) {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		// 尝试读取上一级目录的配置文件
		data, err = os.ReadFile("../config.yaml")
		if err != nil {
			return nil, err
		}
	}

	cfg = &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	// 设置默认值
	if cfg.Server.Port == 0 {
		cfg.Server.Port = 3800
	}
	if cfg.Server.Host == "" {
		cfg.Server.Host = "0.0.0.0"
	}
	if cfg.App.RefreshInterval == 0 {
		cfg.App.RefreshInterval = 60
	}

	return cfg, nil
}

// Get 获取当前配置
func Get() *Config {
	return cfg
}

package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

// Config 存储配置文件的结构
type Config struct {
	ConsistencyAlgorithm string `yaml:"consistency_algorithm"`
	Gossip               struct {
		Port   int      `yaml:"port"`
		Peers  []string `yaml:"peers"`
		NodeID string   `yaml:"node_id"`
	} `yaml:"gossip"`
	Database struct {
		Type     string `yaml:"type"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"db_name"`
	} `yaml:"database"`
}

// LoadConfig 从配置文件中加载配置
func LoadConfig(filePath string) (*Config, error) {
	// 打开配置文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("无法打开配置文件: %v", err)
	}
	defer file.Close()

	// 初始化配置结构体
	config := &Config{}

	// 解码 YAML 配置文件
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	// 可选：检查必填项是否存在或有效
	if config.Database.Host == "" || config.Database.Port == 0 {
		return nil, fmt.Errorf("数据库配置缺失必需字段")
	}

	// 返回配置对象
	return config, nil
}

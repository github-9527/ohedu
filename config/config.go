package config

import (
	"github.com/BurntSushi/toml"
)

// LoggerConfig 日志的配置文件结构体
type LoggerConfig struct {
	Level      string `toml:"level"`       // 日志级别
	File       string `toml:"file"`        // 日志路径，包含文件名，例如：./logs/info.log
	MaxSize    int    `toml:"max_size"`    // 文件大小限制，单位MB
	MaxBackups int    `toml:"max_backups"` // 最大保留日志文件数量
	Compress   bool   `toml:"compress"`    // 是否压缩处理
}

type SystemConfig struct {
	Logger LoggerConfig `toml:"logger"`
	MySQL  struct {
		DBSource            string `toml:"db_source"`
		PullLogTableName    string `toml:"pull_log_table_name"`
		StatisticsTableName string `toml:"statistics_table_name"`
	} `toml:"mysql"`
	Oracle struct {
		DBSource         string `toml:"db_source"`
		PreDistTableName string `toml:"pre_dist_table_name"`
	} `toml:"oracle"`
	Regions []struct {
		File string `toml:"file"`
	} `toml:"regions"`
}

var Config SystemConfig

//输入路径如 config/config.toml
func InitConfig(path string) error {

	_, err := toml.DecodeFile(path, &Config)
	if err != nil {
		return err
	}
	return nil
}

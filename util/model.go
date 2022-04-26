package util

import (
	"fmt"
	"time"

	"gopkg.in/ini.v1"
)

// LoggerConfig 日志的配置文件结构体
type LoggerConfig struct {
	Level      string `toml:"level"`       // 日志级别
	File       string `toml:"file"`        // 日志路径，包含文件名，例如：./logs/info.log
	MaxSize    int    `toml:"max_size"`    // 文件大小限制，单位MB
	MaxBackups int    `toml:"max_backups"` // 最大保留日志文件数量
	Compress   bool   `toml:"compress"`    // 是否压缩处理
}

type Config struct {
	Logger LoggerConfig `toml:"logger"`
	MySQL  struct {
		DBSource            string `toml:"db_source"`
		PullLogTableName    string `toml:"pull_log_table_name"`
		StatisticsTableName string `toml:"statistics_table_name"`
		ErrorTableName      string `toml:"error_table_name"`
	} `toml:"mysql"`
	Oracle struct {
		DBSource         string `toml:"db_source"`
		PreDistTableName string `toml:"pre_dist_table_name"`
	} `toml:"oracle"`
	Regions []struct {
		File string `toml:"file"`
	} `toml:"regions"`
}

type Region struct {
	Region     string `toml:"region"`
	Concurrent int    `toml:"concurrent"`
	Interval   int    `toml:"interval"`
	BackTime   int    `toml:"back_time"`
	CornSpec   string `toml:"corn_spec"`
	Disable    bool   `toml:"disable"`
	Source     []struct {
		SubRegion    string      `toml:"sub_region"`
		DBName       string      `toml:"db_name"`
		Disable      bool        `toml:"disable"`
		MainSource   string      `toml:"main_source"`
		BackupSource string      `toml:"backup_source"`
		Tables       []TableFile `toml:"tables"`
	} `toml:"source"`
}

type TableFile struct {
	File      string    `toml:"file"`
	TableName string    `toml:"table_name"`
	StartTime time.Time `toml:"start_time"`
	Patch     Table     `toml:"patch"`
}

type Table struct {
	// TODO: 使用拼接前缀
	TableName             string `toml:"table_name"`
	TimeCol               string `toml:"time_col"`
	TimeCamelCase         string `toml:"time_camel_case"`
	OrganizationCol       string `toml:"organization_col"`
	OrganizationColCase   string `toml:"organization_case"`
	OrganizationCamelCase string `toml:"organization_camel_case"`

	SrcTableName string    `toml:"src_table_name"`
	StartTime    time.Time `toml:"start_time"`

	Columns []struct {
		Name          string `toml:"name"`
		NameCase      string `toml:"name_case"`
		NameCamelCase string `toml:"name_camel_case"`
		TypeName      string `toml:"type_name"`
		Length        int    `toml:"length"`
		Precision     int    `toml:"precision"`
		Scale         int    `toml:"scale"`
		NullAble      bool   `toml:"null_able"`
		PrimaryKey    bool   `toml:"primary_key"`
	} `toml:"columns"`
}

var Cfg *ini.File

func init() {
	var err error
	Cfg, err = ini.Load("./config.ini")
	if err != nil {
		fmt.Println(err)
		return
	}
}

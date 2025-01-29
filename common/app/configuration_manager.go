package app

import "github.com/denizdoganinsider/kpi_project/common/mysql"

type ConfigurationManager struct {
	MySqlConfig mysql.Config
}

func NewConfigurationManager() *ConfigurationManager {
	MySqlConfig := getMySqlConfig()
	return &ConfigurationManager{
		MySqlConfig: MySqlConfig,
	}
}

func getMySqlConfig() mysql.Config {
	return mysql.Config{
		Host:                  "localhost",
		Port:                  "3306",
		UserName:              "root",
		Password:              "root",
		DbName:                "kpidb",
		MaxConnections:        "10",
		MaxConnectionIdleTime: "30s",
	}
}

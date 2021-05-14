package config

type MysqlConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	DBName   string `mapstructure:"dbname"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type DatabaseConfig struct {
	MysqlConfig MysqlConfig `mapstructure:"mysql"`
}

type ServiceConfig struct {
	DatabaseConfig DatabaseConfig `mapstructure:"database"`
}

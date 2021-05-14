package config

type UserSrvConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type JWTConfig struct {
	SigningKey string `mapstructure:"key"`
}

type SMTPServerConfig struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
}

type RedisConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type ServiceConfig struct {
	Name             string           `mapstructure:"name"`
	Port             int              `mapstructure:"port"`
	UserSrvConfig    UserSrvConfig    `mapstructure:"user_srv"`
	JWTInfo          JWTConfig        `mapstructure:"jwt"`
	SMTPServerConfig SMTPServerConfig `mapstructure:"smtp_server"`
	RedisConfig      RedisConfig      `mapstructure:"redis"`
}

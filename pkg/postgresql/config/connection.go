package config

// DBConfig представляет конфигурацию для подключения к базе данных PostgreSQL
type DBConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
	SSLMode  string `json:"sslmode"`
}

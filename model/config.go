package model

type Config struct {
	Database DatabaseConfig `json:"database"`
}

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Name     string `json:"name"`
	User     string `json:"user"`
	Database string `json:"database"`
	Password string `json:"password"`
}
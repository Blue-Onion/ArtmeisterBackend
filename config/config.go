package config

import (
	"database/sql"
	"errors"
	"os"

	"github.com/Blue-Onion/RestApi-Go/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Config struct {
	DbUrl     string
	Port      string
	JWTSecert string
}
type ApiConfig struct {
	UserRepo database.UserRepository
}

func LoadConfig() *Config {
	godotenv.Load()
	dbUrl := os.Getenv("DATABASE_URL")
	Port := os.Getenv("PORT")
	Jwt := os.Getenv("JWT_SECERT")

	return &Config{
		DbUrl:     dbUrl,
		Port:      Port,
		JWTSecert: Jwt,
	}

}
func DbQuries() (*ApiConfig, error) {
	apiConfig := &ApiConfig{}
	config := LoadConfig()
	conn, err := sql.Open("postgres", config.DbUrl)
	if err != nil {
		return nil, err
	}
	query := database.New(conn)
	if query == nil {
		return nil, errors.New("Connection Failed")
	}
	apiConfig.UserRepo = query
	return apiConfig, nil
}

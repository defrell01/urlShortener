package configs

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Db    DbConfig
	Auth  AuthConfig
	Redis ReadisConfig
}

type AuthConfig struct {
	Secret string
}

type DbConfig struct {
	Dsn string
}

type ReadisConfig struct {
	Addr     string
	Password string
	DB       uint32
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using default config")
	}

	redisDb, err := strconv.ParseUint(os.Getenv("REDIS_DB"), 10, 32)
	if err != nil {
		log.Println("Error parsing redis db number")
	}

	return &Config{
		Db: DbConfig{
			Dsn: os.Getenv("DSN"),
		},
		Auth: AuthConfig{
			Secret: os.Getenv("SECRET"),
		},
		Redis: ReadisConfig{
			Addr:     os.Getenv("REDIS_ADDR"),
			Password: os.Getenv("REDIS_PWD"),
			DB:       uint32(redisDb),
		},
	}
}

package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"time"
)

type Config struct {
	Server      Server
	Debug       Debug
	Selectel    Selectel
	Redis       Redis
	PostgresUrl string
	PostgresCtx time.Duration
}

type Server struct {
	Host string `env:"SERVER_HOST" env-default:"0.0.0.0"`
	Port int    `env:"SERVER_PORT" env-default:"8082"`
}

type Selectel struct {
	AccessKey  string `env:"SELECTEL_ACCESS_KEY" env-default:""`
	SecretKey  string `env:"SELECTEL_SECRET_KEY" env-default:""`
	BucketName string `env:"SELECTEL_BUCKET_NAME" env-default:""`
	Region     string `env:"SELECTEL_REGION" env-default:""`
	Endpoint   string `env:"SELECTEL_ENDPOINT" env-default:""`
}

type Redis struct {
	Address      string        `env:"REDIS_ADDRESS" env-default:""`
	Password     string        `env:"REDIS_PASSWORD" env-default:""`
	DB           int           `env:"REDIS_DB"`
	CtxTime      time.Duration `env:"REDIS_CTX" env-default:"5s"`
	DialTimeout  time.Duration `env:"REDIS_DIAL_TIMEOUT" env-default:""`
	ReadTimeout  time.Duration `env:"REDIS_READ_TIMEOUT" env-default:""`
	WriteTimeout time.Duration `env:"REDIS_WRITE_TIMEOUT" env-default:""`
}

type Debug struct {
	DebugLogger bool `env:"DEBUG_LOGGER" env-default:"false"`
}
type Postgres struct {
	Host     string `env:"PG_HOST"`
	Port     int    `env:"PG_PORT"`
	Dbname   string `env:"PG_DBNAME"`
	User     string `env:"PG_USER"`
	Password string `env:"PG_USER_PASS"`
}

func (s *Config) GetPostgresUrl(p *Postgres) {
	s.PostgresUrl = fmt.Sprintf("host=%s dbname=%s user=%s password=%s port=%d", p.Host, p.Dbname, p.User, p.Password, p.Port)
}

func NewConfig() *Config {
	cfg := new(Config)
	postgres := new(Postgres)
	err := godotenv.Load()
	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		log.Fatal("Error loading config: ", err)
	}
	err = cleanenv.ReadEnv(postgres)
	if err != nil {
		log.Fatal("Error loading Postgres data: ", err)
	}
	cfg.GetPostgresUrl(postgres)
	return cfg
}

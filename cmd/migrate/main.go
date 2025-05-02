package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
	"log"
)

type Migration struct {
	MigrationPath string `env:"MIGRATION_PATH"`
	PGDBNAME      string `env:"PG_DBNAME"`
	PGUser        string `env:"PG_USER"`
	PGPassword    string `env:"PG_USER_PASS"`
	PGHost        string `env:"PG_HOST"`
	PGPort        string `env:"PG_PORT"`
}

const (
	driver = "pgx"
)

func main() {
	migration := new(Migration)
	var command string
	flag.StringVar(&command, "command", "up", "command for migration")
	flag.Parse()
	err := godotenv.Load("C:\\Golang_social_project\\VK_posts\\.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	err = cleanenv.ReadEnv(migration)
	if err != nil {
		log.Fatal("Problem with check env")
	}
	ctx := context.Background()
	db, err := goose.OpenDBWithDriver(driver, fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", migration.PGHost, migration.PGPort, migration.PGUser, migration.PGPassword, migration.PGDBNAME))
	if err != nil {
		log.Fatal(err)
	}
	err = goose.RunContext(ctx, command, db, migration.MigrationPath)
	if err != nil {
		log.Fatal(err.Error())
	}
}

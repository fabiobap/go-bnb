package dbrepo

import (
	"database/sql"

	"github.com/fabiobap/go-bnb/internal/config"
	"github.com/fabiobap/go-bnb/internal/repository"
)

type postgrsDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgrsDBRepo{
		App: a,
		DB:  conn,
	}
}

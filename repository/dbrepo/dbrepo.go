package dbrepo

import (
	"database/sql"
	"github.com/zahnah/study-app/internal/config"
	"github.com/zahnah/study-app/repository"
)

type postgresDbRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(db *sql.DB, app *config.AppConfig) repository.DatabaseRepo {
	return &postgresDbRepo{
		App: app,
		DB:  db,
	}
}

package dbrepo

import (
	"database/sql"
	"github.com/zahnah/study-app/internal/config"
	"github.com/zahnah/study-app/repository"
)

func NewPostgresRepo(db *sql.DB, app *config.AppConfig) repository.DatabaseRepo {
	return &postgresDbRepo{
		App: app,
		DB:  db,
	}
}

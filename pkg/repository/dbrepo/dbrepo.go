package dbrepo

import (
	"database/sql"

	"github.com/ptamarov/go-cards/pkg/config"
	"github.com/ptamarov/go-cards/pkg/repository"
)

// All structs that satisfy DatabaseRepository go here

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepository {
	return &postgresDBRepo{App: a, DB: conn}
}

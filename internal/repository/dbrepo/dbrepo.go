package dbrepo

import (
	"database/sql"

	"github.com/ptamarov/go-cards/internal/config"
	"github.com/ptamarov/go-cards/internal/repository"
)

// All structs that satisfy DatabaseRepository go here

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

type testDBRepo struct {
	App *config.AppConfig
}

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepository {
	return &postgresDBRepo{App: a, DB: conn}
}

func NewTestingRepo(a *config.AppConfig) repository.DatabaseRepository {
	return &testDBRepo{App: a}
}

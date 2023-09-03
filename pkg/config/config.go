package config

import (
	"database/sql"
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
	"github.com/google/uuid"
	"github.com/ptamarov/go-cards/app/card"
)

type AppConfig struct {
	UseCache        bool
	InProduction    bool
	TemplateCache   map[string]*template.Template
	InfoLog         *log.Logger
	ErrorLog        *log.Logger
	Session         *scs.SessionManager
	DataBase        *sql.DB
	UserID          uuid.UUID
	UserData        UserCache
	GetTranslations bool
}

type UserCache struct {
	UserName       string
	AuthKey        string
	TargetLanguage string
	Progress       int
	CardData       card.MemoryCard
	Count          int
}

package config

import (
	"database/sql"
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
	"github.com/google/uuid"
)

type AppConfig struct {
	UseCache        bool
	InProduction    bool
	TemplateCache   map[string]*template.Template
	InfoLog         *log.Logger
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
	CardData       map[string]string
	Count          int
}

package config

import (
	"html/template"
	"log"
	"time"

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
	UserData        UserCache
	GetTranslations bool
	NotFresh        bool
	Time            time.Time
}

type UserCache struct {
	UserID         uuid.UUID
	UserName       string
	AuthKey        string
	TargetLanguage string
	LastAnswer     string
	Progress       int
	CardData       card.MemoryCard
	Count          int
}

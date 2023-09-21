package config

import (
	"html/template"
	"log"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/ptamarov/go-cards/app/user"
)

type AppConfig struct {
	UseCache        bool
	InProduction    bool
	TemplateCache   map[string]*template.Template
	InfoLog         *log.Logger
	ErrorLog        *log.Logger
	Session         *scs.SessionManager
	User            user.User
	GetTranslations bool
	NotFresh        bool
	Time            time.Time
	AuthKey         string
	AlreadyAnswered bool
}

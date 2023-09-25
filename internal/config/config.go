package config

import (
	"html/template"
	"log"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/ptamarov/go-cards/app/user"
	"github.com/ptamarov/go-cards/internal/models"
)

type AppConfig struct {
	UseTemplateCache bool
	InProduction     bool
	TemplateCache    map[string]*template.Template
	HomeDataCache    *models.TemplateData
	SummaryDataCache *models.TemplateData
	CorrectCache     *models.TemplateData
	InfoLog          *log.Logger
	ErrorLog         *log.Logger
	Session          *scs.SessionManager
	User             user.User
	GetTranslations  bool
	NotFresh         bool
	Time             time.Time
	AuthKey          string
	AlreadyAnswered  bool
}

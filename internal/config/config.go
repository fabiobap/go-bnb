package config

import (
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
	"github.com/fabiobap/go-bnb/internal/models"
)

// holds application cache
type AppConfig struct {
	TemplateCache map[string]*template.Template
	UseCache      bool
	InProduction  bool
	Session       *scs.SessionManager
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	MailChan      chan models.MailData
}

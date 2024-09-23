package config

import (
	"html/template"

	"github.com/alexedwards/scs/v2"
)

// holds application cache
type AppConfig struct {
	TemplateCache map[string]*template.Template
	UseCache      bool
	InProduction  bool
	Session       *scs.SessionManager
}

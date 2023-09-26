package configs

import (
	"os"
	"strings"
	"sync"
)

type AppConfig struct {
	name    string
	version string
	env     string
	host    string
	schemes []string
}

var (
	app     *AppConfig
	appOnce sync.Once
)

func (ac *AppConfig) Name() string {
	return ac.name
}

func (ac *AppConfig) Version() string {
	return ac.version
}

func (ac *AppConfig) Env() string {
	return ac.env
}

func (ac *AppConfig) Host() string {
	return ac.host
}

func (ac *AppConfig) Schemes() []string {
	return ac.schemes
}

func App() *AppConfig {
	appOnce.Do(func() {
		app = &AppConfig{
			name:    "mcash-finance-notification",
			version: "1.0.0",
			env:     os.Getenv("ENV"),
			host:    os.Getenv("HOST"),
		}

		trimCsv := strings.TrimSpace(os.Getenv("SCHEMES"))
		lowerCsv := strings.ToLower(trimCsv)
		splittedSchemes := strings.Split(lowerCsv, ",")
		appSchemes := UniqueStrings(splittedSchemes)

		for _, appScheme := range appSchemes {
			if appScheme == "http" || appScheme == "https" {
				app.schemes = append(app.schemes, appScheme)
			}
		}

		if len(app.schemes) == 0 {
			app.schemes = []string{"http"}
		}
	})
	return app
}

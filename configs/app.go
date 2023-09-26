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
	storage string
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

func (ac *AppConfig) StoragePath() string {
	return ac.storage
}

func App() *AppConfig {
	appOnce.Do(func() {
		app = &AppConfig{
			name:    "mcash-finance-consolidation-core",
			version: "1.0.0",
			env:     os.Getenv("ENV"),
			host:    os.Getenv("HOST"),
			storage: os.Getenv("STORAGE_DIRECTORY_PATH"),
		}

		trimCsv := strings.TrimSpace(os.Getenv("SCHEMES"))
		lowerCsv := strings.ToLower(trimCsv)
		splitedSchemes := strings.Split(lowerCsv, ",")
		appSchemes := UniqueStrings(splitedSchemes)

		for _, appScheme := range appSchemes {
			if appScheme == "http" || appScheme == "https" {
				app.schemes = append(app.schemes, appScheme)
			}
		}

		if len(app.schemes) <= 0 {
			app.schemes = []string{"http"}
		}

	})
	return app
}

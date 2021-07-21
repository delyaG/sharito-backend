package configs

import (
	"backend/internal/infra/http"
	"backend/internal/infra/postgres"
	"backend/internal/infra/security"
	"backend/pkg/logging"
	"github.com/jessevdk/go-flags"
	"os"
)

type Config struct {
	Logger   *logging.Config  `group:"Logger args" namespace:"logger" env-namespace:"SHARITO_LOGGER"`
	HTTP     *http.Config     `group:"HTTP args" namespace:"http" env-namespace:"SHARITO_HTTP"`
	Postgres *postgres.Config `group:"Postgres args" namespace:"postgres" env-namespace:"SHARITO_POSTGRES"`
	Security *security.Config `group:"Security args" namespace:"security" env-namespace:"SHARITO_SECURITY"`
}

func Parse() (*Config, error) {
	var config Config
	p := flags.NewParser(&config, flags.HelpFlag|flags.PassDoubleDash)

	_, err := p.ParseArgs(os.Args)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

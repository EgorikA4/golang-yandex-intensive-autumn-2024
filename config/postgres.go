package config

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/EgorikA4/golang-yandex-intensive-autumn-2024/pkg/utils"
)

type PostgresConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
}

var (
	postgresCfg  *PostgresConfig
	postgresOnce sync.Once
)

func LoadPostgresConfig() error {
	var err error
	postgresOnce.Do(func() {
		requiredVars := []string{
			"PG_HOST",
			"PG_PASSWORD",
			"PG_USERNAME",
			"PG_DBNAME",
		}

		if err = utils.CheckEnvVars(requiredVars); err != nil {
			return
		}

		postgresCfg = &PostgresConfig{
			Host:     os.Getenv("PG_HOST"),
			Username: os.Getenv("PG_USERNAME"),
			Password: os.Getenv("PG_PASSWORD"),
			DBName:   os.Getenv("PG_DBNAME"),
		}

		postgresPort := os.Getenv("PG_PORT")
		if _, err := strconv.Atoi(postgresPort); err != nil {
			err = fmt.Errorf("PG_PORT should be integer")
			return
		}
		postgresCfg.Port = postgresPort
	})
	return err
}

package storage

import (
	"context"
	"fmt"
	"sync"

	"github.com/EgorikA4/golang-yandex-intensive-autumn-2024/config"
    "github.com/EgorikA4/golang-yandex-intensive-autumn-2024/internal/consts/db_queries"
    "github.com/jackc/pgx/v5"
)

type PostgresDatabase struct {
    Conn *pgx.Conn
    Ctx context.Context
}

var (
    postgresInstance *PostgresDatabase
    postgresInitOnce sync.Once
)

func InitPostgresDatabase() (*PostgresDatabase, error) {
    var err error
    postgresInitOnce.Do(func() {
        cfg := config.GetPostgresConfig()
        url := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

        postgresInstance = &PostgresDatabase{
            Ctx: context.Background(),
        }

        var conn *pgx.Conn
        conn, err = pgx.Connect(
            postgresInstance.Ctx,
            url,
        )
        if err != nil {
            return 
        }

        postgresInstance.Conn = conn

        _, err = conn.Exec(postgresInstance.Ctx, db_queries.CREATE_EXPRESSION_TABLE)
        if err != nil {
            return
        }
    })
    return postgresInstance, err
}

func (pg *PostgresDatabase) Close() {
    if pg.Conn != nil {
        pg.Conn.Close(pg.Ctx)
    }
}

func GetPostgresInstance() *PostgresDatabase {
    return postgresInstance
}

package storage

import (
	"context"
	"fmt"
	"sync"

	"github.com/EgorikA4/golang-yandex-intensive-autumn-2024/config"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type MemgraphDatabase struct {
    Driver neo4j.DriverWithContext
    Ctx context.Context
}

var (
    memgraphInstance *MemgraphDatabase
    memgraphInitOnce sync.Once
)

func InitMemgraphDatabase() (*MemgraphDatabase, error) {
    var err error
    memgraphInitOnce.Do(func() {
        cfg := config.GetMemgraphConfig()
        uri := fmt.Sprintf("bolt://%s:%s", cfg.Host, cfg.Port)

        var driver neo4j.DriverWithContext
        driver, err = neo4j.NewDriverWithContext(
            uri,
            neo4j.BasicAuth(cfg.Username, cfg.Password, ""),
        )
        if err != nil {
            return 
        }

        memgraphInstance = &MemgraphDatabase{
            Driver: driver,
            Ctx: context.Background(),
        }

        // err = memgraphInstance.Driver.VerifyConnectivity(memgraphInstance.Ctx)
        // if err != nil {
        //    return
        // }
    })
    return memgraphInstance, err
}

func (mdb *MemgraphDatabase) Session() neo4j.SessionWithContext {
    return mdb.Driver.NewSession(mdb.Ctx, neo4j.SessionConfig{})
}

func (mdb *MemgraphDatabase) Close() {
    if mdb.Driver != nil {
        mdb.Driver.Close(mdb.Ctx)
    }
}

func GetDBInstance() *MemgraphDatabase {
    return memgraphInstance
}

package postgres

import (
	"context"
	"fmt"
	"github.com/dvdxa/add-to-favorites/internal/configs"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectToPostgres(cfg *configs.Config) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", cfg.Host, cfg.Port, cfg.Username, cfg.Dbname, cfg.Password, cfg.SSLMode)
	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, err
	}
	//b, err := os.ReadFile("./internal/database/schema/migrations.up.sql")
	//if err != nil {
	//	return nil, err
	//}
	//migrationQuery := string(b)
	//_, err = pool.Exec(context.Background(), migrationQuery)
	//if err != nil {
	//	return nil, err
	//}

	return pool, nil
}

package repos

// Для взаимодействия с Базой данных

import (
	"ProjectGo/backend/internal/config"
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ConnRepo struct {
	conn *pgxpool.Pool
}

func ConnUrlRepos(cfg *config.Config) *ConnRepo {
	pool, err := pgxpool.New(context.TODO(), fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName))
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}
	return &ConnRepo{conn: pool}
}

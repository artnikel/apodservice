package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/artnikel/apodservice/internal/config"
	"github.com/artnikel/apodservice/internal/repository"
	"github.com/artnikel/apodservice/internal/worker"
	"github.com/jackc/pgx/v5/pgxpool"
)

func connectPostgres(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	cfgPostgres, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}
	dbpool, err := pgxpool.NewWithConfig(ctx, cfgPostgres)
	if err != nil {
		return nil, err
	}
	return dbpool, nil
}

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	dbpool, errPool := connectPostgres(ctx, cfg.ConnectionString)
	if errPool != nil {
		log.Fatalf("could not construct the pool: %v", errPool)
	}
	pgclient := repository.NewPgClient(dbpool)
	defer dbpool.Close()
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			apod, err := worker.GetApod(cfg)
			if err != nil {
				fmt.Printf("error recieve apod: %v\n", err)
			} else {
				err := pgclient.ApodCreate(ctx, apod)
				if err != nil {
					fmt.Printf("error saving apod in database: %v\n", err)
				} 
			}
		}
	}
}

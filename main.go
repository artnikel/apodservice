// Package main is an entry point to application
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/artnikel/apodservice/internal/config"
	"github.com/artnikel/apodservice/internal/constants"
	"github.com/artnikel/apodservice/internal/handler"
	"github.com/artnikel/apodservice/internal/repository"
	"github.com/artnikel/apodservice/internal/service"
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
		log.Fatalf("failed to parse config %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	dbpool, errPool := connectPostgres(ctx, cfg.ConnectionString)
	if errPool != nil {
		// nolint gocritic
		log.Fatalf("could not construct the pool: %v", errPool)
	}
	pgclient := repository.NewPgClient(dbpool)
	defer dbpool.Close()
	go func() {
		ticker := time.NewTicker(constants.WorkFrequency)
		defer ticker.Stop()
		for {
			for range ticker.C {
				apod, err := worker.GetApod(cfg)
				if err != nil {
					fmt.Printf("error receiving apod: %v\n", err)
				} else {
					err := pgclient.ApodCreate(ctx, apod)
					if err != nil {
						fmt.Printf("error saving apod in the database: %v\n", err)
					}
				}
			}
		}
	}()
	apodSvc := service.NewApodService(pgclient)
	apodHndl := handler.NewApodHandler(apodSvc)

	mux := http.NewServeMux()
	mux.HandleFunc("/list", apodHndl.GetAll)
	mux.HandleFunc("/today", apodHndl.GetToday)
	mux.HandleFunc("/bydate", apodHndl.GetByDate)

	fs := http.FileServer(http.Dir("./storage"))
	mux.Handle("/storage/", http.StripPrefix("/storage/", fs))

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  constants.RWTimeout,
		WriteTimeout: constants.RWTimeout,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("failed to start http server %v", err)
	}
}

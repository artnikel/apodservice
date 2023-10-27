// Package main is an entry point to application
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/artnikel/apodservice/docs"
	"github.com/artnikel/apodservice/internal/config"
	"github.com/artnikel/apodservice/internal/constants"
	"github.com/artnikel/apodservice/internal/handler"
	"github.com/artnikel/apodservice/internal/repository"
	"github.com/artnikel/apodservice/internal/service"
	"github.com/artnikel/apodservice/internal/worker"
	"github.com/jackc/pgx/v5/pgxpool"
	httpSwagger "github.com/swaggo/http-swagger"
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

// @title NASA APOD API
// @version 1.0
// @description API with methods for getting APOD.

// nolint funlen
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
		currentDate := time.Now().UTC().Format(constants.DateLayout)
		parsedCurrentDate, err := time.Parse(constants.DateLayout, currentDate)
		if err != nil {
			fmt.Printf("error parsing current date %v\n", err)
		}
		apod, err := pgclient.ApodGetByDate(ctx, parsedCurrentDate)
		if err != nil {
			fmt.Printf("apodGetByDate %v\n", err)
		}
		if apod == nil {
			apod, err := worker.GetApodByKey(cfg)
			if err != nil {
				fmt.Printf("error receiving apod: %v\n", err)
			} else {
				err := pgclient.ApodCreate(ctx, apod)
				if err != nil {
					fmt.Printf("error saving apod in the database: %v\n", err)
				}
			}
		}
		ticker := time.NewTicker(constants.WorkFrequency)
		defer ticker.Stop()
		for {
			for range ticker.C {
				apod, err := worker.GetApodByKey(cfg)
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
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	fs := http.FileServer(http.Dir("./storage"))
	mux.Handle("/storage/", http.StripPrefix("/storage/", fs))

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  constants.ServerTimeout,
		WriteTimeout: constants.ServerTimeout,
	}

	stopped := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sigint
		ctx, cancel := context.WithTimeout(context.Background(), constants.ServerTimeout)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("http server shutdown error %v", err)
		}
		close(stopped)
	}()

	log.Printf("starting HTTP server on :%s", cfg.Port)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("http server not listening %v", err)
	}

	<-stopped
}

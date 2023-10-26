package repository

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"testing"

	"github.com/artnikel/apodservice/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ory/dockertest/v3"
)

var (
	client *PgClient
	cfg    *config.Config
)

func SetupTestPgx() (*pgxpool.Pool, func(), error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, nil, fmt.Errorf("could not construct pool: %w", err)
	}
	resource, err := pool.Run("postgres", "latest", []string{
		"POSTGRES_USER=" + cfg.ApodUser,
		"POSTGRES_PASSWORD=" + cfg.ApodPassword,
		"POSTGRES_DB=" + cfg.ApodDB})
	if err != nil {
		return nil, nil, fmt.Errorf("could not start resource: %w", err)
	}
	err = RunMigrations(resource.GetPort(cfg.ApodDBPort + "/tcp"))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to run migrations: %w", err)
	}
	dbURL := fmt.Sprintf("postgres://"+cfg.ApodUser+":"+cfg.ApodPassword+"@localhost:%s/"+cfg.ApodDB,
		resource.GetPort(cfg.ApodDBPort+"/tcp"))
	cfg, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse dbURL: %w", err)
	}
	dbpool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect pgxpool: %w", err)
	}
	cleanup := func() {
		dbpool.Close()
		pool.Purge(resource)
	}
	return dbpool, cleanup, nil
}

func RunMigrations(port string) error {
	cmd := exec.Command("flyway",
		"-url=jdbc:postgresql://localhost:"+port+"/"+cfg.ApodDB,
		"-user="+cfg.ApodUser,
		"-password="+cfg.ApodPassword,
		"-locations=filesystem:../../migrations",
		"-connectRetries=10",
		"migrate")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	return nil
}

func TestMain(m *testing.M) {
	var err error
	cfg, err = config.New()
	if err != nil {
		log.Fatal(err)
	}
	dbpool, cleanupPgx, err := SetupTestPgx()
	if err != nil {
		cleanupPgx()
		log.Fatalf("could not construct the pool: %v", err)
	}
	client = NewPgClient(dbpool)
	exitVal := m.Run()
	cleanupPgx()
	os.Exit(exitVal)
}

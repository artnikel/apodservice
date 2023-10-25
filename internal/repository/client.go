// Package repository is a lower level of project
package repository

import "github.com/jackc/pgx/v5/pgxpool"

// PgClient represents the PostgreSQL repository implementation.
type PgClient struct {
	pool *pgxpool.Pool
}

// NewPgClient creates and returns a new instance of PgClient, using the provided pgxpool.Pool.
func NewPgClient(pool *pgxpool.Pool) *PgClient {
	return &PgClient{pool: pool}
}

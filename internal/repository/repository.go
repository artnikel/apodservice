package repository

import (
	"context"
	"fmt"

	"github.com/artnikel/apodservice/internal/model"
)

// ApodCreate insert new row of APOD in database
func (c *PgClient) ApodCreate(ctx context.Context, apod *model.APOD) error {
	var count int
	err := c.pool.QueryRow(ctx, "SELECT COUNT(id) FROM apod WHERE date = $1", apod.Date).Scan(&count)
	if err != nil {
		return fmt.Errorf("querryRow %w", err)
	}
	if count != 0 {
		return fmt.Errorf("the product name is already exist")
	}
	_, err = c.pool.Exec(ctx, `INSERT INTO apod 
	(id, copyright, date, explanation, media_type, service_version, title, url) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		apod.ID, apod.Copyright, apod.ParsedDate, apod.Explanation,
		apod.MediaType, apod.ServiceVersion, apod.Title, apod.URL)
	if err != nil {
		return fmt.Errorf("exec %w", err)
	}
	return nil
}

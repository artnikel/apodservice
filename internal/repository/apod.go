package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/artnikel/apodservice/internal/constants"
	"github.com/artnikel/apodservice/internal/model"
)

// ApodCreate insert new row of APOD in database
func (r *PgClient) ApodCreate(ctx context.Context, apod *model.APOD) error {
	var count int
	err := r.pool.QueryRow(ctx, "SELECT COUNT(id) FROM apod WHERE date = $1", apod.Date).Scan(&count)
	if err != nil {
		return fmt.Errorf("querryRow %w", err)
	}
	if count != 0 {
		return fmt.Errorf("apod is already exist")
	}
	_, err = r.pool.Exec(ctx, `INSERT INTO apod 
	(copyright, date, explanation, media_type, service_version, title, url) 
	VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		apod.Copyright, apod.ParsedDate, apod.Explanation,
		apod.MediaType, apod.ServiceVersion, apod.Title, apod.URL)
	if err != nil {
		return fmt.Errorf("exec %w", err)
	}
	return nil
}

// ApodGetAll returns a list of all APOD`s in database
func (r *PgClient) ApodGetAll(ctx context.Context) ([]*model.APOD, error) {
	apods := make([]*model.APOD, 0)
	rows, err := r.pool.Query(ctx, `SELECT copyright, date, explanation,
	 media_type,service_version, title, url FROM apod`)
	if err != nil {
		return nil, fmt.Errorf("query %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		apod := &model.APOD{}
		err = rows.Scan(&apod.Copyright, &apod.ParsedDate, &apod.Explanation, &apod.MediaType,
			&apod.ServiceVersion, &apod.Title, &apod.URL)
		if err != nil {
			return nil, fmt.Errorf("scan %w", err)
		}
		apod.Date = apod.ParsedDate.Format(constants.DateLayout)
		apods = append(apods, apod)
	}
	return apods, nil
}

// ApodGetByDate return one row of APOD by date
func (r *PgClient) ApodGetByDate(ctx context.Context, date time.Time) (*model.APOD, error) {
	apod := &model.APOD{}
	err := r.pool.QueryRow(ctx, `SELECT copyright, explanation, media_type,
	 service_version, title, url FROM apod WHERE date=$1`, date).
		Scan(&apod.Copyright, &apod.Explanation, &apod.MediaType,
			&apod.ServiceVersion, &apod.Title, &apod.URL)
	if err != nil {
		return nil, fmt.Errorf("query %w", err)
	}
	return apod, nil
}

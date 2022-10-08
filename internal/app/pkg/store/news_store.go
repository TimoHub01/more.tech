package store

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type postgresStore struct {
	pool *pgxpool.Pool
}

func NewPostgresStore(pool *pgxpool.Pool) *postgresStore {
	return &postgresStore{pool}
}

func (s *postgresStore) CreateNew(ctx context.Context, topic, date, link, text string) (New, error) {
	sql := "INSERT INTO news (link,topic,text,date) VALUES ($1,$2,$3,$4) RETURNING id"
	params := []interface{}{
		link,  // 1
		topic, // 2
		text,  //3
		date,  // 4
	}
	row := s.pool.QueryRow(ctx, sql, params...)
	var id string
	err := row.Scan(&id)
	if err != nil {
		return New{}, err
	}
	return New{
		ID:    id,
		topic: topic,
		text:  text,
		date:  date,
		link:  link,
	}, nil
}

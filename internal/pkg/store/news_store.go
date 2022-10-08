package store

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type postgresStore struct {
	pool *pgxpool.Pool
}

func NewPostgresStore(pool *pgxpool.Pool) *postgresStore {
	return &postgresStore{pool}
}

func scanNew(row pgx.Row) (New, error) {
	var n New
	err := row.Scan(&n.ID, &n.Topic, &n.Text, &n.Date, &n.Link)
	if err != nil {
		return New{}, err
	}
	return n, nil
}
func scanNews(rows pgx.Rows) ([]New, error) {
	var news []New
	var n New
	for rows.Next() {
		err := rows.Scan(&n.ID, &n.Link, &n.Topic, &n.Text, &n.Date)
		news = append(news, n)
		if err != nil {
			return []New{}, fmt.Errorf("failed to select users from db %w", err)
		}
	}
	return news, nil
}

func (s *postgresStore) CreateNew(ctx context.Context, topic, date, link, text string) (New, error) {
	sql := "INSERT INTO news (Link,Topic,Text,Date) VALUES ($1,$2,$3,$4) RETURNING id"
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
		Topic: topic,
		Text:  text,
		Date:  date,
		Link:  link,
	}, nil
}

func (s *postgresStore) GetNews(ctx context.Context) ([]New, error) {
	sql := "SELECT id,Link,Topic,Text,Date FROM news order by Link LIMIT 3;"
	rows, err := s.pool.Query(ctx, sql)
	if err != nil {
		return []New{}, err
	}
	news, err := scanNews(rows)
	if err != nil {
		return []New{}, err
	}
	return news, nil
}

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

type Trend struct {
	Trend string
	Link  string
	Topic string
	Date  string
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

func scanTrends(rows pgx.Rows) ([]Trend, error) {
	var trends []Trend
	var n Trend
	for rows.Next() {
		err := rows.Scan(&n.Link, &n.Topic, &n.Date, &n.Trend)
		trends = append(trends, n)
		if err != nil {
			return []Trend{}, fmt.Errorf("failed to select users from db %w", err)
		}
	}
	return trends, nil
}

func (s *postgresStore) CreateNew(ctx context.Context, topic, date, link, text string) (New, error) {
	sql := "INSERT INTO news (Link,Topic,Text,Date, src) VALUES ($1,$2,$3,$4, 'web') RETURNING id"
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

func (s *postgresStore) GetNews(ctx context.Context, role int) ([]New, error) {
	sql := "SELECT id,link,topic,text,date,DISTINCT(text_cluster) FROM news_new WHERE role == ($1) and src == 'web' order by relevant_score LIMIT 3;"
	rows, err := s.pool.Query(ctx, sql, role)
	if err != nil {
		return []New{}, err
	}
	news, err := scanNews(rows)
	if err != nil {
		return []New{}, err
	}
	return news, nil
}

func (s *postgresStore) GetTrend(ctx context.Context) ([]Trend, error) {
	sql := "SELECT link,topic,date,DISTINCT(trend) FROM news_new WHERE role == ($1) and src == 'web' order by relevant_score;"
	rows, err := s.pool.Query(ctx, sql)
	if err != nil {
		return []Trend{}, err
	}
	trends, err := scanTrends(rows)
	if err != nil {
		return []Trend{}, err
	}
	return trends, nil
}

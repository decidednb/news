package postgres

import (
	"context"
	"news/pkg/storage"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

// New - конструктор, conn - строка подключения к БД
func New(conn string) (*Storage, error) {
	db, err := pgxpool.Connect(context.Background(), conn)

	if err != nil {
		return nil, err
	}

	s := Storage{
		db: db,
	}

	return &s, nil
}

// ConnClose - закрывает соединение к базе данных
func (s *Storage) Close() {
	s.db.Close()
}

// StoreNews - создание массива задач
func (s *Storage) StoreNews(news []storage.Post) error {

	for _, p := range news {
		_, err := s.db.Exec(context.Background(), `
		INSERT INTO news(title, content, pub_time, link)
		VALUES ($1, $2, $3, $4)`,
			p.Title,
			p.Content,
			p.PubTime,
			p.Link)

		if err != nil {
			return err
		}
	}

	return nil
}

// LastNews - возвращает n-количество последних новостей
// (по дате публикации) из базы данных.
func (s *Storage) LastNews(n int) ([]storage.Post, error) {
	if n == 0 {
		n = 10
	}

	rows, err := s.db.Query(context.Background(), `
		SELECT id, title, content, pub_time, link FROM news
		ORDER BY pub_time DESC
		LIMIT $1
		`, n)

	if err != nil {
		return nil, err
	}

	var news []storage.Post
	for rows.Next() {
		var p storage.Post
		err = rows.Scan(
			&p.ID,
			&p.Title,
			&p.Content,
			&p.PubTime,
			&p.Link,
		)
		if err != nil {
			return nil, err
		}
		news = append(news, p)
	}
	return news, rows.Err()
}

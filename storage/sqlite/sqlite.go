package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sunriseex/tgbot-notifier/storage"
)

type Storage struct {
	db *sql.DB
}

// New creates new SQLite storage
func New(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("cant open sqlite database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("cant connect sqlite database: %w", err)
	}
	return &Storage{db: db}, nil
}

// Save pages from user to storage.
func (s *Storage) Save(ctx context.Context, p *storage.Page) error {
	query := `INSERT INTO pages (url,user_name) VALUES (?,?)`
	if _, err := s.db.ExecContext(ctx, query, p.URL, p.UserName); err != nil {
		return fmt.Errorf("can't save page: %w", err)
	}
	return nil
}

// PickRandom picks random page from storage.
func (s *Storage) PickRandom(ctx context.Context, userName string) (*storage.Page, error) {
	query := `SELECT url FROM pages WHERE user_name =? ORDER BY RANDOM() LIMIT 1`
	var url string

	err := s.db.QueryRowContext(ctx, query, userName).Scan(&url)
	if err == sql.ErrNoRows {
		return nil, storage.ErrNoSavedPages
	}

	if err != nil {
		return nil, fmt.Errorf("cant pick random page: %w", err)
	}
	return &storage.Page{
		URL:      url,
		UserName: userName,
	}, nil
}

// Remove page from storage.
func (s Storage) Remove(ctx context.Context, p *storage.Page) error {
	q := `DELETE FROM pages WHERE url =? AND user_name =?`
	if _, err := s.db.ExecContext(ctx, q, p.URL, p.UserName); err != nil {
		return fmt.Errorf("can't remove page: %w", err)
	}
	return nil
}

// IsExists check if page already exists.
func (s Storage) IsExists(ctx context.Context, p *storage.Page) (bool, error) {
	q := `SELECT COUNT(*) FROM pages WHERE url =? AND user_name =?`
	var count int

	if err := s.db.QueryRowContext(ctx, q, p.URL, p.UserName).Scan(&count); err != nil {
		return false, fmt.Errorf("can't check if page exists: %w", err)
	}
	return count > 0, nil
}

func (s *Storage) Init(ctx context.Context) error {
	q := `CREATE TABLE IF NOT EXISTS pages (url TEXT, user_name TEXT)`
	if _, err := s.db.ExecContext(ctx, q); err != nil {
		return fmt.Errorf("can't create table: %w", err)
	}
	return nil
}

package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
	"github.com/popvaleks/url-shortener/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func New(dbPath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS url(
		id INTEGER PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL);
	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUrl(inputUrl string, alias string) (int64, error) {
	const op = "storage.sqlite.SaveUrl"

	stmt, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(inputUrl, alias)

	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
			return 0, fmt.Errorf("%s: %w", op, sqliteErr)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetUrl(alias string) (string, error) {
	const op = "storage.sqlite.GetUrl"

	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias = ?")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	var url string

	err = stmt.QueryRow(alias).Scan(&url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrUrlNotFound
		}

		return "", fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return url, nil
}

func (s *Storage) DeleteUrl(alias string) error {
	const op = "storage.sqlite.DeleteUrl"

	stmt, err := s.db.Prepare("DELETE FROM url WHERE alias = ?")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	result, err := stmt.Exec(alias)
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: get rows affected: %w", op, err)
	}

	if rowsAffected == 0 {
		return storage.ErrUrlNotFound
	}

	return nil
}

func (s *Storage) GetAllUrls() (map[string]string, error) {
	const op = "storage.sqlite.GetAllUrls"
	stmt, err := s.db.Prepare("SELECT url, alias FROM url")

	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	urlMap := make(map[string]string)

	for rows.Next() {
		var url, alias string
		if err := rows.Scan(&url, &alias); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		urlMap[alias] = url
	}

	return urlMap, nil
}

func (s *Storage) UpdateUrl(url, alias string) (string, error) {
	const op = "storage.sqlite.updateUrl"

	fmt.Println(url, alias)

	var exists bool
	r := s.db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM url WHERE alias = :alias)",
		sql.Named("alias", alias),
	)

	err := r.Scan(&exists)
	if err != nil {
		return "", err
	}

	if exists == false {
		return "", storage.ErrAliasNotFound
	}

	result, err := s.db.Exec("UPDATE url SET url = :url WHERE alias = :alias",
		sql.Named("url", url),
		sql.Named("alias", alias),
	)

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		return "", storage.ErrAliasNotFound
	}

	return alias, nil
}

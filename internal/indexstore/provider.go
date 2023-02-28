package indexstore

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"

	database "database/sql"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/mattn/go-sqlite3"
)

const createTable = `
CREATE VIRTUAL TABLE IF NOT EXISTS files 
USING FTS5 (filename, body);
`

func formatDatabaseName(indexName string) string {
	return fmt.Sprintf("index_%s.db", indexName)
}

type Storage struct {
	db *database.DB
	mu sync.Mutex
}

func New(baseDir, indexName string) (*Storage, error) {
	fileName := filepath.Join(baseDir, formatDatabaseName(indexName))
	db, err := database.Open("sqlite3", fileName)
	if err != nil {
		return nil, fmt.Errorf("storage: unable to open a storage: %v", err)
	}

	if _, err := db.Exec(createTable); err != nil {
		return nil, fmt.Errorf("storage: unable to init tables: %v", err)
	}

	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

type File struct {
	Path    string
	Content string
}

func (s *Storage) InsertFiles(ctx context.Context, files []File) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	query := sq.Insert("files").Columns("filename", "body")

	for _, f := range files {
		query = query.Values(f.Path, f.Content)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("storage: unable to compose a SQL consolequery: %v", err)
	}

	if _, err := s.db.ExecContext(ctx, sql, args...); err != nil {
		return fmt.Errorf("storage: unble to insert files to the index: %v", err)
	}

	return nil
}

type SearchResult struct {
	FileName string
	Content  string
}

func (s *Storage) SearchAll(ctx context.Context, query string) ([]SearchResult, error) {
	q := sq.Select(
		"filename",
		"highlight(files, 1, '_hl_start_', '_hl_end_') AS body",
	).
		From("files").
		Where(sq.Expr("body MATCH ?", query)).
		OrderBy("rank")

	return s.search(ctx, q)
}

func (s *Storage) Search(ctx context.Context, query string, page, perPage uint64) ([]SearchResult, error) {
	searchQuery := sq.Select(
		"filename",
		"highlight(files, 1, '_hl_start_', '_hl_end_') AS body",
	).
		From("files").
		Where(sq.Expr("body MATCH ?", query)).
		OrderBy("rank").
		Limit(perPage).
		Offset(perPage * page)

	return s.search(ctx, searchQuery)
}

func (s *Storage) search(ctx context.Context, q sq.SelectBuilder) ([]SearchResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	sql, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("storage: unable to compose a SQL consolequery: %v", err)
	}

	rows, err := s.db.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("storage: unable to get a result from the storage: %v", err)
	}

	defer rows.Close()

	var results []SearchResult
	for rows.Next() {
		var result SearchResult

		if err := rows.Scan(&result.FileName, &result.Content); err != nil {
			return nil, fmt.Errorf("storage: unable to scan a row: %v", err)
		}

		results = append(results, result)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("storage: returned rows contain an error: %v", err)
	}

	return results, nil
}

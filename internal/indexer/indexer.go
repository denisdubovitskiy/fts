package indexer

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"unicode/utf8"

	"github.com/denisdubovitskiy/fts/internal/indexstore"
)

type Storage interface {
	InsertFiles(ctx context.Context, files []indexstore.File) error
}

type Indexer struct {
	storage   Storage
	batchSize int
	path      string
}

func New(storage Storage, path string) *Indexer {
	return &Indexer{
		storage:   storage,
		batchSize: 50,
		path:      path,
	}
}
func (i *Indexer) Run(ctx context.Context) error {
	var files []indexstore.File

	walkErr := filepath.Walk(i.path, func(path string, info fs.FileInfo, e error) error {
		if e != nil {
			return e
		}

		if info.IsDir() {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		if !IsText(content) {
			return nil
		}

		files = append(files, indexstore.File{
			Path:    path,
			Content: string(content),
		})

		if len(files) >= i.batchSize {
			for _, f := range files {
				log.Printf("indexer: indexing file %s\n", f.Path)
			}
			if err := i.storage.InsertFiles(ctx, files); err != nil {
				return fmt.Errorf("indexer: unable to insert files: %v", err)
			}
		}

		return nil
	})

	if walkErr != nil {
		return walkErr
	}

	if len(files) > 0 {
		for _, f := range files {
			log.Printf("indexer: indexing file %s\n", f.Path)
		}
		if err := i.storage.InsertFiles(ctx, files); err != nil {
			return err
		}
	}

	return nil
}

// IsText reports whether a significant prefix of s looks like correct UTF-8;
// that is, if it is likely that s is human-readable text.
func IsText(s []byte) bool {
	const max = 1024 // at least utf8.UTFMax
	if len(s) > max {
		s = s[0:max]
	}
	for i, c := range string(s) {
		if i+utf8.UTFMax > len(s) {
			// last char may be incomplete - ignore
			break
		}
		if c == 0xFFFD || c < ' ' && c != '\n' && c != '\t' && c != '\f' {
			// decoding error or control character - not a text file
			return false
		}
	}
	return true
}

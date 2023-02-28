package consoleindex

import (
	"context"
	"fmt"
	"github.com/denisdubovitskiy/fts/internal/indexer"
	"github.com/denisdubovitskiy/fts/internal/indexstore"
)

func NewProvider(baseDir, indexName string) *Provider {
	return &Provider{
		baseDir:   baseDir,
		indexName: indexName,
	}
}

type Provider struct {
	indexName string
	baseDir   string
}

func (p *Provider) Run(ctx context.Context, path string) error {
	storage, err := indexstore.New(p.baseDir, p.indexName)
	if err != nil {
		return err
	}
	defer storage.Close()

	index := indexer.New(storage, path)
	if err := index.Run(ctx); err != nil {
		return fmt.Errorf("fts: unable to run indexer: %v", err)
	}

	return nil
}

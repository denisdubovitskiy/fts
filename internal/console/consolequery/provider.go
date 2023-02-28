package consolequery

import (
	"context"
	"fmt"
	"strings"

	"github.com/denisdubovitskiy/fts/internal/indexstore"
	"github.com/denisdubovitskiy/fts/internal/textutil/snippeter"
	"github.com/fatih/color"
)

var accent = color.New(color.FgBlue)

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

func (p *Provider) Run(ctx context.Context, query string) error {
	storage, err := indexstore.New(p.baseDir, p.indexName)
	if err != nil {
		return err
	}
	defer storage.Close()

	results, err := storage.SearchAll(ctx, query)
	if err != nil {
		return fmt.Errorf("fts: unable to get search results: %v", err)
	}

	for _, r := range results {
		fmt.Println(r.FileName)

		snippets := snippeter.ParseToSnippets(r.Content, snippeter.Options{
			HighlightStart: "<b>",
			HighlightEnd:   "</b>",

			SourceStart: "_hl_start_",
			SourceEnd:   "_hl_end_",
		})

		for _, snippet := range snippets {
			startIdx := strings.Index(snippet, "<b>")
			if startIdx < 0 {
				continue // TODO
			}
			endIdx := strings.Index(snippet, "</b>")
			if endIdx < 0 {
				continue // TODO
			}

			fmt.Printf(
				"%s%s%s\n\n",
				snippet[:startIdx],
				accent.Sprint(snippet[startIdx+len("<b>"):endIdx]),
				snippet[endIdx+len("</b>"):],
			)
		}
	}

	return nil
}

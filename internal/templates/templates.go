package templates

import (
	"embed"
	"html/template"
	"io"
	"path/filepath"
)

//go:embed html
var html embed.FS

var all = template.
	Must(
		template.
			New("").
			Funcs(template.FuncMap{
				"fileurl": func(path string) string {
					return filepath.Join("/files", path)
				},
			}).
			ParseFS(html, "html/*.gohtml", "html/style.css"),
	)

type IndexData struct {
	SearchText     string
	InputFilesRoot string
	IndexName      string
}

func RenderIndex(wr io.Writer, data IndexData) error {
	return all.ExecuteTemplate(wr, "index.gohtml", data)
}

type SearchResult struct {
	Snippets []template.HTML
	FileName string
}

type SearchResultsData struct {
	SearchResults []SearchResult
	SearchText    string
	IndexName     string
}

func RenderSearchResults(wr io.Writer, data SearchResultsData) error {
	return all.ExecuteTemplate(wr, "results.gohtml", data)
}

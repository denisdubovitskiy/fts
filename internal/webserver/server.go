package webserver

import (
	"context"
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/denisdubovitskiy/fts/internal/indexstore"
	"github.com/denisdubovitskiy/fts/internal/templates"
	"github.com/denisdubovitskiy/fts/internal/textutil/snippeter"
	"github.com/microcosm-cc/bluemonday"
)

type Storage interface {
	Search(ctx context.Context, query string, page, perPage uint64) ([]indexstore.SearchResult, error)
}

type Server struct {
	storage Storage
	opts    Options
	mux     http.Handler
}

type Options struct {
	Addr           string
	IndexName      string
	InputFilesRoot string

	TagHighlightStart string
	TagHighlightEnd   string

	HighlightMarkerStart string
	HighlightMarkerEnd   string
}

func (o *Options) getAddr() string {
	if o.Addr != "" {
		return o.Addr
	}

	return "127.0.0.1:8081"
}

func New(storage Storage, opts Options) *Server {
	mux := http.NewServeMux()

	s := &Server{
		storage: storage,
		opts:    opts,
		mux:     mux,
	}

	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		s.renderIndexTemplate(writer)
		return
	})
	mux.HandleFunc("/search", func(writer http.ResponseWriter, request *http.Request) {
		s.renderSearchResults(writer, request)
		return
	})

	mux.HandleFunc("/files/", func(writer http.ResponseWriter, request *http.Request) {
		path := strings.TrimPrefix(request.URL.Path, "/files")
		content, err := os.ReadFile(path)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		writer.Write(content)
	})

	return s
}

func (s *Server) Run(ctx context.Context) error {
	server := &http.Server{
		Addr:    s.opts.getAddr(),
		Handler: s.mux,
	}
	go func() {
		<-ctx.Done()
		log.Printf("server: shutdown signal received")
		server.Shutdown(context.Background())
	}()

	log.Printf("server: listening at %s\n", s.opts.getAddr())

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	log.Printf("server: stopped")

	return nil
}

func (s *Server) renderIndexTemplate(writer http.ResponseWriter) {
	templates.RenderIndex(writer, templates.IndexData{
		InputFilesRoot: s.opts.InputFilesRoot,
		IndexName:      s.opts.IndexName,
	})
}

func (s *Server) renderSearchResults(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query().Get("query")
	page := parseInt(request.URL.Query().Get("page"))
	if page > 0 {
		page--
	}
	perPage := parseInt(request.URL.Query().Get("per_page"))

	found, err := s.storage.Search(request.Context(), strings.ToLower(query), page, perPage)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}

	policy := bluemonday.StrictPolicy()

	var searchResults []templates.SearchResult
	for _, result := range found {

		snippets := snippeter.ParseToSnippets(policy.Sanitize(result.Content), snippeter.Options{
			HighlightStart: s.opts.TagHighlightStart,
			HighlightEnd:   s.opts.TagHighlightEnd,

			// Returned by the database highlighter
			SourceStart: s.opts.HighlightMarkerStart,
			SourceEnd:   s.opts.HighlightMarkerEnd,
		})

		searchResults = append(searchResults, templates.SearchResult{
			FileName: result.FileName,
			Snippets: convertToHTML(snippets),
		})
	}

	templates.RenderSearchResults(writer, templates.SearchResultsData{
		SearchResults: searchResults,
		SearchText:    query,
		IndexName:     s.opts.IndexName,
	})
}

func parseInt(s string) uint64 {
	if n, err := strconv.ParseUint(s, 10, 64); err == nil {
		return n
	}
	return 0
}

func convertToHTML(snippets []string) []template.HTML {
	result := make([]template.HTML, len(snippets))
	for i, s := range snippets {
		result[i] = template.HTML(s)
	}
	return result
}

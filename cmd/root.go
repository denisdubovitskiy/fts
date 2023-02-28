package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/denisdubovitskiy/fts/internal/console/consoleindex"
	"github.com/denisdubovitskiy/fts/internal/console/consolequery"
	"github.com/denisdubovitskiy/fts/internal/indexstore"
	"github.com/denisdubovitskiy/fts/internal/webserver"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	SilenceUsage: true,
	Use:          "fts",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().String("index-name", "default", "index name (defaults to `default`)")
	rootCmd.PersistentFlags().String("directory", "./", "directory to store database files (defaults to ./)")

	rootCmd.AddCommand(indexCommand)
	rootCmd.AddCommand(queryCommand)

	webCommand.PersistentFlags().String("addr", "127.0.0.1:8082", "an address to listen (defaults to 127.0.0.1:8082)")
	rootCmd.AddCommand(webCommand)
}

var indexCommand = &cobra.Command{
	Use:        "index",
	ArgAliases: []string{"path"},
	Args:       checkDocumentsPath,
	RunE: func(cmd *cobra.Command, args []string) error {
		documentsPath := args[0]

		provider := consoleindex.NewProvider(
			cmd.Flag("directory").Value.String(),
			cmd.Flag("index-name").Value.String(),
		)

		return provider.Run(cmd.Context(), documentsPath)
	},
}

var queryCommand = &cobra.Command{
	Use:        "query",
	ArgAliases: []string{"path"},
	Args:       cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]

		provider := consolequery.NewProvider(
			cmd.Flag("directory").Value.String(),
			cmd.Flag("index-name").Value.String(),
		)

		return provider.Run(cmd.Context(), query)
	},
}

var webCommand = &cobra.Command{
	Use:        "web",
	Args:       checkDocumentsPath,
	ArgAliases: []string{"path"},
	RunE: func(cmd *cobra.Command, args []string) error {
		documentsPath := args[0]
		storage, err := indexstore.New(
			cmd.Flag("directory").Value.String(),
			cmd.Flag("index-name").Value.String(),
		)
		if err != nil {
			return err
		}
		defer storage.Close()

		server := webserver.New(storage, webserver.Options{
			Addr:           cmd.Flag("addr").Value.String(),
			IndexName:      cmd.Flag("index-name").Value.String(),
			InputFilesRoot: documentsPath,

			TagHighlightStart: `<b class="highlight">`,
			TagHighlightEnd:   `</b>`,

			HighlightMarkerStart: `_hl_start_`,
			HighlightMarkerEnd:   `_hl_end_`,
		})

		ctx, cancel := signal.NotifyContext(cmd.Context(), syscall.SIGTERM, syscall.SIGINT)
		defer cancel()

		return server.Run(ctx)
	},
}

func checkDocumentsPath(cmd *cobra.Command, args []string) error {
	if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
		return err
	}

	stat, err := os.Stat(args[0])
	if err != nil {
		return fmt.Errorf("fts: unable to stat a given path: %v", err)
	}

	if !stat.IsDir() {
		return fmt.Errorf("fts: %s does not appear to be a directory", args[0])
	}

	return nil
}

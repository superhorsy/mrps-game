package main

import (
	"errors"
	"log"

	"github.com/spf13/cobra"

	"mrps-game/internal/app"
)

func main() {
	if err := NewRootCommand().Execute(); err != nil {
		log.Fatal(err)
	}
}

func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "game",
		Short: "Start game server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("choose setup mode")
		},
	}
	rootCmd.AddCommand(NewGameServer())
	return rootCmd
}

func NewGameServer() *cobra.Command {
	var dsn string
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start game server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.Start(dsn)
		},
	}
	cmd.Flags().StringVar(&dsn, "dsn", "", "Postgres DSN")
	return cmd
}

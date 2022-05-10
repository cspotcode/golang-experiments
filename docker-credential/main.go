package main

import (
	"fmt"
	"os"

	"github.com/cspotcode/golang-experiments/docker-credential/cli"
	"github.com/cspotcode/golang-experiments/docker-credential/commands/store"
	"github.com/cspotcode/golang-experiments/docker-credential/log"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func main() {
	logger := log.CreateLogger()
	cli := cli.Cli{
		Logger: logger,
	}
	var rootCmd = &cobra.Command{
		Use:   "docker-credential",
		Short: "CLI for manipulating docker CLI's stored credentials",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.Errorf("you must specify a subcommand")
		},
	}
	rootCmd.PersistentFlags().BoolVarP(&logger.VerboseEnabled, "verbose", "v", false, "verbose output")

	// NOTE docker cli plugins are deprecated
	pluginMetadataCmd := &cobra.Command{
		Use: "docker-cli-plugin-metadata",
		RunE: func(cmd *cobra.Command, args []string) error {
			println(`{` +
				`"SchemaVersion":"0.1.0",` +
				`"Vendor":"cspotcode",` +
				`"Version":"0.0.0",` +
				`"ShortDescription":"Store credentials into docker credential store",` +
				`"URL":"https://github.com/cspotcode/golang-experiments/docker-credential"` +
				`}`)
			return nil
		},
	}
	rootCmd.AddCommand(pluginMetadataCmd)
	rootCmd.AddCommand(store.CreateStoreCmd(&cli))

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

package main

import (
	"fmt"
	"os"
	"github.com/cspotcode/golang-experiments/docker-credential-cli/commands/store"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "docker-credential",
		Short: "CLI for manipulating docker CLI's stored credentials",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.Errorf("you must specify a subcommand")
		},
	}

	// NOTE docker cli plugins are deprecated
	pluginMetadataCmd := &cobra.Command{
		Use: "docker-cli-plugin-metadata",
		RunE: func(cmd *cobra.Command, args []string) error {
			println(`{` +
				`"SchemaVersion":"0.1.0",` +
				`"Vendor":"cspotcode",` +
				`"Version":"0.0.0",` +
				`"ShortDescription":"Store credentials into docker credential store",` +
				`"URL":"https://github.com/cspotcode/golang-experiments/docker-credential-cli"` +
				`}`)
			return nil
		},
	}
	rootCmd.AddCommand(pluginMetadataCmd)
	rootCmd.AddCommand(store.CreateStoreCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

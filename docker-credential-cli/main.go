package main

import (
	"fmt"
	"os"

	dockerconfig "github.com/docker/cli/cli/config"
	dockerconfigtypes "github.com/docker/cli/cli/config/types"
	"github.com/docker/docker/api/types"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type isFileStore interface {
	IsFileStore() bool
	GetFilename() string
}

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

	storeCmdArgs := storeCmdArgs{}
	storeCmd := &cobra.Command{
		Use:   "store",
		Short: "store credentials",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return store(storeCmdArgs)
		},
		Annotations: map[string]string{"version": "1.31"},
	}
	storeCmd.Flags().StringVarP(&storeCmdArgs.username, "username", "u", "", "username, for username & password auth")
	storeCmd.Flags().StringVarP(&storeCmdArgs.password, "password", "p", "", "password, for username & password auth")
	storeCmd.Flags().StringVarP(&storeCmdArgs.identityToken, "token", "t", "", "identity token, for token auth")
	storeCmd.Flags().StringVarP(&storeCmdArgs.serverAddress, "address", "a", "", "docker registry address")
	rootCmd.AddCommand(storeCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type storeCmdArgs struct {
	username      string
	password      string
	identityToken string
	serverAddress string
}

func store(args storeCmdArgs) error {
	println("storing credentials for", args.serverAddress)
	config := dockerconfig.LoadDefaultConfigFile(os.Stderr)

	authConfig := types.AuthConfig{}
	authConfig.ServerAddress = args.serverAddress
	authConfig.IdentityToken = args.identityToken
	authConfig.Username = args.username
	authConfig.Password = args.password

	creds := config.GetCredentialsStore(authConfig.ServerAddress)

	// store, isDefault := creds.(isFileStore)
	if err := creds.Store(dockerconfigtypes.AuthConfig(authConfig)); err != nil {
		return errors.Errorf("Error saving credentials: %v", err)
	}
	return nil
}

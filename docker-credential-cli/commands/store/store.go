package store

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

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

type storeCmdArgs struct {
	username      string
	password      string
	passwordStdin bool
	// identityToken string
	serverAddress string
}

func CreateStoreCmd() *cobra.Command {
	opts := storeCmdArgs{}
	cmd := &cobra.Command{
		Use:   "store <server address>",
		Short: "store credentials",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.serverAddress = args[0]
			if opts.password != "" {
				fmt.Fprintln(os.Stderr, "WARNING! Using --password via the CLI is insecure. Use --password-stdin.")
				if opts.passwordStdin {
					return errors.New("--password and --password-stdin are mutually exclusive")
				}
			}
			if opts.username == "" {
				return errors.New("Must provide --username")
			}
			if opts.passwordStdin {
				contents, err := ioutil.ReadAll(os.Stdin)
				if err != nil {
					return err
				}

				opts.password = strings.TrimSuffix(strings.TrimSuffix(string(contents), "\n"), "\r")
			}

			return Store(opts)
		},
	}
	cmd.Flags().StringVarP(&opts.username, "username", "u", "", "username")
	cmd.Flags().StringVarP(&opts.password, "password", "p", "", "password")
	cmd.Flags().BoolVarP(&opts.passwordStdin, "password-stdin", "", false, "Take the password from stdin")
	// cmd.Flags().StringVarP(&storeCmdArgs.identityToken, "token", "t", "", "identity token, for token auth")

	cmd.AddCommand(CreateStoreAwsCmd())
	return cmd
}

func Store(args storeCmdArgs) error {
	// println("storing credentials for", args.serverAddress)
	config := dockerconfig.LoadDefaultConfigFile(os.Stderr)

	authConfig := types.AuthConfig{}
	authConfig.ServerAddress = args.serverAddress
	// authConfig.IdentityToken = args.identityToken
	authConfig.Username = args.username
	authConfig.Password = args.password

	creds := config.GetCredentialsStore(authConfig.ServerAddress)

	// store, isDefault := creds.(isFileStore)
	if err := creds.Store(dockerconfigtypes.AuthConfig(authConfig)); err != nil {
		return errors.Errorf("Error saving credentials: %v", err)
	}
	return nil
}

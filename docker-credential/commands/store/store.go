package store

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/cspotcode/golang-experiments/docker-credential-cli/cli"
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

type StoreCmdArgs struct {
	cli           *cli.Cli
	username      string
	password      string
	passwordStdin bool
	// identityToken string
	serverAddress string
}

func CreateStoreCmd(cli *cli.Cli) *cobra.Command {
	opts := StoreCmdArgs{
		cli: cli,
	}
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
	cmd.Flags().BoolVar(&opts.passwordStdin, "password-stdin", false, "Take the password from stdin")
	// cmd.Flags().StringVarP(&StoreCmdArgs.identityToken, "token", "t", "", "identity token, for token auth")

	cmd.AddCommand(CreateStoreAwsCmd(cli))
	return cmd
}

func Store(args StoreCmdArgs) error {
	logger := args.cli.Logger
	config := dockerconfig.LoadDefaultConfigFile(os.Stderr)

	logger.Verbose("Storing credentials for " + args.serverAddress)
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
	logger.Verbose("Stored credentials for " + args.serverAddress)
	return nil
}

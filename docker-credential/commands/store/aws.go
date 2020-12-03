package store

import (
	"encoding/base64"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/cspotcode/golang-experiments/docker-credential/cli"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type storeAwsCmdArgs struct {
	cli       *cli.Cli
	regions   []string
	accountid string
}

func CreateStoreAwsCmd(cli *cli.Cli) *cobra.Command {
	opts := storeAwsCmdArgs{
		cli: cli,
	}
	cmd := &cobra.Command{
		Use:   "aws <account>",
		Short: "generate and store credentials for AWS ECR",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.accountid = args[0]

			return StoreAws(opts)
		},
	}
	cmd.Flags().StringSliceVar(&opts.regions, "regions", []string{}, "regions, comma-separated.  Omit to attempt auth against all regions")
	return cmd
}

type regionAuthToken struct {
	token  string
	region string
}

func StoreAws(opts storeAwsCmdArgs) error {
	logger := opts.cli.Logger
	useAllRegions := false
	if len(opts.regions) == 0 {
		useAllRegions = true
	}
	regions := opts.regions
	if useAllRegions {
		ec2Session := session.Must(session.NewSession(&aws.Config{
			Region: aws.String("us-east-1"),
		}))
		ec2Svc := ec2.New(ec2Session)
		allRegions, err := ec2Svc.DescribeRegions(&ec2.DescribeRegionsInput{})
		if err != nil {
			return errors.Errorf("Error listing all regions: %v", err)
		}
		regions = make([]string, len(allRegions.Regions))
		for i, region := range allRegions.Regions {
			regions[i] = *region.RegionName
		}
	}

	results := make(chan regionAuthToken)
	errs := make(chan error)
	done, wg := taskSplit(len(regions))
	for _, region := range regions {
		go func(region string) {
			auth, err := func() (*regionAuthToken, error) {
				logger.Verbose("Getting authentication token for " + region)
				ecrSession := session.Must(session.NewSession(&aws.Config{
					Region: aws.String(region),
				}))
				svc := ecr.New(ecrSession)
				authorizationToken, err := svc.GetAuthorizationToken(&ecr.GetAuthorizationTokenInput{
					RegistryIds: []*string{&opts.accountid},
				})
				if err != nil {
					return nil, errors.Errorf("Error generating ecr credentials for region %s: %v", region, err)
				}
				logger.Verbose("Got authentication token for " + region)
				return &regionAuthToken{
					token:  *authorizationToken.AuthorizationData[0].AuthorizationToken,
					region: region,
				}, nil
			}()
			if err != nil {
				errs <- err
			} else {
				results <- *auth
			}
			wg.Done()
		}(region)
	}
Loop:
	for {
		select {
		case <-done:
			break Loop

		case auth := <-results:
			// Intentionally call `Store()` outside of goroutines to avoid potential concurrency issues with go
			// credential store (I don't know if I need to care about this)
			decoded, err := base64.StdEncoding.DecodeString((auth.token))
			if err != nil {
				return err
			}
			err2 := Store(StoreCmdArgs{
				cli:           opts.cli,
				username:      "AWS",
				password:      string(decoded),
				serverAddress: opts.accountid + ".dkr.ecr." + auth.region + ".amazonaws.com",
			})
			if err2 != nil {
				return err2
			}

		case err := <-errs:
			return err
		}
	}

	return nil
}

func taskSplit(taskCount int) (chan int, *sync.WaitGroup) {
	var wg sync.WaitGroup
	wg.Add(taskCount)
	doneChannel := make(chan int)
	go func() {
		wg.Wait()
		doneChannel <- 0
	}()
	return doneChannel, &wg
}

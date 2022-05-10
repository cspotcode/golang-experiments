package store

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type storeAwsCmdArgs struct {
	region    string
	accountid string
}

func CreateStoreAwsCmd() *cobra.Command {
	opts := storeAwsCmdArgs{}
	cmd := &cobra.Command{
		Use:   "aws <account>",
		Short: "generate and store credentials for AWS ECR",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.accountid = args[0]

			return StoreAws(opts)
		},
	}
	cmd.Flags().StringVarP(&opts.region, "region", "", "", "region (pass all or omit to auth all regions)")
	return cmd
}

func StoreAws(opts storeAwsCmdArgs) error {
	useAllRegions := false
	if opts.region == "" || opts.region == "all" {
		useAllRegions = true
	}
	regions := []string{opts.region}
	if useAllRegions {
		ec2Session := session.Must(session.NewSession(&aws.Config{
			Region: aws.String("us-east-1"),
		}))
		//ec2Session := session.Must(session.NewSession())
		ec2Svc := ec2.New(ec2Session)
		allRegions, err := ec2Svc.DescribeRegions(&ec2.DescribeRegionsInput{})
		if err != nil {
			return errors.Errorf("Error listing all regions: %v", err)
		}
		println(allRegions)
		regions = []string{"us-west-2"}
	}

	for _, region := range regions {
		ecrSession := session.Must(session.NewSession(&aws.Config{
			Region: aws.String(region),
		}))
		svc := ecr.New(ecrSession)
		authorizationToken, err := svc.GetAuthorizationToken(&ecr.GetAuthorizationTokenInput{
			RegistryIds: []*string{&opts.accountid},
		})
		if err != nil {
			return errors.Errorf("Error generating ecr credentials: %v", err)
		}
		println(authorizationToken)
	}
	return nil
}

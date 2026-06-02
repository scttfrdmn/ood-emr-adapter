package cmd

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/scttfrdmn/ood-emr-adapter/internal/awscfg"
	"github.com/spf13/cobra"
)

var (
	region        string
	applicationID string
)

var (
	assumeRoleArn     string // #78: per-user role to assume (empty = instance role)
	assumeRoleExtID   string
	assumeRoleSession string
)

var version = "dev" // overridden at release time via -ldflags -X .../cmd.version

var rootCmd = &cobra.Command{
	Version: version,
	Use:     "ood-emr-adapter",
	Short:   "OOD compute adapter for Amazon EMR Serverless",
	Long:    "Translates Open OnDemand job submissions to Amazon EMR Serverless API calls.",
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&region, "region", "us-east-1", "AWS region")
	rootCmd.PersistentFlags().StringVar(&applicationID, "application-id", "", "EMR Serverless application ID")
}

// #78: per-user cross-account AssumeRole flags. Empty (default) = use the OOD instance role.
func init() {
	pf := rootCmd.PersistentFlags()
	pf.StringVar(&assumeRoleArn, "assume-role-arn", "", "IAM role ARN to assume for AWS calls (empty = use the instance role)")
	pf.StringVar(&assumeRoleExtID, "assume-role-external-id", "", "sts:ExternalId for the assumed-role trust policy")
	pf.StringVar(&assumeRoleSession, "assume-role-session-name", "", "RoleSessionName for the assumed role (e.g. the OOD username)")
}

// awsOptions builds the AWS config options from the root flags (region + optional AssumeRole).
func awsOptions(ctx context.Context) []func(*config.LoadOptions) error {
	return awscfg.LoadOptions(ctx, awscfg.Options{
		Region:        region,
		AssumeRoleARN: assumeRoleArn,
		ExternalID:    assumeRoleExtID,
		SessionName:   assumeRoleSession,
	})
}

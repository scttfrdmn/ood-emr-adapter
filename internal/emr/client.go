// Package emr wraps the Amazon EMR Serverless API for the OOD adapter.
package emr

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/emrserverless"
	"github.com/aws/aws-sdk-go-v2/service/emrserverless/types"
)

// Client wraps the AWS EMR Serverless client.
type Client struct {
	svc    *emrserverless.Client
	region string
}

// New creates an EMR Serverless client using the default AWS credential chain.
func New(ctx context.Context, region string) (*Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("load AWS config: %w", err)
	}
	return &Client{svc: emrserverless.NewFromConfig(cfg), region: region}, nil
}

// JobRunSpec holds the parameters for an EMR Serverless job run.
type JobRunSpec struct {
	ApplicationID         string
	ExecutionRoleArn      string
	EntryPoint            string
	EntryPointArgs        []string
	SparkSubmitParameters string
	Env                   map[string]string
	JobName               string
}

// StartJobRun submits an EMR Serverless job run and returns the job run ID.
func (c *Client) StartJobRun(ctx context.Context, spec JobRunSpec) (string, error) {
	sparkProps := &types.SparkSubmit{
		EntryPoint: aws.String(spec.EntryPoint),
	}
	if len(spec.EntryPointArgs) > 0 {
		sparkProps.EntryPointArguments = spec.EntryPointArgs
	}
	if spec.SparkSubmitParameters != "" {
		sparkProps.SparkSubmitParameters = aws.String(spec.SparkSubmitParameters)
	}

	input := &emrserverless.StartJobRunInput{
		ApplicationId:    aws.String(spec.ApplicationID),
		ExecutionRoleArn: aws.String(spec.ExecutionRoleArn),
		JobDriver: &types.JobDriverMemberSparkSubmit{
			Value: *sparkProps,
		},
	}
	if spec.JobName != "" {
		input.Name = aws.String(spec.JobName)
	}
	if len(spec.Env) > 0 {
		envMap := make(map[string]string, len(spec.Env))
		for k, v := range spec.Env {
			envMap[k] = v
		}
		input.ConfigurationOverrides = &types.ConfigurationOverrides{
			MonitoringConfiguration: &types.MonitoringConfiguration{},
		}
		_ = envMap // env passed via SparkSubmitParameters in practice
	}

	out, err := c.svc.StartJobRun(ctx, input)
	if err != nil {
		return "", fmt.Errorf("emr StartJobRun: %w", err)
	}
	return aws.ToString(out.JobRunId), nil
}

// GetJobRun returns the current detail of an EMR Serverless job run.
func (c *Client) GetJobRun(ctx context.Context, applicationID, jobRunID string) (*emrserverless.GetJobRunOutput, error) {
	out, err := c.svc.GetJobRun(ctx, &emrserverless.GetJobRunInput{
		ApplicationId: aws.String(applicationID),
		JobRunId:      aws.String(jobRunID),
	})
	if err != nil {
		return nil, fmt.Errorf("emr GetJobRun: %w", err)
	}
	return out, nil
}

// CancelJobRun cancels an EMR Serverless job run.
func (c *Client) CancelJobRun(ctx context.Context, applicationID, jobRunID string) error {
	_, err := c.svc.CancelJobRun(ctx, &emrserverless.CancelJobRunInput{
		ApplicationId: aws.String(applicationID),
		JobRunId:      aws.String(jobRunID),
	})
	if err != nil {
		return fmt.Errorf("emr CancelJobRun: %w", err)
	}
	return nil
}

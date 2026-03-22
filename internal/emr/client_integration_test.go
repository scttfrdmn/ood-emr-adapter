//go:build integration

package emr_test

import (
	"context"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	awsemr "github.com/aws/aws-sdk-go-v2/service/emrserverless"
	substrate "github.com/scttfrdmn/substrate"

	. "github.com/scttfrdmn/ood-emr-adapter/internal/emr"
)

// substrateEMRClient returns a raw EMR Serverless SDK client pointed at the substrate server.
func substrateEMRClient(t *testing.T, endpointURL string) *awsemr.Client {
	t.Helper()
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion("us-east-1"),
		config.WithBaseEndpoint(endpointURL),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("test", "test", "")),
	)
	if err != nil {
		t.Fatalf("config: %v", err)
	}
	return awsemr.NewFromConfig(cfg)
}

// createTestApplication creates an EMR Serverless application in substrate and
// returns its application ID.
func createTestApplication(t *testing.T, ctx context.Context, raw *awsemr.Client) string {
	t.Helper()
	out, err := raw.CreateApplication(ctx, &awsemr.CreateApplicationInput{
		Name:         aws.String("ood-test-spark"),
		ReleaseLabel: aws.String("emr-7.0.0"),
		Type:         aws.String("SPARK"),
	})
	if err != nil {
		t.Fatalf("CreateApplication: %v", err)
	}
	appID := aws.ToString(out.ApplicationId)
	t.Logf("created EMR Serverless application: %s", appID)
	return appID
}

// TestStartGetCancelJobRun_Substrate exercises the full EMR Serverless job run
// lifecycle (StartJobRun → GetJobRun → CancelJobRun) against the substrate emulator.
func TestStartGetCancelJobRun_Substrate(t *testing.T) {
	ts := substrate.StartTestServer(t)
	t.Setenv("AWS_ENDPOINT_URL", ts.URL)
	t.Setenv("AWS_ACCESS_KEY_ID", "test")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "test")

	ctx := context.Background()
	raw := substrateEMRClient(t, ts.URL)
	appID := createTestApplication(t, ctx, raw)

	// Build the adapter client — picks up AWS_ENDPOINT_URL from the environment.
	client, err := New(ctx, "us-east-1")
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	spec := JobRunSpec{
		ApplicationID:    appID,
		ExecutionRoleArn: "arn:aws:iam::123456789012:role/EMRServerlessExecutionRole",
		EntryPoint:       "s3://my-bucket/scripts/analysis.py",
		EntryPointArgs:   []string{"--input", "s3://my-bucket/data/"},
		JobName:          "ood-integration-test",
	}

	// StartJobRun
	runID, err := client.StartJobRun(ctx, spec)
	if err != nil {
		t.Fatalf("StartJobRun: %v", err)
	}
	if runID == "" {
		t.Fatal("expected non-empty job run ID")
	}
	t.Logf("started job run: %s", runID)

	// GetJobRun
	detail, err := client.GetJobRun(ctx, appID, runID)
	if err != nil {
		t.Fatalf("GetJobRun: %v", err)
	}
	if detail == nil || detail.JobRun == nil {
		t.Fatal("GetJobRun: got nil job run")
	}
	t.Logf("job run state: %s", detail.JobRun.State)

	// CancelJobRun
	err = client.CancelJobRun(ctx, appID, runID)
	if err != nil {
		t.Fatalf("CancelJobRun: %v", err)
	}
	t.Log("job run cancelled successfully")
}

// TestGetJobRun_NotFound_Substrate verifies that GetJobRun returns an error
// for a run ID that was never created.
func TestGetJobRun_NotFound_Substrate(t *testing.T) {
	ts := substrate.StartTestServer(t)
	t.Setenv("AWS_ENDPOINT_URL", ts.URL)
	t.Setenv("AWS_ACCESS_KEY_ID", "test")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "test")

	ctx := context.Background()
	raw := substrateEMRClient(t, ts.URL)
	appID := createTestApplication(t, ctx, raw)

	client, err := New(ctx, "us-east-1")
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	_, err = client.GetJobRun(ctx, appID, "00000000000000000000000000000000")
	if err == nil {
		t.Fatal("expected error for non-existent job run, got nil")
	}
	if !strings.Contains(err.Error(), "emr") {
		t.Logf("error (acceptable): %v", err)
	}
}

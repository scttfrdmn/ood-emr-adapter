//go:build integration

package cmd

import (
	"context"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	awsemr "github.com/aws/aws-sdk-go-v2/service/emrserverless"
	substrate "github.com/scttfrdmn/substrate"
)

// runEMRCmd drives rootCmd with the given args, optionally feeding stdinData
// into os.Stdin. It captures os.Stdout and returns the trimmed output.
func runEMRCmd(t *testing.T, stdinData string, args []string) string {
	t.Helper()

	// Redirect stdin if caller supplied input.
	if stdinData != "" {
		stdinR, stdinW, err := os.Pipe()
		if err != nil {
			t.Fatalf("os.Pipe (stdin): %v", err)
		}
		origStdin := os.Stdin
		os.Stdin = stdinR
		t.Cleanup(func() { os.Stdin = origStdin })

		if _, err := io.WriteString(stdinW, stdinData); err != nil {
			t.Fatalf("write stdin pipe: %v", err)
		}
		stdinW.Close()
	}

	// Redirect stdout so we can capture command output.
	stdoutR, stdoutW, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe (stdout): %v", err)
	}
	origStdout := os.Stdout
	os.Stdout = stdoutW
	t.Cleanup(func() { os.Stdout = origStdout })

	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true
	rootCmd.SetArgs(args)
	if err := rootCmd.Execute(); err != nil {
		t.Logf("rootCmd.Execute error (may be expected): %v", err)
	}

	stdoutW.Close()
	os.Stdout = origStdout

	out, err := io.ReadAll(stdoutR)
	if err != nil {
		t.Fatalf("read stdout pipe: %v", err)
	}
	return strings.TrimSpace(string(out))
}

// TestCLISubmitStatusDelete_EMR_Substrate exercises the full EMR Serverless
// adapter CLI lifecycle (submit → status → delete) against the substrate
// emulator. An EMR Serverless application is created via the raw SDK before
// the CLI commands are exercised, as the adapter requires a pre-existing
// application.
func TestCLISubmitStatusDelete_EMR_Substrate(t *testing.T) {
	ts := substrate.StartTestServer(t)
	t.Setenv("AWS_ENDPOINT_URL", ts.URL)
	t.Setenv("AWS_ACCESS_KEY_ID", "test")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "test")

	// Create a raw EMR Serverless client pointed at the substrate server.
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
		config.WithBaseEndpoint(ts.URL),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("test", "test", "")),
	)
	if err != nil {
		t.Fatalf("config.LoadDefaultConfig: %v", err)
	}
	rawClient := awsemr.NewFromConfig(cfg)

	// Create the EMR Serverless application that the CLI will target.
	createOut, err := rawClient.CreateApplication(ctx, &awsemr.CreateApplicationInput{
		Name:         aws.String("ood-cli-spark"),
		ReleaseLabel: aws.String("emr-7.0.0"),
		Type:         aws.String("SPARK"),
	})
	if err != nil {
		t.Fatalf("CreateApplication: %v", err)
	}
	appID := aws.ToString(createOut.ApplicationId)
	t.Logf("created EMR Serverless application: %s", appID)

	jobSpec := `{"entry_point":"s3://bucket/job.py","execution_role_arn":"arn:aws:iam::123456789012:role/EMRRole","job_name":"cli-emr-test"}`

	// Submit
	compoundID := runEMRCmd(t, jobSpec, []string{"submit", "--region", "us-east-1", "--application-id", appID})
	if compoundID == "" {
		t.Fatal("submit: expected non-empty compound ID, got empty string")
	}
	if !strings.Contains(compoundID, appID) {
		t.Errorf("submit: compound ID %q does not contain application ID %q", compoundID, appID)
	}
	t.Logf("submitted EMR job run: %s", compoundID)

	// Status
	statusOut := runEMRCmd(t, "", []string{"status", "--region", "us-east-1", compoundID})
	if statusOut == "" {
		t.Fatal("status: expected non-empty output")
	}
	t.Logf("status output: %s", statusOut)
	if !strings.Contains(statusOut, "success") &&
		!strings.Contains(statusOut, "running") &&
		!strings.Contains(statusOut, "completed") {
		t.Errorf("status: output %q does not contain a recognised terminal state", statusOut)
	}

	// Delete (CancelJobRun)
	runEMRCmd(t, "", []string{"delete", "--region", "us-east-1", compoundID})
	t.Logf("deleted EMR job run: %s", compoundID)
}

package cmd

import (
	"testing"

	emrtypes "github.com/aws/aws-sdk-go-v2/service/emrserverless/types"
)

func TestEmrStateToOod(t *testing.T) {
	tests := []struct {
		state    emrtypes.JobRunState
		expected string
	}{
		{emrtypes.JobRunStateSubmitted, "queued"},
		{emrtypes.JobRunStatePending, "queued"},
		{emrtypes.JobRunStateScheduled, "queued"},
		{emrtypes.JobRunStateQueued, "queued"},
		{emrtypes.JobRunStateRunning, "running"},
		{emrtypes.JobRunStateSuccess, "completed"},
		{emrtypes.JobRunStateFailed, "failed"},
		{emrtypes.JobRunStateCancelling, "cancelled"},
		{emrtypes.JobRunStateCancelled, "cancelled"},
		{emrtypes.JobRunState("UNKNOWN_STATE"), "undetermined"},
	}

	for _, tt := range tests {
		t.Run(string(tt.state), func(t *testing.T) {
			got := emrStateToOod(tt.state)
			if got != tt.expected {
				t.Errorf("emrStateToOod(%q) = %q, want %q", tt.state, got, tt.expected)
			}
		})
	}
}

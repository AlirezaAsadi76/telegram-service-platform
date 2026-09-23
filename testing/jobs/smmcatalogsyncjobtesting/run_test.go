package smmcatalogsyncjobtesting

import (
	"context"
	"errors"
	"telegram-service-platform/scheduler/jobs/smmcatalogsyncjob"
	"testing"
)

func TestJob_Run(t *testing.T) {
	tests := []struct {
		name       string
		serviceErr error
		wantErr    bool
		wantCalls  int
	}{
		{
			name:      "success",
			wantCalls: 1,
		},
		{
			name: "product service error",
			serviceErr: errors.New(
				"catalog sync failed",
			),
			wantErr:   true,
			wantCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &fakeProductService{
				err: tt.serviceErr,
			}

			job := smmcatalogsyncjob.New(
				service,
			)

			err := job.Run(
				context.Background(),
			)

			if tt.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf(
					"unexpected error: %v",
					err,
				)
			}

			if service.calls != tt.wantCalls {
				t.Fatalf(
					"expected %d SyncSMMServices calls, got %d",
					tt.wantCalls,
					service.calls,
				)
			}
		})
	}
}

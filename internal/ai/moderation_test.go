package ai

  import (
        "math"
        "strings"
        "testing"

        "github.com/zyj/my-blog/pkg/constant"
  )

  func TestValidateModerationResult(t *testing.T) {
        tests := []struct {
                name    string
                result  ModerationResult
                wantErr bool
        }{
                {
                        name: "approved",
                        result: ModerationResult{
                                Status:     constant.ModerationApproved,
                                Confidence: 0.8,
                        },
                },
                {
                        name: "rejected",
                        result: ModerationResult{
                                Status:     constant.ModerationRejected,
                                Confidence: 0.6,
                        },
                },
                {
                        name: "manual review",
                        result: ModerationResult{
                                Status:     constant.ModerationManualReview,
                                Confidence: 0,
                        },
                },
                {
                        name: "invalid status",
                        result: ModerationResult{
                                Status:     constant.ModerationPending,
                        },
                        wantErr: true,
                },
                {
                        name: "negative confidence",
                        result: ModerationResult{
                                Status:     constant.ModerationApproved,
                                Confidence: -0.1,
                        },
                        wantErr: true,
                },
                {
                        name: "confidence greater than one",
                        result: ModerationResult{
                                Status:     constant.ModerationApproved,
                                Confidence: 1.1,
                        },
                        wantErr: true,
                },
                {
                        name: "nan confidence",
                        result: ModerationResult{
                                Status:     constant.ModerationApproved,
                                Confidence: math.NaN(),
                        },
                        wantErr: true,
                },
                {
                        name: "too many categories",
                        result: ModerationResult{
                                Status:     constant.ModerationApproved,
                                Categories: make([]string, 11),
                        },
                        wantErr: true,
                },
                {
                        name: "reason too long",
                        result: ModerationResult{
                                Status: constant.ModerationApproved,
                                Reason: strings.Repeat("a", 501),
                        },
                        wantErr: true,
                },
        }

        for _, tt := range tests {
                t.Run(tt.name, func(t *testing.T) {
                        err := ValidateModerationResult(tt.result)
                        if (err != nil) != tt.wantErr {
                                t.Fatalf("unexpected error state: %v", err)
                        }
                })
        }
  }
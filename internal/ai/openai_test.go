package ai

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/zyj/my-blog/pkg/constant"
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}

func TestOpenAIProviderModerate(t *testing.T) {
	provider := NewOpenAIProvider("https://provider.test", "test-key", "test-model")
	provider.HTTPClient = &http.Client{Transport: roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/chat/completions" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Fatalf("unexpected authorization: %s", r.Header.Get("Authorization"))
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Fatalf("unexpected content type: %s", r.Header.Get("Content-Type"))
		}

		var requestBody openAIChatRequest
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			t.Fatal(err)
		}
		if requestBody.Model != "test-model" {
			t.Fatalf("unexpected model: %s", requestBody.Model)
		}
		if len(requestBody.Messages) != 2 {
			t.Fatalf("unexpected messages: %#v", requestBody.Messages)
		}
		if !strings.Contains(requestBody.Messages[1].Content, "test content") {
			t.Fatalf("unexpected user content: %s", requestBody.Messages[1].Content)
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"choices":[{"message":{"role":"assistant","content":"{\"status\":\"approved\",\"categories\":[\"safe\"],\"confidence\":0.95,\"reason\":\"approved\"}"}}]}`)),
			Header:     make(http.Header),
		}, nil
	})}

	result, err := provider.Moderate(context.Background(), ModerationInput{
		CommentID:  1,
		Content:    "test content",
		TargetType: constant.TargetComment,
		TargetID:   2,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != constant.ModerationApproved || result.Confidence != 0.95 {
		t.Fatalf("unexpected result: %#v", result)
	}
	if len(result.Categories) != 1 || result.Categories[0] != "safe" {
		t.Fatalf("unexpected categories: %#v", result.Categories)
	}
}

func TestOpenAIProviderModerateHTTPError(t *testing.T) {
	provider := NewOpenAIProvider("https://provider.test", "test-key", "test-model")
	provider.HTTPClient = &http.Client{Transport: roundTripperFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       io.NopCloser(strings.NewReader("provider failed")),
			Header:     make(http.Header),
		}, nil
	})}

	_, err := provider.Moderate(context.Background(), ModerationInput{Content: "test content"})
	if err == nil || !strings.Contains(err.Error(), "status 500") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestOpenAIProviderModerateEmptyChoices(t *testing.T) {
	provider := NewOpenAIProvider("https://provider.test", "test-key", "test-model")
	provider.HTTPClient = &http.Client{Transport: roundTripperFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"choices":[]}`)),
			Header:     make(http.Header),
		}, nil
	})}

	_, err := provider.Moderate(context.Background(), ModerationInput{Content: "test content"})
	if err == nil || !strings.Contains(err.Error(), "no choices") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewOpenAIProviderFromEnv(t *testing.T) {
	t.Setenv(constant.EnvKeyAIModerationBaseURL, "https://env-provider.test/v1")
	t.Setenv(constant.EnvKeyAIModerationAPIKey, "env-key")
	t.Setenv(constant.EnvKeyAIModerationModel, "env-model")

	provider := NewOpenAIProviderFromEnv()

	if provider.BaseURL != "https://env-provider.test/v1" {
		t.Fatalf("unexpected base url: %s", provider.BaseURL)
	}
	if provider.APIKey != "env-key" {
		t.Fatalf("unexpected api key: %s", provider.APIKey)
	}
	if provider.Model != "env-model" {
		t.Fatalf("unexpected model: %s", provider.Model)
	}
	if provider.HTTPClient == nil {
		t.Fatal("expected http client")
	}
}

func TestOpenAIProviderValidateConfig(t *testing.T) {
	tests := []struct {
		name     string
		provider *OpenAIProvider
		wantErr  string
	}{
		{
			name:    "nil provider",
			wantErr: "openai provider is not configured",
		},
		{
			name: "nil http client",
			provider: &OpenAIProvider{
				BaseURL: "https://provider.test",
				APIKey:  "test-key",
				Model:   "test-model",
			},
			wantErr: "openai provider is not configured",
		},
		{
			name: "missing base url",
			provider: &OpenAIProvider{
				APIKey:     "test-key",
				Model:      "test-model",
				HTTPClient: http.DefaultClient,
			},
			wantErr: "openai provider configuration is incomplete",
		},
		{
			name: "complete",
			provider: NewOpenAIProvider(
				"https://provider.test",
				"test-key",
				"test-model",
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.provider.ValidateConfig()
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil || err.Error() != tt.wantErr {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

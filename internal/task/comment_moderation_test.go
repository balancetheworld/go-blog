package task

import (
	"encoding/json"
	"testing"
)

func TestNewCommentModerationTask(t *testing.T) {
	task, err := NewCommentModerationTask(12)
	if err != nil {
		t.Fatal(err)
	}
	if task.Type() != TypeCommentModeration {
		t.Fatalf("unexpected task type: %s", task.Type())
	}

	var payload CommentModerationPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		t.Fatal(err)
	}
	if payload.AITaskID != 12 {
		t.Fatalf("unexpected task payload: %#v", payload)
	}
}

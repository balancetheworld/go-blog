package constant

type AITaskType string

const (
	AITaskCommentModeration AITaskType = "comment_moderation"
)

type AITaskStatus string

const (
	AITaskQueued     AITaskStatus = "queued"
	AITaskProcessing AITaskStatus = "processing"
	AITaskSucceeded  AITaskStatus = "succeeded"
	AITaskFailed     AITaskStatus = "failed"
	AITaskDead       AITaskStatus = "dead"
)

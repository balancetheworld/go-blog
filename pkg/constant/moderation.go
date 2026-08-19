package constant

type ModerationStatus string

const (
	ModerationPending      ModerationStatus = "pending"
	ModerationApproved     ModerationStatus = "approved"
	ModerationRejected     ModerationStatus = "rejected"
	ModerationManualReview ModerationStatus = "manual_review"
	ModerationHidden       ModerationStatus = "hidden"
)

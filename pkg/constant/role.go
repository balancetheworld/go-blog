package constant

type RoleApplicationStatus string

const (
	RoleCodeMember = "member"
	RoleCodeEditor = "editor"
	RoleCodeAdmin  = "admin"
)

const (
	RoleApplicationPending  RoleApplicationStatus = "pending"
	RoleApplicationApproved RoleApplicationStatus = "approved"
	RoleApplicationRejected RoleApplicationStatus = "rejected"
)

package constant

type PostVisibility string

const (
	PostVisibilityPublic  PostVisibility = "public"
	PostVisibilityRoles   PostVisibility = "roles"
	PostVisibilityPrivate PostVisibility = "private"
)

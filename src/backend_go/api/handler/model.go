package handler

const (
	ActionDelete  = "delete"
	ActionRestore = "restore"
)

// TokenResponse is a struct representing the response from a token generation request
type TokenResponse struct {
	Token string `json:"token"`
}

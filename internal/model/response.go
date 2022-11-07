package model

// TODO generic
type CreateUserResponse struct {
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

type GetByIdResponse struct {
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
	*User
}

type MatchResponse struct {
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

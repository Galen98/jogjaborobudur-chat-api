package dto

type CreateChatSessionRequest struct {
	Token       string `json:"token"`
	Session     string `json:"session"`
	ProductId   uint   `json:"product_id"`
	Thumbnail   string `json:"thumbnail"`
	ProductName string `json:"product_name"`
}

type CheckChatTokenRequest struct {
	Session   string `json:"session"`
	ProductId uint   `json:"product_id"`
}

type CheckChatTokenResponse struct {
	Token    string `json:"token"`
	IsActive bool   `json:"is_active"`
}

type GetAllSessionUserRequest struct {
	Session string `json:"session"`
}

type GetAllSessionUserResponse struct {
	ProductName string `json:"product_name"`
	Thumbnail   string `json:"thumbnail"`
	Token       string `json:"token"`
	UpdatedAt   string `json:"updated_at"`
}

type GetAllSessionAdminResponse struct {
	UserID      uint   `json:"user_id"`
	UserName    string `json:"user_name"`
	ProductName string `json:"product_name"`
	Thumbnail   string `json:"thumbnail"`
	UpdatedAt   string `json:"updated_at"`
	Token       string `json:"token"`
}

type GetChatSessionRequest struct {
	SessionID string `json:"session_id"`
	ProductID uint   `json:"product_id"`
}

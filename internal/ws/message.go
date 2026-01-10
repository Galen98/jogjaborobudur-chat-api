package ws

type WSBase struct {
	Type string `json:"type"`
}

type WSChatMessage struct {
	Type       string `json:"type"`
	Token      string `json:"token"`
	SenderType string `json:"sender_type"`
	Message    string `json:"message"`
}

type WSSessionUpdate struct {
	Type        string `json:"type"`
	Token       string `json:"token"`
	UpdatedAt   string `json:"updated_at"`
	IsRead      bool   `json:"is_read"`
	IsReadAdmin bool   `json:"is_read_admon"`
}

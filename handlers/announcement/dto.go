package announcement

type CreateChannelReq struct {
	EntityId    uint   `json:"entity_id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

type SSEClient struct {
	Channel chan string
}

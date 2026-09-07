package telegramutils

import (
	"encoding/json"
	"fmt"
	"telegram-service-platform/entity/telegramentity"
)

func ParseTelegramUser(userJSON string) (*telegramentity.User, error) {
	if userJSON == "" {
		return nil, fmt.Errorf("user data is empty in init data")
	}

	var user telegramentity.User
	if err := json.Unmarshal([]byte(userJSON), &user); err != nil {
		return nil, fmt.Errorf("failed to parse telegram user json: %w", err)
	}

	return &user, nil
}

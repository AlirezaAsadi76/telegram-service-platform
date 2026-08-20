package mainhandler

import (
	"fmt"
	"strings"
)

const (
	platformSplitMode int = iota + 2
	CategorySplitMode
	ServiceSplitMode
)

func (h *Handler) callbackQueryData(callbackData string, splitMode int) (data []string, err error) {
	parts := strings.Split(callbackData, ":")
	if len(parts) < splitMode+1 {
		return nil, fmt.Errorf("داده نامعتبر")
	}
	return parts, nil
}

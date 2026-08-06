package telegramproduct

import (
	"fmt"
	"telegram-service-platform/adapter/fzrcards"
)

type Adapter struct {
	client *fzrcards.FzrClient
}

func New(client *fzrcards.FzrClient) Adapter {
	return Adapter{
		client: client,
	}
}

func (a *Adapter) createURL(path string) string {
	return fmt.Sprintf("%s/%s", a.client.BaseURL(), path)
}

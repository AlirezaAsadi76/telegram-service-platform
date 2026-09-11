package helpers

import "strings"

func BuildPaymentURL(StartPayURL string, authority string) string {
	return strings.TrimRight(StartPayURL, "/") + "/" + authority
}

package smmentity

type Platform struct {
	Name PlatformType
}

type PlatformType string

const (
	TelegramPlatform  PlatformType = "Telegram"
	WhatsappPlatform  PlatformType = "Whatsapp"
	InstagramPlatform PlatformType = "Instagram"
	TickTockPlatform  PlatformType = "TickTock"
	XxPlatform        PlatformType = "twitter(X)"
)

func (p PlatformType) String() string {
	return string(p)
}

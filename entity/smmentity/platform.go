package smmentity

type Platform struct {
	Name PlatformType
}

type PlatformType string

const (
	TelegramPlatform  PlatformType = "telegram"
	WhatsappPlatform  PlatformType = "whatsapp"
	InstagramPlatform PlatformType = "instagram"
	TickTockPlatform  PlatformType = "tickTock"
	XxPlatform        PlatformType = "twitter(X)"
)

func (p PlatformType) String() string {
	return string(p)
}

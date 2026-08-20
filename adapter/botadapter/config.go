package botadapter

type Config struct {
	Token string `koanf:"token"`
	Debug bool   `koanf:"debug"`
}

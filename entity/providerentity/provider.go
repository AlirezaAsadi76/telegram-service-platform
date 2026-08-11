package providerentity

type Provider struct {
	ID       uint64
	Name     string       // "justanotherpanel", "caspersmm"
	Type     ProviderType // SMM / PAYMENT / EXCHANGE
	BaseURL  string
	APIKey   string
	Config   map[string]any // تنظیمات اضافی
	Priority int            // اولویت (کمتر = اولویت بالاتر)
	IsActive bool
}

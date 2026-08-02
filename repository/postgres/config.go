package postgres

type DBConfig struct {
	Host     string `koanf:"host"`
	Port     uint16 `koanf:"port"`
	User     string `koanf:"user"`
	Password string `koanf:"password"`
	Database string `koanf:"database"`
	SSLMode  string `koanf:"sslmode"`
}

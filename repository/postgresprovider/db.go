package postgresprovider

import (
	"encoding/json"
	"telegram-service-platform/entity/providerentity"
	"telegram-service-platform/repository/postgres"
)

type DB struct {
	Pool *postgres.DB
}

func New(pool *postgres.DB) *DB {
	return &DB{Pool: pool}
}

func scanProvider(row postgres.Scanner) (providerentity.Provider, error) {
	provider := providerentity.Provider{}

	var config []byte

	var _ []byte
	err := row.Scan(
		&provider.ID,
		&provider.Name,
		&provider.Type,
		&provider.BaseURL,
		&provider.APIKey,
		&config,
		&provider.Priority,
		&provider.IsActive,
	)

	if len(config) > 0 {

		err = json.Unmarshal(
			config,
			&provider.Config,
		)

	}
	return provider, err

}

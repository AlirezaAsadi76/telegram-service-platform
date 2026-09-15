package postgrescheckout

import "telegram-service-platform/repository/postgres"

type DB struct {
	transactionProvider postgres.TransactionProvider
}

func New(transactionProvider postgres.TransactionProvider) *DB {
	return &DB{
		transactionProvider: transactionProvider,
	}
}

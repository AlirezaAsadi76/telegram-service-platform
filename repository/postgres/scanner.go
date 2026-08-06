package postgres

type Scanner interface {
	Scan(...interface{}) error
}

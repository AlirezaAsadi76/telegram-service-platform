package postgres

func (d DB) Close() {
	d.db.Close()
}

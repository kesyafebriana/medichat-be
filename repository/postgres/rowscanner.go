package postgres

type RowScanner interface {
	Scan(dests ...any) error
}

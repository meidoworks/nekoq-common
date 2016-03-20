package database

type Database interface {
	Close() error

	GetTable(string) (Table, error)

	GetTx() (Tx, error)
}

type Table interface {
}

type Tx interface {
	Get(Table, []byte) ([]byte, error)
	Set(table Table, key, value []byte) error
}

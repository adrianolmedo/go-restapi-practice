package postgres

import (
	"fmt"

	config "github.com/adrianolmedo/genesis"
)

// Storage represents all repositories.
type Storage struct {
	User     User
	Product  Product
	Customer Customer
	Invoice  Invoice
}

// NewStorage init postgres database conecton, build the Storage and return
// it its pointer.
func NewStorage(dbcfg config.DB) (*Storage, error) {
	if dbcfg.Engine == "" {
		return nil, fmt.Errorf("database engine not selected")
	}

	if dbcfg.Engine == "postgres" {
		db, err := newDB(dbcfg)
		if err != nil {
			return nil, fmt.Errorf("postgres: %v", err)
		}

		return &Storage{
			User:     User{db: db},
			Product:  Product{db: db},
			Customer: Customer{db: db},
			Invoice: Invoice{
				db:     db,
				header: InvoiceHeader{db: db},
				items:  InvoiceItem{db: db},
			},
		}, nil
	}

	return nil, fmt.Errorf("database engine '%s' not implemented", dbcfg.Engine)
}

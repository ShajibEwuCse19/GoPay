package repo

import (
	"github.com/ShajibEwuCse19/gopay-core/domain"
	"github.com/ShajibEwuCse19/gopay-core/services/customer"
	"github.com/jmoiron/sqlx"
)

type customerRepo struct {
	db *sqlx.DB
}

func NewCustomerRepo(db *sqlx.DB) services.CustomerRepo {
	return &customerRepo{
		db: db,
	}
}

func (r *customerRepo) Create(c domain.Customer) (*domain.Customer, error) {
	query := `
	INSERT INTO customers (phone_number, email, full_name, pin_hash)
	VALUES (:phone_number, :email, :full_name, :pin_hash)
	RETURNING *
	`

	rows, err := r.db.NamedQuery(query, c)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customer domain.Customer
	if rows.Next() {
		err = rows.StructScan(&customer)
		if err != nil {
			return nil, err
		}
	}
	return &customer, nil
}

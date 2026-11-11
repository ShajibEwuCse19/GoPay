package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ShajibEwuCse19/gopay-core/domain"
	"github.com/ShajibEwuCse19/gopay-core/services/customer"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// pgUniqueViolation is the PostgreSQL error code for a unique constraint violation.
const pgUniqueViolation = "23505"

type customerRepo struct {
	db *sqlx.DB
}

func NewCustomerRepo(db *sqlx.DB) services.CustomerRepo {
	return &customerRepo{
		db: db,
	}
}

func (r *customerRepo) Create(ctx context.Context, c domain.Customer) (*domain.Customer, error) {
	query := `
	INSERT INTO customers (phone_number, email, full_name, pin_hash)
	VALUES (:phone_number, :email, :full_name, :pin_hash)
	RETURNING *
	`

	rows, err := r.db.NamedQueryContext(ctx, query, c)
	if err != nil {
		return nil, mapCreateError(err)
	}
	defer rows.Close()

	// rows.Next() returns false on error too (e.g. ctx cancelled or timed out while reading)
	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, mapCreateError(err)
		}
		return nil, sql.ErrNoRows
	}

	var customer domain.Customer
	if err := rows.StructScan(&customer); err != nil {
		return nil, err
	}
	return &customer, nil
}

// mapCreateError turns driver errors into domain errors the upper layers understand.
func mapCreateError(err error) error {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) && pqErr.Code == pgUniqueViolation {
		return domain.ErrCustomerAlreadyExists
	}
	return err
}

package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ShajibEwuCse19/gopay-core/domain"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// pgUniqueViolation is the PostgreSQL error code for a unique constraint violation.
const pgUniqueViolation = "23505"

type CustomerRepo struct {
	db *sqlx.DB
}

func NewCustomerRepo(db *sqlx.DB) *CustomerRepo {
	return &CustomerRepo{
		db: db,
	}
}

func (r *CustomerRepo) Create(ctx context.Context, c domain.Customer) (*domain.Customer, error) {
	query := `
	INSERT INTO customers (phone_number, email, full_name, pin_hash)
	VALUES (:phone_number, :email, :full_name, :pin_hash)
	RETURNING id, phone_number, email, full_name, pin_hash, kyc_status, status, created_at, updated_at
	`

	// creating rows in db and execute INSERT query using c (domain.Customer)
	rows, err := r.db.NamedQueryContext(ctx, query, c)
	if err != nil {
		return nil, mapDBError(err)
	}
	defer rows.Close()

	// rows.Next() returns false on error too (e.g. ctx cancelled or timed out while reading)
	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, mapDBError(err)
		}
		return nil, sql.ErrNoRows
	}

	// saving rows into Customer from database
	// [KycStatus, Status, CreatedAt, etc. auto-generated fields will be updated in the domain]
	var customer domain.Customer
	if err := rows.StructScan(&customer); err != nil {
		return nil, err
	}
	return &customer, nil
}

// mapDBError turns driver errors into domain errors the upper layers understand.
func mapDBError(err error) error {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case pgUniqueViolation:
			return domain.ErrConflict
		}
	}
	return err
}

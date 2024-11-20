package repository

import "github.com/jackc/pgtype"

type Product struct {
	ID    pgtype.Int4   `db:"id"`
	Name  pgtype.Text   `db:"name"`
	Price pgtype.Float8 `db:"price"`
}

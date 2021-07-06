package postgres

import (
	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Readiness() error {
	return r.db.Ping()
}

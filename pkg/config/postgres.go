package config

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type pgUserRepo struct {
	db *sqlx.DB
}

// NewPostgresUserRepo returns a Postgres backed UserRepo implementation.
func NewPostgresUserRepo(db *sqlx.DB) (UserRepo, error) {
	return &pgUserRepo{
		db: db,
	}, nil
}

func (r *pgUserRepo) Get(baseName, id string) (*UserConfig, error) {
	return nil, fmt.Errorf("pgUserRepo.Get() not implemented")
}

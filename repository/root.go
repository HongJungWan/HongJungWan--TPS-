package repository

import (
	"database/sql"
	"go_chat/config"

	_ "github.com/go-sql-driver/mysql"
)

type Repository struct {
	cfg *config.Config

	db *sql.DB
}

func NewRepository(cfg *config.Config) (*Repository, error) {
	repository := &Repository{cfg: cfg}

	var err error
	if repository.db, err = sql.Open(cfg.DB.Database, cfg.DB.URL); err != nil {
		return nil, err
	}

	return repository, nil
}

package repository

import (
	"database/sql"
	"go_chat/config"
	"go_chat/types/schema"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Repository struct {
	cfg *config.Config

	db *sql.DB
}

const (
	room       = "chatting.room"
	chat       = "chatting.chat"
	serverInfo = "chatting.serverInfo"
)

func NewRepository(cfg *config.Config) (*Repository, error) {
	repository := &Repository{cfg: cfg}

	var err error
	if repository.db, err = sql.Open(cfg.DB.Database, cfg.DB.URL); err != nil {
		return nil, err
	}

	return repository, nil
}

func (s *Repository) Room(name string) (*schema.Room, error) {
	domain := new(schema.Room)
	queryString := query([]string{})

	err := s.db.QueryRow(queryString, name).Scan(
		&domain.ID,
		&domain.Name,
		&domain.CreateAt,
		&domain.UpdateAt,
	)
	return domain, err
}

func query(queryString []string) string {
	return strings.Join(queryString, " ") + ";"
}

package repository

import (
	"database/sql"
	"errors"
	"go_chat/config"
	"go_chat/types/schema"
	"log"
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

func (repository *Repository) InsertChatting(user, message, roomName string) error {
	log.Println("Insert Chatting Using WSS", "from", user, "message", message, "room", roomName)
	_, err := repository.db.Exec("INSERT INTO chatting.chat(room, name, message) VALUES(?,?,?)", roomName, user, message)

	return err
}

func (repository *Repository) GetChatList(roomName string) ([]*schema.Chat, error) {
	queryString := query([]string{"SELECT * FROM", chat, "WHERE room = ? ORDER BY `when` DESC LIMIT 10"})

	if cursor, err := repository.db.Query(queryString, roomName); err != nil {
		return nil, err
	} else {
		defer cursor.Close()

		var result []*schema.Chat

		for cursor.Next() {
			d := new(schema.Chat)

			if err = cursor.Scan(&d.ID, &d.Room, &d.Name, &d.Message, &d.When); err != nil {
				return nil, err
			} else {
				result = append(result, d)
			}
		}

		if len(result) == 0 {
			return []*schema.Chat{}, nil
		} else {
			return result, nil
		}
	}
}

func (repository *Repository) RoomList() ([]*schema.Room, error) {
	// TODO: 페이지네이션 기능 추가 필요
	queryString := query([]string{"SELECT * FROM", room})

	if cursor, err := repository.db.Query(queryString); err != nil {
		return nil, err
	} else {
		defer cursor.Close()

		var result []*schema.Room

		for cursor.Next() {
			domain := new(schema.Room)

			if err = cursor.Scan(&domain.ID, &domain.Name, &domain.CreateAt, &domain.UpdateAt); err != nil {
				return nil, err
			} else {
				result = append(result, domain)
			}
		}

		if len(result) == 0 {
			return []*schema.Room{}, nil
		} else {
			return result, nil
		}
	}
}

func (repository *Repository) MakeRoom(roomName string) error {
	_, err := repository.db.Exec("INSERT INTO chatting.room(name) VALUES(?)", roomName)
	return err
}

func (repository *Repository) Room(roomName string) (*schema.Room, error) {
	domain := new(schema.Room)
	queryString := query([]string{"SELECT * FROM", room, "WHERE name = ?"})

	err := repository.db.QueryRow(queryString, roomName).
		Scan(&domain.ID, &domain.Name, &domain.CreateAt, &domain.UpdateAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return domain, err
}

func query(queryString []string) string {
	return strings.Join(queryString, " ") + ";"
}

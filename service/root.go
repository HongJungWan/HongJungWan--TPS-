package service

import (
	"go_chat/repository"
	"go_chat/types/schema"
	"log"
)

type Service struct {
	repository *repository.Repository
}

func NewService(repository *repository.Repository) *Service {
	service := &Service{repository: repository}
	return service
}

func (service *Service) EnterRoom(roomName string) ([]*schema.Chat, error) {
	if response, err := service.repository.GetChatList(roomName); err != nil {
		log.Println("Failed To Get EnterRoom", "err", err.Error())
		return nil, err
	} else {
		return response, nil
	}
}

func (service *Service) RoomList() ([]*schema.Room, error) {
	if res, err := service.repository.RoomList(); err != nil {
		log.Println("Failed To Get All Room List", "err", err.Error())
		return nil, err
	} else {
		return res, nil
	}
}

func (service *Service) MakeRoom(name string) error {
	if err := service.repository.MakeRoom(name); err != nil {
		log.Println("Failed To Make New Room", "err", err.Error())
		return err
	} else {
		return nil
	}
}

func (service *Service) Room(name string) (*schema.Room, error) {
	if res, err := service.repository.Room(name); err != nil {
		log.Println("Failed To Get Room ", "err", err.Error())
		return nil, err
	} else {
		return res, nil
	}
}

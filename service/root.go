package service

import "go_chat/repository"

type Service struct {
	repository *repository.Repository
}

func NewService(repository *repository.Repository) *Service {
	service := &Service{repository: repository}
	return service
}

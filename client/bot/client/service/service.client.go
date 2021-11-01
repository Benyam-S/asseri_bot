package service

import (
	"errors"
	"regexp"

	"github.com/Benyam-S/asseri/client/bot"
	"github.com/Benyam-S/asseri/client/bot/client"
)

// Service is a type that defines a client service
type Service struct {
	clientRepo client.IClientRepository
}

// NewClientService is a function that returns a new client service
func NewClientService(clientRepository client.IClientRepository) client.IService {
	return &Service{clientRepo: clientRepository}
}

// AddClient is a method that adds a new client to the system
func (service *Service) AddClient(newClient *bot.Client) error {
	err := service.clientRepo.Create(newClient)
	if err != nil {
		return errors.New("unable to add new client")
	}

	return nil
}

// FindClient is a method that find and return a client that matches the identifier value
func (service *Service) FindClient(identifier string) (*bot.Client, error) {

	empty, _ := regexp.MatchString(`^\s*$`, identifier)
	if empty {
		return nil, errors.New("no client found")
	}

	client, err := service.clientRepo.Find(identifier)
	if err != nil {
		return nil, errors.New("no client found")
	}
	return client, nil
}

// UpdateClient is a method that updates a client in the system
func (service *Service) UpdateClient(client *bot.Client) error {
	err := service.clientRepo.Update(client)
	if err != nil {
		return errors.New("unable to update client")
	}

	return nil
}

// DeleteClient is a method that deletes a client from the system
func (service *Service) DeleteClient(identifier string) (*bot.Client, error) {
	client, err := service.clientRepo.Delete(identifier)
	if err != nil {
		return nil, errors.New("unable to delete client")
	}

	return client, nil
}

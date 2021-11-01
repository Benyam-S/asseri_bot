package client

import (
	"github.com/Benyam-S/asseri/client/bot"
)

// IService is an interface that defines all the service methods of a bot.Client struct
type IService interface {
	AddClient(newClient *bot.Client) error
	FindClient(identifier string) (*bot.Client, error)
	UpdateClient(client *bot.Client) error
	DeleteClient(identifier string) (*bot.Client, error)
}

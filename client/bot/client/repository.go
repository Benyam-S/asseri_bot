package client

import (
	"github.com/Benyam-S/asseri/client/bot"
)

// IClientRepository is an interface that defines all the repository methods of a bot.Client struct
type IClientRepository interface {
	Create(newClient *bot.Client) error
	Find(identifier string) (*bot.Client, error)
	Update(client *bot.Client) error
	Delete(identifier string) (*bot.Client, error)
}

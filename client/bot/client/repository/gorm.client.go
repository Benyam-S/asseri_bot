package repository

import (
	"github.com/Benyam-S/asseri/client/bot"
	"github.com/Benyam-S/asseri/client/bot/client"
	"github.com/jinzhu/gorm"
)

// ClientRepository is a type that defines a client repository type
type ClientRepository struct {
	conn *gorm.DB
}

// NewClientRepository is a function that creates a new client repository type
func NewClientRepository(connection *gorm.DB) client.IClientRepository {
	return &ClientRepository{conn: connection}
}

// Create is a method that adds a new client to the database
func (repo *ClientRepository) Create(newClient *bot.Client) error {
	err := repo.conn.Create(newClient).Error
	if err != nil {
		return err
	}
	return nil
}

// Find is a method that returns a client that matches the provided identifier.
// In Find() telegram_id or user_id can either be used as a key.
func (repo *ClientRepository) Find(identifier string) (*bot.Client, error) {
	client := new(bot.Client)
	err := repo.conn.Model(client).
		Where("telegram_id = ? || user_id = ?", identifier, identifier).
		First(client).Error

	if err != nil {
		return nil, err
	}
	return client, nil
}

// Update is a method that updates a certain client entries in the database
func (repo *ClientRepository) Update(client *bot.Client) error {

	prevClient := new(bot.Client)
	err := repo.conn.Model(prevClient).Where("user_id = ?", client.UserID).First(prevClient).Error

	if err != nil {
		return err
	}

	err = repo.conn.Save(client).Error
	if err != nil {
		return err
	}
	return nil
}

// Delete is a method that deletes a certain client from the database using telegram_id or user_id.
func (repo *ClientRepository) Delete(identifier string) (*bot.Client, error) {
	client := new(bot.Client)
	err := repo.conn.Model(client).Where("telegram_id = ? || user_id = ?",
		identifier, identifier).First(client).Error

	if err != nil {
		return nil, err
	}

	repo.conn.Delete(client)
	return client, nil
}

package repository

import (
	"strings"

	"github.com/Benyam-S/asseri/client/bot"
	"github.com/Benyam-S/asseri/client/bot/tempuser"
	"github.com/jinzhu/gorm"
)

// TempUserRepository is a type that defines a temporary user repository type
type TempUserRepository struct {
	conn *gorm.DB
}

// NewTempUserRepository is a function that creates a new temporary user repository type
func NewTempUserRepository(connection *gorm.DB) tempuser.ITempUserRepository {
	return &TempUserRepository{conn: connection}
}

// Create is a method that adds a new temporary user to the database
func (repo *TempUserRepository) Create(newTempUser *bot.TempUser) error {

	err := repo.conn.Create(newTempUser).Error
	if err != nil {
		return err
	}
	return nil
}

// Find is a method that finds a certain temporary user from the database using an identifier,
// also Find() uses telegram_id and phone_number as a key for selection
func (repo *TempUserRepository) Find(identifier string) (*bot.TempUser, error) {

	modifiedIdentifier := identifier
	splitIdentifier := strings.Split(identifier, "")
	if splitIdentifier[0] == "0" {
		modifiedIdentifier = "+251" + strings.Join(splitIdentifier[1:], "")
	}

	tempUser := new(bot.TempUser)
	err := repo.conn.Model(tempUser).
		Where("telegram_id = ? || phone_number = ?", identifier, modifiedIdentifier).
		First(tempUser).Error

	if err != nil {
		return nil, err
	}
	return tempUser, nil
}

// Update is a method that updates a certain temporary user entries in the database
func (repo *TempUserRepository) Update(tempUser *bot.TempUser) error {

	prevTempUser := new(bot.TempUser)
	err := repo.conn.Model(prevTempUser).Where("telegram_id = ?", tempUser.TelegramID).
		First(prevTempUser).Error

	if err != nil {
		return err
	}

	/* --------------------------- can change layer if needed --------------------------- */
	tempUser.CreatedAt = prevTempUser.CreatedAt
	/* -------------------------------------- end --------------------------------------- */

	err = repo.conn.Save(tempUser).Error
	if err != nil {
		return err
	}
	return nil
}

// UpdateValue is a method that updates a certain temporary user single column value in the database
func (repo *TempUserRepository) UpdateValue(tempUser *bot.TempUser, columnName string, columnValue interface{}) error {

	prevTempUser := new(bot.TempUser)
	err := repo.conn.Model(prevTempUser).Where("telegram_id = ?", tempUser.TelegramID).First(prevTempUser).Error

	if err != nil {
		return err
	}

	err = repo.conn.Model(bot.TempUser{}).Where("telegram_id = ?", tempUser.TelegramID).
		Update(map[string]interface{}{columnName: columnValue}).Error
	if err != nil {
		return err
	}
	return nil
}

// Delete is a method that deletes a certain temporary user from the database using an identifier.
// In Delete() telegram_id is only used as an key
func (repo *TempUserRepository) Delete(identifier string) (*bot.TempUser, error) {
	tempUser := new(bot.TempUser)
	err := repo.conn.Model(tempUser).Where("telegram_id = ?", identifier).First(tempUser).Error

	if err != nil {
		return nil, err
	}

	repo.conn.Delete(tempUser)
	return tempUser, nil
}

package tempuser

import "github.com/Benyam-S/asseri/client/bot"

// ITempUserRepository is an interface that defines all the repository methods of a bot.TempUser struct
type ITempUserRepository interface {
	Create(newTempUser *bot.TempUser) error
	Find(identifier string) (*bot.TempUser, error)
	Update(tempUser *bot.TempUser) error
	UpdateValue(tempUser *bot.TempUser, columnName string, columnValue interface{}) error
	Delete(identifier string) (*bot.TempUser, error)
}

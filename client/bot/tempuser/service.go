package tempuser

import (
	"github.com/Benyam-S/asseri/client/bot"
	"github.com/Benyam-S/asseri/entity"
)

// IService is an interface that defines all the service methods of a bot.TempUser struct
type IService interface {
	AddTempUser(newTempUser *bot.TempUser) error
	ValidateTempUserProfile(tempUser *bot.TempUser) entity.ErrMap
	FindTempUser(identifier string) (*bot.TempUser, error)
	// SearchTempUsers(key, pagination string, extra ...string) []*bot.TempUser
	// AllTempUsers(pagination string) []*bot.TempUser
	UpdateTempUser(tempUser *bot.TempUser) error
	UpdateTempUserSingleValue(telegramID, columnName string, columnValue interface{}) error
	DeleteTempUser(telegramID string) (*bot.TempUser, error)
}

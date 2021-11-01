package service

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/Benyam-S/asseri/client/bot"
	"github.com/Benyam-S/asseri/client/bot/tempuser"
	"github.com/Benyam-S/asseri/common"
	"github.com/Benyam-S/asseri/entity"
	"github.com/Benyam-S/asseri/user"
	"github.com/nyaruka/phonenumbers"
)

// Service is a type that defines a temporary user service
type Service struct {
	tempUserRepo tempuser.ITempUserRepository
	userRepo     user.IUserRepository
	commonRepo   common.ICommonRepository
}

// NewTempUserService is a function that returns a new temporary user service
func NewTempUserService(tempUserRepository tempuser.ITempUserRepository, userRepository user.IUserRepository,
	commonRepository common.ICommonRepository) tempuser.IService {
	return &Service{tempUserRepo: tempUserRepository, userRepo: userRepository,
		commonRepo: commonRepository}
}

// AddTempUser is a method that adds a new temporary user to the system
func (service *Service) AddTempUser(newTempUser *bot.TempUser) error {
	err := service.tempUserRepo.Create(newTempUser)
	if err != nil {
		return errors.New("unable to add new temporary user")
	}

	return nil
}

// ValidateTempUserProfile is a method that validates a temporary user profile.
// It checks if the temporary user has a valid entries or not and return map of errors if any.
// Also it will add country code to the phone number value if not included: default country code +251
func (service *Service) ValidateTempUserProfile(tempUser *bot.TempUser) entity.ErrMap {

	errMap := make(map[string]error)
	validUserName, _ := regexp.MatchString(`^\w[\w\s]*$`, tempUser.UserName)

	// Removing all whitespaces
	phoneNumber := strings.Join(strings.Fields(tempUser.PhoneNumber), "")

	// Checking for local phone number
	isLocalPhoneNumber, _ := regexp.MatchString(`^0\d{9}$`, phoneNumber)

	if isLocalPhoneNumber {
		phoneNumberSlice := strings.Split(phoneNumber, "")
		if phoneNumberSlice[0] == "0" {
			phoneNumberSlice = phoneNumberSlice[1:]
			internationalPhoneNumber := "+251" + strings.Join(phoneNumberSlice, "")
			phoneNumber = internationalPhoneNumber
		}
	} else {
		// Making the phone number international if it is not local
		if len(phoneNumber) != 0 && string(phoneNumber[0]) != "+" {
			phoneNumber = "+" + phoneNumber
		}
	}

	parsedPhoneNumber, _ := phonenumbers.Parse(phoneNumber, "")
	validPhoneNumber := phonenumbers.IsValidNumber(parsedPhoneNumber)

	if !validUserName {
		errMap["user_name"] = errors.New(`user name should have at least one character and ` +
			`contain only alpha numeric value`)
	} else if len(tempUser.UserName) > 255 {
		errMap["user_name"] = errors.New(`user name should not be longer than 255 characters`)
	}

	if !validPhoneNumber {
		errMap["phone_number"] = errors.New("invalid phonenumber used")
	} else {
		// If a valid phone number is provided, adjust the phone number to fit the database
		// Stored in +251900010197 format
		phoneNumber = fmt.Sprintf("+%d%d", parsedPhoneNumber.GetCountryCode(),
			parsedPhoneNumber.GetNationalNumber())

		tempUser.PhoneNumber = phoneNumber
	}

	if tempUser.Category != entity.UserCategoryasseri &&
		tempUser.Category != entity.UserCategoryAgent &&
		tempUser.Category != entity.UserCategoryJobSeeker {
		errMap["category"] = errors.New("invalid category selected")
	}

	// Checking if the temporary user exists
	prevProfile, err := service.tempUserRepo.Find(tempUser.TelegramID)

	// Meaning a new temporary user is being added
	if err != nil {
		if validPhoneNumber && (!service.commonRepo.IsUnique("phone_number", tempUser.PhoneNumber, "users") ||
			!service.commonRepo.IsUnique("phone_number", tempUser.PhoneNumber, "temp_users")) {
			errMap["phone_number"] = errors.New("phone number already exists")
		}
	} else {

		if validPhoneNumber && prevProfile.PhoneNumber != tempUser.PhoneNumber {
			if !service.commonRepo.IsUnique("phone_number", tempUser.PhoneNumber, "users") ||
				!service.commonRepo.IsUnique("phone_number", tempUser.PhoneNumber, "temp_users") {
				errMap["phone_number"] = errors.New("phone number already exists")
			}
		}
	}

	if len(errMap) > 0 {
		return errMap
	}

	return nil
}

// FindTempUser is a method that find and return a temporary user that matchs the identifier value
func (service *Service) FindTempUser(identifier string) (*bot.TempUser, error) {

	empty, _ := regexp.MatchString(`^\s*$`, identifier)
	if empty {
		return nil, errors.New("no temporary user found")
	}

	tempUser, err := service.tempUserRepo.Find(identifier)
	if err != nil {
		return nil, errors.New("no temporary user found")
	}
	return tempUser, nil
}

// UpdateTempUser is a method that updates a temporary user in the system
func (service *Service) UpdateTempUser(tempUser *bot.TempUser) error {
	err := service.tempUserRepo.Update(tempUser)
	if err != nil {
		return errors.New("unable to update temporary user")
	}

	return nil
}

// UpdateTempUserSingleValue is a method that updates a single column entry of a temporary user
func (service *Service) UpdateTempUserSingleValue(telegramID, columnName string, columnValue interface{}) error {
	tempUser := bot.TempUser{TelegramID: telegramID}
	err := service.tempUserRepo.UpdateValue(&tempUser, columnName, columnValue)
	if err != nil {
		return errors.New("unable to update temporary user")
	}

	return nil
}

// DeleteTempUser is a method that deletes a temporary user from the system
func (service *Service) DeleteTempUser(telegramID string) (*bot.TempUser, error) {
	tempUser, err := service.tempUserRepo.Delete(telegramID)
	if err != nil {
		return nil, errors.New("unable to delete temporary user")
	}

	return tempUser, nil
}

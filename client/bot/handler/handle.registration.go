package handler

import (
	"strconv"
	"strings"

	"github.com/Benyam-S/asseri/client/bot"
	"github.com/Benyam-S/asseri/entity"
	"github.com/Benyam-S/asseri/tools"
)

// HandleRegistration is a method that handles the whole registration process
func (handler *TelegramBotHandler) HandleRegistration(update *bot.Update, tempUser *bot.TempUser) {
	switch tempUser.Status {
	case bot.RegistrationStatusInit:
		handler.HandleRegistrationName(update, tempUser)
	case bot.RegistrationStatusUserName:
		handler.HandleRegistrationPhone(update, tempUser)
	case bot.RegistrationStatusPhoneNumber:
		handler.HandleRegistrationCategory(update, tempUser)
	}
}

// HandleRegistrationInit is a method that initiate bot client registration
func (handler *TelegramBotHandler) HandleRegistrationInit(update *bot.Update) {
	tempUser := new(bot.TempUser)
	tempUser.TelegramID = strconv.FormatInt(update.Message.User.ID, 10)
	err := handler.tuService.AddTempUser(tempUser)
	if err != nil {
		keyboard := bot.CreateReplyKeyboard(true, true, []string{"🏁 Start"})
		bot.SendReplyToTelegramChat(update.Message.Chat.ID,
			"❌ Error unable to initiate registration!", keyboard)
		return
	}

	bot.SendReplyToTelegramChat(update.Message.Chat.ID, "Welcome 👋 to አሰሪ, please register first!")
	bot.SendReplyToTelegramChat(update.Message.Chat.ID, "Enter company/individual name")
}

// HandleRegistrationName is a method that handles user name registration
func (handler *TelegramBotHandler) HandleRegistrationName(update *bot.Update, tempUser *bot.TempUser) {
	tempUser.UserName = strings.TrimSpace(update.Message.Text)
	tempUser.Status = bot.RegistrationStatusUserName

	errMap := handler.tuService.ValidateTempUserProfile(tempUser)
	if errMap["user_name"] != nil {
		bot.SendReplyToTelegramChat(update.Message.Chat.ID, tools.ToSentenceCase(errMap["user_name"].Error()))
		bot.SendReplyToTelegramChat(update.Message.Chat.ID, "Enter company/individual name")
		return
	}

	err := handler.tuService.UpdateTempUser(tempUser)
	if err != nil {
		bot.SendReplyToTelegramChat(update.Message.Chat.ID, "❌ Error unable to register username!")
		bot.SendReplyToTelegramChat(update.Message.Chat.ID, "Re-enter your name")
		return
	}

	keyboard := bot.CreateReplyKeyboardWExtra(true, false, []bot.ReplyKeyboardButton{{Text: "Add 📱", RequestContact: true}})
	bot.SendReplyToTelegramChat(update.Message.Chat.ID, "Add your phonenumber, use 'Add 📱' button to add your phone number", keyboard)
}

// HandleRegistrationPhone is a method that handles phonenumber registration
func (handler *TelegramBotHandler) HandleRegistrationPhone(update *bot.Update, tempUser *bot.TempUser) {

	phoneNumber := update.Message.Contact.PhoneNumber
	tempUser.PhoneNumber = strings.TrimSpace(phoneNumber)
	tempUser.Status = bot.RegistrationStatusPhoneNumber

	errMap := handler.tuService.ValidateTempUserProfile(tempUser)
	if errMap["phone_number"] != nil {
		bot.SendReplyToTelegramChat(update.Message.Chat.ID, tools.ToSentenceCase(errMap["phone_number"].Error()))

		keyboard := bot.CreateReplyKeyboardWExtra(true, false, []bot.ReplyKeyboardButton{{Text: "Add 📱", RequestContact: true}})
		bot.SendReplyToTelegramChat(update.Message.Chat.ID, "Add your phonenumber, use 'Add 📱' button to add your phone number", keyboard)
		return
	}

	err := handler.tuService.UpdateTempUser(tempUser)
	if err != nil {
		bot.SendReplyToTelegramChat(update.Message.Chat.ID, "❌ Error unable to register phonenumber!")

		keyboard := bot.CreateReplyKeyboardWExtra(true, false, []bot.ReplyKeyboardButton{{Text: "Add 📱", RequestContact: true}})
		bot.SendReplyToTelegramChat(update.Message.Chat.ID, "Re-enter your phonenumber", keyboard)
		return
	}

	keyboard := bot.CreateReplyKeyboard(true, true, []string{"አሰሪ", "Job Seeker"}, []string{"Agent"})
	bot.SendReplyToTelegramChat(update.Message.Chat.ID, "You wish to be categorized as ?", keyboard)
}

// HandleRegistrationCategory is a method that handles user category registration
func (handler *TelegramBotHandler) HandleRegistrationCategory(update *bot.Update, tempUser *bot.TempUser) {
	category := strings.TrimSpace(update.Message.Text)

	if category == "አሰሪ" {
		category = entity.UserCategoryasseri

		// We changed it to lower for formality purpose
	} else if strings.ToLower(category) == "job seeker" {
		category = entity.UserCategoryJobSeeker
	} else if strings.ToLower(category) == "agent" {
		category = entity.UserCategoryAgent
	}

	tempUser.Category = category
	tempUser.Status = bot.RegistrationStatusCategory

	errMap := handler.tuService.ValidateTempUserProfile(tempUser)
	if errMap["category"] != nil {
		bot.SendReplyToTelegramChat(update.Message.Chat.ID, tools.ToSentenceCase(errMap["category"].Error()))
		bot.SendReplyToTelegramChat(update.Message.Chat.ID, "You wish to be categorized as ?")
		return
	}

	user := tempUser.ToUser()
	err := handler.urService.AddUser(user)
	if err != nil {
		keyboard := bot.CreateReplyKeyboard(true, true, []string{"አሰሪ", "Job Seeker"}, []string{"Agent"})
		bot.SendReplyToTelegramChat(update.Message.Chat.ID, "❌ Error unable to add new user!")
		bot.SendReplyToTelegramChat(update.Message.Chat.ID, "Re-select category", keyboard)
		return
	}

	client := &bot.Client{UserID: user.ID, TelegramID: tempUser.TelegramID}
	err = handler.clService.AddClient(client)
	if err != nil {
		// Removing the added user this their is no client link
		handler.urService.DeleteUser(user.ID)

		keyboard := bot.CreateReplyKeyboard(true, true, []string{"አሰሪ", "Job Seeker"}, []string{"Agent"})
		bot.SendReplyToTelegramChat(update.Message.Chat.ID, "❌ Error unable to add new user!")
		bot.SendReplyToTelegramChat(update.Message.Chat.ID, "Re-select category", keyboard)
		return
	}

	handler.tuService.DeleteTempUser(tempUser.TelegramID)

	menu := bot.MainMenuW
	if user.Category == entity.UserCategoryJobSeeker {
		menu = bot.MainMenuWO
	}

	bot.SendReplyToTelegramChat(update.Message.Chat.ID,
		"Congratulations 🎉 you have been successfully registered!", menu)
}

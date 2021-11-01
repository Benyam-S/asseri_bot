package handler

import (
	"fmt"
	"strings"

	"github.com/Benyam-S/asseri/client/bot"
	"github.com/Benyam-S/asseri/entity"
	"github.com/Benyam-S/asseri/tools"
)

// HandleViewProfile is a method that handles the profile viewing process
func (handler *TelegramBotHandler) HandleViewProfile(update *bot.Update, user *entity.User) {

	category := user.Category
	if user.Category == entity.UserCategoryasseri {
		category = "አሰሪ"
	}

	userProfile := fmt.Sprintf(
		"<b>Name</b>:   %s\n"+
			"<b>Category</b>:   %s\n"+
			"<b>Phonenumber</b>:   %s\n\n",
		strings.Title(strings.ToLower(user.UserName)), category, user.PhoneNumber)

	profileMenu := bot.CreateReplyKeyboard(true, false, []string{"🔧 Update Profile", "🔙 Main Menu"})
	bot.SendReplyToTelegramChat(update.Message.Chat.ID, userProfile, profileMenu)
}

// HandleInitUpdateProfile is a method that initiates the profile updating process
func (handler *TelegramBotHandler) HandleInitUpdateProfile(update *bot.Update, user *entity.User) {

	backMenu := bot.CreateReplyKeyboard(true, false, []string{"↖️ Skip", "🔙 Main Menu"})
	bot.SendReplyToTelegramChat(update.Message.Chat.ID, "Enter new name", backMenu)
}

// HandleUpdateName is a method that handles user name updating process
func (handler *TelegramBotHandler) HandleUpdateName(update *bot.Update, user *entity.User) bool {
	user.UserName = strings.TrimSpace(update.Message.Text)

	errMap := handler.urService.ValidateUserProfile(user)
	if errMap["user_name"] != nil {
		bot.SendReplyToTelegramChat(update.Message.Chat.ID, tools.ToSentenceCase(errMap["user_name"].Error()))
		bot.SendReplyToTelegramChat(update.Message.Chat.ID, "Enter new name")
		return false
	}

	err := handler.urService.UpdateUser(user)
	if err != nil {
		bot.SendReplyToTelegramChat(update.Message.Chat.ID, "❌ Error unable to update name!")
		bot.SendReplyToTelegramChat(update.Message.Chat.ID, "Re-enter new name")
		return false
	}

	backMenu := bot.CreateReplyKeyboard(true, false, []string{"↖️ Skip", "🔙 Main Menu"})
	bot.SendReplyToTelegramChat(update.Message.Chat.ID, "Enter new phonenumber", backMenu)
	return true
}

// HandleUpdatePhone is a method that handles phonenumber updating process
func (handler *TelegramBotHandler) HandleUpdatePhone(update *bot.Update, user *entity.User) bool {
	user.PhoneNumber = strings.TrimSpace(update.Message.Text)

	errMap := handler.urService.ValidateUserProfile(user)
	if errMap["phone_number"] != nil {
		bot.SendReplyToTelegramChat(update.Message.Chat.ID, tools.ToSentenceCase(errMap["phone_number"].Error()))
		bot.SendReplyToTelegramChat(update.Message.Chat.ID, "Enter new phonenumber")
		return false
	}

	err := handler.urService.UpdateUser(user)
	if err != nil {
		bot.SendReplyToTelegramChat(update.Message.Chat.ID, "❌ Error unable to update phonenumber!")
		bot.SendReplyToTelegramChat(update.Message.Chat.ID, "Re-enter new phonenumber")
		return false
	}

	keyboard := bot.CreateReplyKeyboard(true, false, []string{"አሰሪ", "Job Seeker", "Agent"},
		[]string{"↖️ Skip", "🔙 Main Menu"})
	bot.SendReplyToTelegramChat(update.Message.Chat.ID, "Wish to change category?", keyboard)
	return true
}

// HandleUpdateCategory is a method that handles user category updating process
func (handler *TelegramBotHandler) HandleUpdateCategory(update *bot.Update, user *entity.User) bool {
	category := strings.TrimSpace(update.Message.Text)

	if category == "አሰሪ" {
		category = entity.UserCategoryasseri

		// We changed it to lower for formality purpose
	} else if strings.ToLower(category) == "job seeker" {
		category = entity.UserCategoryJobSeeker
	} else if strings.ToLower(category) == "agent" {
		category = entity.UserCategoryAgent
	}

	user.Category = category

	errMap := handler.urService.ValidateUserProfile(user)
	if errMap["category"] != nil {
		bot.SendReplyToTelegramChat(update.Message.Chat.ID, tools.ToSentenceCase(errMap["category"].Error()))
		bot.SendReplyToTelegramChat(update.Message.Chat.ID, "Wish to change category?")
		return false
	}

	err := handler.urService.UpdateUser(user)
	if err != nil {
		bot.SendReplyToTelegramChat(update.Message.Chat.ID, "❌ Error unable to update category!")
		bot.SendReplyToTelegramChat(update.Message.Chat.ID, "Re-select category")
		return false
	}

	bot.SendReplyToTelegramChat(update.Message.Chat.ID,
		"Congratulations 🎉 you have successfully update your profile!")
	return true
}

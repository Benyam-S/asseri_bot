package handler

import (
	"strings"

	"github.com/Benyam-S/asseri/client/bot"
	"github.com/Benyam-S/asseri/entity"
	"github.com/Benyam-S/asseri/tools"
)

// HandlePromptFeedback is a method that initiate feedback receiving process
func (handler *TelegramBotHandler) HandlePromptFeedback(update *bot.Update) {
	backMenu := bot.CreateReplyKeyboard(true, false, []string{"🔙 Main Menu"})
	bot.SendReplyToTelegramChat(update.Message.Chat.ID, "Please write your feedback in concise and short way",
		backMenu)
}

// HandleReceiveFeedback is a method that handles feedback receiving process
func (handler *TelegramBotHandler) HandleReceiveFeedback(update *bot.Update, user *entity.User) bool {
	feedback := new(entity.Feedback)
	feedback.Comment = strings.TrimSpace(update.Message.Text)
	feedback.UserID = user.ID

	errMap := handler.fdService.ValidateFeedback(feedback)

	if errMap["comment"] != nil {
		bot.SendReplyToTelegramChat(update.Message.Chat.ID, tools.ToSentenceCase(errMap["comment"].Error()))
		bot.SendReplyToTelegramChat(update.Message.Chat.ID, "Please write your feedback in concise and short way")
		return false
	}

	err := handler.fdService.AddFeedback(feedback)
	if err != nil {
		bot.SendReplyToTelegramChat(update.Message.Chat.ID, "❌ Error unable to add your feedback!")
		bot.SendReplyToTelegramChat(update.Message.Chat.ID, "Re-enter your feedback")
		return false
	}

	bot.SendReplyToTelegramChat(update.Message.Chat.ID, "Thank you 😁 for your feedback!"+
		" We will do our best to satisfy your requests.")
	return true
}

package handler

import "github.com/Benyam-S/asseri/client/bot"

// HandleSettings is a method that handles settings menu viewing
func (handler *TelegramBotHandler) HandleSettings(update *bot.Update) {

	settingsMenu := bot.CreateReplyKeyboard(true, false, []string{"👥 Profile", "🗣️ Feedback"},
		[]string{"🔙 Main Menu"})
	bot.SendReplyToTelegramChat(update.Message.Chat.ID, "Choose preference", settingsMenu)
}

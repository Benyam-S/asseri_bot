package handler

import (
	"net/http"
	"strconv"

	"github.com/Benyam-S/asseri/client/bot"
	"github.com/Benyam-S/asseri/entity"
)

// HandleWebHook sends a message back to the chat with a punchline starting by the message provided by the user.
func (handler *TelegramBotHandler) HandleWebHook(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	update, ok := ctx.Value(entity.Key("update_info")).(*bot.Update)

	if !ok {
		handler.logger.LogFileError("Unable to get parsed update.", entity.BotLogFile)
		return
	}

	telegramID := strconv.FormatInt(update.Message.User.ID, 10)

	// This is used for call back query response so as to identify the user
	if telegramID == "0" {
		telegramID = strconv.FormatInt(update.CallbackQuery.User.ID, 10)
	}

	client, err := handler.clService.FindClient(telegramID)

	// The telegram user hasn't been registered yet
	if err != nil {

		// If not registered, check temporary database
		tempUser, err := handler.tuService.FindTempUser(telegramID)

		// If not found in temporary database then initiate registration process
		if err != nil {
			handler.HandleRegistrationInit(update)
			return
		}

		// Complete registration process and exit
		handler.HandleRegistration(update, tempUser)
		return
	}

	// Error checking is not needed here since if the user is deleted the client is automatically deleted too
	user, _ := handler.urService.FindUser(client.UserID)
	handler.HandleClient(update, user)
}

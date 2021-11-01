package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Benyam-S/asseri/client/bot"
	"github.com/Benyam-S/asseri/entity"
)

// ParseRequest is a middleware that parse an incoming request update value
func (handler *TelegramBotHandler) ParseRequest(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		update := new(bot.Update)
		err := json.NewDecoder(r.Body).Decode(update)
		if err != nil {
			handler.logger.LogFileError(fmt.Sprintf("Unable to parse update, %s", err.Error()), entity.BotLogFile)
			return
		}

		// Adding the update information to the context
		ctx := r.Context()
		ctx = context.WithValue(ctx, entity.Key("update_info"), update)
		r = r.WithContext(ctx)

		next(w, r)
	}
}

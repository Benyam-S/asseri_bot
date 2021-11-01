package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Benyam-S/asseri/client/bot"
	"github.com/Benyam-S/asseri/entity"
	"github.com/Benyam-S/asseri/tools"
	"github.com/gorilla/mux"
)

// HandleApprovalResult is a handler func that handles a request for sending approval result to telegram user
func (handler *TelegramBotHandler) HandleApprovalResult(w http.ResponseWriter, r *http.Request) {

	jobID := mux.Vars(r)["id"]

	job, err := handler.jbService.FindJob(jobID)
	if err != nil {
		output, _ := json.MarshalIndent(map[string]string{"error": err.Error()}, "", "\t")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	user, err := handler.urService.FindUser(job.Employer)
	if err != nil {
		// Since a job doesn't necessarily need to be owned by a user
		w.WriteHeader(http.StatusOK)
		return
	}

	client, err := handler.clService.FindClient(user.ID)
	if err == nil {
		handler.ProcessJobResult(job, user, client, w, r)
	}
}

// HandlePushNotificationToChannel is a handler func that handles a request for pushing notification to the channel
func (handler *TelegramBotHandler) HandlePushNotificationToChannel(w http.ResponseWriter, r *http.Request) {

	jobID := mux.Vars(r)["id"]

	job, err := handler.jbService.FindJob(jobID)
	if err != nil {
		output, _ := json.MarshalIndent(map[string]string{"error": err.Error()}, "", "\t")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	handler.PushNotificationToChannel(job, w, r)
}

// HandlePushNotificationToSubscribers is a handler func that handles a request for pushing notification for subscribers
func (handler *TelegramBotHandler) HandlePushNotificationToSubscribers(w http.ResponseWriter, r *http.Request) {

	jobID := mux.Vars(r)["id"]

	job, err := handler.jbService.FindJob(jobID)
	if err != nil {
		output, _ := json.MarshalIndent(map[string]string{"error": err.Error()}, "", "\t")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	handler.PushNotificationToSubscribers(job, w, r)

}

// ProcessJobResult is a method that process a job and send the need reply to the telegram bot
func (handler *TelegramBotHandler) ProcessJobResult(job *entity.Job, user *entity.User,
	client *bot.Client, w http.ResponseWriter, r *http.Request) {

	var statusString string
	var postToChat string

	if job.Status == entity.JobStatusOpened {
		statusString = "------------- <b>Approved</b> -------------\n\n"
	} else if job.Status == entity.JobStatusDecelined {
		statusString = "------------- <b>Declined</b> -------------\n\n"
	} else if job.Status == entity.JobStatusClosed {
		statusString = "------------- <b>Closed</b> -------------\n\n"
	} else {
		output, _ := json.MarshalIndent(map[string]string{"error": "unable to perform operation"}, "", "\t")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	postToChat = fmt.Sprintf(
		"%s"+"<b>Job Title</b>:  %s\n\n"+
			"<b>Job Type</b>:  %s\n"+
			"<b>Gender</b>:  %s\n"+
			"<b>Education level</b>:  %s\n"+
			"<b>Experience</b>:  %s\n"+
			"<b>Contact Type</b>:  %s\n\n"+
			"<b>Description</b>:  %s\n\n"+"#%s\n\n"+"%s",
		statusString, job.Title, job.Type, bot.GetGender(job.Gender),
		job.EducationLevel, job.Experience, job.ContactType, job.Description,
		tools.ChangeSpaceToUnderscore(job.Sector), statusString)

	chatID, _ := strconv.ParseInt(client.TelegramID, 10, 64)
	value, err := bot.SendReplyToTelegramChat(chatID, postToChat)
	if err != nil {
		handler.logger.LogFileError(string(err.Error()), entity.BotLogFile)
		output, _ := json.MarshalIndent(map[string]string{"error": err.Error()}, "", "\t")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	botResponse := new(BotResponse)
	json.Unmarshal([]byte(value), botResponse)

	if !botResponse.Ok && botResponse.ErrorCode == 429 {
		output, _ := json.MarshalIndent(map[string]string{"error": "retry"}, "", "\t")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}
}

// PushNotificationToChannel is a method that pushes job alert notifications to channel
func (handler *TelegramBotHandler) PushNotificationToChannel(job *entity.Job,
	w http.ResponseWriter, r *http.Request) {

	var contact string
	var employer string
	var inlineKeyboard string

	if job.Status != entity.JobStatusOpened &&
		job.Status != entity.JobStatusClosed {
		return
	}

	if job.PostType == entity.PostCategoryUser {

		user, err := handler.urService.FindUser(job.Employer)
		if err != nil {
			output, _ := json.MarshalIndent(map[string]string{"error": err.Error()}, "", "\t")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(output)
			return
		}

		employer = user.UserName
		client, _ := handler.clService.FindClient(user.ID)

		// Means via telegram account
		if job.ContactType == handler.cmService.GetValidContactTypes()[0] {

			var chatID string
			if client != nil {
				chatID = client.TelegramID
			}
			var getChatURL string = os.Getenv("api_access_point") +
				os.Getenv("bot_api_token") + "/getChat?chat_id=" + chatID

			result, _ := http.Get(getChatURL)

			type UserProfile struct {
				UserName string `json:"username"`
			}

			chatResponse := &struct {
				Result UserProfile `json:"result"`
			}{}

			json.NewDecoder(result.Body).Decode(chatResponse)
			if chatResponse.Result.UserName != "" {
				contact = "<b>Contact</b>: @" + chatResponse.Result.UserName + "\n\n"
			} else {
				contact = "<b>Contact</b>: " + strings.ReplaceAll(user.PhoneNumber, "+251", "0") + "\n\n"
			}

		} else if job.ContactType == handler.cmService.GetValidContactTypes()[1] {
			inlineKeyboard = bot.CreateInlineKeyboard([]bot.InlineKeyboardButton{
				{Text: "🔗 Apply", URL: os.Getenv("bot_url") + "?start=" + "apply_" + job.ID},
			})
		}
	} else if job.PostType == entity.PostCategoryInternal {
		employer = job.Employer
		contact = "<b>Contact</b>: " + job.ContactInfo + "\n\n"

	} else if job.PostType == entity.PostCategoryExternal {
		employer = job.Employer
		emptyLink, _ := regexp.MatchString(`^\s*$`, job.Link)
		if !emptyLink {
			inlineKeyboard = bot.CreateInlineKeyboard([]bot.InlineKeyboardButton{
				{Text: "Tap to view details / ዝርዝሩን ለማየት ይሄንን ይጫኑ", URL: job.Link},
			})
		}

	}

	pack := &bot.StructuredPackage{Employer: employer, Contact: contact}
	postToChannel := bot.BuildNotification(job, pack)

	if job.Status == entity.JobStatusClosed {
		postToChannel = "------------- <b>Closed</b> -------------\n\n" +
			postToChannel +
			"------------- <b>Closed</b> -------------\n\n"
		inlineKeyboard = ""
	}

	// Posting to telegram channel if opened
	value, err := bot.PostToTelegramChannel(postToChannel, inlineKeyboard)
	if err != nil {
		handler.logger.LogFileError(err.Error(), entity.BotLogFile)
		output, _ := json.MarshalIndent(map[string]string{"error": err.Error()}, "", "\t")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	botResponse := new(BotResponse)
	json.Unmarshal([]byte(value), botResponse)

	if !botResponse.Ok && botResponse.ErrorCode == 429 {
		output, _ := json.MarshalIndent(map[string]string{"error": "retry"}, "", "\t")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}
}

// PushNotificationToSubscribers is a method that pushes job alert notifications to subscribers
func (handler *TelegramBotHandler) PushNotificationToSubscribers(job *entity.Job,
	w http.ResponseWriter, r *http.Request) {

	var contact string
	var employer string
	var inlineKeyboard string

	if job.Status != entity.JobStatusOpened {
		return
	}

	if job.PostType == entity.PostCategoryUser {

		user, err := handler.urService.FindUser(job.Employer)
		if err != nil {
			output, _ := json.MarshalIndent(map[string]string{"error": err.Error()}, "", "\t")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(output)
			return
		}

		employer = user.UserName
		client, _ := handler.clService.FindClient(user.ID)

		// Means via telegram account
		if job.ContactType == handler.cmService.GetValidContactTypes()[0] {

			var chatID string
			if client != nil {
				chatID = client.TelegramID
			}
			var getChatURL string = os.Getenv("api_access_point") +
				os.Getenv("bot_api_token") + "/getChat?chat_id=" + chatID

			result, _ := http.Get(getChatURL)

			type UserProfile struct {
				UserName string `json:"username"`
			}

			chatResponse := &struct {
				Result UserProfile `json:"result"`
			}{}

			json.NewDecoder(result.Body).Decode(chatResponse)
			if chatResponse.Result.UserName != "" {
				contact = "@" + chatResponse.Result.UserName + "\n\n"
			} else {
				contact = "<b>Contact</b>: " + strings.ReplaceAll(user.PhoneNumber, "+251", "0") + "\n\n"
			}

		} else if job.ContactType == handler.cmService.GetValidContactTypes()[1] {
			inlineKeyboard = bot.CreateInlineKeyboard([]bot.InlineKeyboardButton{
				{Text: "🔗 Apply", URL: os.Getenv("bot_url") + "?start=" + "apply_" + job.ID},
			})
		}
	} else if job.PostType == entity.PostCategoryInternal {
		employer = job.Employer
		contact = "<b>Contact</b>: " + job.ContactInfo + "\n\n"

	} else if job.PostType == entity.PostCategoryExternal {
		employer = job.Employer
		emptyLink, _ := regexp.MatchString(`^\s*$`, job.Link)
		if !emptyLink {
			inlineKeyboard = bot.CreateInlineKeyboard([]bot.InlineKeyboardButton{
				{Text: "Tap to view details / ዝርዝሩን ለማየት ይሄንን ይጫኑ", URL: job.Link},
			})
		}

	}

	pack := &bot.StructuredPackage{Employer: employer, Contact: contact}
	postToSubscribers := "------------- <b>Subscription</b> -------------\n\n" +
		bot.BuildNotification(job, pack)

	subscribers := handler.sbService.FindSubscriptionMatch(job.Sector, job.Type, job.EducationLevel, job.Experience)
	for _, subscriber := range subscribers {
		subscriberClient, err := handler.clService.FindClient(subscriber.UserID)
		if err != nil || job.Employer == subscriber.UserID {
			continue
		}

		chatID, _ := strconv.ParseInt(subscriberClient.TelegramID, 10, 64)

		newRequest := new(entity.ChannelRequest)
		newRequest.ChatID = chatID
		newRequest.Value = postToSubscribers
		newRequest.Extra = inlineKeyboard

		handler.pq.AddToQueue(newRequest)
	}

	handler.pushChan <- entity.StartPush
}

// HandlePushRequest is a method that handles push notification sending process to bot handler
func (handler *TelegramBotHandler) HandlePushRequest() {

	for range handler.pushChan {

		for {

			i := 0
			requestCount := 0

			if len(handler.pq.GetQueue()) == 0 {
				break
			}

			for {

				if i >= len(handler.pq.GetQueue()) {
					break
				}

				request := handler.pq.GetQueue()[i]

				if requestCount > 15 {
					time.Sleep(time.Second * 15)
					requestCount = 0
				}

				value, _ := bot.SendReplyToTelegramChat(request.ChatID, request.Value, request.Extra)

				botResponse := new(BotResponse)
				json.Unmarshal([]byte(value), botResponse)

				if !botResponse.Ok && botResponse.ErrorCode == 429 {
					i++
					requestCount++
					continue
				}

				handler.pq.RemoveFromQueueWithIndex(i)
				requestCount++
			}
		}

	}
}

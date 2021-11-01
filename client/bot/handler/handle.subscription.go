package handler

import (
	"fmt"
	"strconv"

	"github.com/Benyam-S/asseri/client/bot"
	"github.com/Benyam-S/asseri/entity"
	"github.com/Benyam-S/asseri/tools"
)

// HandleJobSubscription is a method that handles job subscription menu and registered subscription viewing
func (handler *TelegramBotHandler) HandleJobSubscription(update *bot.Update, user *entity.User) {

	// Remove any invalid subscription
	subscriptions := handler.sbService.FindMultipleSubscriptions(user.ID)
	for _, subscription := range subscriptions {
		errMap := handler.sbService.ValidateSubscription(subscription)
		if len(errMap) > 0 {
			handler.sbService.DeleteSubscription(subscription.ID)
		}
	}

	subscriptionMenu := bot.CreateReplyKeyboard(true, false,
		[]string{"➕ Add Subscription", "📝 Edit Subscriptions"}, []string{"🔙 Main Menu"})
	bot.SendReplyToTelegramChat(update.Message.Chat.ID, "Choose option", subscriptionMenu)
}

// HandleEditJobSubscriptions is a method that handles registered subscription viewing and editing process
func (handler *TelegramBotHandler) HandleEditJobSubscriptions(update *bot.Update, user *entity.User) {
	subscriptions := handler.sbService.FindMultipleSubscriptions(user.ID)

	if len(subscriptions) == 0 {
		bot.SendReplyToTelegramChat(update.Message.Chat.ID, "You haven't subscribed for a job!"+
			" please subscribe using the add subscription button.")
		return
	}

	for _, subscription := range subscriptions {

		errMap := handler.sbService.ValidateSubscription(subscription)
		if len(errMap) == 0 {
			reply := fmt.Sprintf(
				"<b>Job Subscription</b>\n\n"+
					"<b>Job Type</b>:  %s\n"+
					"<b>Sector</b>:  %s\n"+
					"<b>Education Level</b>:  %s\n"+
					"<b>Experience</b>:  %s\n\n",
				subscription.Type, subscription.Sector, subscription.EducationLevel, subscription.Experience)

			removeButton := bot.CreateInlineKeyboard([]bot.InlineKeyboardButton{
				{Text: "🗑️ Remove", CallbackData: "subscription/remove/" + subscription.ID},
			})
			bot.SendReplyToTelegramChat(update.Message.Chat.ID, reply, removeButton)
		}
	}
}

// HandleInitAddSubscriptionSector is a method that initiates the job subscription sector adding process
func (handler *TelegramBotHandler) HandleInitAddSubscriptionSector(client *bot.Client) {

	chatID, _ := strconv.ParseInt(client.TelegramID, 10, 64)
	validJobSectors := handler.cmService.GetValidJobSectorsForSubscription()
	validJobSectorButtons := [][]bot.InlineKeyboardButton{}
	row := []bot.InlineKeyboardButton{}

	for index, validJobSector := range validJobSectors {

		if index%2 == 0 {
			row = []bot.InlineKeyboardButton{}
			row = append(row,
				bot.InlineKeyboardButton{Text: validJobSector.Name,
					CallbackData: "subscription/add/sector/" + validJobSector.ID})
		} else {
			row = append(row,
				bot.InlineKeyboardButton{Text: validJobSector.Name,
					CallbackData: "subscription/add/sector/" + validJobSector.ID})
			validJobSectorButtons = append(validJobSectorButtons, row)
		}

		// Adding the last row if the index ends with even
		if index%2 == 0 && index == len(validJobSectors)-1 {
			validJobSectorButtons = append(validJobSectorButtons, row)
		}

	}

	backMenu := bot.CreateReplyKeyboard(true, false, []string{"🔙 Main Menu"})
	bot.SendReplyToTelegramChat(chatID, "Select job sector", backMenu)

	validJobSectorMenu := bot.CreateInlineKeyboard(validJobSectorButtons...)
	bot.SendReplyToTelegramChat(chatID,
		"<b>The following are the valid job sectors avaliable</b>", validJobSectorMenu)
}

// HandleInitAddSubscriptionType is a method that shows the valid job types avaliable for subscription
func (handler *TelegramBotHandler) HandleInitAddSubscriptionType(subscriptionID string, client *bot.Client) {

	chatID, _ := strconv.ParseInt(client.TelegramID, 10, 64)
	ValidJobTypes := handler.cmService.GetValidJobTypesForSubscription()
	validJobTypesButtons := [][]bot.InlineKeyboardButton{}

	for _, validJobType := range ValidJobTypes {
		validJobTypesButtons = append(validJobTypesButtons, []bot.InlineKeyboardButton{
			{Text: validJobType.Name, CallbackData: "subscription/" + subscriptionID + "/add/type/" + validJobType.ID},
		})
	}

	backMenu := bot.CreateReplyKeyboard(true, false, []string{"🔙 Main Menu"})
	bot.SendReplyToTelegramChat(chatID, "Select job type", backMenu)

	validJobTypeMenu := bot.CreateInlineKeyboard(validJobTypesButtons...)
	bot.SendReplyToTelegramChat(chatID,
		"<b>The following are the valid job types avaliable</b>", validJobTypeMenu)
}

// HandleInitAddSubscriptionEducationLevel is a method that shows the valid education levels avaliable for subscription
func (handler *TelegramBotHandler) HandleInitAddSubscriptionEducationLevel(subscriptionID string, client *bot.Client) {

	chatID, _ := strconv.ParseInt(client.TelegramID, 10, 64)
	validEducationLevels := handler.cmService.GetValidEducationLevelsForSubscription()
	validEducationLevelsButtons := [][]bot.InlineKeyboardButton{}
	row := []bot.InlineKeyboardButton{}

	for index, validEducationLevel := range validEducationLevels {

		if index%2 == 0 {
			row = []bot.InlineKeyboardButton{}
			row = append(row,
				bot.InlineKeyboardButton{
					Text: validEducationLevel.Name, CallbackData: "subscription/" + subscriptionID + "/add/education_level/" + validEducationLevel.ID})
		} else {
			row = append(row,
				bot.InlineKeyboardButton{
					Text: validEducationLevel.Name, CallbackData: "subscription/" + subscriptionID + "/add/education_level/" + validEducationLevel.ID})
			validEducationLevelsButtons = append(validEducationLevelsButtons, row)
		}

		// Adding the last row if the index ends with even
		if index%2 == 0 && index == len(validEducationLevels)-1 {
			validEducationLevelsButtons = append(validEducationLevelsButtons, row)
		}

	}

	backMenu := bot.CreateReplyKeyboard(true, false, []string{"🔙 Main Menu"})
	bot.SendReplyToTelegramChat(chatID, "Select education level", backMenu)

	validEducationLevelMenu := bot.CreateInlineKeyboard(validEducationLevelsButtons...)
	bot.SendReplyToTelegramChat(chatID,
		"<b>The following are the valid education levels avaliable</b>", validEducationLevelMenu)
}

// HandleInitAddSubscriptionExperience is a method that shows the valid work experience avaliable for subscription
func (handler *TelegramBotHandler) HandleInitAddSubscriptionExperience(subscriptionID string, client *bot.Client) {

	chatID, _ := strconv.ParseInt(client.TelegramID, 10, 64)
	validExperiences := handler.cmService.GetValidWorkExperiencesForSubscription()
	validExperiencesButtons := [][]bot.InlineKeyboardButton{}
	row := []bot.InlineKeyboardButton{}

	for index, validExperience := range validExperiences {

		if index%2 == 0 {
			row = []bot.InlineKeyboardButton{}
			row = append(row,
				bot.InlineKeyboardButton{
					Text: validExperience, CallbackData: "subscription/" + subscriptionID + "/add/experience/" + validExperience})
		} else {
			row = append(row,
				bot.InlineKeyboardButton{
					Text: validExperience, CallbackData: "subscription/" + subscriptionID + "/add/experience/" + validExperience})
			validExperiencesButtons = append(validExperiencesButtons, row)
		}

		// Adding the last row if the index ends with even
		if index%2 == 0 && index == len(validExperiences)-1 {
			validExperiencesButtons = append(validExperiencesButtons, row)
		}

	}

	backMenu := bot.CreateReplyKeyboard(true, false, []string{"🔙 Main Menu"})
	bot.SendReplyToTelegramChat(chatID, "Select work experience", backMenu)

	validExperienceMenu := bot.CreateInlineKeyboard(validExperiencesButtons...)
	bot.SendReplyToTelegramChat(chatID,
		"<b>The following are the valid work experiences avaliable</b>", validExperienceMenu)
}

// AddSubscriptionSector is a method that handles job subscription sector adding process
func (handler *TelegramBotHandler) AddSubscriptionSector(jobSectorID string, user *entity.User, client *bot.Client) {

	chatID, _ := strconv.ParseInt(client.TelegramID, 10, 64)
	subscription := new(entity.Subscription)
	subscription.UserID = user.ID
	subscription.Sector = jobSectorID

	errMap := handler.sbService.ValidateSubscription(subscription)
	if errMap["sector"] != nil {
		bot.SendReplyToTelegramChat(chatID, "❌ "+tools.ToSentenceCase(errMap["sector"].Error()))
		handler.HandleInitAddSubscriptionSector(client)
		return
	}

	err := handler.sbService.AddSubscription(subscription)
	if err != nil {
		bot.SendReplyToTelegramChat(chatID, "❌ Error unable to add job subscription sector!")
		handler.HandleInitAddSubscriptionSector(client)
		return
	}

	handler.HandleInitAddSubscriptionType(subscription.ID, client)
}

// AddSubscriptionType is a method that handles job subscription type adding process
func (handler *TelegramBotHandler) AddSubscriptionType(subscriptionID, jobTypeID string, user *entity.User, client *bot.Client) int {

	chatID, _ := strconv.ParseInt(client.TelegramID, 10, 64)
	subscription, err := handler.sbService.FindSubscription(subscriptionID)

	// subscription.Type != "" is used so we can't edit existing subscription
	if err != nil || subscription.Type != "" {
		subscriptionMenu := bot.CreateReplyKeyboard(true, false,
			[]string{"➕ Add Subscription", "📝 Edit Subscriptions"}, []string{"🔙 Main Menu"})
		bot.SendReplyToTelegramChat(chatID, "Oops 😳 something terribly went wrong!")
		bot.SendReplyToTelegramChat(chatID, "Choose option", subscriptionMenu)
		return bot.SubscriptionNotFound
	}

	subscription.Type = jobTypeID

	errMap := handler.sbService.ValidateSubscription(subscription)
	if errMap["type"] != nil {
		bot.SendReplyToTelegramChat(chatID, "❌ "+tools.ToSentenceCase(errMap["type"].Error()))
		handler.HandleInitAddSubscriptionType(subscriptionID, client)
		return bot.SubscriptionError
	}

	if errMap["error"] != nil {
		bot.SendReplyToTelegramChat(chatID, "❌ "+tools.ToSentenceCase(errMap["error"].Error()))
		handler.HandleInitAddSubscriptionType(subscriptionID, client)
		return bot.SubscriptionError
	}

	err = handler.sbService.UpdateSubscription(subscription)
	if err != nil {
		bot.SendReplyToTelegramChat(chatID, "❌ Error unable to add job subscription type!")
		handler.HandleInitAddSubscriptionType(subscriptionID, client)
		return bot.SubscriptionError
	}

	handler.HandleInitAddSubscriptionEducationLevel(subscription.ID, client)
	return bot.SubscriptionModified
}

// AddSubscriptionEducationLevel is a method that handles job subscription education level adding process
func (handler *TelegramBotHandler) AddSubscriptionEducationLevel(subscriptionID, educationLevelID string, user *entity.User, client *bot.Client) int {

	chatID, _ := strconv.ParseInt(client.TelegramID, 10, 64)
	subscription, err := handler.sbService.FindSubscription(subscriptionID)

	// subscription.EducationLevel != "" is used so we can't edit existing subscription
	if err != nil || subscription.EducationLevel != "" {
		subscriptionMenu := bot.CreateReplyKeyboard(true, false,
			[]string{"➕ Add Subscription", "📝 Edit Subscriptions"}, []string{"🔙 Main Menu"})
		bot.SendReplyToTelegramChat(chatID, "Oops 😳 something terribly went wrong!")
		bot.SendReplyToTelegramChat(chatID, "Choose option", subscriptionMenu)
		return bot.SubscriptionNotFound
	}

	subscription.EducationLevel = educationLevelID

	errMap := handler.sbService.ValidateSubscription(subscription)
	if errMap["education_level"] != nil {
		bot.SendReplyToTelegramChat(chatID, "❌ "+tools.ToSentenceCase(errMap["education_level"].Error()))
		handler.HandleInitAddSubscriptionEducationLevel(subscriptionID, client)
		return bot.SubscriptionError
	}

	if errMap["error"] != nil {
		bot.SendReplyToTelegramChat(chatID, "❌ "+tools.ToSentenceCase(errMap["error"].Error()))
		handler.HandleInitAddSubscriptionEducationLevel(subscriptionID, client)
		return bot.SubscriptionError
	}

	err = handler.sbService.UpdateSubscription(subscription)
	if err != nil {
		bot.SendReplyToTelegramChat(chatID, "❌ Error unable to add education level for the job subscription!")
		handler.HandleInitAddSubscriptionEducationLevel(subscriptionID, client)
		return bot.SubscriptionError
	}

	handler.HandleInitAddSubscriptionExperience(subscription.ID, client)
	return bot.SubscriptionModified
}

// AddSubscriptionExperience is a method that handles job subscription work experience adding process
func (handler *TelegramBotHandler) AddSubscriptionExperience(subscriptionID, experience string, user *entity.User, client *bot.Client) int {

	chatID, _ := strconv.ParseInt(client.TelegramID, 10, 64)
	subscription, err := handler.sbService.FindSubscription(subscriptionID)

	// subscription.Experience != "" is used so we can't edit existing subscription
	if err != nil || subscription.Experience != "" {
		subscriptionMenu := bot.CreateReplyKeyboard(true, false,
			[]string{"➕ Add Subscription", "📝 Edit Subscriptions"}, []string{"🔙 Main Menu"})
		bot.SendReplyToTelegramChat(chatID, "Oops 😳 something terribly went wrong!")
		bot.SendReplyToTelegramChat(chatID, "Choose option", subscriptionMenu)
		return bot.SubscriptionNotFound
	}

	subscription.Experience = experience

	errMap := handler.sbService.ValidateSubscription(subscription)
	if errMap["experience"] != nil {
		bot.SendReplyToTelegramChat(chatID, "❌ "+tools.ToSentenceCase(errMap["experience"].Error()))
		handler.HandleInitAddSubscriptionExperience(subscriptionID, client)
		return bot.SubscriptionError
	}

	if errMap["error"] != nil {
		bot.SendReplyToTelegramChat(chatID, "❌ "+tools.ToSentenceCase(errMap["error"].Error()))
		handler.HandleInitAddSubscriptionExperience(subscriptionID, client)
		return bot.SubscriptionError
	}

	err = handler.sbService.UpdateSubscription(subscription)
	if err != nil {
		bot.SendReplyToTelegramChat(chatID, "❌ Error unable to add work experience for job subscription!")
		handler.HandleInitAddSubscriptionExperience(subscriptionID, client)
		return bot.SubscriptionError
	}

	reply := fmt.Sprintf(
		"<b>Job Subscription</b>\n\n"+
			"<b>Job Type</b>:  %s\n"+
			"<b>Sector</b>:  %s\n"+
			"<b>Education Level</b>:  %s\n"+
			"<b>Experience</b>:  %s\n\n",
		subscription.Type, subscription.Sector, subscription.EducationLevel, subscription.Experience)

	bot.SendReplyToTelegramChat(chatID,
		"Congratulations 🎉 you have successfully added new job subscription!")
	bot.SendReplyToTelegramChat(chatID, reply)

	subscriptionMenu := bot.CreateReplyKeyboard(true, false,
		[]string{"➕ Add Subscription", "📝 Edit Subscriptions"}, []string{"🔙 Main Menu"})
	bot.SendReplyToTelegramChat(chatID, "Choose option", subscriptionMenu)
	return bot.SubscriptionModified
}

// RemoveSubscription is a method that removes a certain job subscription of a user
func (handler *TelegramBotHandler) RemoveSubscription(subscriptionID string, user *entity.User) (string, error) {
	subscription, err := handler.sbService.DeleteSubscription(subscriptionID)
	if err != nil {
		return "🙁 Oops! unable to remove the subscription", err
	}

	reply := fmt.Sprintf(
		"------------- <b>Removed</b> -------------\n\n"+
			"<b>Job Subscription</b>\n\n"+
			"<b>Job Type</b>:  %s\n"+
			"<b>Sector</b>:  %s\n\n"+
			"------------- <b>Removed</b> -------------\n\n",
		subscription.Type, subscription.Sector)

	return reply, nil
}

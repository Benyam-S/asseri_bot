package handler

import (
	"regexp"
	"strings"

	"github.com/Benyam-S/asseri/client/bot"
	"github.com/Benyam-S/asseri/entity"
	emoji "github.com/tmdvs/Go-Emoji-Utils"
)

// HandleClient is a method that handles the whole client handling process
func (handler *TelegramBotHandler) HandleClient(update *bot.Update, user *entity.User) {
	command := emoji.RemoveAll(update.Message.Text)
	action := emoji.RemoveAll(update.CallbackQuery.Data)

	client, err := handler.clService.FindClient(user.ID)
	if err != nil {
		return
	}

	// First check for action
	if handler.HandleCallBackAction(action, update, user, client) {
		return
	}
	handler.HandleCommand(command, update, user, client)
}

// RegisterPreviousCommand is a method that register the previous command on the client
func (handler *TelegramBotHandler) RegisterPreviousCommand(command string, client *bot.Client) error {
	client.PrevCommand = command
	return handler.clService.UpdateClient(client)
}

// HandleShowMainMenu is a method that shows the main menu
func (handler *TelegramBotHandler) HandleShowMainMenu(update *bot.Update, user *entity.User) {
	menu := bot.MainMenuW
	if user.Category == entity.UserCategoryJobSeeker {
		menu = bot.MainMenuWO
	}

	bot.SendReplyToTelegramChat(update.Message.Chat.ID, "Select option", menu)
}

// HandleCallBackAction is a method that handles a callback action sent as a response and
// returns wether the action has been handled or not
func (handler *TelegramBotHandler) HandleCallBackAction(action string, update *bot.Update,
	user *entity.User, client *bot.Client) bool {

	switch client.PrevCommand {
	case "Opened":
		if strings.HasPrefix(action, "job/close/") {
			jobID := action[len("job/close/"):]
			reply, err := handler.CloseJob(jobID)

			if err != nil {
				bot.AnswerToTelegramCallBack(update.CallbackQuery.ID, reply)
			} else {
				bot.AnswerToTelegramCallBack(update.CallbackQuery.ID, "")
				bot.SendReplyToTelegramChat(update.CallbackQuery.User.ID, reply)
				bot.PostToTelegramChannel(reply)
			}

			return true
		}

	case "Job Subscriptions":
		if strings.HasPrefix(action, "subscription/remove/") {
			subscriptionID := action[len("subscription/remove/"):]
			reply, err := handler.RemoveSubscription(subscriptionID, user)

			if err != nil {
				bot.AnswerToTelegramCallBack(update.CallbackQuery.ID, reply)
			} else {
				bot.AnswerToTelegramCallBack(update.CallbackQuery.ID, "")
				bot.SendReplyToTelegramChat(update.CallbackQuery.User.ID, reply)
			}

			return true
		}

	case "Add Subscription":

		if strings.HasPrefix(action, "subscription/add/sector/") {
			jobSector := action[len("subscription/add/sector/"):]

			handler.AddSubscriptionSector(jobSector, user, client)

			bot.AnswerToTelegramCallBack(update.CallbackQuery.ID, "")

			handler.RegisterPreviousCommand("Add Subscription Sector", client)
			return true
		}

	case "Add Subscription Sector":
		reg := regexp.MustCompile(`^subscription/.+/add/type/.+$`)
		if reg.MatchString(action) {
			subscriptionID := action[len("subscription/"):strings.Index(action, "/add/type/")]
			jobType := action[strings.Index(action, "/add/type/")+len("/add/type/"):]

			typeAdded := handler.AddSubscriptionType(subscriptionID, jobType, user, client)
			if typeAdded == bot.SubscriptionNotFound {
				handler.RegisterPreviousCommand("Job Subscriptions", client)

			} else if typeAdded == bot.SubscriptionModified {
				handler.RegisterPreviousCommand("Add Subscription Type", client)
			}

			bot.AnswerToTelegramCallBack(update.CallbackQuery.ID, "")
			return true
		}

	case "Add Subscription Type":
		reg := regexp.MustCompile(`^subscription/.+/add/education_level/.+$`)
		if reg.MatchString(action) {
			subscriptionID := action[len("subscription/"):strings.Index(action, "/add/education_level/")]
			educationLevel := action[strings.Index(action, "/add/education_level/")+len("/add/education_level/"):]

			educationLevelAdded := handler.AddSubscriptionEducationLevel(subscriptionID, educationLevel, user, client)
			if educationLevelAdded == bot.SubscriptionNotFound {
				handler.RegisterPreviousCommand("Job Subscriptions", client)

			} else if educationLevelAdded == bot.SubscriptionModified {
				handler.RegisterPreviousCommand("Add Subscription Education Level", client)
			}

			bot.AnswerToTelegramCallBack(update.CallbackQuery.ID, "")
			return true
		}

	case "Add Subscription Education Level":
		reg := regexp.MustCompile(`^subscription/.+/add/experience/.+$`)
		if reg.MatchString(action) {
			subscriptionID := action[len("subscription/"):strings.Index(action, "/add/experience/")]
			experience := action[strings.Index(action, "/add/experience/")+len("/add/experience/"):]

			experienceAdded := handler.AddSubscriptionExperience(subscriptionID, experience, user, client)
			if experienceAdded == bot.SubscriptionNotFound || experienceAdded == bot.SubscriptionModified {
				handler.RegisterPreviousCommand("Job Subscriptions", client)
			}

			bot.AnswerToTelegramCallBack(update.CallbackQuery.ID, "")
			return true
		}

	default:
		// Callback that doesn't depend on pervious command
		if strings.HasPrefix(action, "job/view/") {
			jobID := action[len("job/view/"):]
			reply := handler.HandleViewJobDetail(jobID, update.CallbackQuery.User.ID)
			bot.AnswerToTelegramCallBack(update.CallbackQuery.ID, reply)
			return true
		}

	}

	return false
}

// HandleCommand is a method that handles the whole command routing process
func (handler *TelegramBotHandler) HandleCommand(command string, update *bot.Update,
	user *entity.User, client *bot.Client) {

	// Handling sub menu commands
	switch client.PrevCommand {
	case "Manage Jobs", "Pending", "Opened", "Closed", "Declined":
		switch command {
		case "Pending":
			handler.HandlePendingJobs(update, user)
			handler.RegisterPreviousCommand(command, client)
			return
		case "Opened":
			handler.HandleOpenedJobs(update, user)
			handler.RegisterPreviousCommand(command, client)
			return
		case "Closed":
			handler.HandleClosedJobs(update, user)
			handler.RegisterPreviousCommand(command, client)
			return
		case "Declined":
			handler.HandleDeclinedJobs(update, user)
			handler.RegisterPreviousCommand(command, client)
			return
		}

	case "Job Subscriptions":
		switch command {
		case "Add Subscription":
			handler.HandleInitAddSubscriptionSector(client)
			handler.RegisterPreviousCommand(command, client)
			return
		case "Edit Subscriptions":
			handler.HandleEditJobSubscriptions(update, user)
			// No previous command registration required
			return
		}

	case "Settings":
		switch command {
		case "Profile":
			handler.HandleViewProfile(update, user)
			handler.RegisterPreviousCommand(command, client)
			return
		case "Feedback":
			handler.HandlePromptFeedback(update)
			handler.RegisterPreviousCommand(command, client)
			return
		}

	case "Feedback":
		switch command {
		case "Main Menu":
			handler.HandleShowMainMenu(update, user)
			handler.RegisterPreviousCommand(command, client)
			return
		default:
			received := handler.HandleReceiveFeedback(update, user)
			if received {
				handler.HandleShowMainMenu(update, user)
				handler.RegisterPreviousCommand("Main Menu", client)
			}
			return
		}

	case "Profile":
		switch command {
		case "Update Profile":
			handler.HandleInitUpdateProfile(update, user)
			handler.RegisterPreviousCommand(command, client)
			return
		}

	case "Update Profile":
		switch command {
		case "Skip":
			backMenu := bot.CreateReplyKeyboard(true, false, []string{"↖️ Skip", "🔙 Main Menu"})
			bot.SendReplyToTelegramChat(update.Message.Chat.ID, "Enter new phonenumber", backMenu)
			handler.RegisterPreviousCommand("Update Name", client)
			return
		case "Main Menu":
			handler.HandleShowMainMenu(update, user)
			handler.RegisterPreviousCommand(command, client)
			return
		default:
			updated := handler.HandleUpdateName(update, user)
			if updated {
				handler.RegisterPreviousCommand("Update Name", client)
			}
			return
		}

	case "Update Name":
		switch command {
		case "Skip":
			keyboard := bot.CreateReplyKeyboard(true, false, []string{"አሰሪ", "Job Seeker", "Agent"},
				[]string{"↖️ Skip", "🔙 Main Menu"})
			bot.SendReplyToTelegramChat(update.Message.Chat.ID, "Wish to change category?", keyboard)
			handler.RegisterPreviousCommand("Update Phonenumber", client)
			return
		case "Main Menu":
			handler.HandleShowMainMenu(update, user)
			handler.RegisterPreviousCommand(command, client)
			return
		default:
			updated := handler.HandleUpdatePhone(update, user)
			if updated {
				handler.RegisterPreviousCommand("Update Phonenumber", client)
			}
			return
		}

	case "Update Phonenumber":
		switch command {
		case "Skip":
			bot.SendReplyToTelegramChat(update.Message.Chat.ID,
				"Congratulations 🎉 you have successfully update your profile!")
			handler.HandleViewProfile(update, user)
			handler.RegisterPreviousCommand("Profile", client)
			return
		case "Main Menu":
			handler.HandleShowMainMenu(update, user)
			handler.RegisterPreviousCommand(command, client)
			return
		default:
			updated := handler.HandleUpdateCategory(update, user)
			if updated {
				handler.HandleViewProfile(update, user)
				handler.RegisterPreviousCommand("Profile", client)
			}
			return
		}
	}

	// Applying Process
	if strings.Contains(client.PrevCommand, "Apply ") {

		if command == "Cancel Application" {
			handler.HandleShowMainMenu(update, user)
			handler.RegisterPreviousCommand(command, client)
			return
		}

		jobID := client.PrevCommand[len("Apply "):]
		err := handler.HandleApplyForJob(update, jobID, user)
		if err == nil || err.Error() == "unable to apply for the job" {
			handler.HandleShowMainMenu(update, user)
			handler.RegisterPreviousCommand(command, client)
			return
		}
	}

	// Handling the main menu commands
	switch command {
	case "Main Menu":
		handler.HandleShowMainMenu(update, user)
		handler.RegisterPreviousCommand(command, client)
		return
	case "Post Job":
		handler.HandlePostJob(update, user)
		handler.RegisterPreviousCommand(command, client)
		return
	case "Manage Jobs":
		handler.HandleMangeJobs(update)
		handler.RegisterPreviousCommand(command, client)
		return
	case "Job Subscriptions":
		handler.HandleJobSubscription(update, user)
		handler.RegisterPreviousCommand(command, client)
		return
	case "Settings":
		handler.HandleSettings(update)
		handler.RegisterPreviousCommand(command, client)
		return
	}

	// Star command for job application and get menu flow
	if strings.Contains(command, "/start apply_") {
		jobID := command[len("/start apply_"):]
		if handler.HandleInitApplyForJob(jobID, user, update.Message.Chat.ID) {
			handler.RegisterPreviousCommand("Apply "+jobID, client)
		} else {
			handler.HandleShowMainMenu(update, user)
			handler.RegisterPreviousCommand(command, client)
			return
		}

	} else if command == "/start" {
		handler.HandleShowMainMenu(update, user)
		handler.RegisterPreviousCommand(command, client)
		return
	}
}

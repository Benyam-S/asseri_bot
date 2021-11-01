package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/Benyam-S/asseri/client/bot"
	"github.com/Benyam-S/asseri/entity"
	"github.com/Benyam-S/asseri/tools"
	"github.com/google/uuid"
)

// HandlePostJob is a method that sends the job posting website url
func (handler *TelegramBotHandler) HandlePostJob(update *bot.Update, user *entity.User) {

	if user.Category == entity.UserCategoryJobSeeker {
		bot.SendReplyToTelegramChat(update.Message.Chat.ID, "🙁 Oops! Can't perform operation for job seeker.")
		return
	}

	// storing the job posting access token
	accessToken := uuid.Must(uuid.NewRandom())
	handler.store.Add(accessToken.String(), user.ID)

	bot.SendReplyToTelegramChat(update.Message.Chat.ID, fmt.Sprintf(`please follow the following link to post a job. 
	 https://www.asseri.net/job/post.html?employer_id=%s&access_token=%s`, user.ID, accessToken))
}

// HandleMangeJobs is a method that handles the job managing process
func (handler *TelegramBotHandler) HandleMangeJobs(update *bot.Update) {

	jobStatusMenu := bot.CreateReplyKeyboard(true, false, []string{"⌛ Pending", "📖 Opened"},
		[]string{"📕 Closed", "🚫 Declined"}, []string{"🔙 Main Menu"})
	bot.SendReplyToTelegramChat(update.Message.Chat.ID, "Choose status", jobStatusMenu)
}

// HandleViewJobDetail is a method that enables user to view a certain job details
func (handler *TelegramBotHandler) HandleViewJobDetail(jobID string, chatID int64) string {

	var inlineKeyboard string
	job, err := handler.jbService.FindJob(jobID)
	if err != nil {
		return "😳 Oops! unable to view job detail."
	}

	reply := fmt.Sprintf(
		"<b>Job Title</b>:  %s\n\n"+
			"<b>Job Type</b>:  %s\n"+
			"<b>Gender</b>:  %s\n"+
			"<b>Education level</b>:  %s\n"+
			"<b>Experience</b>:  %s\n"+
			"<b>Contact Type</b>:  %s\n\n"+
			"<b>Description</b>:  %s\n\n"+"#%s",
		job.Title, job.Type, bot.GetGender(job.Gender),
		job.EducationLevel, job.Experience, job.ContactType,
		job.Description, tools.ChangeSpaceToUnderscore(job.Sector))

	if job.Status == entity.JobStatusOpened {
		inlineKeyboard = bot.CreateInlineKeyboard([]bot.InlineKeyboardButton{
			{Text: "❌ Close", CallbackData: "job/close/" + job.ID},
		})
	}

	bot.SendReplyToTelegramChat(chatID, reply, inlineKeyboard)
	return ""
}

// HandlePendingJobs is a method that shows all the pending jobs sent by the user that are waiting approval
func (handler *TelegramBotHandler) HandlePendingJobs(update *bot.Update, user *entity.User) {

	jobs := handler.jbService.FindMultipleJobs(user.ID)
	pendingJobs := make([]*entity.Job, 0)

	for _, job := range jobs {
		if job.Status == entity.JobStatusPending {
			pendingJobs = append(pendingJobs, job)
		}
	}

	if len(pendingJobs) == 0 {
		bot.SendReplyToTelegramChat(update.Message.Chat.ID, "You don't have any pending"+
			" job waiting for approval.")
		return
	}

	for _, pendingJob := range pendingJobs {
		reply := fmt.Sprintf(
			"<b>Job Title</b>:  %s\n\n"+
				"<b>Job Type</b>:  %s\n"+
				"<b>Gender</b>:  %s\n"+
				"<b>Education level</b>:  %s\n"+
				"<b>Experience</b>:  %s\n"+
				"<b>Contact Type</b>:  %s\n\n"+
				"<b>Description</b>:  %s\n\n"+"#%s",
			pendingJob.Title, pendingJob.Type, bot.GetGender(pendingJob.Gender),
			pendingJob.EducationLevel, pendingJob.Experience, pendingJob.ContactType,
			pendingJob.Description, tools.ChangeSpaceToUnderscore(pendingJob.Sector))

		bot.SendReplyToTelegramChat(update.Message.Chat.ID, reply)
	}
}

// HandleOpenedJobs is a method that shows all the opened jobs sent by the user
func (handler *TelegramBotHandler) HandleOpenedJobs(update *bot.Update, user *entity.User) {

	jobs := handler.jbService.FindMultipleJobs(user.ID)
	openedJobs := make([]*entity.Job, 0)

	for _, job := range jobs {
		if job.Status == entity.JobStatusOpened {
			openedJobs = append(openedJobs, job)
		}
	}

	if len(openedJobs) == 0 {
		bot.SendReplyToTelegramChat(update.Message.Chat.ID, "You don't have any opened"+
			" job waiting for an applier.")
		return
	}

	for _, openedJob := range openedJobs {
		reply := fmt.Sprintf(
			"<b>Job Title</b>:  %s\n\n"+
				"<b>Job Type</b>:  %s\n"+
				"<b>Gender</b>:  %s\n"+
				"<b>Education level</b>:  %s\n"+
				"<b>Experience</b>:  %s\n"+
				"<b>Contact Type</b>:  %s\n\n"+
				"<b>Description</b>:  %s\n\n"+"#%s",
			openedJob.Title, openedJob.Type, bot.GetGender(openedJob.Gender),
			openedJob.EducationLevel, openedJob.Experience, openedJob.ContactType,
			openedJob.Description, tools.ChangeSpaceToUnderscore(openedJob.Sector))

		inlineKeyboard := bot.CreateInlineKeyboard([]bot.InlineKeyboardButton{
			{Text: "❌ Close", CallbackData: "job/close/" + openedJob.ID},
		})
		bot.SendReplyToTelegramChat(update.Message.Chat.ID, reply, inlineKeyboard)
	}
}

// HandleClosedJobs is a method that shows all the closed jobs of a certain user
func (handler *TelegramBotHandler) HandleClosedJobs(update *bot.Update, user *entity.User) {

	jobs := handler.jbService.FindMultipleJobs(user.ID)
	closedJobs := make([]*entity.Job, 0)

	for _, job := range jobs {
		if job.Status == entity.JobStatusClosed {
			closedJobs = append(closedJobs, job)
		}
	}

	if len(closedJobs) == 0 {
		bot.SendReplyToTelegramChat(update.Message.Chat.ID, "You don't have any closed job.")
		return
	}

	for _, closedJob := range closedJobs {
		reply := fmt.Sprintf(
			"------------- <b>Closed</b> -------------\n\n"+
				"<b>Job Title</b>:  %s\n\n"+
				"<b>Job Type</b>:  %s\n"+
				"<b>Gender</b>:  %s\n"+
				"<b>Education level</b>:  %s\n"+
				"<b>Experience</b>:  %s\n"+
				"<b>Contact Type</b>:  %s\n\n"+
				"<b>Description</b>:  %s\n\n"+"#%s\n\n"+
				"------------- <b>Closed</b> -------------\n\n",
			closedJob.Title, closedJob.Type, bot.GetGender(closedJob.Gender),
			closedJob.EducationLevel, closedJob.Experience, closedJob.ContactType,
			closedJob.Description, tools.ChangeSpaceToUnderscore(closedJob.Sector))

		bot.SendReplyToTelegramChat(update.Message.Chat.ID, reply)
	}
}

// HandleDeclinedJobs is a method that shows all the declined jobs of a certain user
func (handler *TelegramBotHandler) HandleDeclinedJobs(update *bot.Update, user *entity.User) {

	jobs := handler.jbService.FindMultipleJobs(user.ID)
	declinedJobs := make([]*entity.Job, 0)

	for _, job := range jobs {
		if job.Status == entity.JobStatusDecelined {
			declinedJobs = append(declinedJobs, job)
		}
	}

	if len(declinedJobs) == 0 {
		bot.SendReplyToTelegramChat(update.Message.Chat.ID, "You don't have any declined job.")
		return
	}

	for _, declinedJob := range declinedJobs {
		reply := fmt.Sprintf(
			"------------- <b>Declined</b> -------------\n\n"+
				"<b>Job Title</b>:  %s\n\n"+
				"<b>Job Type</b>:  %s\n"+
				"<b>Gender</b>:  %s\n"+
				"<b>Education level</b>:  %s\n"+
				"<b>Experience</b>:  %s\n"+
				"<b>Contact Type</b>:  %s\n\n"+
				"<b>Description</b>:  %s\n\n"+"#%s\n\n"+
				"------------- <b>Declined</b> -------------\n\n",
			declinedJob.Title, declinedJob.Type, bot.GetGender(declinedJob.Gender),
			declinedJob.EducationLevel, declinedJob.Experience, declinedJob.ContactType,
			declinedJob.Description, tools.ChangeSpaceToUnderscore(declinedJob.Sector))

		bot.SendReplyToTelegramChat(update.Message.Chat.ID, reply)
	}
}

// CloseJob is a method that closes a certain job so no one can apply for the job
func (handler *TelegramBotHandler) CloseJob(jobID string) (string, error) {
	job, err := handler.jbService.FindJob(jobID)
	if err != nil {
		return "🙁 Oops! unable to close the job", err
	}

	if job.Status != entity.JobStatusOpened {
		return "🙁 Oops! unable to close the job", errors.New("unable to close unopened job")
	}

	job.Status = entity.JobStatusClosed
	err = handler.jbService.UpdateJob(job)
	if err != nil {
		return "🙁 Oops! unable to close the job", err
	}

	reply := fmt.Sprintf(
		"------------- <b>Closed</b> -------------\n\n"+
			"<b>Job Title</b>:  %s\n\n"+
			"<b>Job Type</b>:  %s\n"+
			"<b>Gender</b>:  %s\n"+
			"<b>Education level</b>:  %s\n"+
			"<b>Experience</b>:  %s\n"+
			"<b>Contact Type</b>:  %s\n\n"+
			"<b>Description</b>:  %s\n\n"+"#%s\n\n"+
			"------------- <b>Closed</b> -------------\n\n",
		job.Title, job.Type, bot.GetGender(job.Gender),
		job.EducationLevel, job.Experience, job.ContactType,
		job.Description, tools.ChangeSpaceToUnderscore(job.Sector))

	return reply, nil
}

// HandleInitApplyForJob is a method that prompt user to send cv
func (handler *TelegramBotHandler) HandleInitApplyForJob(jobID string, user *entity.User, chatID int64) bool {

	if _, errMsg := handler.IsJobApplicable(jobID, user); errMsg != "" {
		bot.SendReplyToTelegramChat(chatID, errMsg)
		return false
	}

	if handler.jaService.JobApplicationExists(jobID, user.ID) {
		bot.SendReplyToTelegramChat(chatID, "❌ You have already applied for the job")
		return false
	}

	cancelMenu := bot.CreateReplyKeyboard(true, false,
		[]string{"🔙 Cancel Application"})
	bot.SendReplyToTelegramChat(chatID, "Send CV (*PDF format only)", cancelMenu)
	return true
}

// HandleApplyForJob is a method that enables user to apply for a job using their cv
func (handler *TelegramBotHandler) HandleApplyForJob(update *bot.Update, jobID string, user *entity.User) error {

	job, errMsg := handler.IsJobApplicable(jobID, user)
	if errMsg != "" {
		bot.SendReplyToTelegramChat(update.Message.Chat.ID, errMsg)
		return errors.New("unable to apply for the job")
	}

	// Verifying the file with type
	file := update.Message.Document
	if file.Type != "application/pdf" {
		bot.SendReplyToTelegramChat(update.Message.Chat.ID, "Please send *.pdf file only")
		return errors.New("invalid format")
	}

	client, err := handler.clService.FindClient(job.Employer)
	if err != nil {
		bot.SendReplyToTelegramChat(update.Message.Chat.ID, "🙁 Unable to apply for the job")
		return errors.New("unable to apply for the job")
	}

	jobApplication := new(entity.JobApplication)
	jobApplication.JobID = jobID
	jobApplication.JobSeekerID = user.ID

	err = handler.jaService.AddJobApplication(jobApplication)
	if err != nil {
		bot.SendReplyToTelegramChat(update.Message.Chat.ID, "🙁 Unable to apply for the job")
		return errors.New("unable to apply for the job")
	}

	chatID, _ := strconv.ParseInt(client.TelegramID, 10, 64)
	applyCaption := fmt.Sprintf(
		"------------- <b>Application</b> -------------\n\n"+
			"<b>Job Title</b>:  %s\n\n"+
			"<b>Description</b>:  %s\n\n", job.Title, job.Description)

	inlineKeyboard := bot.CreateInlineKeyboard([]bot.InlineKeyboardButton{
		{Text: "👀 Job Details", CallbackData: "job/view/" + job.ID},
	})

	type TelegramResponse struct {
		Ok bool `json:"ok"`
	}

	response, err := bot.SendDocumentToTelegramChat(chatID, file.ID, applyCaption, inlineKeyboard)
	result := new(TelegramResponse)
	json.Unmarshal([]byte(response), result)

	if err != nil || !result.Ok {
		handler.jaService.DeleteJobApplication(jobID, user.ID)
		bot.SendReplyToTelegramChat(update.Message.Chat.ID, "😳 Oops! something went wrong, re-apply again.")
		return errors.New("application not completed")
	}

	bot.SendReplyToTelegramChat(update.Message.Chat.ID,
		"🎉 Your application has been sent to the employer. Good Luck!")
	return nil
}

// IsJobApplicable is a method that identify the applicability of a job
func (handler *TelegramBotHandler) IsJobApplicable(jobID string, user *entity.User) (*entity.Job, string) {

	job, err := handler.jbService.FindJob(jobID)
	if err != nil {
		return nil, "🙁 Unable to apply for the job"
	}

	if job.Status == entity.JobStatusClosed {
		return nil, "🙁 Sorry the job has been closed"

	} else if job.Status != entity.JobStatusOpened {
		return nil, "🙁 Unable to apply for the job"
	}

	if job.Employer == user.ID {
		return nil, "❌ You can't apply for your own job"
	}

	return job, ""
}

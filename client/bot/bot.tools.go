package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/Benyam-S/asseri/entity"
	"github.com/Benyam-S/asseri/tools"
)

// SendReplyToTelegramChat sends a reply to the Telegram chat identified by its chat Id
func SendReplyToTelegramChat(chatID int64, reply ...string) (string, error) {

	text := ""
	replyMarkup := ""

	if len(reply) > 0 {
		text = reply[0]
	}

	if len(reply) > 1 {
		replyMarkup = reply[1]
	}

	var telegramAPI string = os.Getenv("api_access_point") + os.Getenv("bot_api_token") + "/sendMessage"
	response, err := http.PostForm(
		telegramAPI,
		url.Values{
			"chat_id":      {strconv.FormatInt(chatID, 10)},
			"text":         {text},
			"reply_markup": {replyMarkup},
			"parse_mode":   {"html"},
		})

	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	var bodyBytes, errRead = ioutil.ReadAll(response.Body)
	if errRead != nil {
		return "", err
	}

	return string(bodyBytes), nil
}

// SendDocumentToTelegramChat sends a document to the Telegram chat identified by its chat Id
func SendDocumentToTelegramChat(chatID int64, fileID string, reply ...string) (string, error) {

	caption := ""
	replyMarkup := ""

	if len(reply) > 0 {
		caption = reply[0]
	}

	if len(reply) > 1 {
		replyMarkup = reply[1]
	}

	var telegramAPI string = os.Getenv("api_access_point") + os.Getenv("bot_api_token") + "/sendDocument"
	response, err := http.PostForm(
		telegramAPI,
		url.Values{
			"chat_id":      {strconv.FormatInt(chatID, 10)},
			"document":     {fileID},
			"caption":      {caption},
			"parse_mode":   {"html"},
			"reply_markup": {replyMarkup},
		})

	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	var bodyBytes, errRead = ioutil.ReadAll(response.Body)
	if errRead != nil {
		return "", err
	}

	return string(bodyBytes), nil
}

// PostToTelegramChannel posts a certain content to a telegram channel
func PostToTelegramChannel(post ...string) (string, error) {

	text := ""
	replyMarkup := ""

	if len(post) > 0 {
		text = post[0]
	}

	if len(post) > 1 {
		replyMarkup = post[1]
	}

	var telegramAPI string = os.Getenv("api_access_point") + os.Getenv("bot_api_token") + "/sendMessage"
	response, err := http.PostForm(
		telegramAPI,
		url.Values{
			"chat_id":      {os.Getenv("channel_name")},
			"text":         {text},
			"reply_markup": {replyMarkup},
			"parse_mode":   {"html"},
		})

	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	var bodyBytes, errRead = ioutil.ReadAll(response.Body)
	if errRead != nil {
		return "", err
	}

	return string(bodyBytes), nil
}

// AnswerToTelegramCallBack sends a reply to the Telegram call back request identified by the query id
func AnswerToTelegramCallBack(queryID string, text string) (string, error) {

	var telegramAPI string = os.Getenv("api_access_point") + os.Getenv("bot_api_token") + "/answerCallbackQuery"
	response, err := http.PostForm(
		telegramAPI,
		url.Values{
			"callback_query_id": {queryID},
			"text":              {text},
		})

	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	var bodyBytes, errRead = ioutil.ReadAll(response.Body)
	if errRead != nil {
		return "", err
	}

	return string(bodyBytes), nil
}

// CreateReplyKeyboard is a function that creates a reply keyboard from set of parameters
func CreateReplyKeyboard(resizeKeyboard, oneTimeKeyboard bool, keyboardButtons ...[]string) string {

	buttonRows := make([][]*ReplyKeyboardButton, 0)

	for _, keyboardRow := range keyboardButtons {

		row := make([]*ReplyKeyboardButton, 0)
		for _, keyboardText := range keyboardRow {
			button := new(ReplyKeyboardButton)
			button.Text = keyboardText
			row = append(row, button)
		}

		buttonRows = append(buttonRows, row)
	}

	keyboard := ReplyKeyboardMarkup{
		Keyboard:        buttonRows,
		ResizeKeyboard:  resizeKeyboard,
		OneTimeKeyboard: oneTimeKeyboard,
	}

	keyboardS, _ := json.Marshal(&keyboard)
	return string(keyboardS)
}

// CreateReplyKeyboardWExtra is a function that creates a reply keyboard from set of parameters with extra capabilities
func CreateReplyKeyboardWExtra(resizeKeyboard, oneTimeKeyboard bool, keyboardButtons ...[]ReplyKeyboardButton) string {

	buttonRows := make([][]*ReplyKeyboardButton, 0)

	for _, keyboardRow := range keyboardButtons {

		row := make([]*ReplyKeyboardButton, 0)
		for _, keyboardButton := range keyboardRow {
			button := new(ReplyKeyboardButton)
			button.Text = keyboardButton.Text
			button.RequestContact = keyboardButton.RequestContact
			row = append(row, button)
		}

		buttonRows = append(buttonRows, row)
	}

	keyboard := ReplyKeyboardMarkup{
		Keyboard:        buttonRows,
		ResizeKeyboard:  resizeKeyboard,
		OneTimeKeyboard: oneTimeKeyboard,
	}

	keyboardS, _ := json.Marshal(&keyboard)
	return string(keyboardS)
}

// CreateInlineKeyboard is a function that creates an inline keyboard from set of parameters for a chat
func CreateInlineKeyboard(keyboardButtons ...[]InlineKeyboardButton) string {

	buttonRows := make([][]*InlineKeyboardButton, 0)

	for _, keyboardRow := range keyboardButtons {

		row := make([]*InlineKeyboardButton, 0)
		for _, keyboardButton := range keyboardRow {
			button := new(InlineKeyboardButton)
			button.Text = keyboardButton.Text
			button.URL = keyboardButton.URL
			button.CallbackData = keyboardButton.CallbackData
			row = append(row, button)
		}

		buttonRows = append(buttonRows, row)
	}

	keyboard := InlineKeyboardMarkup{
		InlineKeyboard: buttonRows,
	}

	keyboardS, _ := json.Marshal(&keyboard)
	return string(keyboardS)
}

// GetGender is a function tha get the appropriate gender value for a given gender acronym
func GetGender(gender string) string {

	switch gender {
	case "M":
		gender = "Male"
	case "F":
		gender = "Female"
	case "B":
		fallthrough
	default:
		gender = "Both"
	}

	return gender
}

// BuildNotification is a function that builds a notification and keyboard from given job
func BuildNotification(job *entity.Job, pack *StructuredPackage) string {

	var notification string
	var jobSector = ""
	var educationLevel = ""
	var jobType = ""
	var jobTypes = make([]string, 0)
	var jobSectors = strings.Split(strings.TrimSpace(job.Sector), ",")

	for _, jobS := range jobSectors {
		jobSector += fmt.Sprintf("#%s    ", tools.ChangeSpaceToUnderscore(jobS))
	}

	for _, jobT := range strings.Split(strings.TrimSpace(job.Type), ",") {
		if strings.ToLower(jobT) != "other" {
			jobTypes = append(jobTypes, jobT)
		}
	}

	if strings.ToLower(job.EducationLevel) != "other" {
		educationLevel = job.EducationLevel
	}

	jobType = strings.Join(jobTypes, ", ")

	emptyJobType, _ := regexp.MatchString(`^\s*$`, jobType)
	if !emptyJobType {
		jobType = fmt.Sprintf("<b>Job Type</b>:  %s\n", jobType)
	} else {
		jobType = ""
	}

	emptyEducationLevel, _ := regexp.MatchString(`^\s*$`, educationLevel)
	if !emptyEducationLevel {
		educationLevel = fmt.Sprintf("<b>Education level</b>:  %s\n", educationLevel)
	} else {
		educationLevel = ""
	}

	switch job.PostType {
	case entity.PostCategoryUser:
		notification = fmt.Sprintf(
			"<b>Job Title</b>:  %s\n\n"+
				"<b>አሰሪ</b>:  %s\n\n"+
				jobType+
				"<b>Gender</b>:  %s\n"+
				educationLevel+
				"<b>Experience</b>:  %s\n\n"+
				"<b>Description</b>:  %s\n\n%s"+"%s\n\n"+
				"@asseri_bot         @asseri_bot\n\n",
			job.Title, pack.Employer, GetGender(job.Gender), job.Experience,
			job.Description, pack.Contact, jobSector)
	case entity.PostCategoryInternal:
		notification = fmt.Sprintf(
			"<b>Job Title</b>:  %s\n\n"+
				"<b>አሰሪ</b>:  %s\n\n"+
				jobType+
				"<b>Gender</b>:  %s\n"+
				educationLevel+
				"<b>Experience</b>:  %s\n\n"+
				"<b>Description</b>:  %s\n\n%s"+"%s\n\n"+
				"@asseri_bot         @asseri_bot\n\n",
			job.Title, pack.Employer, GetGender(job.Gender), job.Experience,
			job.Description, pack.Contact, jobSector)
	case entity.PostCategoryExternal:
		notification = fmt.Sprintf(
			"<b>Job Title</b>:  %s\n\n"+
				"<b>አሰሪ</b>:  %s\n\n"+
				jobType+
				educationLevel+
				"<b>Experience</b>:  %s\n\n"+
				"<b>Description</b>:  %s\n\n"+"%s\n\n"+
				"@asseri_bot         @asseri_bot\n\n",
			job.Title, pack.Employer, job.Experience, job.Description,
			jobSector)
	}

	return notification
}

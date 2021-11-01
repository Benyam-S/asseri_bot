package bot

import (
	"time"

	"github.com/Benyam-S/asseri/entity"
)

// TempUser is a struct that holds the temporary user data before registration
type TempUser struct {
	TelegramID  string `gorm:"primary_key; unique; not null"`
	UserName    string
	PhoneNumber string
	Category    string
	Status      int64 // Used to identify the state of the registration process
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Client is a struct that defines the relation between user and Telegram bot user
type Client struct {
	UserID      string `gorm:"primary_key; unique; not null"`
	TelegramID  string `gorm:"unique; not null"`
	PrevCommand string
}

// TableName overrides the table name used by Client to `bot_clients`
func (Client) TableName() string {
	return "bot_clients"
}

// Update is a Telegram object that the handler receives every time an user interacts with the bot.
type Update struct {
	UpdateID      int64         `json:"update_id"`
	Message       Message       `json:"message"`
	CallbackQuery CallbackQuery `json:"callback_query"`
}

// Message is a Telegram object that can be found inside an update.
type Message struct {
	Text     string    `json:"text"`
	Chat     Chat      `json:"chat"`
	User     TUser     `json:"from"`
	Document TDocument `json:"document"`
	Contact  TContact  `json:"contact"`
}

// CallbackQuery is a Telegram object that can be found inside an update.
type CallbackQuery struct {
	ID   string `json:"id"`
	Data string `json:"data"`
	User TUser  `json:"from"`
}

// Chat indicates the conversation to which the message belongs.
type Chat struct {
	ID int64 `json:"id"`
}

// TUser is a Telegram user object
type TUser struct {
	ID           int64  `json:"id"`
	LanguageCode string `json:"language_code"`
}

// TDocument is a Telegram document object
type TDocument struct {
	ID       string `json:"file_id"`
	UniqueID string `json:"file_unique_id"`
	Name     string `json:"file_name"`
	Type     string `json:"mime_type"`
}

// TContact is a Telegram contact object
type TContact struct {
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
}

// ReplyKeyboardMarkup is a struct that represents a reply to form Telegram keyboard
type ReplyKeyboardMarkup struct {
	Keyboard        [][]*ReplyKeyboardButton `json:"keyboard"`
	ResizeKeyboard  bool                     `json:"resize_keyboard"`
	OneTimeKeyboard bool                     `json:"one_time_keyboard"`
}

// ReplyKeyboardButton is a struct that represents a Telegram reply keyboard button
type ReplyKeyboardButton struct {
	Text           string `json:"text"`
	RequestContact bool   `json:"request_contact"`
}

// InlineKeyboardMarkup is a struct that represents an inline keyboard for a reply chat
type InlineKeyboardMarkup struct {
	InlineKeyboard [][]*InlineKeyboardButton `json:"inline_keyboard"`
}

// InlineKeyboardButton is a struct that represents a Telegram inline keyboard button
type InlineKeyboardButton struct {
	Text         string `json:"text"`
	URL          string `json:"url"`
	CallbackData string `json:"callback_data"`
}

// ToUser is a method that converts TempUser to entity.User
func (tempUser *TempUser) ToUser() *entity.User {
	user := new(entity.User)

	user.UserName = tempUser.UserName
	user.PhoneNumber = tempUser.PhoneNumber
	user.Category = tempUser.Category

	return user
}

// StructuredPackage is a type that holds all the structured and modified entities ready for consumption
type StructuredPackage struct {
	Employer string
	Contact  string
}

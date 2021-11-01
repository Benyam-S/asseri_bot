package handler

import (
	"github.com/Benyam-S/asseri/client/bot/client"
	"github.com/Benyam-S/asseri/client/bot/tempuser"
	"github.com/Benyam-S/asseri/common"
	"github.com/Benyam-S/asseri/feedback"
	"github.com/Benyam-S/asseri/job"
	"github.com/Benyam-S/asseri/jobapplication"
	"github.com/Benyam-S/asseri/log"
	"github.com/Benyam-S/asseri/subscription"
	"github.com/Benyam-S/asseri/tools"
	"github.com/Benyam-S/asseri/user"
)

// TelegramBotHandler is a struct that defines a telegram bot handler
type TelegramBotHandler struct {
	tuService tempuser.IService
	clService client.IService
	urService user.IService
	jbService job.IService
	jaService jobapplication.IService
	sbService subscription.IService
	fdService feedback.IService
	cmService common.IService
	logger    *log.Logger
	store     tools.IStore
	pushChan  chan string
	pq        common.IPushQueue
}

// BotResponse is a type that defines a bot response message
type BotResponse struct {
	Ok        bool  `json:"ok"`
	ErrorCode int64 `json:"error_code"`
}

// NewTelegramBotHandler is a function that returns a new telegram bot handler
func NewTelegramBotHandler(tempUserService tempuser.IService, clientService client.IService,
	userService user.IService, jobService job.IService,
	jobApplicationService jobapplication.IService, subscriptionService subscription.IService,
	feedbackService feedback.IService, commonService common.IService, store tools.IStore,
	pushChannel chan string, pushQueue common.IPushQueue, log *log.Logger) *TelegramBotHandler {
	return &TelegramBotHandler{
		tuService: tempUserService, clService: clientService, urService: userService,
		jbService: jobService, jaService: jobApplicationService, sbService: subscriptionService,
		fdService: feedbackService, cmService: commonService, pq: pushQueue, store: store,
		pushChan: pushChannel, logger: log}
}

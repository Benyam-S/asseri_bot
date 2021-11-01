package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Benyam-S/asseri/client/bot"
	"github.com/Benyam-S/asseri/client/bot/handler"
	"github.com/Benyam-S/asseri/entity"
	"github.com/Benyam-S/asseri/log"
	"github.com/Benyam-S/asseri/tools"
	"github.com/go-redis/redis"

	cmRepository "github.com/Benyam-S/asseri/common/repository"
	cmService "github.com/Benyam-S/asseri/common/service"

	jaRepository "github.com/Benyam-S/asseri/jobapplication/repository"
	jaService "github.com/Benyam-S/asseri/jobapplication/service"

	jbRepository "github.com/Benyam-S/asseri/job/repository"
	jbService "github.com/Benyam-S/asseri/job/service"

	urRepository "github.com/Benyam-S/asseri/user/repository"
	urService "github.com/Benyam-S/asseri/user/service"

	fdRepository "github.com/Benyam-S/asseri/feedback/repository"
	fdService "github.com/Benyam-S/asseri/feedback/service"

	sbRepository "github.com/Benyam-S/asseri/subscription/repository"
	sbService "github.com/Benyam-S/asseri/subscription/service"

	tuRepository "github.com/Benyam-S/asseri/client/bot/tempuser/repository"
	tuService "github.com/Benyam-S/asseri/client/bot/tempuser/service"

	clRepository "github.com/Benyam-S/asseri/client/bot/client/repository"
	clService "github.com/Benyam-S/asseri/client/bot/client/service"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

var (
	configFilesDir string
	redisClient    *redis.Client
	mysqlDB        *gorm.DB
	sysConfig      SystemConfig
	err            error
	botHandler     *handler.TelegramBotHandler
)

// SystemConfig is a type that defines a server system configuration file
type SystemConfig struct {
	RedisClient         map[string]string `json:"redis_client"`
	MysqlClient         map[string]string `json:"mysql_client"`
	BotDomainAddres     string            `json:"bot_domain_address"`
	BotClientServerPort string            `json:"bot_client_server_port"`
	ServerLogFile       string            `json:"server_log_file"`
	BotLogFile          string            `json:"bot_log_file"`
}

// initServer initialize the web server for takeoff
func initServer() {

	// Reading data from config.server.json file and creating the systemconfig  object
	sysConfigDir := filepath.Join(configFilesDir, "/config.server.json")
	sysConfigData, _ := ioutil.ReadFile(sysConfigDir)

	// Reading data from config.asseri.json file
	asseriConfig := make(map[string]interface{})
	asseriConfigDir := filepath.Join(configFilesDir, "/config.asseri.json")
	asseriConfigData, _ := ioutil.ReadFile(asseriConfigDir)

	err = json.Unmarshal(sysConfigData, &sysConfig)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(asseriConfigData, &asseriConfig)
	if err != nil {
		panic(err)
	}

	apiAccessPoint, ok1 := asseriConfig["api_access_point"].(string)
	botAPIToken, ok2 := asseriConfig["bot_api_token"].(string)
	channelName, ok3 := asseriConfig["channel_name"].(string)
	botURL, ok4 := asseriConfig["bot_url"].(string)

	if !ok1 || !ok2 || !ok3 || !ok4 {
		panic(errors.New("unable to parse asseri config data"))
	}

	// Setting environmental variables so they can be used any where on the application
	os.Setenv("config_files_dir", configFilesDir)
	os.Setenv("bot_domain_address", sysConfig.BotDomainAddres)
	os.Setenv("bot_client_server_port", sysConfig.BotClientServerPort)
	os.Setenv("server_log_file", sysConfig.ServerLogFile)
	os.Setenv("bot_log_file", sysConfig.BotLogFile)

	os.Setenv("api_access_point", apiAccessPoint)
	os.Setenv("bot_api_token", botAPIToken)
	os.Setenv("channel_name", channelName)
	os.Setenv("bot_url", botURL)

	// Initializing the database with the needed tables and values
	initDB()

	path, _ := os.Getwd()

	jobApplicationRepo := jaRepository.NewJobApplicationRepository(mysqlDB)
	jobRepo := jbRepository.NewJobRepository(mysqlDB)
	userRepo := urRepository.NewUserRepository(mysqlDB)
	subscriptionRepo := sbRepository.NewSubscriptionRepository(mysqlDB)
	feedbackRepo := fdRepository.NewFeedbackRepository(mysqlDB)
	commonRepo := cmRepository.NewCommonRepository(mysqlDB)

	commonService := cmService.NewCommonService(commonRepo)
	userService := urService.NewUserService(userRepo, jobRepo, jobApplicationRepo, commonRepo)
	jobService := jbService.NewJobService(jobRepo, userRepo, commonService)
	jobApplicationService := jaService.NewJobApplicationService(jobApplicationRepo, commonRepo)
	subscriptionService := sbService.NewSubscriptionService(subscriptionRepo, commonService)
	feedbackService := fdService.NewFeedbackService(feedbackRepo, userRepo)

	// Creating push channel and queue
	pushChannel := make(chan string, 1000)
	pushQueue := cmService.NewPushQueue()
	logger := &log.Logger{ServerLogFile: filepath.Join(path, "../../log", os.Getenv("server_log_file")),
		BotLogFile: filepath.Join(path, "../../log", os.Getenv("bot_log_file"))}

	// ----- Bot level init -----
	tempUserRepo := tuRepository.NewTempUserRepository(mysqlDB)
	clientRepo := clRepository.NewClientRepository(mysqlDB)

	tempUserService := tuService.NewTempUserService(tempUserRepo, userRepo, commonRepo)
	clientService := clService.NewClientService(clientRepo)

	// ----- Creating store -----
	store := tools.NewRedisStore(redisClient)

	botHandler = handler.NewTelegramBotHandler(tempUserService, clientService, userService,
		jobService, jobApplicationService, subscriptionService, feedbackService,
		commonService, store, pushChannel, pushQueue, logger)
}

// initDB initialize the database for takeoff
func initDB() {

	redisDB, err := strconv.ParseInt(sysConfig.RedisClient["database"], 0, 0)
	redisClient = redis.NewClient(&redis.Options{
		Addr:     sysConfig.RedisClient["address"] + ":" + sysConfig.RedisClient["port"],
		Password: sysConfig.RedisClient["password"], // no password set
		DB:       int(redisDB),                      // use default DB
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	})

	if err != nil {
		panic(err)
	}

	mysqlDB, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		sysConfig.MysqlClient["user"], sysConfig.MysqlClient["password"],
		sysConfig.MysqlClient["address"], sysConfig.MysqlClient["port"], sysConfig.MysqlClient["database"]))

	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to the database: mysql @GORM")

	// Creating and Migrating tables from the structures
	mysqlDB.AutoMigrate(&entity.Feedback{})
	mysqlDB.AutoMigrate(&entity.Subscription{})
	mysqlDB.AutoMigrate(&entity.JobApplication{})
	mysqlDB.AutoMigrate(&entity.Job{})
	mysqlDB.AutoMigrate(&entity.User{})

	// ----- Bot level database -----
	mysqlDB.AutoMigrate(&bot.TempUser{})
	mysqlDB.AutoMigrate(&bot.Client{})

	// Setting foreign key constraint
	mysqlDB.Model(&entity.JobApplication{}).AddForeignKey("job_seeker_id", "users(id)", "CASCADE", "CASCADE")
	mysqlDB.Model(&entity.JobApplication{}).AddForeignKey("job_id", "jobs(id)", "CASCADE", "CASCADE")
	mysqlDB.Model(&entity.Feedback{}).AddForeignKey("user_id", "users(id)", "SET NULL", "CASCADE")
	mysqlDB.Model(&entity.Subscription{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")

	// ----- Bot level constraint -----
	mysqlDB.Model(&bot.Client{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
}

func main() {
	configFilesDir = "/asseri_bot/config"

	// Initializing the server
	initServer()
	defer mysqlDB.Close()

	router := mux.NewRouter()

	router.HandleFunc("/approval/result/{id}", botHandler.HandleApprovalResult).Methods("GET")
	router.HandleFunc("/push/notification/channel/{id}", botHandler.HandlePushNotificationToChannel).Methods("GET")
	router.HandleFunc("/push/notification/subscriber/{id}", botHandler.HandlePushNotificationToSubscribers).Methods("GET")

	router.HandleFunc("/", tools.MiddlewareFactory(botHandler.HandleWebHook, botHandler.ParseRequest))

	go func() {
		botHandler.HandlePushRequest()
	}()

	http.ListenAndServeTLS(":"+os.Getenv("bot_client_server_port"),
		filepath.Join(configFilesDir, "/server.pem"),
		filepath.Join(configFilesDir, "/server.key"), router)
}

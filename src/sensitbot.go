package main

import (
	// internal package
	"./controllers"
	"./db"
	//"./models"
	"./utils"

	// standard packages
	"fmt"
	//"io/ioutil"
	"log"
	"log/syslog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	// external packages
	"github.com/BurntSushi/toml"
	"github.com/julienschmidt/httprouter"
	"github.com/tucnak/telebot"
	"golang.org/x/oauth2"
)

// Program structures.
//  Define Start and Stop methods.
type program struct {
	exit chan struct{}
}

type tomlConfig struct {
	Title     string
	Owner     utils.OwnerInfo
	Sensitbot utils.Sensitbot
	Telegram  utils.Telegram
	Sensit    utils.Sensit
	Mongodb   utils.Mongodb
}

var ConfigFile string = "sensitbot.toml"
//var logger service.Logger
var config tomlConfig

var bot *telebot.Bot


func botMessages(uc *controllers.UserController, messages chan telebot.Message, bot *telebot.Bot) {

	for message := range messages {

		// a virer, enregistrer que les users qui font la demande d'oauth2 de sensit
		if !uc.MongoFindByUsername(message.Sender.Username) {
			uc.MongoCreateUser(&message.Sender)
		}

		//fmt.Println("uc username : ", uc.MongoGetUser())
		//uc.MongoGetUser()

		// telegram user's ID as oauth2 state
		state := strconv.Itoa(message.Sender.ID)
		println("service : ", message.IsService())
		println("personnal : ", message.IsPersonal())

		if len(message.Photo) != 0 {
			println("message photo")

		} else if message.Audio.FileSize != 0 {
			println("message audio")
		} else if message.Document.FileSize != 0 {
			println("message document")
		} else if message.Sticker.FileSize != 0 {
			println("message sticker")
		} else if message.Video.FileSize != 0 {
			println("message video")
		} else if message.Location.Latitude != 0 {
			println("message location")
		} else if message.Contact.UserID != 0 {
			println("message contact")
		} else if len(message.Text) != 0 {

			switch message.Text {

			case "/auth":
				url := utils.OAuth2_config.AuthCodeURL(state, oauth2.AccessTypeOffline)
				//redirect := "&redirect_uri=http://9b6268cd.ngrok.io/sensit/" + id + "/"
				//			msg := "Please authorize sigfoxbot on your sensit's account : " + "<a href=\"" + url + "\">authorize</a>"
				msg := "Please authorize sensitbot on your sensit.io account : " + `<a href="` + url + `">authorize</a>`

				//dest := telebot.User{ID: teleid_i}
				uc.UserSendMessage(message.Origin(), msg, nil)

				//url := utils.OauthConf.AuthCodeURL("state", oauth2.AccessTypeOffline)
				//bot.SendMessage(message.Chat,
				//	"Please authorize sigfoxbot on your sensit's account : "+url, nil)
			case "/start":
				msg := "/auth to connect your device (see /help)"
				uc.UserSendMessage(message.Origin(), msg, nil)

			case "/help":
				msg := `@sensitBot is a bot which allow to receive notifications from your sens'it device. <a href="https://www.sensit.io">sensit</a>`
				if uc.UserGetToken(state) == true {
					//controllers.UserSendMessage(bot, &message, msg)
					bot.SendMessage(message.Chat, msg, &telebot.SendOptions{
						ParseMode: "HTML",
						ReplyMarkup: telebot.ReplyMarkup{
							ForceReply: false,
							Selective:  true,

							CustomKeyboard: [][]string{
								[]string{"/devices", "/notifications"},
								[]string{"/version", "/account"},
							},
						},
					},
					)
				} else {
					msg += "\n/auth (to connect to your sens'it device)"
					uc.UserSendMessage(message.Origin(), msg, nil)
				}

			case "/back":
				msg := `I'm waiting your order Master.`
				if uc.UserGetToken(state) == true {
					//controllers.UserSendMessage(bot, &message, msg)
					bot.SendMessage(message.Chat, msg, &telebot.SendOptions{
						ParseMode: "HTML",
						ReplyMarkup: telebot.ReplyMarkup{
							ForceReply: false,
							Selective:  true,

							CustomKeyboard: [][]string{
								[]string{"/devices", "/notifications"},
								[]string{"/version", "/account"},
							},
						},
					},
					)
				} else {
					msg += "\n/auth (to connect to your sens'it device)"
					uc.UserSendMessage(message.Origin(), msg, nil)
				}

			case "/devices":
				arr, err := uc.GetDevices(state)
				//uc.UserSendMessage(message.Origin(), msg, nil)

				if len(err) != 0 {
					uc.UserSendMessage(message.Origin(), err, nil)
				} else {
					// do not count "/back" from len(arr)
					var nb int
					if len(arr) == 1 {
						nb = 0
					} else {
						nb = len(arr) - 1
					}
					msg := strconv.Itoa(nb) + " devices"
					uc.UserSendMessage(message.Origin(), msg, arr)
				}



			case "/account":
				msg := uc.GetAccount(state)
				uc.UserSendMessage(message.Origin(), msg, nil)

			case "/notifications":
				arr, err := uc.GetNotifications(state)
				//var arr [][]string = [][]string{[]string{"notification 15143", "notification 15187"}, []string{"notif 3", "notif 4"}, []string{"/back"}}

				//fmt.Println("arr : ", arr)

				if len(err) != 0 {
					uc.UserSendMessage(message.Origin(), err, nil)
				} else {
					// do not count "/back" from len(arr)
					var nb int
					if len(arr) == 1 {
						nb = 0
					} else {
						nb = len(arr) - 1
					}
					msg := strconv.Itoa(nb) + " notifications"
					uc.UserSendMessage(message.Origin(), msg, arr)
				}

			case "/version":
				msg := "@sensitBot <b>" + config.Sensitbot.Version + "</b>. Made by @fredix."
				uc.UserSendMessage(message.Origin(), msg, nil)

			case "/setcallback":
				url := uc.SetCallbackURL(state)
				msg := "callback url : " + url
				uc.UserSendMessage(message.Origin(), msg, nil)

			case "/image":
				pix, err := telebot.NewFile("golang.jpg")
				if err != nil {
					//return err
					fmt.Println("err : " + err.Error())
				}
				image := telebot.Photo{File: pix}

				// Next time you send &audio, telebot won't issue
				// an upload, but would re-use existing file.
				err = bot.SendPhoto(message.Chat, &image, nil)

			case "/ping":
				var arr [][]string = [][]string{[]string{"1", "2", "3"}}
				uc.UserSendMessage(message.Origin(), "pong", arr)

			case "/html":
				msg := `<b>bold</b>, <strong>bold</strong>
				<i>italic</i>, <em>italic</em>
				<a href="URL">inline URL</a>
				<code>inline fixed-width code</code>
				<pre>pre-formatted fixed-width code block</pre>`

				uc.UserSendMessage(message.Origin(), msg, nil)

			default:

				if strings.HasPrefix(message.Text, "notification") == true {

					if strings.Contains(message.Text, "delete") {
						msg_arr := strings.Split(message.Text, " ")
						if len(msg_arr) == 3 {
							id, err := strconv.Atoi(msg_arr[2])
							if err != nil {
								fmt.Println("ERROR on notification id : ", err)
								uc.UserSendMessage(message.Origin(), err.Error(), nil)
							} else {
								fmt.Println("notification id : ", id)
								msg := uc.DeleteNotificationByID(state, id)
								var arr [][]string = [][]string{[]string{"/back"}}
								uc.UserSendMessage(message.Origin(), msg, arr)
							}
						}
					} else {
						id, err := strconv.Atoi(strings.Replace(message.Text, "notification ", "", -1))
						if err != nil {
							fmt.Println("ERROR on notification id : ", err)
							uc.UserSendMessage(message.Origin(), err.Error(), nil)
						} else {
							fmt.Println("notification id : ", id)
							msg := uc.GetNotificationByID(state, id)

							var arr [][]string = [][]string{[]string{"notification delete " + strconv.Itoa(id), "/back"}}
							uc.UserSendMessage(message.Origin(), msg, arr)
						}
					}
				} else if strings.HasPrefix(message.Text, "device") == true {


					if strings.Contains(message.Text, "notification") {
						msg_arr := strings.Split(message.Text, " ")
						if len(msg_arr) == 4 {
							deviceid := msg_arr[2]
							sensorid := msg_arr[3]
							fmt.Println("device id : ", deviceid)
							fmt.Println("sensor id : ", sensorid)
							//msg := uc.DeleteNotificationByID(state, id)

							msg := uc.AddNotification(state, deviceid, sensorid, "plop")

							var arr [][]string = [][]string{[]string{"/back"}}
							uc.UserSendMessage(message.Origin(), msg, arr)
						}
					} else {
						id, err := strconv.Atoi(strings.Replace(message.Text, "device ", "", -1))
						if err != nil {
							fmt.Println("ERROR on device id : ", err)
							uc.UserSendMessage(message.Origin(), err.Error(), nil)
						} else {
							fmt.Println("device id : ", id)
							//msg := uc.GetDeviceByID(state, id)

							msg, arr, _ := uc.GetDeviceByID(state, id)

							//var arr [][]string = [][]string{[]string{"/back"}}
							//uc.UserSendMessage(message.Origin(), msg, arr)

							uc.UserSendMessage(message.Origin(), msg, arr)
						}
					}
					
				} else {
					msg := "/help"
					uc.UserSendMessage(message.Origin(), msg, nil)
				}
			}
		}
	}
}

func botQueries(uc *controllers.UserController, queries chan telebot.Query, bot *telebot.Bot) {
	for query := range queries {

		// a virer, enregistrer que les users qui font la demande d'oauth2 de sensit
		//if !uc.MongoFindByUsername(query.From.Username) {
		//	uc.MongoCreateUser(&query.From. message.Sender)
		//}

		//fmt.Println("uc username : ", uc.MongoGetUser())
		//uc.MongoGetUser()

		// telegram user's ID as oauth2 state
		state := strconv.Itoa(query.From.ID)

		var text string

		if uc.UserGetToken(state) == true {
			//controllers.UserSendMessage(bot, &message, msg)
			text = "AUTH OK"
		} else {
			text = "\n/auth (to connect to your sens'it device)"
		}

		fmt.Println("--- new query ---")
		fmt.Println("from:", query.From.Username)
		fmt.Println("text:", query.Text)

		// Create an article (a link) object to show in our results.
		article := &telebot.InlineQueryResultArticle{
			Title: "Choose a notification",
			//URL:   "https://github.com/tucnak/telebot",
			InputMessageContent: &telebot.InputTextMessageContent{
				Text:           text,
				DisablePreview: false,
			},
		}

		// Build the list of results. In this instance, just our 1 article from above.
		results := []telebot.InlineQueryResult{article}

		// Build a response object to answer the query.
		response := telebot.QueryResponse{
			Results:    results,
			IsPersonal: true,
		}

		// And finally send the response.
		if err := bot.AnswerInlineQuery(&query, &response); err != nil {
			fmt.Println("Failed to respond to query:", err)
		}
	}
}


func sigfoxBotApi(bot *telebot.Bot, r *httprouter.Router, uc *controllers.UserController) {

	/*

	   file, err := os.Open("file.go") // For read access.
	   if err != nil {
	       log.Fatal(err)
	   }
	   count, err := file.Read(data)
	   if err != nil {
	       log.Fatal(err)
	   }
	   fmt.Printf("read %d bytes: %q\n", count, data[:count])
	*/

	utils.Oauth2Config(config.Sensit)

	//uc := controllers.NewUserController(db.GetSession(config.Mongodb.Url), config.Mongodb.Database)

	// Get sensit code
	//r.GET("/user/:id", uc.GetUserById)
	//r.GET("/sensit/:id/*code", uc.GetSensitCode)
	r.GET("/sensit/*code", uc.GetSensitCode)
	r.POST("/sensit/", uc.PostSensitCallback)

	// Get a user resource
	//r.GET("/user/:id", uc.GetUserById)
	r.GET("/user/:nickname", uc.GetUserByNickname)
	r.GET("/token/:token", uc.GetUserByToken)

	// Create a new user
	r.POST("/user", uc.CreateUser)

	// Remove an existing user
	r.DELETE("/user/:id", uc.RemoveUser)

	// Fire up the server
	http.ListenAndServe(config.Sensitbot.Url+":"+strconv.Itoa(config.Sensitbot.Port), r)

}

func main() {
	logwriter, e := syslog.New(syslog.LOG_NOTICE, "sensitbot")
    if e == nil {
        log.SetOutput(logwriter)
    }


	if len(os.Args) > 1 {
		ConfigFile = os.Args[1]
	}


	// Do work here

	if _, err := toml.DecodeFile(ConfigFile, &config); err != nil {
		panic(fmt.Sprintf("%s", err))
	}

	fmt.Println("Title: ", config.Title)
	fmt.Printf("Owner: %s (%s, %s)\n",
		config.Owner.Name, config.Owner.Org, config.Owner.DOB)

	//fmt.Printf("sensit: (Oauth_url %s, Token_url %s Token %s), telegram : (Token %s)",
	//		config.Sensit.Oauth_url, config.Sensit.Token_url, config.Sensit.Token, config.Telegram.Token)

	// boltdb to store json payload
	// and maybe other stuffs
	//    boltDb := gputils.NewBoltDb(config.Boltdb.Dbfile, config.Graylog)
	//boltDb := utils.NewBoltDb(config.Boltdb.Dbfile)

	//    boltDb.Resend("graylog")

	//        done := make(chan bool)

	// Get a UserController instance

	// Instantiate a new router
	r := httprouter.New()

	var err error
	bot, err = telebot.NewBot(config.Telegram.Token)
	if err != nil {
		return
	}

	uc := controllers.NewUserController(bot, db.GetSession(config.Mongodb.Url, config.Mongodb.Database, config.Mongodb.Dbowner, config.Mongodb.Dbpass), config.Mongodb.Database, config.Sensitbot.CallbackUrl)

	go sigfoxBotApi(bot, r, uc)

	//var bot_user *telebot.Bot

	bot.Messages = make(chan telebot.Message)
	bot.Queries = make(chan telebot.Query)

	//	bot.Listen(messages, 1*time.Second)

	go botMessages(uc, bot.Messages, bot)
	go botQueries(uc, bot.Queries, bot)

	bot.Start(1 * time.Second)


}

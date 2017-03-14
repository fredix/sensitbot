package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	//"log"
	"bytes"
	"errors"
	"net/http"
	"reflect"
	"strconv"

	"../models"
	"../utils"

	"github.com/julienschmidt/httprouter"
	//"github.com/nu7hatch/gouuid"
	"github.com/tucnak/telebot"
	"golang.org/x/net/context"
	//"golang.org/x/oauth2"
	"golang.org/x/oauth2"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type (
	// UserController represents the controller for operating on the User resource
	UserController struct {
		bot         *telebot.Bot
		session     *mgo.Session
		db          string
		callbackurl string
		//models.User
	}
)

type ErrorData struct {
	Err     int    `json:"err"`
	Details string `json:"details"`
}

/*
100, "Error database access"
101, "Wrong Content-Type header : only accept application/json"
102, "Too many login/activation attempts"
300, "Device error"
400, "Account error"
500, "Notification error"
*/

/*
func UserSendNotifications(bot *telebot.Bot, message *telebot.Message, token string) {
	client := utils.OauthConf.Client(oauth2.NoContext, token)
	resp, err := client.Get("https://api.sensit.io/v2/notifications")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Reading the body
	raw, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println("raw body error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Unmarshalling the JSON of the Profile
	var notifications map[string]interface{}
	if err := json.Unmarshal(raw, &notifications); err != nil {
		fmt.Println("json Unmarshal error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("notifications : ", notifications)

	url := utils.OauthConf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	//message := "Visit the URL for the auth dialog: " + url
	bot.SendMessage(message.Chat,
		"Sensit's notifications : "+notifications, nil)

}
*/

// NewUserController provides a reference to a UserController with provided mongo session
func NewUserController(bot *telebot.Bot, s *mgo.Session, db string, callbackurl string) *UserController {
	return &UserController{bot, s, db, callbackurl}
}

/*func (uc UserController) MongoGetUser() {
	fmt.Println("plop")

	fmt.Println("MongoGetUser : ", uc.User.LastName)
	//return uc.mgo_user.Username
}*/

func (uc UserController) UserSendMessage(dest telebot.User, msg string, keyboard [][]string) {

	/*
		fmt.Println("len keyboard : ", len(keyboard))
		fmt.Println("keyboard : ", keyboard)
		fmt.Println("type keyboard : ", reflect.TypeOf(keyboard))
	*/

	if len(keyboard) == 0 {
		uc.bot.SendMessage(&dest, msg, &telebot.SendOptions{
			ParseMode: "HTML",
		},
		)
	} else {
		uc.bot.SendMessage(&dest, msg, &telebot.SendOptions{
			ParseMode: "HTML",
			ReplyMarkup: telebot.ReplyMarkup{
				ForceReply: true,
				Selective:  true,

				CustomKeyboard: keyboard,

				/*CustomKeyboard: [][]string{
					[]string{"notification 15143", "notification 15187"},
					[]string{"notif 3", "notif 4"},
					[]string{"/back"},
				},
				*/

			},
		},
		)

	}
}

func (uc UserController) MongoFindByUsername(username string) bool {

	fmt.Println(username)
	// Stub user
	//u := models.User{}

	u := models.User{}
	//	uc.User = models.User{}

	// Fetch user
	if err := uc.session.DB(uc.db).C("users").Find(bson.M{"username": username}).One(&u); err != nil {
		fmt.Println("user not found")
		return false
	}
	//uc.User = u
	//reflect.Copy(uc.User.Username, u.Username)
	fmt.Println("mgo username :", u.LastName)

	// Marshal provided interface into JSON structure
	//uj, _ := json.Marshal(u)

	// Write content-type, statuscode, payload
	//fmt.Println(uj)
	return true
}



func (uc UserController) MongoFindDevice(serial string) bool {

	fmt.Println(serial)
	// Stub user
	//u := models.User{}

	d := models.Device{}
	//	uc.User = models.User{}

	// Fetch user
	if err := uc.session.DB(uc.db).C("devices").Find(bson.M{"serial": serial}).One(&d); err != nil {
		fmt.Println("device not found")
		return false
	}
	//uc.User = u
	//reflect.Copy(uc.User.Username, u.Username)
	fmt.Println("mgo DeviceId :", d.DeviceId)

	// Marshal provided interface into JSON structure
	//uj, _ := json.Marshal(u)

	// Write content-type, statuscode, payload
	//fmt.Println(uj)
	return true
}



func (uc UserController) MongoGetByTeleId(tele_id string) (user models.User, err error) {

	fmt.Println("MongoGetByTeleId : ", tele_id)
	// Stub user
	//u := models.User{}

	id, _ := strconv.Atoi(tele_id)

	u := models.User{}
	//	uc.User = models.User{}

	// Fetch user
	if err := uc.session.DB(uc.db).C("users").Find(bson.M{"tele_id": id}).One(&u); err != nil {
		fmt.Println("MongoGetByTeleId error : ", err)
		if err.Error() == "EOF" {
			uc.session.Refresh()
		}
		err := errors.New("mongodb refresh session...")
		return u, err
	}
	//uc.User = u
	//reflect.Copy(uc.User.Username, u.Username)
	fmt.Println("MongoGetByTeleId username :", u.Username)
	fmt.Println("MongoGetByTeleId username ID : ", u.Id)
	fmt.Println("MongoGetByTeleId username accesstoken : ", u.Token.AccessToken)
	fmt.Println("MongoGetByTeleId username refreshtoken : ", u.Token.RefreshToken)
	// Marshal provided interface into JSON structure
	//uj, _ := json.Marshal(u)

	// Write content-type, statuscode, payload
	//fmt.Println(uj)
	return u, nil
}

func (uc UserController) MongoUpdateToken(id bson.ObjectId, token *oauth2.Token) bool {
	/*
		if !bson.IsObjectIdHex(id) {
			fmt.Println("ID error")
			return false
		}
	*/

	// Grab id
	//oid := bson.ObjectIdHex(id)

	fmt.Println("MongoUpdateToken : ", token)
	// Stub user
	//u := models.User{}

	u := models.User{}

	// Fetch user
	//	if err := uc.session.DB(uc.db).C("users").Find(bson.M{"tele_id": id}).One(&u); err != nil {
	if err := uc.session.DB(uc.db).C("users").Find(bson.M{"_id": id}).One(&u); err != nil {

		fmt.Println("MongoGetById : user not found")
		return false
	}

	c := uc.session.DB(uc.db).C("users")
	colQuerier := bson.M{"_id": id}
	change := bson.M{"$set": bson.M{"token": token}}

	err := c.Update(colQuerier, change)
	if err != nil {
		panic(err)
	}

	// Marshal provided interface into JSON structure
	//uj, _ := json.Marshal(u)

	// Write content-type, statuscode, payload
	//fmt.Println(uj)
	return true
}

// CreateUser creates a new user resource
func (uc UserController) MongoCreateUser(sender *telebot.User) {
	// Stub an user to be populated from the body

	fmt.Println("firstname : ", sender.FirstName)
	fmt.Println("LastName : ", sender.LastName)
	fmt.Println("Username : ", sender.Username)
	fmt.Println("ID : ", sender.ID)
	fmt.Println("Destination : ", sender.Destination())

	u := models.User{}

	// Add an Id
	u.Id = bson.NewObjectId()
	u.TeleId = sender.ID
	u.Username = sender.Username
	u.FirstName = sender.FirstName
	u.LastName = sender.LastName

	// Add an uuid
	/*
		u4, err := uuid.NewV4()
		if err != nil {
			log.Fatal(err)
		}
		u.Token = u4.String()
	*/

	// Write the user to mongo
	uc.session.DB(uc.db).C("users").Insert(u)

	// Marshal provided interface into JSON structure
	//uj, _ := json.Marshal(u)

}



// CreateUser creates a new user resource
func (uc UserController) MongoCreateDevice(userid bson.ObjectId, device models.DeviceData) {
	// Stub an user to be populated from the body

	fmt.Println("user Id : ", userid)
	fmt.Println("device : ", device)

	d := models.Device{}

/*
	Id         bson.ObjectId `json:"id" bson:"_id"`
		AccountId  bson.ObjectId `json:"account_id" bson:"account_id"`
		DeviceId   int           `json:"device_id" bson:"device_id"`
		Serial     string        `json:"serial_number" bson:"serial"`
		Model      int        	 `json:"device_model" bson:"model"`
		Battery    string        `json:"battery" bson:"battery"`
		Mode 	   int           `json:"mode" bson:"mode"`
		ActivationDate string `json:"activation_date" bson:"activation_date"`
		LastCommDate   string `json:"last_comm_date" bson:"modlast_comm_datee"`
		LastConfigDate string `json:"last_config_date" bson:"last_config_date"`
*/


	// Add an Id
	d.Id = bson.NewObjectId()
	d.AccountId = userid
	d.DeviceId = device.Id
	d.Serial = device.SerialNumber
	d.Model = device.DeviceModel
	d.Battery = device.Battery
	d.Mode = device.Mode
	d.ActivationDate = device.ActivationDate
	d.LastCommDate = device.LastCommDate
	d.LastConfigDate = device.LastConfigDate


	// Add an uuid
	/*
		u4, err := uuid.NewV4()
		if err != nil {
			log.Fatal(err)
		}
		u.Token = u4.String()
	*/

	// Write user's device
	uc.session.DB(uc.db).C("devices").Insert(d)

	fmt.Println("Create Device : ", d)

	// Marshal provided interface into JSON structure
	//uj, _ := json.Marshal(u)

}





func (uc UserController) GetSensitCode(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	//nickname := p.ByName("nickname")
	//fmt.Println("id : ", p.ByName("id"))

	state := r.FormValue("state")
	fmt.Println("state : ", state)

	user, err := uc.MongoGetByTeleId(state)
	if err != nil {
		fmt.Println("GetSensitCode : user not found")
		return
	}

	code := r.FormValue("code")
	fmt.Println("code : ", code)
	ctx := context.Background()

	// remplacer
	println("GetSensitCode context : ", ctx)

	/*
		token, err := utils.OAuth2_config.Exchange(ctx, code)
		if err != nil {
			fmt.Printf("OauthConf.Exchange() failed with '%s'\n", err)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		fmt.Println("token : ", token)
		fmt.Println("token AccessToken : ", token.AccessToken)
		fmt.Println("token TokenType : ", token.TokenType)
		fmt.Println("token RefreshToken : ", token.RefreshToken)
		fmt.Println("token Expiry : ", token.Expiry)
		fmt.Println("token type : ", reflect.TypeOf(token))
		uc.MongoUpdateToken(user.Id, token)
	*/
	// remplacer par
	client, _ := uc.UserOAuthClient(ctx, code, user)

	/*
		client := utils.OAuth2_config.Client(ctx, token)
		resp, err := client.Get("https://api.sensit.io/v2/account")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Reading the body
		raw, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			fmt.Println("raw body error")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Unmarshalling the JSON of the Profile
		var profile map[string]interface{}
		if err := json.Unmarshal(raw, &profile); err != nil {
			fmt.Println("json Unmarshal error")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	*/

	/*
	   		POST {endpoint}/callback
	   {
	     url: "https://myCallbackURL",
	     header_name: "Authorization",
	     header_value: "Bearer my_token_value",
	   }
	*/

	// update callback URL to user's sensit account
	//var jsonStr = []byte(`{url: "http://8cd069a5.eu.ngrok.io/sensit/", header_name: "TeleId", header_value: ` + state + `}`)
	var jsonStr = []byte(`{url: "` + uc.callbackurl + `", header_name: "TeleId", header_value: ` + state + `}`)

	//client := utils.OAuth2_config.Client(ctx, token)
	resp, err := client.Post("https://api.sensit.io/v2/callback", "application/json; charset=utf-8", bytes.NewBuffer(jsonStr))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Reading the body
	raw, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println("raw body error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Unmarshalling the JSON
	var callback_user map[string]interface{}
	if err := json.Unmarshal(raw, &callback_user); err != nil {
		fmt.Println("json Unmarshal error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("callback_user : ", callback_user)
	fmt.Println("account updated : ", user.TeleId)

	var redirection = `<html>
<head>
<title>redirect</title>
<meta http-equiv="refresh" content="1; url=https://telegram.me/sensitbot" />
</head>
<body>
<a href="https://telegram.me/sensitbot">redirect to telegram</a>
</body>
</html>`

	fmt.Fprintf(w, redirection)

	teleid_i, err := strconv.Atoi(state)
	dest := telebot.User{ID: teleid_i}
	res := "@sensitBot is now linked to your sens'it account ! /help \n"
	uc.bot.SendMessage(dest, res, nil)
}

// https://gist.github.com/jfcote87/89eca3032cd5f9705ba3
// UserOauthClient returns an oauth2 client for a specific user
func (uc UserController) UserOAuthClient(ctx context.Context, code string, user models.User) (client *http.Client, err error) {
	fmt.Println("UserOAuthClient")
	var userToken = &user.Token
	//if userToken, err = getCachedToken(userId); err != nil {
	fmt.Println("UserOAuthClient : ", userToken.AccessToken)

	if (len(userToken.AccessToken) == 0) || (len(code) != 0) {
	//if len(code) != 0 {
		// if token for user is not cached then go through oauth2 flow
		if userToken, err = uc.newUserToken(ctx, code, user); err != nil {
			return nil, err
		}
	}
	if !userToken.Valid() { // if user token is expired
		fmt.Println("UserOAuthClient user token is expired")
		userToken = &oauth2.Token{RefreshToken: userToken.RefreshToken}
	}

	return utils.OAuth2_config.Client(ctx, userToken), nil
}

func (uc UserController) newUserToken(ctx context.Context, code string, user models.User) (*oauth2.Token, error) {
	fmt.Println("newUserToken")
	if len(code) == 0 {
		err := errors.New("please /auth")
		return nil, err
	}

	token, err := utils.OAuth2_config.Exchange(ctx, code)
	if err != nil {
		fmt.Printf("OauthConf.Exchange() failed with '%s'\n", err)
		return nil, err
	}
	fmt.Println("token : ", token)
	fmt.Println("token AccessToken : ", token.AccessToken)
	fmt.Println("token TokenType : ", token.TokenType)
	fmt.Println("token RefreshToken : ", token.RefreshToken)
	fmt.Println("token Expiry : ", token.Expiry)
	fmt.Println("token type : ", reflect.TypeOf(token))

	uc.MongoUpdateToken(user.Id, token)

	return token, nil
}

func (uc UserController) SetCallbackURL(user_teleid string) string {

	var jsonStr = []byte(`{url: "` + uc.callbackurl + `", header_name: "TeleId", header_value: ` + user_teleid + `}`)

	var jsonres map[string]interface{}
	if err0 := json.Unmarshal(jsonStr, &jsonres); err0 != nil {
		fmt.Println("SetCallbackURL json Unmarshal error : ", err0)
	}
	fmt.Println("SetCallbackURL json : ", jsonres)

	user, err := uc.MongoGetByTeleId(user_teleid)
	if err != nil {
		fmt.Println("SetCallbackURL : user not found")
		err := errors.New("user not found")
		return err.Error()
	}

	ctx := context.Background()
	println("SetCallbackURL context : ", ctx)

	//client := utils.OAuth2_config.Client(ctx, &user.Token)
	client, err := uc.UserOAuthClient(ctx, "", user)
	if err != nil {
		return err.Error()
	}

	//client := utils.OAuth2_config.Client(ctx, token)
	resp, err := client.Post("https://api.sensit.io/v2/callback", "application/json; charset=utf-8", bytes.NewBuffer(jsonStr))
	if err != nil {
		return err.Error()
	}
	// Reading the body
	raw, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println("raw body error")
		return err.Error()
	}

	// Unmarshalling the JSON
	var callback_user map[string]interface{}
	if err := json.Unmarshal(raw, &callback_user); err != nil {
		fmt.Println("json Unmarshal error")
		return err.Error()
	}

	fmt.Println("callback_user : ", callback_user)
	fmt.Println("account updated : ", user.TeleId)

	return uc.callbackurl
}

func (uc UserController) AddNotification(user_teleid string, device_id string, sensor_id string, message string) string {

	//	var jsonStr = []byte(`{url: "` + uc.callbackurl + `", header_name: "TeleId", header_value: ` + user_teleid + `}`)

	var jsonStr = []byte(`{
  template: "` + message + `",
  trigger: {
    id_device: "` + device_id + `",
    "id_sensor": "` + sensor_id + `",
    type: "GENERIC_PUNCTUAL"
  },
  connector: {
    data: "` + uc.callbackurl + `",
    type: "callback"
  },
  mode: 3
}`)

	//	map[results:1 data:map[template:http://bc40648d.eu.ngrok.io/sensit/ id:15143 trigger:map[id_device:2056 value: id_sensor:7282 type:GENERIC_PUNCTUAL] data: connector:map[data: type:callback] mode:3] links:map[]]

	var jsonres map[string]interface{}
	if err0 := json.Unmarshal(jsonStr, &jsonres); err0 != nil {
		fmt.Println("AddNotification json Unmarshal error : ", err0)
	}
	fmt.Println("AddNotification json : ", jsonres)

	user, err := uc.MongoGetByTeleId(user_teleid)
	if err != nil {
		fmt.Println("AddNotification : user not found")
		err := errors.New("user not found")
		return err.Error()
	}

	ctx := context.Background()
	println("AddNotification context : ", ctx)

	//client := utils.OAuth2_config.Client(ctx, &user.Token)
	client, err := uc.UserOAuthClient(ctx, "", user)
	if err != nil {
		return err.Error()
	}

	request, err := http.NewRequest("PUT", "https://api.sensit.io/v2/notifications", bytes.NewBuffer(jsonStr))
	if err != nil {
		return err.Error()
	}

	resp, err := client.Do(request)
	if err != nil {
		fmt.Println("request error : ", err)
		return err.Error()
	} else {
		// Reading the body
		raw, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			fmt.Println("raw body error")
			return err.Error()
		}

		fmt.Println("AddNotification raw : ", string(raw))
		//body := []byte{raw}
		//fmt.Println(string(body))

		fmt.Println("The calculated length is:", len(string(raw)), "for the url:", "https://api.sensit.io/v2/notifications")
		fmt.Println("status code : ", resp.StatusCode)
		// 200 : OK
		// 400 : ERROR

		hdr := resp.Header
		for key, value := range hdr {
			fmt.Println("   ", key, ":", value)
		}

		// Unmarshalling the JSON
		//var callback_user map[string]interface{}

		if resp.StatusCode != 200 {
			put_error := ErrorData{}
			if err := json.Unmarshal(raw, &put_error); err != nil {
				fmt.Println("AddNotification json Unmarshal error : ", err)
				//return err.Error()
			}
			//out, _ := json.Marshal(put_error)
			fmt.Println("put_error : ", put_error.Details)
			return "Error " + strconv.Itoa(put_error.Err) + " : " + put_error.Details
		} else {
			var notification map[string]interface{}
			if err := json.Unmarshal(raw, &notification); err != nil {
				fmt.Println("AddNotification json Unmarshal error : ", err)
				//return err.Error()
			}
			fmt.Println("notification : ", notification)
		}

		// OK : notification :  map[results:1 data:map[template:http://bc40648d.eu.ngrok.io/sensit/ id:15120 trigger:map[id_device:2056 value: id_sensor:7282 type:GENERIC_PUNCTUAL] data: connector:map[data: type:callback] mode:3] links:map[]]
		// KO : notification :  map[details:Notification already present:  err:500]

	}
	//	fmt.Println("callback_user : ", callback_user)
	//	fmt.Println("account updated : ", user.TeleId)
	return "added notification"

	//	return uc.callbackurl
}

func (uc UserController) UserGetToken(user_teleid string) bool {
	user, err := uc.MongoGetByTeleId(user_teleid)
	if err != nil {
		fmt.Println("UserGetToken : user not found", err.Error())
		return false
	}
	if len(user.Token.AccessToken) == 0 {
		return false
	} else {
		return true
	}
}

func (uc UserController) GetDevices(user_teleid string) ([][]string, string) {

	//nickname := p.ByName("nickname")
	//fmt.Println("id : ", p.ByName("id"))

	user, err := uc.MongoGetByTeleId(user_teleid)
	if err != nil {
		fmt.Println("GetDevices : user not found")
		err := errors.New("user not found")
		return nil, err.Error()
	}

	ctx := context.Background()
	println("GetDevices context : ", ctx)

	//client := utils.OAuth2_config.Client(ctx, &user.Token)
	client, err := uc.UserOAuthClient(ctx, "", user)
	if err != nil {
				return nil, err.Error()
	}

	resp, err := client.Get("https://api.sensit.io/v2/devices")
	if err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, "failed"
	}
	// Reading the body
	raw, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println("raw body error")
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, "failed"
	}

	//	devices :  map[results:1 data:[map[last_comm_date:2016-09-16T07:19Z battery:78 last_config_date:2016-01-25T00:13Z serial_number:C6E06
	// mode:1 sensors:[] id:2056 device_model:2 activation_date:2015-08-25T19:41Z]] links:map[last:/v2/devices?page=1 first:/v2/devices?page=1]]

/*
	type DeviceData struct {
		Id             string `json:"id"`
		DeviceModel    string `json:"device_model"`
		SerialNumber   string `json:"serial_number"`
		Battery        int    `json:"battery"`
		LastConfigDate string `json:"last_config_date"`
		Mode           int    `json:"mode"`
		ActivationDate string `json:"activation_date"`
		LastCommDate   string `json:"last_comm_date"`
	}

	type Devices struct {
		Results int          `json:"results"`
		Data    []DeviceData `json:"data"`
	}
*/

	// Unmarshalling the JSON of the Profile
	//var devices map[string]interface{}
	devices := models.Devices{}

	if err := json.Unmarshal(raw, &devices); err != nil {
		fmt.Println("json Unmarshal error : ", err)
		//http.Error(w, err.Error(), http.StatusInternalServerError)
			return nil, "json unmarshal failed : " + err.Error()
	}

	fmt.Println("TypeOf : ", reflect.TypeOf(devices))
	fmt.Println("devices : ", devices)

	// devices :  map[results:1 data:[map[id:2056 device_model:2 battery:78 last_config_date:2016-01-25T00:13Z sensors:[]
	//	activation_date:2015-08-25T19:41Z last_comm_date:2016-09-24T10:09Z serial_number:C6E06 mode:1]]
	//	links:map[first:/v2/devices?page=1 last:/v2/devices?page=1]]

	var res2 string
	for key, _ := range devices.Data {
		fmt.Println("data : ", devices.Data[key].SerialNumber)

		res2 = "nb devices : <b>" + strconv.Itoa(devices.Results) + "</b>\n"
		res2 += "Id : <b>" + devices.Data[key].Id + "</b>\n"
		res2 += "Serial number : <b>" + devices.Data[key].SerialNumber + "</b>\n"
		res2 += "Device model : <b>" + devices.Data[key].DeviceModel + "</b>\n"
		res2 += "Battery : <b>" + strconv.Itoa(devices.Data[key].Battery) + "</b>\n"
		res2 += "Mode : <b>" + strconv.Itoa(devices.Data[key].Mode) + "</b>\n"
		res2 += "Activation date : <b>" + devices.Data[key].ActivationDate + "</b>\n"
		res2 += "Last config date : <b>" + devices.Data[key].LastConfigDate + "</b>\n"
		res2 += "Last comm date : <b>" + devices.Data[key].LastCommDate + "</b>\n"



		if !uc.MongoFindDevice(devices.Data[key].SerialNumber) {
			fmt.Println("Device not found, serial : ", devices.Data[key].SerialNumber)
			uc.MongoCreateDevice(user.Id, devices.Data[key])
		}

	}


	fmt.Println("devices : ", devices)

	var res [][]string

	//var arr [][]string = [][]string{[]string{"notification 15143", "notification 15187"}, []string{"notif 3", "notif 4"}, []string{"/back"}}

	for key, _ := range devices.Data {
		res = append(res, []string{"device " + devices.Data[key].Id})
	}
	res = append(res, []string{"/back"})

	return res, ""
}




func (uc UserController) GetDeviceByID(user_teleid string, device_id int)  (string, [][]string, string)  {

	//nickname := p.ByName("nickname")
	//fmt.Println("id : ", p.ByName("id"))

	user, err := uc.MongoGetByTeleId(user_teleid)
	if err != nil {
		fmt.Println("GetDevices : user not found")
		err := errors.New("user not found")
		return "failed", nil, err.Error()
	}

	ctx := context.Background()
	println("GetDevices context : ", ctx)

	//client := utils.OAuth2_config.Client(ctx, &user.Token)
	client, err := uc.UserOAuthClient(ctx, "", user)
	if err != nil {
		return "failed", nil, err.Error()
	}

	id := strconv.Itoa(device_id)
	resp, err := client.Get("https://api.sensit.io/v2/devices/" + id)
	if err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
			return "failed", nil, err.Error()
	}
	// Reading the body
	raw, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println("raw body error")
		//http.Error(w, err.Error(), http.StatusInternalServerError)
			return "failed", nil, err.Error()
	}


	//	devices :  map[results:1 data:[map[last_comm_date:2016-09-16T07:19Z battery:78 last_config_date:2016-01-25T00:13Z serial_number:C6E06
	// mode:1 sensors:[] id:2056 device_model:2 activation_date:2015-08-25T19:41Z]] links:map[last:/v2/devices?page=1 first:/v2/devices?page=1]]

/*
	type DeviceData struct {
		Id             string `json:"id"`
		DeviceModel    string `json:"device_model"`
		SerialNumber   string `json:"serial_number"`
		Battery        int    `json:"battery"`
		LastConfigDate string `json:"last_config_date"`
		Mode           int    `json:"mode"`
		ActivationDate string `json:"activation_date"`
		LastCommDate   string `json:"last_comm_date"`
	}

	type Devices struct {
		Results int          `json:"results"`
		Data    []DeviceData `json:"data"`
	}
*/

	// Unmarshalling the JSON of the Profile
	//var devices map[string]interface{}
	devices := models.DevicesSensors{}

	if err := json.Unmarshal(raw, &devices); err != nil {
		fmt.Println("json Unmarshal error : ", err)
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Device raw : ", string(raw))
		return "failed", nil, err.Error()
	}

	fmt.Println("TypeOf : ", reflect.TypeOf(devices))
	fmt.Println("devices : ", devices)

	// devices :  map[results:1 data:[map[id:2056 device_model:2 battery:78 last_config_date:2016-01-25T00:13Z sensors:[]
	//	activation_date:2015-08-25T19:41Z last_comm_date:2016-09-24T10:09Z serial_number:C6E06 mode:1]]
	//	links:map[first:/v2/devices?page=1 last:/v2/devices?page=1]]

	var msg string
	var res [][]string

	msg += "Id : <b>" + devices.Data.Id + "</b>\n"
	msg += "Serial number : <b>" + devices.Data.SerialNumber + "</b>\n"
	msg += "Device model : <b>" + devices.Data.DeviceModel + "</b>\n"
	msg += "Battery : <b>" + strconv.Itoa(devices.Data.Battery) + "</b>\n"
	msg += "Mode : <b>" + strconv.Itoa(devices.Data.Mode) + "</b>\n"
	msg += "Activation date : <b>" + devices.Data.ActivationDate + "</b>\n"
	msg += "Last config date : <b>" + devices.Data.LastConfigDate + "</b>\n"
	msg += "Last comm date : <b>" + devices.Data.LastCommDate + "</b>\n"


	for key, _ := range devices.Data.SensorData {
		msg += "Sensor ID : <b>" + devices.Data.SensorData[key].Id + "</b>\n"
		msg += "SensorType : <b>" + devices.Data.SensorData[key].SensorType + "</b>\n"

		res = append(res, []string{"device notification " + devices.Data.Id + " " + devices.Data.SensorData[key].Id})
	}

	res = append(res, []string{"/back"})


	//var arr [][]string = [][]string{[]string{"notification 15143", "notification 15187"}, []string{"notif 3", "notif 4"}, []string{"/back"}}

/*	for key, _ := range notifications.Data {
		res = append(res, []string{"notification " + notifications.Data[key].Id})
	}
	res = append(res, []string{"/back"})
*/
	return msg, res, ""

}


func (uc UserController) GetAccount(user_teleid string) string {

	//nickname := p.ByName("nickname")
	//fmt.Println("id : ", p.ByName("id"))

	user, err := uc.MongoGetByTeleId(user_teleid)
	if err != nil {
		fmt.Println("GetAccount : user not found")
		err := errors.New("user not found")
		return err.Error()
	}

	ctx := context.Background()
	println("GetAccount context : ", ctx)

	//client := utils.OAuth2_config.Client(ctx, &user.Token)
	client, err := uc.UserOAuthClient(ctx, "", user)
	if err != nil {
		return err.Error()
	}

	resp, err := client.Get("https://api.sensit.io/v2/account")
	if err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		return "failed"
	}
	// Reading the body
	raw, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println("raw body error")
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		return "failed"
	}

	//	devices :  map[results:1 data:[map[last_comm_date:2016-09-16T07:19Z battery:78 last_config_date:2016-01-25T00:13Z serial_number:C6E06
	// mode:1 sensors:[] id:2056 device_model:2 activation_date:2015-08-25T19:41Z]] links:map[last:/v2/devices?page=1 first:/v2/devices?page=1]]

	type AccountData struct {
		Civility     string `json:"civility"`
		FirstName    string `json:"first_name"`
		LastName     string `json:"last_name"`
		Email        string `json:"email"`
		Phone        string `json:"phone"`
		Company      string `json:"company"`
		Country      string `json:"country"`
		ZipCode      string `json:"zip_code"`
		Address      string `json:"address"`
		AddressExtra string `json:"address_extra"`
		Profile      string `json:"profile"`
		Lang         string `json:"lang"`
		Town         string `json:"town"`
		PhoneIdZone  string `json:"phone_id_zone"`
		Id           string `json:"id"`
		CreationDate string `json:"creation_date"`
	}

	type Account struct {
		Results int         `json:"results"`
		Data    AccountData `json:"data"`
	}

	// Unmarshalling the JSON of the Profile
	//var devices map[string]interface{}
	account := Account{}
	//var account2 map[string]interface{}

	if err := json.Unmarshal(raw, &account); err != nil {
		fmt.Println("json Unmarshal error : ", err)
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		return "failed"
	}

	fmt.Println("TypeOf : ", reflect.TypeOf(account))

	fmt.Println("account : ", account)

	fmt.Println("Email : ", account.Data.Email)

	res := "Civility : <b>" + account.Data.Civility + "</b>\n"
	res += "First name : <b>" + account.Data.FirstName + "</b>\n"
	res += "Last name : <b>" + account.Data.LastName + "</b>\n"
	res += "Email : <b>" + account.Data.Email + "</b>\n"
	res += "Phone : <b>" + account.Data.Phone + "</b>\n"
	res += "Company : <b>" + account.Data.Company + "</b>\n"
	res += "Country : <b>" + account.Data.Country + "</b>\n"
	res += "Zip code : <b>" + account.Data.ZipCode + "</b>\n"
	res += "Address : <b>" + account.Data.Address + "</b>\n"
	res += "Address extra : <b>" + account.Data.AddressExtra + "</b>\n"
	res += "Profile : <b>" + account.Data.Profile + "</b>\n"
	res += "Lang : <b>" + account.Data.Lang + "</b>\n"
	res += "Town : <b>" + account.Data.Town + "</b>\n"
	res += "Phone id zone : <b>" + account.Data.PhoneIdZone + "</b>\n"
	res += "Id : <b>" + account.Data.Id + "</b>\n"
	res += "Creation date : <b>" + account.Data.CreationDate + "</b>\n"

	return res
}

func (uc UserController) GetNotifications(user_teleid string) ([][]string, string) {

	//nickname := p.ByName("nickname")
	//fmt.Println("id : ", p.ByName("id"))

	user, err := uc.MongoGetByTeleId(user_teleid)
	if err != nil {
		fmt.Println("GetNotifications : user not found")
		//err := errors.New("user not found")
		return nil, err.Error()
	}

	ctx := context.Background()
	println("GetNotifications context : ", ctx)

	//client := utils.OAuth2_config.Client(ctx, &user.Token)
	client, err := uc.UserOAuthClient(ctx, "", user)
	if err != nil {
		return nil, err.Error()
	}

	resp, err := client.Get("https://api.sensit.io/v2/notifications")
	if err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, "failed"
	}
	// Reading the body
	raw, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println("raw body error")
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, "failed"
	}

	type TriggerData struct {
		IdDevice string `json:"id_device"`
		IdSensor string `json:"id_sensor"`
		Value    string `json:"value"`
		Type     string `json:"type"`
	}

	type ConnectorData struct {
		Data string `json:"data"`
		Type string `json:"type"`
	}

	type NotificationData struct {
		Id        string        `json:"id"`
		Template  string 		`json:"template"`
		Data      string        `json:"data"`
		Mode      int           `json:"mode"`
		Trigger   TriggerData   `json:"trigger"`
		Connector ConnectorData `json:"connector"`
	}

	type Notification struct {
		Results int                `json:"results"`
		Data    []NotificationData `json:"data"`
	}

	// Unmarshalling the JSON of the Profile
	//var notifications map[string]interface{}
	notifications := Notification{}
	//var account2 map[string]interface{}

	if err := json.Unmarshal(raw, &notifications); err != nil {
		fmt.Println("json Unmarshal error : ", err)
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, "json Unmarshal error : " + err.Error()
	}

	fmt.Println("notification : ", notifications)

	fmt.Println("notification raw : ", string(raw))


	var res [][]string

	//var arr [][]string = [][]string{[]string{"notification 15143", "notification 15187"}, []string{"notif 3", "notif 4"}, []string{"/back"}}

	for key, _ := range notifications.Data {
		res = append(res, []string{"notification " + notifications.Data[key].Id})
	}
	res = append(res, []string{"/back"})

	return res, ""
}

func (uc UserController) GetNotificationByID(user_teleid string, notification_id int) string {

	//nickname := p.ByName("nickname")
	//fmt.Println("id : ", p.ByName("id"))

	user, err := uc.MongoGetByTeleId(user_teleid)
	if err != nil {
		fmt.Println("GetNotifications : user not found")
		//err := errors.New("user not found")
		return err.Error()
	}

	ctx := context.Background()
	println("GetNotifications context : ", ctx)

	//client := utils.OAuth2_config.Client(ctx, &user.Token)
	client, err := uc.UserOAuthClient(ctx, "", user)
	if err != nil {
		return err.Error()
	}
	id := strconv.Itoa(notification_id)

	resp, err := client.Get("https://api.sensit.io/v2/notifications/" + id)
	if err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		return "failed"
	}
	// Reading the body
	raw, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println("raw body error")
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		return "failed"
	}

	type TriggerData struct {
		IdDevice string `json:"id_device"`
		IdSensor string `json:"id_sensor"`
		Value    string `json:"value"`
		Type     string `json:"type"`
	}

	type ConnectorData struct {
		Data string `json:"data"`
		Type string `json:"type"`
	}

	type NotificationData struct {
		Id        string        `json:"id"`
		Data      string        `json:"data"`
		Mode      int           `json:"mode"`
		Trigger   TriggerData   `json:"trigger"`
		Connector ConnectorData `json:"connector"`
	}

	type Notification struct {
		Results int              `json:"results"`
		Data    NotificationData `json:"data"`
	}

	// Unmarshalling the JSON of the Profile
	//var notifications map[string]interface{}
	notifications := Notification{}
	//var account2 map[string]interface{}

	if err := json.Unmarshal(raw, &notifications); err != nil {
		fmt.Println("json Unmarshal error : ", err)
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		return "json Unmarshal error : " + err.Error()
	}

	fmt.Println("notification : ", notifications)

	var res string
	res += "notifications Id : <b>" + notifications.Data.Id + "</b>\n"
	res += "notifications Data : <b>" + notifications.Data.Data + "</b>\n"
	res += "notifications Mode : <b>" + strconv.Itoa(notifications.Data.Mode) + "</b>\n"
	res += "notifications Trigger IdDevice : <b>" + notifications.Data.Trigger.IdDevice + "</b>\n"
	res += "notifications Trigger IdSensor : <b>" + notifications.Data.Trigger.IdSensor + "</b>\n"
	res += "notifications Trigger Value : <b>" + notifications.Data.Trigger.Value + "</b>\n"
	res += "notifications Trigger Type : <b>" + notifications.Data.Trigger.Type + "</b>\n"
	res += "notifications Connector Data : <b>" + notifications.Data.Connector.Data + "</b>\n"
	res += "notifications Connector Type : <b>" + notifications.Data.Connector.Type + "</b>\n\n"

	return res
}

func (uc UserController) DeleteNotificationByID(user_teleid string, notification_id int) string {

	user, err := uc.MongoGetByTeleId(user_teleid)
	if err != nil {
		fmt.Println("DeleteNotificationByID : user not found")
		//err := errors.New("user not found")
		return err.Error()
	}

	ctx := context.Background()
	println("DeleteNotificationByID context : ", ctx)

	//client := utils.OAuth2_config.Client(ctx, &user.Token)
	client, err := uc.UserOAuthClient(ctx, "", user)
	if err != nil {
		return err.Error()
	}
	id := strconv.Itoa(notification_id)

	//	var jsonStr = []byte(`{url: "` + uc.callbackurl + `", header_name: "TeleId", header_value: ` + state + `}`)
	//	request, err := http.NewRequest("DELETE", "https://api.sensit.io/v2/notifications/"+id, bytes.NewBuffer(jsonStr))
	request, err := http.NewRequest("DELETE", "https://api.sensit.io/v2/notifications/"+id, nil)

	if err != nil {
		return err.Error()
	}

	resp, err := client.Do(request)
	if err != nil {
		fmt.Println("request error : ", err)
		return err.Error()
	} else {
		// Reading the body
		raw, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			fmt.Println("raw body error")
			return err.Error()
		}

		fmt.Println("DeleteNotificationByID raw : ", raw)

		// Unmarshalling the JSON of the Profile
		var notifications map[string]interface{}

		if err := json.Unmarshal(raw, &notifications); err != nil {
			fmt.Println("json Unmarshal error : ", err)
			//http.Error(w, err.Error(), http.StatusInternalServerError)
			return "json Unmarshal error : " + err.Error()
		}

		fmt.Println("notification : ", notifications)
		return "notification deleted"
	}

}

func (uc UserController) PostSensitCallback(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	//nickname := p.ByName("nickname")
	//fmt.Println("id : ", p.ByName("id"))

	// callback :  map[id:2056 activation_date:2015-08-25T19:41Z last_comm_date:2016-09-17T08:35Z last_config_date:0002-11-29T23:00Z
	//	serial_number:C6E06 mode:1 sensors:[map[history:[map[signal_level:average data: date:2016-09-17T08:35Z]] id:7282 sensor_type:button
	//	config:map[]] map[id:7280 sensor_type:temperature config:map[period:0] history:[map[signal_level:average data:25.0 date:2016-09-17T08:35Z]]]]
	//	device_model: battery:81]

	teleid := r.Header.Get("Teleid")
	fmt.Println("teleid : ", teleid)
	user, err := uc.MongoGetByTeleId(teleid)
	if err != nil {
		fmt.Println("PostSensitCallback : user not found")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("PostSensitCallback username : ", user.Username)

	// Reading the body
	raw, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		fmt.Println("raw body error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type SensorHistoryData struct {
		SignalLevel string `json:"signal_level"`
		Data        string `json:"data"`
		Date        string `json:"date"`
		DatePeriod  string `json:"date_period"`
	}

	type SensorsData struct {
		Id            string              `json:"id"`
		SensorType    string              `json:"sensor_type"`
		SensorHistory []SensorHistoryData `json:"history"`
	}

	type Callback struct {
		Id             string        `json:"id"`
		SerialNumber   string        `json:"serial_number"`
		Mode           int           `json:"mode"`
		ActivationDate string        `json:"activation_date"`
		LastCommDate   string        `json:"last_comm_date"`
		LastConfigDate string        `json:"last_config_date"`
		Battery        int           `json:"battery"`
		Sensors        []SensorsData `json:"sensors"`
	}

	// Unmarshalling the JSON of the Callback
	//var callback map[string]interface{}

	callback := Callback{}

	if err := json.Unmarshal(raw, &callback); err != nil {
		fmt.Println("json Unmarshal error : ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("callback : ", callback)

	teleid_i, err := strconv.Atoi(teleid)
	dest := telebot.User{ID: teleid_i}

	//dest := telebot.Chat{ID: teleid, Type: "private"}
	res := "Received a callback alert ! \n"
	res += "battery : <b>" + strconv.Itoa(callback.Battery) + "</b>\n"
	res += "mode : <b>" + strconv.Itoa(callback.Mode) + "</b>\n"

	for key, _ := range callback.Sensors {
		res += "Sensor type : <b>" + callback.Sensors[key].SensorType + "</b>\n"
		res += "Sensor id : <b>" + callback.Sensors[key].Id + "</b>\n"

		for key2, _ := range callback.Sensors[key].SensorHistory {
			res += "Sensor history signal level : <b>" + callback.Sensors[key].SensorHistory[key2].SignalLevel + "</b>\n"
			res += "Sensor history data : <b>" + callback.Sensors[key].SensorHistory[key2].Data + "</b>\n"
			res += "Sensor history date : <b>" + callback.Sensors[key].SensorHistory[key2].Date + "</b>\n"
			res += "Sensor history date period : <b>" + callback.Sensors[key].SensorHistory[key2].DatePeriod + "</b>\n\n"
		}
	}

	uc.UserSendMessage(dest, res, nil)
	/*	uc.bot.SendMessage(dest, res, &telebot.SendOptions{
			ParseMode: "HTML",
		},
		)
	*/
}

// GetUser retrieves an individual user resource
func (uc UserController) GetUserById(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Grab id
	id := p.ByName("id")

	// Verify id is ObjectId, otherwise bail
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(404)
		return
	}

	// Grab id
	oid := bson.ObjectIdHex(id)

	// Stub user
	u := models.User{}

	// Fetch user
	//uc.session.DB(name).Login(user, pass).C
	//uc.session.Login(cred)
	if err := uc.session.DB(uc.db).C("users").FindId(oid).One(&u); err != nil {
		//if err := uc.session.DB.C("users").FindId(oid).One(&u); err != nil {

		w.WriteHeader(404)
		return
	}

	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(u)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", uj)
}

// GetUser retrieves an individual user resource
func (uc UserController) GetUserByNickname(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Grab id
	nickname := p.ByName("nickname")

	// Stub user
	u := models.User{}

	// Fetch user
	if err := uc.session.DB(uc.db).C("users").Find(bson.M{"nickname": nickname}).One(&u); err != nil {
		w.WriteHeader(404)
		return
	}

	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(u)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", uj)
}

// GetUser retrieves an individual user resource
func (uc UserController) GetUserByToken(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Grab id
	token := p.ByName("token")

	// Stub user
	u := models.User{}

	// Fetch user
	if err := uc.session.DB(uc.db).C("users").Find(bson.M{"token": token}).One(&u); err != nil {
		w.WriteHeader(404)
		return
	}

	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(u)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", uj)
}

// CreateUser creates a new user resource
func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Stub an user to be populated from the body
	u := models.User{}

	// Populate the user data
	json.NewDecoder(r.Body).Decode(&u)

	// Add an Id
	u.Id = bson.NewObjectId()

	// Add an uuid
	/*u4, err := uuid.NewV4()
	if err != nil {
		log.Fatal(err)
	}
	u.Token = u4.String()
	*/

	// Write the user to mongo
	uc.session.DB(uc.db).C("users").Insert(u)

	// Marshal provided interface into JSON structure
	//uj, _ := json.Marshal(u)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", u.Token)
}

// RemoveUser removes an existing user resource
func (uc UserController) RemoveUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Grab id
	id := p.ByName("id")

	// Verify id is ObjectId, otherwise bail
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(404)
		return
	}

	// Grab id
	oid := bson.ObjectIdHex(id)

	// Remove user
	if err := uc.session.DB(uc.db).C("users").RemoveId(oid); err != nil {
		w.WriteHeader(404)
		return
	}

	// Write status
	w.WriteHeader(200)
}

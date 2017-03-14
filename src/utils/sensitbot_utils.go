package utils

import (
	// standard packages
	//"strconv"
	"fmt"
	"io"
	//"log"
	"os"
	"time"
	//"bytes"
	//"encoding/json"
	//"fmt"
	//"strings"
	"encoding/binary"
	//"io/ioutil"
	"log"
	"net/http"
	//"path/filepath"
	//"time"

	// external packages
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	//"gopkg.in/mgo.v2"
	//"gopkg.in/fsnotify.v1"
	//"github.com/boltdb/bolt"
	//".././models"
)

/*
type Message struct {
	ID              int
	Version         string
	Message  	string `json:"short_message"`
	Host     	string `json:"host"`
	Level 	        int    `json:"level"`
	MessageLog 	string `json:"_log"`
	File     	string `json:"_file"`
	ArchiveDir     	string `json:"_archivedir"`
	Localtime	string `json:"_localtime"`
}


type BoltDb struct {
	Dbfile string `toml:"dbfile"`
	DB         *bolt.DB
	writerChan chan [3]interface{}
	//graylog Graylog
}
*/

type OwnerInfo struct {
	Name string
	Org  string `toml:"organization"`
	DOB  time.Time
}

type Sensit struct {
	Api_url       string `toml:"api_url"`
	Oauth_url     string `toml:"oauth_url"`
	Token_url     string `toml:"token_url"`
	Client_id     string `toml:"client_id"`
	Client_secret string `toml:"client_secret"`
}

type Telegram struct {
	Token string `toml:"token"`
}

type Sensitbot struct {
	Url         string `toml:"url"`
	CallbackUrl string `toml:"callback_url"`
	Port        int    `toml:"port"`
	Version     string `toml:"version"`
}

type Mongodb struct {
	Url      string `toml:"url"`
	Database string `toml:"database"`
	Dbowner  string `toml:"dbowner"`
	Dbpass   string `toml:"dbpass"`
}

// Config is a custom oauth2 config that is used to store the token
// in the datastore
type Config struct {
	*oauth2.Config
}

type DatastoreTokenSource struct {
	config *Config
	source oauth2.TokenSource
	ctx    context.Context
}

/*
var (
	OauthConfOld = &oauth2.Config{
		ClientID:     "PRODLoJHuhsaGLM",
		ClientSecret: "qn2kIGYNxDXWJgQSNsv9NKkghQVYfQRI",
		Scopes:       []string{"read", "write"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.sensit.io/oauth/authorize",
			TokenURL: "https://api.sensit.io/oauth/token",
		},
	}
	// random string for oauth2 API calls to protect against CSRF
	OauthStateString = "prostiprosta"
)
*/

/*
type Config struct {
	// ClientID is the application's ID.
	ClientID string

	// ClientSecret is the application's secret.
	ClientSecret string

	// Endpoint contains the resource server's token endpoint
	// URLs. These are constants specific to each server and are
	// often available via site-specific packages, such as
	// google.Endpoint or github.Endpoint.
	Endpoint oauth2.Endpoint

	// RedirectURL is the URL to redirect users going through
	// the OAuth flow, after the resource owner's URLs.
	RedirectURL string

	// Scope specifies optional requested permissions.
	Scopes []string
}
*/

var OAuth2_config Config

/*
type Config struct {
	*oauth2.Config
	session *mgo.Session
	db      string

	//Mongo MyMongoInterface
}

https://gist.github.com/patrick91/f1d725a985b552261448
https://gist.github.com/agtorre/350c5b4ce0ccebc5ac0f
*/

// StoreToken is called when exchanging the token and saves the token
// in the datastore
func (c *Config) StoreToken(ctx context.Context, token *oauth2.Token) error {
	//log.Infof(ctx, "storing the token")
	log.Println("storing the token")

	// STORE MONGO !!!!!
	//key := datastore.NewKey(ctx, "Tokens", "fitbit", 0, nil)
	//_, err := datastore.Put(ctx, key, token)

	return nil
}

// Exchange is a wrapper around oauth2.config.Exchange and stores the Token
// in the datastore
func (c *Config) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := c.Config.Exchange(ctx, code)

	println("Exchange : ", code)

	if err != nil {
		return nil, err
	}

	/*
		if err := c.StoreToken(ctx, token); err != nil {
			return nil, err
		}
	*/
	return token, nil
}

// Client creates a new client using our custom TokenSource
func (c *Config) Client(ctx context.Context, t *oauth2.Token) *http.Client {
	return oauth2.NewClient(ctx, c.TokenSource(ctx, t))
}

/*


func (c *Config) StoreToken(token *oauth2.Token) error {
	u := models.User{}

	if err := c.session.DB(uc.db).C("users").Find(bson.M{"username": username}).One(&u); err != nil {
		fmt.Println("user not found")
		return false
	}
	//store the token in redis here using c.Redis
}

func (c *Config) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := c.Config.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}
	if err := c.StoreToken(token); err != nil {
		return nil, err
	}
	return token, nil
}

func (c *Config) Client(ctx context.Context, t *Token) *http.Client {
	return oauth2.NewClient(ctx, c.TokenSource(ctx, t))
}

func (c *Config) TokenSource(ctx context.Context, t *oauth2.Token) oauth2.TokenSource {
	rts := &MongoTokenSource{
		source: c.Config.TokenSource(ctx, t),
		config: c,
	}
	return oauth2.ReuseTokenSource(t, rts)
}

type MongoTokenSource struct {
	source oauth2.TokenSource
	config *Config
}

func (t *MongoTokenSource) Token() (*oauth2.Token, error) {
	token, err := t.source.Token()
	if err != nil {
		return nil, err
	}
	if err := t.config.StoreToken(token); err != nil {
		return nil, err
	}
	return token, nil
}*/

// TokenSource uses uses our DatastoreTokenSource as the source token
func (c *Config) TokenSource(ctx context.Context, t *oauth2.Token) oauth2.TokenSource {
	rts := &DatastoreTokenSource{
		source: c.Config.TokenSource(ctx, t),
		config: c,
		ctx:    ctx,
	}

	return oauth2.ReuseTokenSource(t, rts)
}

// Token saves the token in the datastore when it is updated
func (t *DatastoreTokenSource) Token() (*oauth2.Token, error) {
	token, err := t.source.Token()

	if err != nil {
		return nil, err
	}

	if err := t.config.StoreToken(t.ctx, token); err != nil {
		return nil, err
	}

	return token, nil
}

func getToken(ctx context.Context, data *oauth2.Token) error {
	//key := datastore.NewKey(ctx, "Tokens", "fitbit", 0, nil)
	//err := datastore.Get(ctx, key, data)

	// GET FROM MONGO !!!
	return nil
}

/*
func getFitbitConf(ctx context.Context) (*Config, error) {
	var data Settings

	//	err := GetSettings(ctx, &data)

	if err != nil {
		return nil, err
	}

	var fitbitConf = &Config{
		Config: &oauth2.Config{
			ClientID:     data.FitbitClientID,
			ClientSecret: data.FitbitClientSecret,
			Scopes:       []string{"activity", "weight", "profile"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://www.fitbit.com/oauth2/authorize",
				TokenURL: "https://api.fitbit.com/oauth2/token",
			},
		},
	}

	return fitbitConf, nil
}
*/

func Oauth2Config(confsens Sensit) {
	OAuth2_config.Config = &oauth2.Config{
		ClientID:     confsens.Client_id,
		ClientSecret: confsens.Client_secret,
		Scopes:       []string{"read", "write"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  confsens.Oauth_url,
			TokenURL: confsens.Token_url,
		},
	}

	//	OauthConf.ClientID = confsens.Client_id
	//OauthConf.ClientSecret = confsens.Client_secret
	//	OauthConf.Scopes = []string{"read", "write"}
	/*	OauthConf.Endpoint = oauth2.Endpoint{
			AuthURL:  confsens.Oauth_url,
			TokenURL: confsens.Token_url,
		}
	*/
	/*OauthConf = &oauth2.Config{
		ClientID:     confsens.Client_id,
		ClientSecret: confsens.Client_secret,
		Scopes:       []string{"read", "write"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  confsens.Oauth_url,
			TokenURL: confsens.Token_url,
		},
	}*/
	fmt.Println("OauthConf : ", OAuth2_config.Config)
}

func cp(dst, src string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	// no need to check errors on read only file, we already got everything
	// we need from the filesystem, so nothing can go wrong now.
	defer s.Close()
	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(d, s); err != nil {
		d.Close()
		return err
	}
	return d.Close()
}

// itob returns an 8-byte big endian representation of v.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

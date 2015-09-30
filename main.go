package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/url"
	"strconv"

	"github.com/ChimeraCoder/anaconda"
)

func initApi() *anaconda.TwitterApi {
	anaconda.SetConsumerKey(CONSUMER_KEY)
	anaconda.SetConsumerSecret(CONSUMER_SECRET)
	return anaconda.NewTwitterApi(ACCESS_TOKEN, ACCESS_TOKEN_SECRET)
}

var maxApiFavoritesCount = 20

func getFavorites(api *anaconda.TwitterApi) {
	urlValues := url.Values{}
	urlValues.Set("screen_name", SCREEN_NAME)
	urlValues.Set("count", strconv.Itoa(maxApiFavoritesCount))
	favorites, err := api.GetFavorites(urlValues)

	if err != nil {
		log.Println("GetFavorites returned error: %s", err.Error())
		return
	}

	log.Printf("got %d favorites\n", len(favorites))
	for _, tweet := range favorites {
		log.Printf("%s %s", tweet.CreatedAt, tweet.IdStr)
		log.Printf("%s", tweet.Text)
	}

	lastIdStr := favorites[len(favorites)-1].IdStr

	for {
		urlValues.Set("max_id", lastIdStr)
		favorites, err := api.GetFavorites(urlValues)

		if err != nil {
			log.Println("GetFavorites returned error: %s", err.Error())
			break
		}

		log.Printf("got %d favorites\n", len(favorites))
		for _, tweet := range favorites[1:] { //the first tweet after a max_id call is the max_id tweet which we already processed
			log.Printf("%s %s", tweet.CreatedAt, tweet.IdStr)
			log.Printf("%s", tweet.Text)
		}

		if lastIdStr == favorites[len(favorites)-1].IdStr {
			log.Println("read last favorite")
			break
		}

		lastIdStr = favorites[len(favorites)-1].IdStr
	}
}

func main() {
	readTwitterData()
	api := initApi()

	if api.Credentials == nil {
		panic("Twitter Api client has empty (nil) credentials")
	}

	getFavorites(api)
}

func readTwitterData() {
	data, err := ioutil.ReadFile("twitter.json")
	if err != nil {
		panic(err)
	}
	config := &TwitterData{}
	err = json.Unmarshal([]byte(data), &config)
	if err != nil {
		panic(err)
	}

	CONSUMER_KEY = config.ConsumerKey
	CONSUMER_SECRET = config.ConsumerSecret
	ACCESS_TOKEN = config.AccessToken
	ACCESS_TOKEN_SECRET = config.AccessTokenSecret
	SCREEN_NAME = config.ScreenName
}

var CONSUMER_KEY string
var CONSUMER_SECRET string
var ACCESS_TOKEN string
var ACCESS_TOKEN_SECRET string
var SCREEN_NAME string

type TwitterData struct {
	ConsumerKey       string `json:"consumerKey"`
	ConsumerSecret    string `json:"consumerSecret"`
	AccessToken       string `json:"accessToken"`
	AccessTokenSecret string `json:"accessTokenSecret"`
	ScreenName        string `json:"ScreenName"`
}

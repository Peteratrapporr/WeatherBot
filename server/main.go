package bot

import (
	"appengine"
	"appengine/urlfetch"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const URL_OPEN_WEATHER_MAP_API = "http://api.openweathermap.org/data/2.5/weather"
const APP_ID = "58cbb42fff062a563a71946439ea9aea"  //API Key from openweather.org
const Rapporr_Token = "MyHclPbYMuAYjyYqTzgh" //Token from Rapporr API setup screen for Webhook

//OpenWeather JSON data
type Forecast struct {
	Id   float64   `json:"id"`
	Name string    `json:"name"`
	Cod  float64   `json:"cod"`
	Info MainBlock `json:"main"`
}

type MainBlock struct {
	Temp     float64 `json:"temp"`
	Pressure float64 `json:"pressure"`
	humidity float64 `json:"humidity"`
	TempMin  float64 `json:"temp_min"`
	TempMax  float64 `json:"temp_max"`
}

//Rapporr Simple Reply JSON format
type SimpleReply struct {
	Response_type   string    	`json:"response_type"`
	Text     	string    	`json:"text"`
}

func init() {
	http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {

	ctx := appengine.NewContext(r)
	//From data being sent from Rapporr
	/*
		token=ItoB7oEyZIbNmHPfxHQ2GrbC
		team_id=0001
		team_domain=example
		channel_id=xLxysxdhTY23
		channel_name=test
		user_id=0001-1001
		user_name=Steve
		command=/weather
		text=Sydney
	*/

	//get data from call
	token := r.FormValue("token")

	if token != Rapporr_Token {
		w.WriteHeader(400)
	}
	//Stop a posssble infinite loop if process inofmration from itself
	//This could be  the user_id and recormmend for your own system.
	//Have kept this by Name as being used in different demos
	if "Weather Bot" == r.FormValue("user_name") {
		w.WriteHeader(400)
	}

	//City to search on//Shoud add extra testing to make sure this is just a city name and not some other random text
	text := r.FormValue("text")

	ctx.Infof("City %v", text)

	weatherData, err := Search(text, ctx)
	if err != nil {
		ctx.Infof("Failed to get weather %v",err)
		w.WriteHeader(400)
		return
	}

	//setup reply to Rapporr in simple format
	reply := SimpleReply{}

	reply.Response_type = "text"
	reply.Text = "Current Temp for " + text + " is " + fmt.Sprintf("%.2f", weatherData.Info.Temp - 273.15) + "c"

	json, _ := json.Marshal(reply)

	w.Write(json)

}

func Search(city string, ctx appengine.Context) (*Forecast, error) {
	resp, err := GetCurrentWeather(city, ctx)
	//Handle response from search
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var f Forecast
	err = json.Unmarshal(body, &f)

	if err != nil {
		return nil, err
	}

	return &f, nil
}

func GetCurrentWeather(city string, ctx appengine.Context ) (*http.Response, error) {
	//get weather data from openweather.org
	url := fmt.Sprintf(URL_OPEN_WEATHER_MAP_API+"?q=%s&appid=%s", city, APP_ID)

	request, err := http.NewRequest("GET", url, nil)

	client := urlfetch.Client(ctx)
	response, err := client.Do(request)

	if err != nil {
		return response, err
	}
	return response, nil
}
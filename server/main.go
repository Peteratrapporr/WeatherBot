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
const APP_ID = "58cbb42fff062a563a71946439ea9aea"
const Rapporr_Token = "MyHclPbYMuAYjyYqTzgh"

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

type SimpleReply struct {
	Response_type   string    	`json:"response_type"`
	Text     	string    	`json:"text"`
}

func init() {
	http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {

	ctx := appengine.NewContext(r)

	//get data from call
	token := r.FormValue("token")

	if token != Rapporr_Token {
		w.WriteHeader(400)
	}

	if "Weather Bot" == r.FormValue("user_name") {
		w.WriteHeader(400)
	}
	text := r.FormValue("text")

	ctx.Infof("City %v", text)

	weatherData, err := Search(text, ctx)
	if err != nil {
		ctx.Infof("Failed to get weather %v",err)
		return
	}

	reply := SimpleReply{}

	reply.Response_type = "text"
	reply.Text = "Current Temp for " + text + " is " + fmt.Sprintf("%.2f", weatherData.Info.Temp - 273.15) + "c"

	json, _ := json.Marshal(reply)

	w.Write(json)

}

func Search(city string, ctx appengine.Context) (*Forecast, error) {
	resp, err := GetCurrentWeather(city, ctx)

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

	url := fmt.Sprintf(URL_OPEN_WEATHER_MAP_API+"?q=%s&appid=%s", city, APP_ID)

	request, err := http.NewRequest("GET", url, nil)

	client := urlfetch.Client(ctx)
	response, err := client.Do(request)

	if err != nil {
		return response, err
	}
	return response, nil
}
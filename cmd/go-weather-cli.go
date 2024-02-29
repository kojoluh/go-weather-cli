package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

type Weather struct {
	Timezone string `json:"timezone"`
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct {
		TempC     float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`
	Forecast struct {
		Forecastday []struct {
			Hour []struct {
				TimeEpoch int64   `json:"time_epoch"`
				TempC     float64 `json:"temp_c"`
				Condition struct {
					Text string `json:"text"`
				} `json:"condition"`
				ChanceOfRain float64 `json:"chance_of_rain"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func main() {
	godotenv.Load()
	var lat string
	var lon string
	if len(os.Args) > 2 {
		lat = os.Args[1]
		lon = os.Args[2]
	} else {
		lat = os.Getenv("LATITUDE")
		lon = os.Getenv("LONGITUDE")
	}
	url := "https://api.openweathermap.org/data/3.0/onecall?lat=" + lat + "&lon=" + lon + "&appid=" + os.Getenv("WEATHER_API_KEY")
	// url := "http://api.weatherapi.com/v1/forecast.json?q=Iasi&days=1&aqi=no&alerts=no&key=" + os.Getenv("WEATHER_API_KEY")
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	fmt.Println("url: ", url, os.Getenv("LATITUDE"), os.Getenv("LONGITUDE"), os.Getenv("WEATHER_API_KEY"))
	if res.StatusCode != 200 {
		panic("Weather API unavailable")
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		panic("Weather API response invalid")
	}

	fmt.Println("res: ", string(body))

	var weather Weather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		panic(err)
	}

	fmt.Println(weather)

	location, current, hours := weather.Location, weather.Current, weather.Forecast.Forecastday[0].Hour

	fmt.Printf("%s, %s: %.0fC, %s\n",
		location.Name,
		location.Country,
		current.TempC,
		current.Condition.Text,
	)

	for _, hour := range hours {
		date := time.Unix(hour.TimeEpoch, 0)
		if date.Before(time.Now()) {
			continue
		}

		message := fmt.Sprintf("%s - %.0fC, %.0f, %s\n ",
			date.Format("20:20"),
			hour.TempC,
			hour.ChanceOfRain,
			hour.Condition.Text,
		)

		if hour.ChanceOfRain < 40 {
			fmt.Print(message)
		} else {
			color.Red(message)
		}
	}
}

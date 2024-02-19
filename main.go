package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type WeatherData struct {
	Current struct {
		Condition struct {
			Icon string `json:"icon"`
			Text string `json:"text"`
		} `json:"condition"`
		IsDay int     `json:"is_day"`
		TempC float64 `json:"temp_c"`
	} `json:"current"`
	Forecast struct {
		Forecastday []struct {
			Astro struct {
				IsMoonUp int `json:"is_moon_up"`
				IsSunUp  int `json:"is_sun_up"`
			} `json:"astro"`
			Date string `json:"date"`
			Day  struct {
				AvgTempC  float64 `json:"avgtemp_c"`
				Condition struct {
					Code int    `json:"code"`
					Icon string `json:"icon"`
					Text string `json:"text"`
				} `json:"condition"`
				DailyChanceOfRain int     `json:"daily_chance_of_rain"`
				DailyChanceOfSnow int     `json:"daily_chance_of_snow"`
				DailyWillItRain   int     `json:"daily_will_it_rain"`
				DailyWillItSnow   int     `json:"daily_will_it_snow"`
				MaxTempC          float64 `json:"maxtemp_c"`
				MaxWindKph        float64 `json:"maxwind_kph"`
				MinTempC          float64 `json:"mintemp_c"`
			} `json:"day"`
			Hour []struct {
				Condition struct{} `json:"condition"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
	Location struct {
		Country        string  `json:"country"`
		Lat            float64 `json:"lat"`
		Localtime      string  `json:"localtime"`
		LocaltimeEpoch int     `json:"localtime_epoch"`
		Lon            float64 `json:"lon"`
		Name           string  `json:"name"`
		Region         string  `json:"region"`
		TzID           string  `json:"tz_id"`
	} `json:"location"`
}

type City struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Region  string  `json:"region"`
	Country string  `json:"country"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	URL     string  `json:"url"`
}

func main() {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	router.Use(cors.New(config))

	router.GET("/ping", func(context *gin.Context) {
		context.IndentedJSON(http.StatusOK, gin.H{
			"ping": "pong",
		})
	})

	router.GET("/forcast", func(ctx *gin.Context) {
		id := ctx.Query("id")
		qParams := fmt.Sprintf("key=%s&q=id:%s&days=5&aqi=no&alerts=no", viper.GetString("API_KEY"), id)
		url := viper.GetString("BASE_URL") + "/forecast.json?" + qParams
		resp, err := http.Get(url)
		if err != nil {
			ctx.IndentedJSON(http.StatusForbidden, gin.H{
				"message": err.Error(),
			})
		}

		defer resp.Body.Close()

		forcast, err := io.ReadAll(resp.Body)
		if err != nil {
			ctx.IndentedJSON(http.StatusForbidden, gin.H{
				"message": err.Error(),
			})
		}

		var formattedJSON WeatherData
		if err := json.Unmarshal(forcast, &formattedJSON); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.IndentedJSON(http.StatusOK, formattedJSON)
	})

	router.GET("/city", func(ctx *gin.Context) {
		name := ctx.Query("name")
		qParams := fmt.Sprintf("key=%s&q=%s", viper.GetString("API_KEY"), name)
		url := viper.GetString("BASE_URL") + "/search.json?" + qParams

		resp, err := http.Get(url)
		if err != nil {
			ctx.IndentedJSON(http.StatusForbidden, gin.H{
				"message": err.Error(),
			})
		}

		defer resp.Body.Close()

		list, err := io.ReadAll(resp.Body)
		if err != nil {
			ctx.IndentedJSON(http.StatusForbidden, gin.H{
				"message": err.Error(),
			})
		}

		var cityList []City
		if err := json.Unmarshal(list, &cityList); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.IndentedJSON(http.StatusOK, cityList)
	})

	router.Run("0.0.0.0:8000")
}

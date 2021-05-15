// package utils contains functions that populate the structures that contain the reponses from the api.
// It also has the function that format the various print-outs that happend in the discord client.
package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"projectGroup23/structs"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// struct used to retrieved IP location from api
var ipAddress structs.IPLocation

// struct used to retrieved alpha 2 code from api
var alpha2Code []structs.Alpha2Code

// EnvVar reads from .env file, sets environment variables and returns value based on key
func EnvVar(key string) string {
	// Load .env file
	err := godotenv.Load("./.env")
	if err != nil {
		fmt.Println(err)
	}
	return os.Getenv(key)
}

// GetIPLocation returns the city location of the person calling the discord bot.
func GetIPLocation() (string, error) {
	resp, err := http.Get("https://ipwhois.app/json/")
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&ipAddress)
	if err != nil {
		return "", err
	}

	return ipAddress.City, nil
}

// Get2AlphaCode is used to match the city name supplied in newsletter command to actual cities.
func Get2AlphaCode(countryName string) (string, error) {

	country := strings.Title(strings.ToLower(countryName))

	resp, err := http.Get("https://restcountries.eu/rest/v2/name/" + country + "?fullText=true")
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&alpha2Code)
	if err != nil {
		return "", err
	}

	return alpha2Code[0].Alpha2Code, nil
}

// CheckIfSameDate is a helper for the webhook functionality.
func CheckIfSameDate(date, date2 time.Time) bool {
	y, m, d := date.Date()
	y2, m2, d2 := date2.Date()
	return y == y2 && m == m2 && d == d2
}

// returns status code from api
func CheckStatusCodeApi(url string) string {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Errorf("Error in creating request:", err.Error())
	}
	client := &http.Client{}
	res, err := client.Do(request)
	statusCode := 0
	if res != nil {
		statusCode = res.StatusCode
	} else {
		statusCode = http.StatusBadRequest
	}
	return strconv.Itoa(statusCode)
}

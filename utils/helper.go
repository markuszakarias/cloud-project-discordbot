package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"projectGroup23/structs"
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

// GetIPLocation -
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

	fmt.Println(alpha2Code[0].Alpha2Code)

	return alpha2Code[0].Alpha2Code, nil
}

func CheckIfSameDate(date, date2 time.Time) bool {
	y, m, d := date.Date()
	y2, m2, d2 := date2.Date()
	return y == y2 && m == m2 && d == d2
}

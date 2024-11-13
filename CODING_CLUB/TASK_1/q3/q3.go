package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var URL string = "https://api.wheretheiss.at/v1/satellites/25544/positions?timestamps="

func main() {
	var time_entered_as_string string
	var time_req_layout time.Time
	var err error

	for true {
		fmt.Println("Enter the date and time in \"DD/MM/YYYY_Hours:Minutes:Seconds\" format")
		fmt.Scanln(&time_entered_as_string)
		time_entered_as_string = strings.TrimSpace(time_entered_as_string)
		time_req_layout, err = time.Parse("02/01/2006_15:04:05", time_entered_as_string)
		if err == nil {
			break
		}
		fmt.Println("Invalid date/time format.")
	}

	unixTime := time_req_layout.Unix()

	URL += strconv.FormatInt(unixTime, 10)

	body := get_request()
	latitude, longitude := decodeJSON(body)
	fmt.Println("Latitude: ", latitude)
	fmt.Println("Longitude: ", longitude)
}

func get_request() []byte {
	resp, err := http.Get(URL)
	if err != nil {
		fmt.Println("Error encountered")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error encountered")
	}
	return body
}

func decodeJSON(dataJSON []byte) (latitude float64, longitude float64) {
	var decodedJSON []map[string]interface{}
	json.Unmarshal(dataJSON, &decodedJSON)
	latitude, _ = decodedJSON[0]["latitude"].(float64)
	longitude, _ = decodedJSON[0]["longitude"].(float64)
	return
}

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
	var date, time_t string
	fmt.Println("Enter the date in DD/MM/YYYY format('/' are necessary): ")
	fmt.Scanln(&date)
	fmt.Println("Now enter the time in Hours:Minutes:Seconds format(':' are necessary): ")
	fmt.Scanln(&time_t)

	unixTime := convert_time_to_unix_time(date, time_t)

	URL += strconv.FormatInt(unixTime, 10)

	body := get_request()
	latitude, longitude := decodeJSON(body)
	fmt.Printf("The latitude = %v and longitude = %v", latitude, longitude)
}

func convert_time_to_unix_time(date string, time_t string) int64 {
	date_spliced := strings.Split(date, "/")
	time_spliced := strings.Split(time_t, ":")

	t := time.Date(convert_string_to_int(date_spliced[2]), time.Month(convert_string_to_int(date_spliced[1])), convert_string_to_int(date_spliced[0]), convert_string_to_int(time_spliced[0]), convert_string_to_int(time_spliced[1]), convert_string_to_int(time_spliced[2]), 0, time.UTC)
	converted_time := t.Unix()
	return converted_time
}

func convert_string_to_int(s string) (integer int) {
	integer, _ = strconv.Atoi(s)
	return
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

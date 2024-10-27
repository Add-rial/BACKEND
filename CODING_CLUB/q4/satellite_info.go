package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

const url_current_lat_long = "https://api.wheretheiss.at/v1/satellites/25544"

var url_current_country_timezone string = "https://api.wheretheiss.at/v1/coordinates/"
var url_time_lat_long string = "https://api.wheretheiss.at/v1/satellites/25544/positions?timestamps="

func main() {
	//Directory open in vscode is GO
	lat_long := flag.Bool("lat_long", false, "Get current ISS location in latitude and longitude\nExample :go run satellite_info.go -lat_long")
	country_timezone := flag.Bool("country_timezone", false, "Get current ISS location as country code and time zone\ngo run satellite_info.go -country_timezone")
	lat_long_time := flag.String("lat_long_time", "", "Get ISS location at a specific date/time (format: \"DD/MM/YYYY_HH:MM:SS\")\ngo run satellite_info.go -lat_long_time \"26/10/2024_14:30:00\"")

	flag.Parse()

	if len(os.Args) < 2 { //Failsafe to check when the iser enters the flag but not the values
		fmt.Println("Usage:")
		flag.PrintDefaults()
		return
	}

	if *lat_long {
		//For current lat and long
		body := get_request(url_current_lat_long)
		latitude, longitude := decodeJSON_lat_long(body)

		fmt.Println("Latitude: ", latitude)
		fmt.Println("Longitude: ", longitude)
	} else if *country_timezone {
		//get current country and timezone
		body := get_request(url_current_lat_long)
		latitude, longitude := decodeJSON_lat_long(body)

		url_current_country_timezone += strconv.FormatFloat(latitude, 'f', -1, 64) + "," + strconv.FormatFloat(longitude, 'f', -1, 64)

		body = get_request(url_current_country_timezone)
		country_code, timezone_id := decodeJSON_Country_timezone(body)

		fmt.Println("Country code: ", country_code)
		fmt.Println("Timezone id: ", timezone_id)
	} else if *lat_long_time != "" {
		//gets ISS location at time entered by user
		parsedTime, err := time.Parse("02/01/2006_15:04:05", *lat_long_time)
		if err != nil {
			fmt.Println("Invalid date/time format. Usage: \"DD/MM/YYYY_HH:MM:SS\".")
			return
		}

		unixTime := parsedTime.Unix()

		url_time_lat_long += strconv.FormatInt(unixTime, 10)

		body := get_request(url_time_lat_long)
		latitude, longitude := decodeJSON_lat_long_time(body)
		fmt.Println("Latitude: ", latitude)
		fmt.Println("Longitude: ", longitude)
	} else {
		fmt.Println("Invalid command. Use -h or --help for usage information.")
	}
}

func get_request(URL string) []byte {
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

func decodeJSON_lat_long(dataJSON []byte) (latitude float64, longitude float64) {
	var dataDecoded map[string]interface{}

	json.Unmarshal(dataJSON, &dataDecoded)
	latitude, _ = dataDecoded["latitude"].(float64)
	longitude, _ = dataDecoded["longitude"].(float64)
	return
}

func decodeJSON_Country_timezone(dataJSON []byte) (country_code string, timezone_id string) {
	var decodedJSON map[string]interface{}

	json.Unmarshal(dataJSON, &decodedJSON)
	country_code, _ = decodedJSON["country_code"].(string)
	timezone_id, _ = decodedJSON["timezone_id"].(string)
	return
}

func decodeJSON_lat_long_time(dataJSON []byte) (latitude float64, longitude float64) {
	var decodedJSON []map[string]interface{}
	json.Unmarshal(dataJSON, &decodedJSON)
	latitude, _ = decodedJSON[0]["latitude"].(float64)
	longitude, _ = decodedJSON[0]["longitude"].(float64)
	return
}

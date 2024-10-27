package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

const url1 = "https://api.wheretheiss.at/v1/satellites/25544"

var url2 string = "https://api.wheretheiss.at/v1/coordinates/"

func main() {
	//getting the latitude and longitude
	body1 := get_request(url1)
	latitude, longitude := decodeJSON(body1)

	url2 += strconv.FormatFloat(latitude, 'f', -1, 64) + "," + strconv.FormatFloat(longitude, 'f', -1, 64)

	body2 := get_request(url2)
	country_code, timezone_id := getCountry(body2)
	fmt.Printf("The country code is: %v and timezone id is: %v", country_code, timezone_id)
}

func get_request(url string) []byte {
	resp, err := http.Get(url)
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
	var decodedJSON map[string]interface{}

	json.Unmarshal(dataJSON, &decodedJSON)
	latitude, _ = decodedJSON["latitude"].(float64)
	longitude, _ = decodedJSON["longitude"].(float64)
	return
}

func getCountry(dataJSON []byte) (country_code string, timezone_id string) {
	var decodedJSON map[string]interface{}

	json.Unmarshal(dataJSON, &decodedJSON)
	country_code, _ = decodedJSON["country_code"].(string)
	timezone_id, _ = decodedJSON["timezone_id"].(string)
	return
}

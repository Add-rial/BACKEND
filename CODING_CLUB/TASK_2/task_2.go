package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var id_data []string

var branch_codes = map[string]string{
	"cs":         "A7",
	"ece":        "AA",
	"eee":        "A3",
	"eni":        "A8",
	"mech":       "A4",
	"civil":      "A2",
	"phy":        "B5",
	"chem":       "B2",
	"chemical":   "A1",
	"math":       "B4",
	"bio":        "B1",
	"eco":        "B3",
	"pharma":     "A5",
	"manu":       "AB",
	"genstudies": "C2",
	"mnc":        "AD", //Don't forget poor wittle mnc
}

var year_codes = map[string]string{
	"1": "2024",
	"2": "2023",
	"3": "2022",
	"4": "2021",
	"5": "2020",
}

var campus_codes = map[string]string{
	"G": "goa",
	"H": "hyderabad",
	"P": "pilani",
}

func main() {            //sample terminal: go run "TASK_2\task_2.go" TASK_2\data.txt
	fmt.Println("This is a program for task 2 of coding club")

	if(len(os.Args) < 2){
		fmt.Println("Provide at least one argument")
	}

	text_to_data(os.Args[1])

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/favicon.ico" { //If the browser requests for the favicon, we ignore it
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json") //content type is set to json after removing
	//favicon request so that the icon does not convert
	//to json

	if(queryRedirector(w, r.URL.Query(), r.URL.Path) == 0){
		if(pathParameterRedirector(w, r.URL.Path) == 0 ){
			fmt.Fprintf(w, "No valid path or query parameter detected")
		}
	}
}

func text_to_data(file_path string) { //stores values in data.txt to a string slice
	file, err := os.Open(file_path)
	if err != nil {
		fmt.Println("An error has occured")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		id_data = append(id_data, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
}

func queryRedirector(w http.ResponseWriter, queries url.Values, path string) int { //redirects all the queries to the various functions
	if queries.Has("format") && queries.Get("format") == "text" {
		base_url_format_text(w)
	} else if queries.Has("branch") {
		base_url_branch(w, branch_codes[queries.Get("branch")])
	} else if queries.Has("year") {
		base_url_year(w, year_codes[queries.Get("year")])
	} else if queries.Has("campus"){
		base_url_campus(w, queries.Get("campus"))//usage <localhost:8000?campus=P>
	} else if path == "/"{
		base_url(w)
	} else{
		return 0
	}
	return 1
}

func pathParameterRedirector(w http.ResponseWriter, path string) int{
	var id_to_print_index int = -1
	split_path := strings.Split(path, "/")

	if len(split_path) != 2 {
		return 0
	}

	for n, id := range id_data {
		if id[8:12] == strings.TrimSpace(split_path[1]) {
			id_to_print_index = n
		}
	}
	if id_to_print_index == -1 {
		err := map[string]string{
			"error": "ID not found in dataset",
		}
		json.NewEncoder(w).Encode(err)
		return 1
	}

	response := map[string]map[string]string{
		"id": {},
	}
	response["id"]["year"] = reverseLookupMap(year_codes, id_data[id_to_print_index][:4])
	response["id"]["branch"] = reverseLookupMap(branch_codes, id_data[id_to_print_index][4:6])
	response["id"]["campus"] = campus_codes[strings.TrimSpace(id_data[id_to_print_index][12:])]
	response["id"]["email"] = "f" + id_data[id_to_print_index][:4] + id_data[id_to_print_index][8:12] + "@" + response["id"]["campus"] + ".bits-pilani.ac.in"
	response["id"]["id"] = id_data[id_to_print_index]
	response["id"]["uid"] = id_data[id_to_print_index][8:12]
	json.NewEncoder(w).Encode(response)
	return 1
}

func base_url(w http.ResponseWriter) {
	response := map[string][]string{
		"ids": id_data,
	}
	json.NewEncoder(w).Encode(response)
}

func base_url_format_text(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain")

	for _, id := range id_data {
		fmt.Fprintf(w, "%s\n", id)
	}
}

func base_url_branch(w http.ResponseWriter, branch string) {
	var id_data_sorted_branch []string
	var is_there_data bool = false

	for _, id := range id_data {
		if id[4:6] == branch {
			id_data_sorted_branch = append(id_data_sorted_branch, id)
			is_there_data = true
		}
	}

	if(is_there_data){
		response := map[string][]string{
			"ids": id_data_sorted_branch,
		}
		json.NewEncoder(w).Encode(response)
    }else{
		response := map[string]string{
			"error": "There is no ID for given branch",
		}
		json.NewEncoder(w).Encode(response)
	}
}

func base_url_year(w http.ResponseWriter, year string) {
	var id_data_sorted_year []string
	var is_there_data bool = false

	for _, id := range id_data {
		if id[0:4] == year {
			id_data_sorted_year = append(id_data_sorted_year, id)
			is_there_data = true
		}
	}

	if(is_there_data){
		response := map[string][]string{
			"ids": id_data_sorted_year,
		}
		json.NewEncoder(w).Encode(response)
    }else{
		response := map[string]string{
			"error": "There is no ID for given year",
		}
		json.NewEncoder(w).Encode(response)
	}
}

func base_url_campus(w http.ResponseWriter, campus string)  {
	var id_data_sorted_campus []string
	var is_there_data bool = false

	for _, id := range id_data {
		if id[12:] == campus {
			id_data_sorted_campus = append(id_data_sorted_campus, id)
			is_there_data = true
		}
	}

	if(is_there_data){
		response := map[string][]string{
			"ids": id_data_sorted_campus,
		}
		json.NewEncoder(w).Encode(response)
    }else{
		response := map[string]string{
			"error": "There is no ID for given campus",
		}
		json.NewEncoder(w).Encode(response)
	}
}

func reverseLookupMap(m map[string]string, value string) string {
	for k, v := range m {
		if v == value {
			return k
		}
	}
	return ""
}

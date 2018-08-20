package main

import (
		"net/http"
	"./timezone"
	"io"
	"fmt"
	"log"
	"encoding/json"
	"time"
	"strconv"
)

var port string = "7777"
var locations map[string]timezone.Timezone

func main() {
	locations = initCommonTimezons()
	for k := range locations {
		stringRepresentation := fmt.Sprintf("%+v", locations[k])
		fmt.Println("Key:", k, "Value:", stringRepresentation)
	}
	fmt.Printf("Starting backend on port %s\n", port)
	runServer()
}

func runServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", defaultServer)
	mux.HandleFunc("/health", healthServer)
	mux.HandleFunc("/locations", locationsServer)
	mux.HandleFunc("/timediff", timeDiffServer)

	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func logRequest( req *http.Request) {
	fmt.Printf("%s request to %s\n", req.Method, req.RequestURI)
}

func timeDiffServer(w http.ResponseWriter, req *http.Request) {
	logRequest(req)
	locFrom := req.URL.Query().Get("locFrom")
	if locFrom == "" {
		log.Fatal("Missing param `locFrom`")
	}
	locToo := req.URL.Query().Get("locToo")
	if locToo == "" {
		log.Fatal("Missing param `locToo`")
	}

	// collect location from
	locationFrom, ok := locations[locFrom]
	if !ok {
		log.Fatal("Missing param `time`")
	}

	// collect location too
	locationToo, ok := locations[locToo]
	if !ok {
		log.Fatal("Missing param `time`")
	}

	// get current time in UTC
	currentTimeUTC := time.Now().UTC()
	hours := fmt.Sprintf("%02d%02d",
		currentTimeUTC.Hour(), currentTimeUTC.Minute())
	i, _ := strconv.Atoi(hours)
	hoursFrom := i + locationFrom.Offset
	hoursToo := i + locationToo.Offset
	difference := hoursFrom - hoursToo
	if hoursToo > 2400 {
		hoursToo = hoursToo - 2400
	}
	leftPad := ""
	if hoursToo < 1000 {
		leftPad = "0"
	}
	if hoursToo < 100{
		leftPad = "00"
	}
	if hoursToo < 10{
		leftPad = "000"
	}
	hoursFromString := fmt.Sprintf("%v", hoursFrom)
	hoursTooString := fmt.Sprintf("%s%v", leftPad,hoursToo)

	// create data struct
	timeDifference := timezone.TimeDifference{
		LocationFrom: locFrom,
		LocationFromTime: printAdjustedTime(hoursFromString),
		LocationToo: locToo,
		LocationTooTime: printAdjustedTime(hoursTooString),
		TimeDifference: difference,
	}
	// print data struct as json

	timeDifferenceJson, _ := json.Marshal(timeDifference)
	io.WriteString(w, string(timeDifferenceJson))
}

func normalizeDifference(rawDifference int) int {
	if rawDifference < 0 {
		rawDifference = rawDifference + (-rawDifference *2)
	}
	return rawDifference
}

func printAdjustedTime(hours string) string {
	return fmt.Sprintf("%s:%s",  hours[0:2], hours[2:4])
}

func locationsServer(w http.ResponseWriter, req *http.Request) {
	logRequest(req)
	locationsJson, _ := json.Marshal(locations)
	io.WriteString(w, string(locationsJson))
}

func defaultServer(w http.ResponseWriter, req *http.Request) {
	logRequest(req)
	io.WriteString(w, "Timezone Service 0.1.0-alpha")
}

func healthServer(w http.ResponseWriter, req *http.Request) {
	logRequest(req)
	io.WriteString(w, "healthy")
}

func initCommonTimezons() map[string] timezone.Timezone {
	BST := timezone.Timezone {
		Abbreviation: "BST",
		Name:         "British Summer Time",
		Locations:    []string{"Europe"},
		Offset:       100,
	}
	CEST := timezone.Timezone {
		Abbreviation: "CEST",
		Name:         "Central European Summer Time",
		Locations:    []string{"Europe", "Antarctica"},
		Offset:       200,
	}
	EST := timezone.Timezone{
		Abbreviation: "EST",
		Name:         "Eastern Standard Time",
		Locations:    []string {"North America", "Caribbean", "Central America"},
		Offset:       -500,
	}
	IST := timezone.Timezone{
		Abbreviation: "IST",
		Name:         "Irish Standard Time",
		Locations:    []string {"Europe"},
		Offset:       100,
	}
	PDT := timezone.Timezone{
		Abbreviation: "PDT",
		Name:         "Pacific Daylight Time",
		Locations:    []string {"North America"},
		Offset:       -700,
	}

	// San Fransisco - Pacific Daylight Time
	// The Netherlands - CEST (currently)
	// North Carolina - Eastern Daylight Time
	// Ireland - Irish Standard Time
	// UK - BST (British Summer Time)
	var locations map[string] timezone.Timezone
	locations = make(map[string]timezone.Timezone)
	locations["west-coast"] = PDT
	locations["east-coast"] = EST
	locations["aib"] = IST
	locations["uk"] = BST
	locations["home"] = CEST
	return locations
}

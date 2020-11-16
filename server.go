package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type Request struct {
	Host string
}

type Time struct {
	CreationTime int
	LastTime int
}

var storage = map[string]Time{}
var timeout = 10 // unit is minute

func handler(w http.ResponseWriter, r *http.Request) {
	var request Request
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var minutes = getBrowsingTime(request.Host, r.RemoteAddr)
	_, err = fmt.Fprintf(w, "%d", minutes)
	if err != nil {
		log.Fatal(err)
	}
}

func getBrowsingTime(host string, ip string) (minutes int) {
	currentTime := int(time.Now().Unix())
	key := ip+host
	storedTime, ok := storage[key]
	if ok {
		// First we should check if this is a new visit.
		minutes = (currentTime - storedTime.LastTime) / 60
		fmt.Println(minutes)
		if minutes >= timeout {
			// This is a new visit
			storedTime.CreationTime = currentTime
			storedTime.LastTime = currentTime
		} else {
			// Not a new visit, update the last visit time.
			storedTime.LastTime = currentTime
		}
	} else {
		// This is a new visit.
		minutes = 0
		storedTime.CreationTime = currentTime
		storedTime.LastTime = currentTime
	}
	storage[key] = storedTime
	return
}


func main()  {
	http.HandleFunc("/", handler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}
	log.Println("Starting server at port " + port + ".")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type Request struct {
	Host string
}

type Time struct {
	CreationTime int
	LastTime     int
}

var storage = sync.Map{}
var timeout = 10 // unit is minute
var debug bool
var port = flag.String("port", "3000", "specify the listening port")

func pingHandler(w http.ResponseWriter, r *http.Request) {
	var request Request
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ip := strings.Split(r.RemoteAddr, ":")[0]
	var minutes = getBrowsingTime(request.Host, ip)
	_, err = fmt.Fprintf(w, "%d", minutes)
	if err != nil {
		log.Fatal(err)
	}
	if debug {
		log.Println(fmt.Sprintf("New request from address: %s, host: %s, return with %d.", ip, request.Host, minutes))
	}
}

func clearHandler(w http.ResponseWriter, r *http.Request) {
	var request Request
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ip := strings.Split(r.RemoteAddr, ":")[0]
	clearBrowsingTime(request.Host, ip)
	http.StatusText(http.StatusAccepted)
	if debug {
		log.Println(fmt.Sprintf("Clear counter for address: %s, host: %s.", ip, request.Host))
	}
}

func getBrowsingTime(host string, ip string) (minutes int) {
	currentTime := int(time.Now().Unix())
	key := ip + host
	value, ok := storage.Load(key)
	var storedTime Time
	if ok {
		storedTime = value.(Time)
		// First we should check if this is a new visit.
		interval := (currentTime - storedTime.LastTime) / 60
		if interval >= timeout {
			// This is a new visit
			storedTime.CreationTime = currentTime
			storedTime.LastTime = currentTime
		} else {
			// Not a new visit, update the last visit time.
			storedTime.LastTime = currentTime
			minutes = (currentTime - storedTime.CreationTime) / 60
		}
	} else {
		// This is a new visit.
		storedTime = Time{
			CreationTime: currentTime,
			LastTime:     currentTime,
		}
	}
	storage.Store(key, storedTime)
	return
}

func clearBrowsingTime(host string, ip string) {
	currentTime := int(time.Now().Unix())
	key := ip + host
	value, ok := storage.Load(key)
	var storedTime Time
	if ok {
		storedTime = value.(Time)
		storedTime.CreationTime = currentTime
		storedTime.LastTime = currentTime
		storage.Store(key, storedTime)
	}
}

func main() {
	flag.Parse()
	debug = os.Getenv("MODE") == "debug"
	if debug {
		log.Println("Debug mode enabled.")
	}
	if *port == "3000" {
		envPort := os.Getenv("PORT")
		if envPort != "" {
			*port = envPort
		}
	}
	log.Println("Starting server at port " + *port + ".")
	server := http.NewServeMux()
	server.HandleFunc("/clear", clearHandler)
	server.HandleFunc("/", pingHandler)
	c := cors.New(cors.Options{AllowedOrigins: []string{"*"}, AllowOriginFunc: func(origin string) bool {
		return true
	}})
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", *port), c.Handler(server)))
}

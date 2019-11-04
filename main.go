package main

import (
  "encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
  "os"
	"github.com/gorilla/mux"
)


type event struct {
	ID          string `json:"ID"`
	status       string `json:"status"`
	size          string `json:"size"`
}

type allEvents []event

var events = allEvents{
	{
		ID:          "1",
		status:       "Success",
		size:         "82kb",
	},
}


func createEvent(w http.ResponseWriter, r *http.Request) {
	var newEvent event
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}

	json.Unmarshal(reqBody, &newEvent)
	events = append(events, newEvent)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newEvent)
}


func getOneEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]

	for _, singleEvent := range events {
		if singleEvent.ID == eventID {
			json.NewEncoder(w).Encode(singleEvent)
		}
	}
}

func getAllEvents(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(events)
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}


func getImageSize(image string) int,error {
  file, err := os.Open(image)
  if err != nil {
      return nil,err
  }
  defer file.Close()

  stat, err := file.Stat()
  if err != nil {
     return nil,err
  }

  var bytes int
  bytes = stat.Size()
  kilobytes := (bytes / 1024)
  return kilobytes,nil
}

func main() {
  success   := true
  size,err  := getImageSize(os.Args[1])
  if err == nil {
     success = true
  }
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
  router.HandleFunc("/query", createEvent).Methods("POST")
  router.HandleFunc("/events", getAllEvents).Methods("GET")
  router.HandleFunc("/events/{id}", getOneEvent).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}

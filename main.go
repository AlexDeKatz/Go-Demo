package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

type Coaster struct {
	Name         string `json:"name"`
	Manufacturer string `json:"manufacturer"`
	ID           string `json:"id"`
	InPark       string `json:"inPark"`
	Height       int    `json:"height"`
}

type CoasterHandler struct {
	sync.Mutex
	store map[string]Coaster
}

func (c *CoasterHandler) getCoasters(w http.ResponseWriter, r *http.Request) {
	coasters := make([]Coaster, len(c.store))
	c.Lock()
	i := 0
	for _, item := range c.store {
		coasters[i] = item
		i++
	}
	c.Unlock()
	jsonBytes, err := json.Marshal(coasters)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (c *CoasterHandler) getCoasterByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if ok, _ := regexp.Match("^[a-zA-Z]+", []byte(params["id"])); ok == false {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Invalid ID"))
		return
	}
	c.Lock()
	selctedCoaster := c.store[params["id"]]
	c.Unlock()
	jsonBytes, err := json.Marshal(selctedCoaster)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (c *CoasterHandler) addCoaster(w http.ResponseWriter, r *http.Request) {
	var coaster Coaster
	err := json.NewDecoder(r.Body).Decode(&coaster)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	coaster.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	c.Lock()
	c.store[coaster.ID] = coaster
	c.Unlock()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Data saved successfully"))
}

func newCoasterHandler() *CoasterHandler {
	return &CoasterHandler{
		Mutex: sync.Mutex{},
		store: map[string]Coaster{
			"id1": {
				Name:         "Taron",
				InPark:       "PhantasiaLand",
				Height:       30,
				Manufacturer: "Intamin",
				ID:           "id1",
			},
		},
	}
}

func main() {
	router := mux.NewRouter()
	coasterHandler := newCoasterHandler()
	router.HandleFunc("/coasters", coasterHandler.getCoasters).Methods(http.MethodGet)
	router.HandleFunc("/coasters/{id}", coasterHandler.getCoasterByID).Methods(http.MethodGet)
	router.HandleFunc("/coasters", coasterHandler.addCoaster).Methods(http.MethodPost)
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}
}

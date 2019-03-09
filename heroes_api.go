package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Hero : stores the information about a super hero
type Hero struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Quality     string `json:"quality"`
	Rating      int8   `json:"rating"` // Out of 5
}

var keyValueStore = make(map[string]Hero)

// HTTP POST request handler for creating a new hero.
func postHeroHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	var hero Hero
	err := json.NewDecoder(r.Body).Decode(&hero)
	if err != nil {
		// Put 500s Internal Server Error code here
		// and replace panic with a log message.
		panic(err)
	}
	keyValueStore[hero.Name] = hero

	j, err := json.Marshal(keyValueStore[hero.Name])
	if err != nil {
		// Put 500s Internal Server Error code here
		// and replace panic with a log message.
		panic(err)
	}
	log.Printf("Created hero: %s\n", hero.Name)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}

// HTTP PUT request handler to update hero specified in the url.
func updateHeroHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	var heroToUpdate Hero
	vars := mux.Vars(r)
	heroName := vars["hero-name"]
	err := json.NewDecoder(r.Body).Decode(&heroToUpdate)
	if err != nil {
		panic(err)
	}
	if _, ok := keyValueStore[heroName]; ok {
		delete(keyValueStore, heroName)
		keyValueStore[heroName] = heroToUpdate
		log.Printf("Hero named: %s updated successfully!", heroName)
		w.WriteHeader(http.StatusOK)
	} else {
		log.Printf("Hero named: %s could not be updated.\n", heroName)
		w.WriteHeader(http.StatusBadRequest)
	}
}

// HTTP GET request handler to get specific hero specified by name in the url.
func getHeroHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	vars := mux.Vars(r)
	heroName := vars["hero-name"]
	if hero, ok := keyValueStore[heroName]; ok {
		j, err := json.Marshal(hero)
		if err != nil {
			// Put 500s Internal Server Error code here
			// and replace panic with a log message
			panic(err)
		}
		log.Printf("Hero: %s found and served\n", heroName)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusFound)
		w.Write(j)
	} else {
		log.Printf("No Hero named: %s found!\n", heroName)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
	}
}

// HTTP GET request handler to get the list of all the heroes
// present in the keyValueStore.
func getAllHeroesHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	var heroes []Hero
	for _, hero := range keyValueStore {
		heroes = append(heroes, hero)
	}
	j, err := json.Marshal(heroes)
	if err != nil {
		// Put 500s Internal Server Error code here
		// and replace panic with a log message
		panic(err)
	}
	log.Println("Found and served all heroes!!")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusFound)
	w.Write(j)
}

// HTTP DELETE request handler to delete all the heroes from keyValuePair.
func deleteAllHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	keyValueStore = make(map[string]Hero)
	log.Println("Thanos snapped and all the heroes died :(")
	w.WriteHeader(http.StatusOK) // Actually not ok ;)
}

func deleteHeroHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	vars := mux.Vars(r)
	heroName := vars["hero-name"]
	if _, ok := keyValueStore[heroName]; ok {
		delete(keyValueStore, heroName)
		log.Printf("Thanos killed %s :(\n", heroName)
		w.WriteHeader(http.StatusOK) // Actually not ok ;)
	} else {
		log.Printf("%s killed Thanos!\n", heroName)
		w.WriteHeader(http.StatusBadRequest)
	}
}

func main() {
	r := mux.NewRouter().StrictSlash(false) // router

	r.HandleFunc("/api/", getAllHeroesHandler).Methods("GET")
	r.HandleFunc("/api/{hero-name}", getHeroHandler).Methods("GET")
	r.HandleFunc("/api/create", postHeroHandler).Methods("POST")
	r.HandleFunc("/api/update/{hero-name}", updateHeroHandler).Methods("PUT")
	r.HandleFunc("/api/delete", deleteAllHandler).Methods("DELETE")
	r.HandleFunc("/api/delete/{hero-name}", deleteHeroHandler).Methods("DELETE")

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	log.Printf("Listening on port: %s", server.Addr)
	server.ListenAndServe()
}

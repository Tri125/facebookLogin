package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	fb "github.com/huandu/facebook"
	"log"
	"net/http"
	"time"
)

var appId, appSecret, redirectUri string

var globalApp *fb.App

func main() {
	globalApp = fb.New(appId, appSecret)
	// http://localhost:8080/login
	globalApp.RedirectUri = redirectUri
	globalApp.EnableAppsecretProof = true

	r := mux.NewRouter()
	r.HandleFunc("/login", LoginHandler).Methods("PUT")
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./public/"))))

	http.Handle("/", r)

	srv := &http.Server{
		Handler: r,
		Addr:    ":8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var tokenPayload fbPayload
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&tokenPayload)
	if r.Body.Close() != nil || err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	session := globalApp.Session(tokenPayload.FbToken)
	err = session.Validate()
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	res, _ := session.Get("/me", fb.Params{
		"fields": "first_name, last_name, email, locale",
	})

	var user User
	res.Decode(&user)

	jsonResponse, _ := json.Marshal(user)
	w.Header().Add("Location", string(jsonResponse))
	w.WriteHeader(http.StatusCreated)
	return
}

type fbPayload struct {
	FbToken string
}

type User struct {
	FirstName string
	LastName  string
	Locale    string
	Email     string
}

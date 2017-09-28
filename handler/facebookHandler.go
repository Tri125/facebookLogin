package handler

import (
	"encoding/json"
	"net/http"

	"fmt"
	fb "github.com/huandu/facebook"
)

type Env struct {
	FbApp *fb.App
}

type User struct {
	Id        string
	FirstName string
	LastName  string
	Locale    string
	Email     string
}

type FacebookRequest struct {
	Token string
}

func CreateFacebookClient(appID string, appSecret string, redirectUri string, enableAppSecretProof bool) *fb.App {
	fbApp := fb.New(appID, appSecret)
	fbApp.RedirectUri = redirectUri
	fbApp.EnableAppsecretProof = enableAppSecretProof
	return fbApp
}

func (env *Env) FacebookLoginHandler(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()
	requestedFields := queries.Get("fields")
	// TODO: Validate fields + token orr send code 400 bad request
	decoder := json.NewDecoder(r.Body)
	var t FacebookRequest
	err := decoder.Decode(&t)
	defer r.Body.Close()
	if err != nil {
		msg := fmt.Errorf("Invalid payload format.")
		writeResponse(w, http.StatusBadRequest, msg)
		return
	}
	if len(t.Token) == 0 {
		msg := fmt.Errorf("Token not sent.")
		writeResponse(w, http.StatusBadRequest, msg)
		return
	}
	session := env.FbApp.Session(t.Token)
	err = session.Validate()
	if err != nil {
		msg := fmt.Errorf("Session is invalid: %s", err)
		writeResponse(w, http.StatusBadRequest, msg)
		return
	}
	res, err := session.Get("/me", fb.Params{
		"fields": requestedFields,
	})
	if err != nil {
		msg := fmt.Errorf("Error querying graphAPI: %s", err)
		writeResponse(w, http.StatusBadRequest, msg)
		return
	}
	var user User
	res.Decode(&user)
	jsonResponse, _ := json.Marshal(user)
	//w.Header().Add("Location", string(jsonResponse))
	//w.WriteHeader(http.StatusCreated)
	w.Write(jsonResponse)
	return
}

func writeResponse(w http.ResponseWriter, status int, msg error) {
	w.WriteHeader(status)
	w.Write([]byte(msg.Error()))
}

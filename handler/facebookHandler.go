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

type facebookRequest struct {
	token string
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
	// TODO: Get token from post
	// TODO: Validate fields + token orr send code 400 bad request
	session := env.FbApp.Session("token")
	err := session.Validate()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Errorf("Session is invalid: %s", err)
		w.Write([]byte(msg.Error()))
		return
	}
	res, err := session.Get("/me", fb.Params{
		"fields": requestedFields,
	})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Errorf("Error querying graphAPI: %s", err)
		w.Write([]byte(msg.Error()))
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

package cmd

import (
	"github.com/Tri125/facebookLogin/handler"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

const DATA_DIR string = "templates/"

var (
	port             int
	cfgFile          string
	endpoint         string
	timeout          time.Duration
	facebookSettings FacebookSettings
)

type FacebookSettings struct {
	AppID                string
	AppSecret            string
	RedirectUri          string
	EnableAppsecretProof bool
}

func runHTTPServer(r *mux.Router, port int, timeout time.Duration, env *handler.Env) error {
	addrs := ":" + strconv.Itoa(port)
	srv := &http.Server{
		Handler:      r,
		Addr:         addrs,
		WriteTimeout: timeout,
		ReadTimeout:  timeout,
	}
	return srv.ListenAndServe()
}

func apiSink(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	return
}

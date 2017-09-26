package cmd

import "time"

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

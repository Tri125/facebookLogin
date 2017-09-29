// Copyright © 2017 Tristan S.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/Tri125/facebookLogin/handler"
	"github.com/gorilla/mux"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "facebookLogin",
	Short: "facebookLogin is a standalone proxy to make request to facebook graphAPI on the user node.",
	Long: `facebookLogin is a standalone proxy to make request to facebook graphAPI on the user node.

You can use the program to run server a side proxy to query simple data about users who authorized your app.
This is ideal to quickly deploy a server-side login flow if an official SDK isn't available for you.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if len(facebookSettings.AppID) == 0 {
			return fmt.Errorf("AppID is missing.")
		}
		if len(facebookSettings.AppSecret) == 0 {
			return fmt.Errorf("AppSecret is missing.")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		fbApp := handler.CreateFacebookClient(facebookSettings.AppID, facebookSettings.AppSecret, facebookSettings.RedirectUri, facebookSettings.EnableAppsecretProof)
		env := &handler.Env{fbApp}
		return runRoot(env)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.facebookLogin.yaml)")
	RootCmd.PersistentFlags().IntVar(&port, "p", 8080, "Port which the http server will listen for traffic.")
	RootCmd.PersistentFlags().StringVar(&endpoint, "path", "/", "Endpoint where the http server will listen for traffic.")
	RootCmd.PersistentFlags().DurationVar(&timeout, "timeout", 15*time.Second, "Set the WriteTimeout and ReadTimeout value for the http server.")
	RootCmd.PersistentFlags().StringVar(&facebookSettings.AppID, "appID", "", "Set the App ID for your facebook application.")
	RootCmd.PersistentFlags().StringVar(&facebookSettings.AppSecret, "appSecret", "", "Set the App Secret for your facebook application.")
	RootCmd.PersistentFlags().BoolVar(&facebookSettings.EnableAppsecretProof, "proof", true, "Prevents malicious clients from making requests on your behalf if tokens are stolen. "+
		"Enabling the appsecret proof status will verify your graph API calls by generating a secret from your appSecret and the token. "+
		"Make sure to change the setting of your facebook app to require app secret on every calls.")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".facebookLogin" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".facebookLogin")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func runRoot(env *handler.Env) error {
	r := mux.NewRouter()
	r = setProxyRouter(r, endpoint, env)
	log.Printf("Listening on port %v", port)
	return runHTTPServer(r, port, timeout, env)
}

func setProxyRouter(r *mux.Router, path string, env *handler.Env) *mux.Router {
	r.HandleFunc(path, env.FacebookLoginHandler).Methods("POST")
	r.HandleFunc(path, apiSink)
	return r
}

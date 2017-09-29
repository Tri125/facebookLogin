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
	"github.com/Tri125/facebookLogin/data"
	"github.com/Tri125/facebookLogin/handler"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"html/template"
	"log"
	"net/http"
)

// localCmd represents the local command
var localCmd = &cobra.Command{
	Use:   "local",
	Short: "Runs a local website to test the proxy.",
	Long: `Runs a local website to test the proxy.
It's ideal to see a working example on how to call the proxy and test if you correctly configured your facebook app.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fbApp := handler.CreateFacebookClient(facebookSettings.AppID, facebookSettings.AppSecret, facebookSettings.RedirectUri, facebookSettings.EnableAppsecretProof)
		env := &handler.Env{fbApp}
		return runLocal(env)
	},
}

func init() {
	RootCmd.AddCommand(localCmd)
}

func runLocal(env *handler.Env) error {
	r := mux.NewRouter()
	r = setProxyRouter(r, endpoint, env)
	r = setLocalWebsiteRouter(r)
	log.Printf("Listening on port %v", port)
	log.Printf("Open http://localhost:%v/dev in a web browser to test.", port)
	return runHTTPServer(r, port, timeout, env)
}

func setLocalWebsiteRouter(r *mux.Router) *mux.Router {
	r.HandleFunc("/dev", indexHandler)
	r.NotFoundHandler = http.HandlerFunc(redirectToRoot)
	return r
}

func redirectToRoot(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/dev", http.StatusSeeOther)
}

func loadTemplates() *template.Template {
	publicFolder := "./public/"
	err := data.RestoreAsset(publicFolder, DATA_DIR+"index.html")
	if err != nil {
		log.Fatal(err)
	}
	return template.Must(template.ParseGlob(publicFolder + DATA_DIR + "/*.html"))
}

var templates = loadTemplates()

func indexHandler(w http.ResponseWriter, r *http.Request) {
	params := struct {
		AppID string
	}{facebookSettings.AppID}
	// you access the cached templates with the defined name, not the filename
	err := templates.ExecuteTemplate(w, "indexPage", params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Currentbrowsers package attempts to find the most recent
// versions of popular browsers.  This data is then easily
// consumable as an API.
package currentbrowsers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"appengine"
	"appengine/datastore"
)

// Browser contains the necessary information for browser
// type and release version.
type Browser struct {
	//Chrome Desktop, Chrome Android, Chrome iOS
	Type    string `json:"type"`
	Version string `json:"version"`
}

func init() {
	r := mux.NewRouter()
	r.HandleFunc("/", addDefaultHeaders(IndexHandler))
	r.HandleFunc("/check", addDefaultHeaders(CheckHandler))
	http.Handle("/", r)
}

func addDefaultHeaders(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		fn(w, r)
	}
}

// IndexHandler is responsible for listing the most
// recent browsers.
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	q := datastore.NewQuery("browsers")
	var browsers []Browser
	_, err := q.GetAll(c, &browsers)
	if err != nil {
		log.Fatal(err.Error())
	}
	bs, _ := json.Marshal(browsers)
	fmt.Fprintf(w, string(bs))
}

// CheckHandler is responsible for refreshing the list of most
// recent browsers.
func CheckHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	loadChrome(c)
	loadFirefox(c)
	loadSafari(c)
}

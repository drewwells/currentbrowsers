// Currentbrowsers package attempts to find the most recent
// versions of popular browsers.  This data is then easily
// consumable as an API.
package currentbrowsers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"appengine"

	"github.com/gorilla/mux"
	"labix.org/v2/mgo/bson"
)

// Browser contains the necessary information for browser
// type and release version.
type Browser struct {
	//Chrome Desktop, Chrome Android, Chrome iOS
	Type    string `bson:"_id" json:"type"`
	Version string `bson:"version" json:"version"`
}

func init() {
	r := mux.NewRouter()
	r.HandleFunc("/", IndexHandler)
	r.HandleFunc("/check", CheckHandler)
	http.Handle("/", r)
}

// IndexHandler is responsible for listing the most
// recent browsers.
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	sess := session()
	defer sess.Close()
	db := sess.DB("checker")
	c := db.C("browsers")
	bros := make([]Browser, 10)
	err := c.Find(bson.M{}).All(&bros)
	if err != nil {
		log.Fatal(err.Error())
	}
	bs, _ := json.Marshal(bros)
	fmt.Fprintf(w, string(bs))
}

// CheckHandler is responsible for refreshing the list of most
// recent browsers.
func CheckHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	_, _ = loadChrome(c), loadFirefox(c)
}

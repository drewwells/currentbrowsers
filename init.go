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

func FirefoxHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	fmt.Fprintf(w, "%v", loadFirefox(c))
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	sess := Session()
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

func CheckHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	_, _ = loadChrome(c), loadFirefox(c)
}

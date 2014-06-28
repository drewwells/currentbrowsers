package checker

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Browser struct {
	//Chrome Desktop, Chrome Android, Chrome iOS
	Type    string `bson:"_id"`
	Version string `bson:"version"`
}

func init() {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	http.Handle("/", r)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	f := loadChrome()
	for _, v := range f {
		fmt.Fprintf(w, "%s: %s\n", string(v.Type), string(v.Version))
	}
}

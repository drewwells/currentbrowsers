package currentbrowsers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"

	"appengine"
	"appengine/datastore"
	"appengine/urlfetch"
)

var safariURI = `http://support.apple.com/kb/index?page=downloads_search&facet=all&category=&q=safari&locale=en_US`

type sdl struct {
	Title     string
	Thumbnail string
}

type sJSON struct {
	Downloads []sdl
}

func loadSafari(c appengine.Context) {
	client := urlfetch.Client(c)
	resp, err := client.Get(safariURI)
	if err != nil {
		log.Fatal(err.Error())
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err.Error())
	}
	sj := sJSON{}
	err = json.Unmarshal(bytes, &sj)
	if err != nil {
		log.Fatal(err.Error())
	}
	var exs []Browser
	q := datastore.NewQuery("browsers").
		Filter("Type =", "Safari Desktop")
	if _, err := q.GetAll(c, &exs); err != nil {
		log.Fatal(err.Error())
	}
	for _, v := range sj.Downloads {
		if !safariWL(v.Title) {
			continue
		}
		infos := strings.Split(v.Title, " ")
		b := Browser{
			safariNames[v.Title],
			infos[1],
		}
		if len(exs) != 0 &&
			!compareVersions(exs[0].Version, b.Version) {
			continue
		}
		key := datastore.NewKey(
			c,
			"browsers",
			b.Type,
			0,
			nil,
		)
		_, err := datastore.Put(c, key, &b)
		if err != nil {
			log.Fatal(err.Error())
		}

	}
}

var safariNames = map[string]string{
	"Safari 5.0.6 for Leopard": "Safari Desktop",
}

func safariWL(s string) bool {
	for k, _ := range safariNames {
		if s == k {
			return true
		}
	}
	return false
}

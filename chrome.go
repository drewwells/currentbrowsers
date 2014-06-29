package currentbrowsers

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"appengine"
	"appengine/datastore"
	"appengine/urlfetch"
)

type xFE struct {
	Id      string `xml:"id"`
	Title   string `xml:"title"`
	Content []byte `xml:"content"`
}

type xFeed struct {
	XMLName xml.Name `xml:"feed"`
	Entry   []xFE    `xml:"entry"`
	Id      string   `xml:"id"`
}

func loadChrome(c appengine.Context) (browsers []Browser) {
	//bytes, err := ioutil.ReadFile("./default.xml")
	client := urlfetch.Client(c)
	resp, err := client.Get("https://www.blogger.com/feeds/8982037438137564684/posts/default")
	if err != nil {
		log.Fatal(err.Error())
	}
	bytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	f := xFeed{}
	err = xml.Unmarshal(bytes, &f)
	if err != nil {
		log.Fatal(err.Error())
	}
	rVer := regexp.MustCompile(`\d\d\.[^\s]+`)
	for _, v := range f.Entry {
		if rVer.Match(v.Content) &&
			chromeWL(v.Title) {
			b := Browser{
				v.Title,
				string(rVer.Find(v.Content)),
			}
			var exs []Browser
			q := datastore.NewQuery("browsers").
				Filter("Type =", v.Title)
			_, err := q.GetAll(c, &exs)
			if err != nil {
				log.Fatal(err.Error())
			}
			if len(exs) == 0 ||
				compareVersions(exs[0].Version, b.Version) {
				key := datastore.NewKey(
					c,
					"browsers",
					v.Title,
					0,
					nil,
				)
				_, err := datastore.Put(c, key, &b)
				if err != nil {
					log.Fatal(err.Error())
				}
			}
		}
	}
	return browsers
}

var cWL = []string{
	"Stable Channel Update",
	"Chrome for Android Update",
	"Chrome for iOS Update",
}

func chromeWL(s string) bool {
	for _, w := range cWL {
		if s == w {
			return true
		}
	}
	return false
}

func compareVersions(a, b string) bool {
	a = strings.Replace(a, ".", "", 4)
	b = strings.Replace(b, ".", "", 4)
	return a < b
}

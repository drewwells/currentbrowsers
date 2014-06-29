package currentbrowsers

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"appengine"
	"appengine/urlfetch"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type xFE struct {
	Id      string `xml:"id"`
	Title   string `xml:"title"`
	Content []byte `xml:"content"`
}

type xFeed struct {
	XMLName xml.Name `xml:"feed"`
	Entry   []xFE    `xml:"entry"`
	Id      string   `xml:"id" bson:"_id"`
}

func loadChrome(ctx appengine.Context) (browsers []Browser) {
	//bytes, err := ioutil.ReadFile("./default.xml")
	client := urlfetch.Client(ctx)
	resp, err := client.Get("https://www.blogger.com/feeds/8982037438137564684/posts/default")
	if err != nil {
		log.Fatal(err.Error())
	}
	bytes, err := ioutil.ReadAll(resp.Body)

	session, err := mgo.Dial("mongodb://checker:checker1@kahana.mongohq.com:10012/checker")

	if err != nil {
		log.Fatal(err)
	}
	db := session.DB("checker")
	c := db.C("entries")
	defer session.Close()

	f := xFeed{}
	err = xml.Unmarshal(bytes, &f)
	if err != nil {
		log.Fatal(err.Error())
	}
	infs := make([]interface{}, len(f.Entry))
	for i, v := range f.Entry {
		infs[i] = v
	}
	//This will generate a lot of errors, should check
	//if document already exists.
	_ = c.Insert(infs...)
	bc := db.C("browsers")
	rVer := regexp.MustCompile(`\d\d\.[^\s]+`)
	for _, v := range f.Entry {
		if rVer.Match(v.Content) &&
			chromeWL(v.Title) {
			b := Browser{
				v.Title,
				string(rVer.Find(v.Content)),
			}
			ex := Browser{}
			err = bc.Find(bson.M{"_id": v.Title}).One(&ex)
			if compareVersions(ex.Version, b.Version) {
				log.Println("Attempted Insert")
				bc.UpsertId(v.Title, b)
			}
			browsers = append(browsers, b)
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

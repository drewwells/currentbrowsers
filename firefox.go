package currentbrowsers

import (
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"code.google.com/p/go-html-transform/css/selector"
	"code.google.com/p/go-html-transform/h5"
	"labix.org/v2/mgo/bson"

	"appengine"
	"appengine/urlfetch"
)

func loadFirefox(ctx appengine.Context) (browsers []Browser) {
	client := urlfetch.Client(ctx)
	resp, err := client.Get("http://ftp.mozilla.org/pub/mozilla.org/firefox/releases/")
	if err != nil {
		log.Fatal(err.Error())
	}
	tree, err := h5.New(resp.Body)
	if err != nil {
		log.Fatal(err.Error())
	}
	n := tree.Top()
	// Create a Chain object from a CSS selector statement
	chn, err := selector.Selector("td > a")
	if err != nil {
		log.Fatal(err.Error())
	}
	nodes := chn.Find(n)
	tvers := make([]string, len(nodes))
	rVer := regexp.MustCompile(`^\d*[\.\d*]*$`)
	i := 0
	_ = rVer
	for _, v := range nodes {
		var iv int64
		txt := strings.TrimSuffix(v.FirstChild.Data, "/")
		split := strings.Split(txt, ".")
		if iv, err = strconv.ParseInt(split[0], 10, 8); err != nil {
			continue
		}
		if rVer.MatchString(txt) && iv > 10 {
			tvers[i] = txt
			i = i + 1
		}
	}
	vers := make([]string, i)
	copy(vers, tvers)
	sort.Strings(vers)

	sess := Session()
	defer sess.Close()
	col := sess.DB("checker").C("browsers")
	b := Browser{
		Type:    "Firefox Desktop",
		Version: vers[len(vers)-1],
	}
	ex := Browser{}
	err = col.Find(bson.M{"_id": b.Type}).One(&ex)
	if err != nil || compareVersions(ex.Version, b.Version) {
		log.Println("Attempted Insert")
		col.UpsertId(b.Type, b)
		browsers = append(browsers, b)
	} else {
		browsers = append(browsers, ex)
	}

	return browsers
}

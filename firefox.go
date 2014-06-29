package currentbrowsers

import (
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"code.google.com/p/go-html-transform/css/selector"
	"code.google.com/p/go-html-transform/h5"

	"appengine"
	"appengine/datastore"
	"appengine/urlfetch"
)

func loadFirefox(c appengine.Context) (browsers []Browser) {
	client := urlfetch.Client(c)
	resp, err := client.Get("http://ftp.mozilla.org/pub/mozilla.org/firefox/releases/")
	if err != nil {
		log.Fatal(err.Error())
	}
	tree, err := h5.New(resp.Body)
	if err != nil {
		log.Fatal(err.Error())
	}
	n := tree.Top()
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
	//The last anchor in the sorted list will be the latest version
	b := Browser{
		Type:    "Firefox Desktop",
		Version: vers[len(vers)-1],
	}
	var exs []Browser
	q := datastore.NewQuery("browsers").
		Filter("Type =", "Firefox Desktop")
	_, err = q.GetAll(c, &exs)
	if err != nil {
		log.Fatal(err.Error())
	}
	if len(exs) == 0 || compareVersions(exs[0].Version, b.Version) {
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
		browsers = append(browsers, b)
	} else if len(exs) > 0 {
		browsers = append(browsers, exs[0])
	}

	return browsers
}

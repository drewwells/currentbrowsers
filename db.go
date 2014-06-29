package currentbrowsers

import (
	"log"

	"labix.org/v2/mgo"
)

func session() *mgo.Session {
	session, err := mgo.Dial("mongodb://checker:checker1@kahana.mongohq.com:10012/checker")

	if err != nil {
		log.Fatal(err.Error())
	}
	return session
}

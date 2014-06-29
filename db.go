package currentbrowsers

import (
	"log"
	"os"

	"labix.org/v2/mgo"
)

func session() *mgo.Session {
	session, err := mgo.Dial("mongodb://checker:checker1@kahana.mongohq.com:10012/checker")

	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	return session
}

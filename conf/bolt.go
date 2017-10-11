package conf

import (
	"log"

	"github.com/boltdb/bolt"
	"github.com/xlab/closer"
	"github.com/Sirupsen/logrus"
)

func BoltConnect(config *Configuration) (*bolt.DB, error) {

	db, err := bolt.Open(config.DB.Path, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	closer.Bind(func() {
		logrus.Info("Closing database file")
		db.Close()
	})


	return db, nil
}

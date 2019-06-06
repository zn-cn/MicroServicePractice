package common

import (
	"fmt"
	"time"
	"Ethan/MicroServicePractice/config"

	"gopkg.in/mgo.v2"
)

const (
    DEFAULT_HOST = "127.0.0.1"
    DEFAULT_Port = 27017
    DEFAULT_DB = "mongodb"
    DEFAULT_Admin_DB_Name = "admin"
)


// CreateDBSession create the session of db
func CreateDBSession(service string) (*mgo.Session, error) {
    db := config.GetDB(service)
	if db.Host == "" {
		db.Host = DEFAULT_HOST
    }
    if db.Port == 0 {
        db.Port = DEFAULT_Port
    }

    if db.DriverName == "" {
        db.DriverName = DEFAULT_DB
    }
    var dbURL string
    if db.User != "" && db.PW != "" {
	    if db.AdminDBName == "" {
		    db.AdminDBName = DEFAULT_Admin_DB_Name
	    }
	    dbURL = fmt.Sprintf("%s://%s:%s@%s:%d/%s",db.DriverName,  db.User, db.PW, db.Host, db.Port, db.AdminDBName)
    } else {
    	dbURL = fmt.Sprintf("%s://%s:%d",db.DriverName, db.Host, db.Port)
    }

    // TODO
    // switch db.DriverName {
    // case "mongodb":

    // }
	session, err := mgo.DialWithTimeout(dbURL, 10*time.Second)
	if err != nil {
		return nil, err
	}
	session.SetMode(mgo.Monotonic, true)
	// default is 4096
	session.SetPoolLimit(1000)
	return session, nil
}

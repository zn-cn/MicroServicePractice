package dao

import (
	"fmt"
	"time"

	"gopkg.in/mgo.v2"
)

const (
	DEFAULT_HOST = "127.0.0.1"
)

// CreateSession create the session of mongo
func CreateSession(host string, port int32) (*mgo.Session, error) {
	if host == "" {
		host = DEFAULT_HOST
	}
	mongoURL := fmt.Sprintf("mongodb://%s:%d", host, port)
	session, err := mgo.DialWithTimeout(mongoURL, 10*time.Second)
	if err != nil {
		return nil, err
	}
	session.SetMode(mgo.Monotonic, true)
	// default is 4096
	session.SetPoolLimit(1000)
	return session, nil
}

func CreateSessionWithAuth(host string, port int32, user, password, adminDBName string) (*mgo.Session, error) {
	if adminDBName == "" {
		adminDBName = "admin"
	}
	mongoURL := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s", user, password, host, port, adminDBName)
	session, err := mgo.DialWithTimeout(mongoURL, 10*time.Second)
	if err != nil {
		return nil, err
	}
	session.SetMode(mgo.Monotonic, true)
	// default is 4096
	session.SetPoolLimit(1000)
	return session, nil
}

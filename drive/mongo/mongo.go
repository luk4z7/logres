// Logres - Distributed logs system PostgresSQL to MongoDB
// https://github.com/luk4z7/logres for the canonical source repository
//
// Copyright 2017 The Lucas Alves Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// drive
package mongo

import (
	"gopkg.in/mgo.v2"
	"logres/service/model"
	"os"
	"fmt"
	"logres/service/config"
)

const (
	PRODUCTION = "production"
	LOCALHOST  = "localhost"
)

func New() (*mgo.Collection, *mgo.Collection) {
	return GetSession(LOCALHOST, "logger"), GetSession(PRODUCTION, "logger")
}

func session(connection string) (configData model.Config, localhost *mgo.Session, production *mgo.Session) {
	configData = config.GetConfig()
	if os.Args[1] == "--run" {
		if configData.Databaselocal.Host != "" && connection == LOCALHOST {
			localhost, _ = mgo.DialWithInfo(&mgo.DialInfo{
				Addrs:    []string{configData.Databaselocal.Host},
				Username: configData.Databaselocal.Username,
				Password: configData.Databaselocal.Password,
				Database: configData.Databaselocal.Database,
			})
		}

		if configData.Databasecentralized.Host != "" && connection == PRODUCTION {
			production, _ = mgo.DialWithInfo(&mgo.DialInfo{
				Addrs:    []string{configData.Databasecentralized.Host},
				Username: configData.Databasecentralized.Username,
				Password: configData.Databasecentralized.Password,
				Database: configData.Databasecentralized.Database,
			})
		}
	}
	return configData, localhost, production
}

func GetSession(connection string, collection string) *mgo.Collection {
	coll :=  &mgo.Collection{}
	fmt.Println(connection)
	if os.Args[1] == "--run" {
		configData, localhost, production := session(connection)

		if connection == LOCALHOST && configData.Databaselocal.Host != "" {
			return localhost.DB(configData.Databaselocal.Database).C(collection)
		}
		if connection == PRODUCTION && configData.Databasecentralized.Host != "" {
			return production.DB(configData.Databasecentralized.Database).C(collection)
		}
	}
	return coll
}
// Logres - Distributed logs system PostgresSQL to MongoDB
// https://github.com/luk4z7/logres for the canonical source repository
//
// Copyright 2017 The Lucas Alves Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package logger is a layer with methods for persistence and operations with data
package logger

import (
	"github.com/luk4z7/logres/service/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Persist provide the insert on mongo client/production
func Persist(session *mgo.Collection, register model.LoggerPostgreSQL) error {
	if err := session.Insert(register); err != nil {
		return err
	}
	return nil
}

// DeletePerObjectId provide the delete of objects per ObjectId on mongo client/production
func DeletePerObjectId(session *mgo.Collection, id bson.ObjectId) error {
	if err := session.Remove(bson.M{"_id": id}); err != nil {
		return err
	}
	return nil
}

// GetAll provide the query of all data on database on mongo client/production
func GetAll(session *mgo.Collection) ([]model.LoggerPostgreSQL, error) {
	result := []model.LoggerPostgreSQL{}
	err := session.Find(nil).All(&result)
	return result, err
}

// Logres - Distributed logs system PostgresSQL to MongoDB
// https://github.com/luk4z7/logres for the canonical source repository
//
// Copyright 2017 The Lucas Alves Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// service
package model

import "gopkg.in/mgo.v2/bson"

type LoggerPostgreSQL struct {
	LogTime              string        `bson:"log_time"`
	UserName             string        `bson:"user_name"`
	DatabaseName         string        `bson:"database_name"`
	ProcessID            string        `bson:"process_id"`
	ConnectionFrom       string        `bson:"connection_from"`
	SessionID            string        `bson:"session_id"`
	SessionLineNum       string        `bson:"session_line_num"`
	CommandTag           string        `bson:"command_tag"`
	SessionStartTime     string        `bson:"session_start_time"`
	VirtualTransactionID string        `bson:"virtual_transaction_id"`
	TransactionID        string        `bson:"transaction_id"`
	ErrorSeverity        string        `bson:"error_severity"`
	SqlStateCode         string        `bson:"sql_state_code"`
	Message              string        `bson:"message"`
	Detail               string        `bson:"detail"`
	Hint                 string        `bson:"hint"`
	InternalQuery        string        `bson:"internal_query"`
	InternalQueryPos     string        `bson:"internal_query_pos"`
	Context              string        `bson:"context"`
	Query                string        `bson:"query"`
	QueryPos             string        `bson:"query_pos"`
	Location             string        `bson:"location"`
	ApplicationName      string        `bson:"application_name"`
	ID                   bson.ObjectId `bson:"_id,omitempty"`
	Client               string        `bson:"client"`
}

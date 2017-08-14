// Logres - Distributed logs system PostgresSQL to MongoDB
// https://github.com/luk4z7/logres for the canonical source repository
//
// Copyright 2017 The Lucas Alves Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// lib
package error

import (
	"fmt"
	"runtime"
)

type ErrorsAPI struct {
	Errors []Errors `json:"errors"`
	Url    string   `json:"url"`
	Method string   `json:"method"`
}

type Errors struct {
	ParameterName string `json:"parameter_name"`
	Type          string `json:"type"`
	Message       string `json:"message"`
}

type Err struct {
	Name string
}

func (e *Err) Error() string {
	return e.Name
}

// Return panic when error != nil
func Check(e error, m string) {
	if e != nil {
		if m == "" {
			panic(m)
		} else {
			panic(e)
		}
	}
}

func CatchPanic(err *error, functionName string) {
	if r := recover(); r != nil {
		fmt.Printf("%s : PANIC Defered : %v\n", functionName, r)

		// capture the stack trace
		buf := make([]byte, 10000)
		runtime.Stack(buf, false)

		fmt.Printf("%s : Stack Trace : %s", functionName, string(buf))

		if err != nil {
			*err = fmt.Errorf("%v", r)
		}
	} else if err != nil && *err != nil {
		fmt.Printf("%s : ERROR : %v\n", functionName, *err)

		// Capture the stack trace
		buf := make([]byte, 10000)
		runtime.Stack(buf, false)

		fmt.Printf("%s : Stack Trace : %s", functionName, string(buf))
	}
}

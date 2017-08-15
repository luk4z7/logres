// Logres - Distributed logs system PostgresSQL to MongoDB
// https://github.com/luk4z7/logres for the canonical source repository
//
// Copyright 2017 The Lucas Alves Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package logs return instance of gocolorize
package logs

import (
	"github.com/agtorre/gocolorize"
	"log"
	"os"
)

var (
	INFO     *log.Logger // Green log
	WARNING  *log.Logger // yellow log
	CRITICAL *log.Logger // red log
)

// Start initiate the instances of gocolorize
func Start() {
	info := gocolorize.NewColor("green")
	warning := gocolorize.NewColor("yellow")
	critical := gocolorize.NewColor("back+u:red")

	//helper functions to shorten code
	i := info.Paint
	w := warning.Paint
	c := critical.Paint

	INFO = log.New(os.Stdout, i("INFO "), log.Ldate|log.Lmicroseconds|log.Lshortfile)
	WARNING = log.New(os.Stdout, w("WARNING "), log.Ldate|log.Lmicroseconds|log.Lshortfile)
	CRITICAL = log.New(os.Stdout, c("CRITICAL "), log.Ldate|log.Lmicroseconds|log.Lshortfile)
}

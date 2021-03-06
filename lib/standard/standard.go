// Logres - Distributed logs system PostgresSQL to MongoDB
// https://github.com/luk4z7/logres for the canonical source repository
//
// Copyright 2017 The Lucas Alves Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package standard provide helpers methods
package standard

import "os/user"

// InArray check if exist val in array
func InArray(val string, array []string) (exists bool, index int) {
	exists = false
	index = -1
	for i, v := range array {
		if val == v {
			index = i
			exists = true
			return
		}
	}
	return
}

// GetCurrentUser return current user
func GetCurrentUser() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}
	return u.Username, err
}

// IsAlpha check if is character is alpha
func IsAlpha(r string) bool {
	str := []rune(r)
	if len(str) < 1 {
		return false
	}
	return int(str[0]) >= 48 && int(str[0]) <= 57
}

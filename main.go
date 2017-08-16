// Logres - Distributed logs system PostgresSQL to MongoDB
// https://github.com/luk4z7/logres for the canonical source repository
//
// Copyright 2017 The Lucas Alves Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/luk4z7/logres/drive/mongo"
	"github.com/luk4z7/logres/lib/logs"
	"github.com/luk4z7/logres/lib/standard"
	"github.com/luk4z7/logres/service/config"
	"github.com/luk4z7/logres/service/logger"
	"github.com/luk4z7/logres/service/model"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	workers     = runtime.NumCPU()
	local, prod = mongo.New()
	store       = Store{}
	seek        int64
	filename    string
)

// Store store file transactions and mutex for operatiosn on various
// goroutines
type Store struct {
	File        string
	Transaction []string
	sync.Mutex
}

func main() {
	if len(os.Args) == 1 {
		fmt.Printf("usage: %s --[option]\n", filepath.Base(os.Args[0]))
		fmt.Println("To see help text, you can run:")
		fmt.Println()
		fmt.Println("  logres --help")
		os.Exit(2)
	}
	if len(os.Args) > 1 && os.Args[1] == "--help" {
		fmt.Println("Configuration")
		fmt.Println(" logres --config")
		fmt.Println()
		fmt.Println("Running")
		fmt.Println(" logres --run")
		os.Exit(2)
	}
	if len(os.Args) == 2 {
		if os.Args[1] == "--config" {
			config.GenerateConfig()
		}
	}
	if len(os.Args) > 1 && os.Args[1] == "--run" {
		_, err := config.CheckConfig()
		if err != nil {
			os.Exit(1)
		}
		run()
	}
	if len(os.Args) > 1 && os.Args[1] == "--server-http" {
		_, err := config.CheckConfig()
		if err != nil {
			os.Exit(1)
		}
		serverHttp()
	}
}

// Getting configuration and running all methods
func run() {
	logs.Start()

	// Send all data for the centralized database
	go store.push()
	store.Lock()
	defer store.Unlock()

	// Creating the listener
	configData := config.GetConfig()
	watcher(configData)
}

// The watcher() method implements the algorithm used from the fsnotify project
// link: https://github.com/fsnotify/fsnotify (File System Notifications to Go)
// a listener is created for a specific directory when a file is modified
// is triggered the method readLines()
func watcher(configModel model.Config) {
	// Set the client variable
	config.Client = configModel.Client.Name

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan struct{})
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					logs.INFO.Println("Modified file -> ", event.Name)
					// When the file name has not been defined, it is time to
					// use the SetFile() method to add a new file to read.
					if filename == "" {
						store.SetFile(event.Name)
						filename = event.Name
					}
					if filename != "" && filename != event.Name {
						logs.INFO.Println("Reset seek")
						seek = 0
					}
					readLines(event.Name)
				}
			case err := <-watcher.Errors:
				logs.CRITICAL.Println("Error on watcher: ", err)
			}
		}
	}()
	err = watcher.Add(configModel.Pathlog.Name)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

// The readLines() method uses ReadBytes() to read the bytes that are
// read in the file, some are converted to string where they are read
// back by csv.NewReader(), which makes the string parser for the struct
// LoggerPostgreSQL{}, done The Unmarshal data is persisted in the local mongodb.
func readLines(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		logs.CRITICAL.Println("Failed to open the file -> ", err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	if _, err := file.Seek(seek, 0); err != nil {
		logs.CRITICAL.Println("Failed to seek the file -> ", err)
	}
	pos := seek
	for {
		data, err := reader.ReadBytes('\n')
		line := string(data)
		// sum "pos" more total of bytes of "data"
		pos += int64(len(data))
		// the end of file reading
		// setter the value of the bytes for "seek" about sum of "pos"
		if err == io.EOF {
			seek = pos
			break
		}
		r := csv.NewReader(strings.NewReader(line))
		structure := model.LoggerPostgreSQL{}
		_ = config.Unmarshal(r, &structure)
		if structure.LogTime != "" {
			logs.INFO.Println("SCAN -> UserName: " +
				structure.UserName + " :: DataBase: " +
				structure.DatabaseName + " :: VirtualTransactionID: " +
				structure.VirtualTransactionID)

			structure.Client = config.Client
			err = logger.Persist(local, structure)
			if err != nil {
				logs.CRITICAL.Println("Failed to persist structure in the mongodb client -> ", err)
			}
		}
		if err != nil {
			if err != io.EOF {
				logs.CRITICAL.Println("Failed to finish reading the file -> ", err)
			}
			break
		}
	}
}

// The removeLines() method is called by the sync() method every time
// no value is returned of the query of the client database (mongodb),
// the method is used to remove the data from the file, avoiding
// re-reading of the same data.
// The method receives store.getFile() that has the pointer to get
// the file being monitored, always starts from row 1 and in the
// third parameter is passed a negative value, that indicates that
// the number of rows must be queried by the method store.getLines()
// The file is read its contents and through the parameters passed
// using the method skip() to obtain the data that must be escaped
// for writing byte in the file.
func removeLines(fn string, start, n int) (err error) {
	logs.INFO.Println("Clear file -> ", fn)
	if n < 0 {
		n = store.getLines()
	}
	if n == 0 {
		logs.INFO.Println("Nothing to clear")
		seek = 0
		return nil
	}
	logs.INFO.Println("Total lines -> ", n)
	if start < 1 {
		logs.WARNING.Println("Invalid request.  line numbers start at 1.")
	}
	var f *os.File
	if f, err = os.OpenFile(fn, os.O_RDWR, 0); err != nil {
		logs.CRITICAL.Println("Failed to open the file -> ", err)
		return
	}
	defer func() {
		if cErr := f.Close(); err == nil {
			err = cErr
		}
	}()
	var b []byte
	if b, err = ioutil.ReadAll(f); err != nil {
		logs.CRITICAL.Println("Failed to reading the file -> ", err)
		return
	}
	cut, ok := skip(b, start-1)
	if !ok {
		logs.CRITICAL.Printf("less than %d lines -> ", start)
		return
	}
	if n == 0 {
		return nil
	}
	tail, ok := skip(cut, n)
	if !ok {
		logs.CRITICAL.Printf("less than %d lines after line %d ", n, start)
		return
	}
	t := int64(len(b) - len(cut))
	if err = f.Truncate(t); err != nil {
		return
	}
	// Writing in the archive the bytes already with cut removed
	if len(tail) > 0 {
		_, err = f.WriteAt(tail, t)
	}
	return
}

// Get total of line in the log file used into current verification
// through method getFile().
// using MaxScanTokenSize what is the maximum size used to buffer a token
// Scan the file and iterate the variable "counter" and returns it
func (s *Store) getLines() int {
	file, err := os.Open(s.getFile())
	if err != nil {
		panic(err)
	}
	defer file.Close()
	inputReader := bufio.NewReader(file)
	var scanner *bufio.Scanner
	buffer := make([]byte, bufio.MaxScanTokenSize)
	counter := 0
	for scanner == nil || scanner.Err() == bufio.ErrTooLong {
		scanner = bufio.NewScanner(inputReader)
		scanner.Buffer(buffer, 0)
		for scanner.Scan() {
			counter++
		}
	}
	return counter
}

// Push the data of the mongodb client for the mongodb server
// every 15 minutes
func (s *Store) push() {
	for {
		select {
		case <-time.After(time.Second * 15):
			logs.INFO.Println("Initiating push....")
			go s.sync()
		}
	}
}

// Every 15 minutes this method is executed for synchronization
// between the client and server.
// Verifies that there are records in the client database, iterate it and
// push for database centralized.
// At the same time removes data from the mongodb client.
// When the variable "total" is equal zero the method removeLines()
// is activated for clearing of the log file, also store.Transaction
// is set with value null because the file is clean and the
// store.Transaction it should be like this too.
func (s *Store) sync() {
	result, err := logger.GetAll(local)
	if err != nil {
		logs.CRITICAL.Println("Panic for get all objects")
	}
	total := len(result)
	if total != 0 {
		logs.INFO.Println("Total of records: ", total)
	}
	if total == 0 {
		log.Println("Nothing to sync")
		if store.getFile() != "" {
			logs.INFO.Println("Initiating clearing....")
			removeLines(store.getFile(), 1, -1)
		}
		store.Transaction = nil
		return
	}
	for i := 0; i < total; i++ {
		exists, _ := standard.InArray(result[i].VirtualTransactionID, store.Transaction)
		if !exists {
			logs.INFO.Println("PUSH -> UserName: " +
				result[i].UserName + " :: DataBase: " +
				result[i].DatabaseName + " :: VirtualTransactionID: " +
				result[i].VirtualTransactionID)

			store.Transaction = append(store.Transaction, result[i].VirtualTransactionID)
			logger.Persist(prod, result[i])
		}
		// Depending on the amount of data and traffic, goroutines that were
		// first run have already removed the registry, not identifying the
		// registry in the database at the current execution.
		err := logger.DeletePerObjectId(local, result[i].ID)
		if err != nil {
			logs.INFO.Println("ObjectId -> " + result[i].ID.Hex() + " removed on the last goroutine")
		}
	}
}

// SetFile set a file for struct Store
func (s *Store) SetFile(file string) {
	s.File = file
}

// getFile get a file of struct Store
func (s *Store) getFile() string {
	return s.File
}

func skip(b []byte, n int) ([]byte, bool) {
	// checks if "n" is greater than zero and decrements the start value
	for ; n > 0; n-- {
		if len(b) == 0 {
			return nil, false
		}
		// IndexByte returns the index of the first instance of c in s, or -1 if c is not present in s.
		x := bytes.IndexByte(b, '\n')
		// Checks if the number of bytes is less than zero
		// -1 when it is returned from IndexByte when not found c '\n' in s []byte
		if x < 0 {
			// Set the total of bytes for x
			x = len(b)
		} else {
			// Iterates the index value returned by IndexByte
			x++
		}
		// Returns b by escaping a few bytes.
		b = b[x:]
	}
	return b, true
}

func serverHttp() {}

// Logres - Distributed logs system PostgresSQL to MongoDB
// https://github.com/luk4z7/logres for the canonical source repository
//
// Copyright 2017 The Lucas Alves Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// config
package config

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/fatih/structs"
	liberr "github.com/luk4z7/logres/lib/error"
	"github.com/luk4z7/logres/lib/standard"
	"github.com/luk4z7/logres/service/model"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

const (
	fileConfig = ".logres.yaml"
)

// GetOS return the operating system in use
func GetOS() string {
	return runtime.GOOS
}

// GetPathConfigFile only darwin and linux are running the configuration
func GetPathConfigFile() string {
	var file string
	switch GetOS() {
	case "darwin":
		file = "/var/root/" + fileConfig
	case "linux":
		file = "/root/" + fileConfig
	}
	return file
}

// GetConfig using lib yaml.v2 for read the yaml file
// to get the configuration of hosts and return struct Config{}
func GetConfig() model.Config {
	_, err := CheckConfig()
	if err != nil {
		os.Exit(1)
	}
	file := GetPathConfigFile()
	configData := model.Config{}
	openFile, err := os.OpenFile(file, os.O_RDONLY, 0600)
	if err != nil {
		panic(err)
	}
	raw, err := ioutil.ReadAll(openFile)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(raw, &configData)
	if err != nil {
		panic(err)
	}
	return configData
}

// CheckConfig check if exists the config file
func CheckConfig() (*os.File, error) {
	file := GetPathConfigFile()
	openFile, err := os.OpenFile(file, os.O_RDONLY, 0600)
	if err != nil {
		fmt.Println("Unable to locate configure. You can configure by running \"logres --config\"")
	}
	return openFile, err
}

// CreateFileConfig create config file and put the byte data
func CreateFileConfig(data []byte) error {
	file := GetPathConfigFile()
	newFile, err := os.Create(file)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	_, err = newFile.Write(data)
	if err != nil {
		return err
	}
	return nil
}

// GenerateConfig generates the configuration file from some questions that are answered by the user
func GenerateConfig() {
	config := model.Config{}

askLocallyAndCentralizedAgain:
	fmt.Print("Would you like configure the database centralized? [Y/n] ")

	readerOne := bufio.NewReader(os.Stdin)
	askOne, _ := readerOne.ReadString('\n')
	askOne = strings.Trim(askOne, "\n")

	if askOne != "" && (askOne == "y" || askOne == "Y") {
	askCentralized:
		fmt.Println("Enter the value, or press ENTER for the anything")

		host1 := bufio.NewReader(os.Stdin)
		fmt.Print("   Host: ")
		config.Databasecentralized.Host, _ = host1.ReadString('\n')
		config.Databasecentralized.Host = strings.Trim(config.Databasecentralized.Host, "\n")

		username := bufio.NewReader(os.Stdin)
		fmt.Print("   Username: ")
		config.Databasecentralized.Username, _ = username.ReadString('\n')
		config.Databasecentralized.Username = strings.Trim(config.Databasecentralized.Username, "\n")

		password := bufio.NewReader(os.Stdin)
		fmt.Print("   Password: ")
		config.Databasecentralized.Password, _ = password.ReadString('\n')
		config.Databasecentralized.Password = strings.Trim(config.Databasecentralized.Password, "\n")

		database := bufio.NewReader(os.Stdin)
		fmt.Print("   Database: ")
		config.Databasecentralized.Database, _ = database.ReadString('\n')
		config.Databasecentralized.Database = strings.Trim(config.Databasecentralized.Database, "\n")

		field, err := IsNotEmpty(
			&config,
			[]string{
				"Databasecentralized",
			},
			func() map[string][]string {
				return map[string][]string{
					"Databasecentralized": {
						"Host",
						"Username",
						"Database",
					},
				}
			},
		)
		if err != nil {
			fmt.Println("field " + field + " is empty")
		}
		if config.Databasecentralized.Host == "" ||
			config.Databasecentralized.Username == "" ||
			config.Databasecentralized.Database == "" {
			fmt.Print("The information above database centralized is empty, Try again? [Y/n] ")
			centralizedEmpty := bufio.NewReader(os.Stdin)
			centralizedEmptyAnswer, _ := centralizedEmpty.ReadString('\n')
			centralizedEmptyAnswer = strings.Trim(centralizedEmptyAnswer, "\n")
			if centralizedEmptyAnswer != "" && (centralizedEmptyAnswer == "y" || centralizedEmptyAnswer == "Y") {
				goto askCentralized
			}
		}
	}
	fmt.Print("Would you like configure the database locally? [Y/n] ")

	readerTwo := bufio.NewReader(os.Stdin)
	askTwo, _ := readerTwo.ReadString('\n')
	askTwo = strings.Trim(askTwo, "\n")

	if askTwo != "" && (askTwo == "y" || askTwo == "Y") {
	askLocally:
		fmt.Println("Enter the value, or press ENTER for the anything")

		host2 := bufio.NewReader(os.Stdin)
		fmt.Print("   Host: ")
		config.Databaselocal.Host, _ = host2.ReadString('\n')
		config.Databaselocal.Host = strings.Trim(config.Databaselocal.Host, "\n")

		username2 := bufio.NewReader(os.Stdin)
		fmt.Print("   Username: ")
		config.Databaselocal.Username, _ = username2.ReadString('\n')
		config.Databaselocal.Username = strings.Trim(config.Databaselocal.Username, "\n")

		password2 := bufio.NewReader(os.Stdin)
		fmt.Print("   Password: ")
		config.Databaselocal.Password, _ = password2.ReadString('\n')
		config.Databaselocal.Password = strings.Trim(config.Databaselocal.Password, "\n")

		database2 := bufio.NewReader(os.Stdin)
		fmt.Print("   Database: ")
		config.Databaselocal.Database, _ = database2.ReadString('\n')
		config.Databaselocal.Database = strings.Trim(config.Databaselocal.Database, "\n")

		field, err := IsNotEmpty(
			&config,
			[]string{
				"Databaselocal",
			},
			func() map[string][]string {
				return map[string][]string{
					"Databaselocal": {
						"Host",
						"Username",
						"Database",
					},
				}
			},
		)
		if err != nil {
			fmt.Println("field " + field + " is empty")
		}
		if config.Databaselocal.Host == "" ||
			config.Databaselocal.Username == "" ||
			config.Databaselocal.Database == "" {
			fmt.Print("The information above database locally is empty, Try again? [Y/n] ")
			locallyEmpty := bufio.NewReader(os.Stdin)
			locallyEmptyAnswer, _ := locallyEmpty.ReadString('\n')
			locallyEmptyAnswer = strings.Trim(locallyEmptyAnswer, "\n")
			if locallyEmptyAnswer != "" && (locallyEmptyAnswer == "y" || locallyEmptyAnswer == "Y") {
				goto askLocally
			}
		}
	}
	_, errCentralized := IsNotEmpty(
		&config,
		[]string{
			"Databasecentralized",
		},
		func() map[string][]string {
			return map[string][]string{
				"Databasecentralized": {
					"Host",
					"Username",
					"Database",
				},
			}
		},
	)
	_ , errLocally := IsNotEmpty(
		&config,
		[]string{
			"Databaselocal",
		},
		func() map[string][]string {
			return map[string][]string{
				"Databaselocal": {
					"Host",
					"Username",
					"Database",
				},
			}
		},
	)
	if errCentralized != nil && errLocally != nil {
		fmt.Print("The information above database locally and centralized is empty, Try again? [Y/n] ")
		readerLocallyAndCentralized := bufio.NewReader(os.Stdin)
		locallyAndCentralizedEmptyAnswer, _ := readerLocallyAndCentralized.ReadString('\n')
		locallyAndCentralizedEmptyAnswer = strings.Trim(locallyAndCentralizedEmptyAnswer, "\n")

		if locallyAndCentralizedEmptyAnswer != "" && (locallyAndCentralizedEmptyAnswer == "y" ||
			locallyAndCentralizedEmptyAnswer == "Y") {
			goto askLocallyAndCentralizedAgain
		}
	}
	fmt.Print("Enter name for this server: ")
	client := bufio.NewReader(os.Stdin)
	config.Client.Name, _ = client.ReadString('\n')
	config.Client.Name = strings.Trim(config.Client.Name, "\n")

askPathLog:
	fmt.Print("Enter directory for logs scan: ")
	readerPath := bufio.NewReader(os.Stdin)
	config.Pathlog.Name, _ = readerPath.ReadString('\n')
	config.Pathlog.Name = strings.Trim(config.Pathlog.Name, "\n")

	if config.Pathlog.Name == "" {
		fmt.Print("The information above Path log is empty, Try again? [Y/n] ")
		pathlogEmpty := bufio.NewReader(os.Stdin)
		pathlogEmptyAnswer, _ := pathlogEmpty.ReadString('\n')
		pathlogEmptyAnswer = strings.Trim(pathlogEmptyAnswer, "\n")
		if pathlogEmptyAnswer != "" && (pathlogEmptyAnswer == "y" || pathlogEmptyAnswer == "Y") {
			goto askPathLog
		}
	}
	if config.Pathlog.Name != "" {
		if !IsDirectoryExists(config.Pathlog.Name) {
			fmt.Printf("the %s folder doesn't exist \n", config.Pathlog.Name)
			os.Exit(1)
		}
	}
	if config.Pathlog.Name == "" &&
		config.Databaselocal.Host == "" &&
		config.Databaselocal.Username == "" &&
		config.Databaselocal.Database == "" &&
		config.Databaselocal.Password == "" &&
		config.Databasecentralized.Host == "" &&
		config.Databasecentralized.Username == "" &&
		config.Databasecentralized.Database == "" &&
		config.Databasecentralized.Password == "" &&
		config.Client.Name == "" {
		fmt.Println("Nothing to change!")
		os.Exit(1)
	}

	d, err := yaml.Marshal(&config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	err = CreateFileConfig(d)
	if err != nil {
		fmt.Println("Unable to configure file")
		os.Exit(1)
	}
	fmt.Println("File configured!")
}

// IsDirectoryExists verify is directory exists
func IsDirectoryExists(path string) bool {
	if _, err := os.Stat(strings.Trim(path, "\n")); os.IsNotExist(err) {
		return false
	}
	return true
}

// IsNotEmpty using the closure MustBeNotEmpty for check if exist the data on the struct Config
func IsNotEmpty(config *model.Config, Item []string, subItem func() map[string][]string) (string, error) {
	field, err := MustBeNotEmpty(config, func() []string {
		return Item
	}, subItem)
	return field, err
}

// MustBeNotEmpty check the values the map passed
func MustBeNotEmpty(v interface{}, require func() []string, sub func() map[string][]string) (string, error) {
	required := require()
	subitems := sub()

	s := structs.New(v)
	for _, v := range required {
		name := s.Field(v)
		for key, values := range subitems {
			if key == v {
				for i := 0; i < len(values); i++ {
					if name.Field(values[i]).Value() == "" {
						return values[i], &liberr.Err{Name: values[i] + " - Parametro incorreto"}
					}
				}
			}
		}
		if name.Kind() == reflect.String {
			value := name.Value().(string)
			if !standard.IsAlpha(value) || value == "" {
				return v, &liberr.Err{Name: v + " - Parametro incorreto"}
			}
		}
	}
	return "", nil
}

// FieldMismatch
type FieldMismatch struct {
	expected, found int
}

func (e *FieldMismatch) Error() string {
	return "CSV line fields mismatch. Expected " + strconv.Itoa(e.expected) + " found " + strconv.Itoa(e.found)
}

// UnsupportedType
type UnsupportedType struct {
	Type string
}

func (e *UnsupportedType) Error() string {
	return "Unsupported type: " + e.Type
}

// Unmarshal check data on the csv string
func Unmarshal(reader *csv.Reader, v interface{}) error {
	record, err := reader.Read()
	if err != nil {
		return err
	}
	s := reflect.ValueOf(v).Elem()
	if s.NumField()-2 != len(record) {
		return &FieldMismatch{s.NumField(), len(record)}
	}
	for i := 0; i < s.NumField()-2; i++ {
		f := s.Field(i)
		switch f.Type().String() {
		case "string":
			f.SetString(record[i])
		case "bson.ObjectId":
		case "int":
			ival, err := strconv.ParseInt(record[i], 10, 0)
			if err != nil {
				return err
			}
			f.SetInt(ival)
		default:
			return &UnsupportedType{f.Type().String()}
		}
	}
	return nil
}

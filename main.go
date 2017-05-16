package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type Config struct {
	URLs      map[string]int `json:"urls"`
	Name      string         `json:"name"`
	FromEmail string         `json:"from-email"`
	ToEmails  []string       `json:"to-emails"`
	Password  string         `json:"password"`
	Host      string         `json:"host"`
	Port      string         `json:"port"`
}

var (
	config = Config{}
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: " + os.Args[0] + " /path/to/configig/file")
		os.Exit(1)
	}
	loadConfig(os.Args[1])
	sitesToNotifyAbout := ""
	for url, expectedStatus := range config.URLs {
		actualStatus, err := getHTTPStatus(url)
		if err != nil {
			sitesToNotifyAbout += err.Error() + " -- "
		} else if actualStatus != expectedStatus {
			sitesToNotifyAbout += url + " Expected: " + strconv.Itoa(expectedStatus) + " Actual: " + strconv.Itoa(actualStatus) + "\n"
		}
	}
	if sitesToNotifyAbout != "" {
		for _, email := range config.ToEmails {
			sendEmail(email, "Site Monitor - Attention", "Affected site(s):\n"+sitesToNotifyAbout)
		}
	}
}

func loadConfig(filePath string) {
	configData, error := ioutil.ReadFile(filePath)
	panicOnError(error)
	error = json.Unmarshal(configData, &config)
	panicOnError(error)
}

func getHTTPStatus(url string) (status int, err error) {
	resp, err := http.Head(url)
	if err != nil {
		return
	}
	status = resp.StatusCode
	return
}

func panicOnError(error error) {
	if error != nil {
		panic(error)
	}
}

func debug(things ...interface{}) {
	for _, thing := range things {
		fmt.Printf("%+v\n", thing)
	}
	fmt.Println("^^^^^^^^^^^^^^^^^^^^")
}

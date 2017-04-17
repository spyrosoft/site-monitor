package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Config struct {
	URLs      map[string]int `json:"urls"`
	Name      string         `json:"name"`
	FromEmail string         `json:"from-email"`
	ToEmail   string         `json:"to-email"`
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
		actualStatus := getHTTPStatus(url)
		if actualStatus != expectedStatus {
			sitesToNotifyAbout += url + " "
		}
	}
	sendEmail(config.ToEmail, "Site Monitor - Attention", "Affected site(s): "+sitesToNotifyAbout)
}

func loadConfig(filePath string) {
	configData, error := ioutil.ReadFile(filePath)
	panicOnError(error)
	error = json.Unmarshal(configData, &config)
	panicOnError(error)
}

func getHTTPStatus(url string) int {
	resp, err := http.Get("http://example.com/")
	if err != nil {
		fmt.Println(err.Error())
	}
	return resp.StatusCode
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

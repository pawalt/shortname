package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"net/http"

	"github.com/mitchellh/go-homedir"
	"github.com/txn2/txeh"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Sites map[string]string `yaml:"sites"`
}

var (
	hostConfig = Config{}
)

func main() {
	timeConfig()
	http.HandleFunc("/", root)
	err := http.ListenAndServe(":80", nil)
	handleErr(err)
}

func timeConfig() {
	reloadConfig()
	reloadHosts()
	ticker := time.NewTicker(60 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				reloadConfig()
				reloadHosts()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func reloadHosts() {
	hosts, err := txeh.NewHostsDefault()
	handleErr(err)

	hosts.RemoveAddress("127.0.0.1")
	hosts.AddHost("localhost", "127.0.0.1")
	for key := range hostConfig.Sites {
		hosts.AddHost("127.0.0.1", key)
	}

	hosts.Save()
}

func reloadConfig() {
	hd, err := homedir.Dir()
	handleErr(err)
	configLocation := filepath.Join(hd, ".shortnamerc")

	if _, err := os.Stat(configLocation); os.IsNotExist(err) {
		ioutil.WriteFile(configLocation, []byte("sites: {}\n"), 0755)
	}

	configContents, err := ioutil.ReadFile(configLocation)
	handleErr(err)

	err = yaml.Unmarshal(configContents, &hostConfig)
	handleErr(err)
}

func root(w http.ResponseWriter, req *http.Request) {
	if url, ok := hostConfig.Sites[req.Host]; ok {
		http.Redirect(w, req, url, http.StatusMovedPermanently)
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "could not find hostname "+req.Host)
	}
}

func handleErr(e error) {
	if e != nil {
		panic(e)
	}
}

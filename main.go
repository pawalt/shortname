package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"
	"flag"

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
	cfgFile string
)

func main() {
	hd, err := homedir.Dir()
	handleErr(err)
	configLocation := filepath.Join(hd, ".shortnamerc")
	flag.StringVar(&cfgFile, "f", configLocation, "file to load in config from")
	flag.Parse()

	timeConfig()
	http.HandleFunc("/", root)
	log.Println("Starting shortname HTTP server")
	err = http.ListenAndServe(":80", nil)
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
	hosts.AddHost("127.0.0.1", "localhost")
	for key := range hostConfig.Sites {
		hosts.AddHost("127.0.0.1", key)
	}

	hosts.Save()
}

func reloadConfig() {
	if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
		ioutil.WriteFile(cfgFile, []byte("sites: {}\n"), 0755)
	}

	configContents, err := ioutil.ReadFile(cfgFile)
	handleErr(err)

	err = yaml.Unmarshal(configContents, &hostConfig)
	handleErr(err)

	hostConfig.Sites["sn"] = ""
}

func root(w http.ResponseWriter, req *http.Request) {
	if req.Host == "sn" {
		strSites, _ := json.MarshalIndent(hostConfig, "", "  ")
		fmt.Fprintf(w, string(strSites))
	} else if url, ok := hostConfig.Sites[req.Host]; ok {
		http.Redirect(w, req, path.Join(url, req.URL.Path), http.StatusMovedPermanently)
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

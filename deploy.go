package main

import "os"
import "github.com/bgentry/heroku-go"
import hk "github.com/heroku/hk/hkclient"

const (
	PLUGIN_NAME    = "deploy"
	PLUGIN_VERSION = 1
	// PLUGIN_USER_AGENT = "hk-" + PLUGIN_NAME "/1"
)

var client *heroku.Client
var nrc *hk.NetRc

func help() {}

func init() {
	nrc, err := hk.LoadNetRc()
	if err != nil && os.IsNotExist(err) {
		nrc = &hk.NetRc{}
	}

	clients, err := hk.New(nrc, "TODO user agent")

	if err == nil {
		client = clients.Client
	} else {
		// TODO
	}
}

func main() {
	if os.Getenv("HKPLUGINMODE") == "info" {
		help()
		os.Exit(0)
	}

	if len(os.Args) < 2 {
		help()
		os.Exit(1)
	}

	/*
		TODO:
			* Check that we have an APP context or set it
			* tarzip current directory or git root
			* ideally with gitignoring, if possible
			* upload to S3 with an object expiry of ~5min
			* hit build API with that link
			* tail output (if build api has implemented that)
	*/
}

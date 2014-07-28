package main

import (
	"fmt"
	"github.com/bgentry/heroku-go"
	hk "github.com/heroku/hk/hkclient"
	"os"
	"path/filepath"
)

const (
	PLUGIN_NAME    = "deploy"
	PLUGIN_VERSION = 1
	// PLUGIN_USER_AGENT = "hk-" + PLUGIN_NAME "/1"
	ENDPOINT = "https://hk-deploy.herokuapp.com/slot"
)

var client *heroku.Client
var nrc *hk.NetRc

func help() {
	fmt.Println(`hk deploy: Deploy a directory of code to Heroku using the Build API.

Run "hk deploy DIRECTORY" to deploy the specified directory to Heroku.`)
}

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

	if os.Args[1] == "-h" || os.Args[1] == "--help" {
		help()
		os.Exit(0)
	}

	dir := os.Args[1] // TODO: Maybe fallback to CWD or Git root?

	fullPath, _ := filepath.Abs(dir)
	fmt.Printf("Creating .tgz of %s...\n", fullPath)
	tgz := buildTgz(dir)
	fmt.Printf("done (%d bytes)\n", tgz.Len())

	fmt.Print("Requesting upload slot... ")
	slot, err := getUploadSlot()
	if err == nil {
		fmt.Println("done")
	} else {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Print("Uploading .tgz to S3... ")
	if err := upload(&tgz, slot); err == nil {
		fmt.Println("done")
	} else {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// working downloadable link now in slot.DownloadUrl
	fmt.Println("Submitting build with download link... not implemented")
	fmt.Println("Commenting build... not implemented")
}

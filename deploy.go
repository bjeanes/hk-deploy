package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/naaman/hbuild"
)

const (
	PLUGIN_NAME    = "deploy"
	PLUGIN_VERSION = 1
	ENDPOINT       = "https://hk-deploy.herokuapp.com/slot"
	INFO_PREAMBLE  = `%s %d: Deploy code to Heroku using the API`
	HELP_TEXT      = `Usage: hk deploy DIRECTORY

	Deploy the specified directory to Heroku`
)

func help() {
	fmt.Println(HELP_TEXT)
}

func info() {
	fmt.Printf(INFO_PREAMBLE+"\n\n", PLUGIN_NAME, PLUGIN_VERSION)
	help()
	os.Exit(0)
}

func main() {
	if os.Getenv("HKPLUGINMODE") == "info" {
		info()
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

	var source hbuild.Source
	var build hbuild.Build
	var err error

	app := os.Getenv("HKAPP")
	apiKey := os.Getenv("HKPASS")

	fmt.Println(apiKey)

	fmt.Print("Creating source...")
	source, err = hbuild.NewSource(apiKey, app, fullPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("done.")

	fmt.Print("Compressing source...")
	err = source.Compress()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("done.")

	fmt.Print("Uploading source...")
	err = source.Upload()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("done.")

	fmt.Println("Building:")
	build, err = hbuild.NewBuild(apiKey, app, source)
	if err != nil {
		fmt.Println(err)
		return
	}

	io.Copy(os.Stdout, build.Output)
}

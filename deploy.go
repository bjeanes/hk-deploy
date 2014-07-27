package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/bgentry/heroku-go"
	hk "github.com/heroku/hk/hkclient"
	// "io"
	"io/ioutil"
	"os"
	"path/filepath"
)

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

func shouldIgnore(path string) bool {
	// TODO: gitignore-ish rules, if a .gitignore exists?
	return path == ".git"
}

func buildTgz(root string) bytes.Buffer {
	buf := new(bytes.Buffer)
	gz := gzip.NewWriter(buf)
	tw := tar.NewWriter(gz)

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		// TODO: handle incoming err more meaningfully
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		if shouldIgnore(path) {
			// FIXME path may not always be a dir here
			return filepath.SkipDir
		}

		if info.IsDir() {
			return nil
		}

		fmt.Printf("Adding %s (size: %d).\n", path, info.Size())
		hdr, err := tar.FileInfoHeader(info, path)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		hdr.Name = path

		if err = tw.WriteHeader(hdr); err != nil {
			fmt.Println(err.Error())
			return err
		}

		body, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		if _, err = tw.Write(body); err != nil {
			fmt.Println(err.Error())
			return err
		}

		return nil
	})

	if err := tw.Close(); err != nil {
		fmt.Println(err.Error())
	}

	if err := gz.Close(); err != nil {
		fmt.Println(err.Error())
	}

	return *buf
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

	dir := os.Args[1] // TODO: Maybe fallback to CWD or Git root?
	tgz := buildTgz(dir)
	fmt.Printf("%v %d\n", tgz.Bytes(), tgz.Len())
	// fmt.Println(string(tgz.Bytes()))

	/*
		TODO:
			* Check that we have an APP context or set it
			* upload tgz to S3 with an object expiry of ~5min
			* hit build API with that link
			* tail output (if build api has implemented that)
	*/
}

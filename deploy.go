package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bgentry/heroku-go"
	hk "github.com/heroku/hk/hkclient"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
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

		fmt.Printf("  Adding %s (%d bytes)\n", path, info.Size())

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

type S3Upload struct {
	Action, DownloadUrl string
	Fields              map[string]string
}

func getUploadSlot() (u S3Upload, err error) {
	if resp, err := http.Get(ENDPOINT); err == nil {
		defer resp.Body.Close()
		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			json.Unmarshal([]byte(body), &u)
		}
	}

	return
}

func upload(tgz *bytes.Buffer, u S3Upload) (err error) {
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)

	// Set the form fields that the server told us to
	for field, value := range u.Fields {
		w.WriteField(field, value)
	}

	// then attach the actual file (it must be last)
	if writer, err := w.CreateFormFile("file", "code.tgz"); err == nil {
		if _, err = io.Copy(writer, tgz); err != nil {
			return err
		}
	} else {
		return err
	}

	if err = w.Close(); err != nil {
		return
	}

	if req, err := http.NewRequest("POST", u.Action, buf); err == nil {
		req.Header.Add("Content-Type", w.FormDataContentType())
		resp, err := http.DefaultClient.Do(req)
		defer resp.Body.Close()
		if err == nil {
			if resp.StatusCode >= 400 {
				body, _ := ioutil.ReadAll(resp.Body)
				fmt.Println(string(body))
				return errors.New("S3 upload rejected")
			}
		}
	}

	return
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

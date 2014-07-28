package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

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
				return errors.New("S3 upload rejected")
			}
		}
	}

	return
}

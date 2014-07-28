package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

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

package main

import (
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

type structure struct {
	path string
	tmpl *template.Template
}

func run() error {
	b, err := ioutil.ReadFile(".skaffer.yaml")
	if err != nil {
		return err
	}
	data := make(map[string]interface{})
	yaml.Unmarshal(b, &data)

	var templates []structure

	err = filepath.Walk("./tmpl", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		tmpl, err := template.ParseFiles(path)
		if err != nil {
			return err
		}

		item := structure{
			path: strings.TrimPrefix(path, "tmpl/"),
			tmpl: tmpl,
		}

		templates = append(templates, item)

		return nil
	})
	if err != nil {
		return err
	}

	for _, item := range templates {
		directory := filepath.Join("example", filepath.Dir(item.path))
		err := os.MkdirAll(directory, os.ModePerm)
		if err != nil {
			if !os.IsExist(err) {
				return err
			}
		}

		f, err := os.Create(filepath.Join("example", item.path))
		if err != nil {
			return err
		}

		defer f.Close()

		err = item.tmpl.Execute(f, data["values"])
		if err != nil {
			return err
		}
	}
	return nil
}

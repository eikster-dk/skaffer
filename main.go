package main

import (
	"html/template"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
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

type repository struct {
	Repo    string
	Version string
}

type skaffer struct {
	Repository repository
	Template   map[string]interface{}
}

func run() error {
	b, err := ioutil.ReadFile(".skaffer.yaml")
	if err != nil {
		return err
	}

	var skafferTemplate skaffer
	yaml.Unmarshal(b, &skafferTemplate)

	var templates []structure

	_, err = git.PlainClone("./.template", false, &git.CloneOptions{
		URL:      skafferTemplate.Repository.Repo,
		Progress: os.Stdout,
	})
	if err != nil {

		return err
	}
	defer func() {
		os.RemoveAll("./.template")
	}()

	err = filepath.Walk(".template", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.Contains(path, ".git") {
			return nil
		}

		tmpl, err := template.ParseFiles(path)
		if err != nil {
			return err
		}

		item := structure{
			path: strings.TrimPrefix(path, ".template/"),
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

		err = item.tmpl.Execute(f, skafferTemplate.Template)
		if err != nil {
			return err
		}
	}
	return nil
}

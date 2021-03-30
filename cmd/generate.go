package cmd

import (
	"fmt"
	"html/template"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func init() {
	rootCmd.AddCommand(generateCmd)
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "generates the desired template",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		type templateFile struct {
			path string
			tmpl *template.Template
		}

		type repository struct {
			Repo    string
			Version string
		}

		type skafferConfiguration struct {
			Repository repository
			Template   map[string]interface{}
		}

		b, err := ioutil.ReadFile(".skaffer.yaml")
		if err != nil {
			return err
		}

		var skafferConfig skafferConfiguration
		yaml.Unmarshal(b, &skafferConfig)

		var templates []templateFile

		_, err = git.PlainClone("./.template", false, &git.CloneOptions{
			URL:      skafferConfig.Repository.Repo,
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

			item := templateFile{
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

			err = item.tmpl.Execute(f, skafferConfig.Template)
			if err != nil {
				return err
			}
		}

		fmt.Println("------")
		fmt.Println("Template generated")
		return nil
	},
}

package cli

import (
	"errors"
	"fmt"
	"log"
	"path"
	"strings"
	"text/template"

	"github.com/spf13/cobra"

	"tasks/cmd/asynqgen/internal"
)

func init() {
	rootCmd.AddCommand(taskCmd)
	taskCmd.Flags().StringVarP(&Domain, "domain", "d", "", "asynqgen --domain <domain_name>")
	taskCmd.Flags().StringVarP(&Name, "name", "n", "", "asynqgen --name <task_name>")
}

const (
	dirTemplate = "template"
)

var (
	Domain = ""
	Name   = ""
)

var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "`asynqgen task` is an asynq new task scaffolder",
	Long:  `Quickly scaffold a new asynq task - https://github.com/gmhafiz/asynq`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := checkVersion(); err != nil {
			log.Fatal(err)
		}

		domain := strings.TrimSpace(strings.Title(Domain))
		domainLowerCase := strings.TrimSpace(strings.ToLower(Domain))
		name := strings.Title(strings.TrimSpace(Name))

		a := internal.New(internal.Project{
			TaskConstant:      fmt.Sprintf("Type%s%s", domain, name),
			TaskConstantValue: fmt.Sprintf(`"%s:%s"`, Domain, Name),
			TaskName:          Name,
			Domain:            domain,
			DomainLowerCase:   domainLowerCase,
		})

		if err := a.CreateDomainDirectory(domainLowerCase); err != nil {
			log.Fatal(err)
		}

		templateFiles, err := a.Embed.ReadDir(dirTemplate)
		if err != nil {
			return
		}

		for _, tmpl := range templateFiles {
			a.Templates[tmpl.Name()] = &tmpl

			if tmpl.IsDir() {
				continue
			}

			t, err := template.ParseFS(a.Embed, path.Join(dirTemplate, tmpl.Name()))
			if err != nil {
				log.Fatal(err)
			}

			file, err := a.CreateFile(tmpl.Name())
			if err != nil {
				if errors.Is(err, internal.ErrFileExists) || errors.Is(err, internal.ErrBlacklisted) {
					continue
				}
				log.Fatal(err)
			}

			if err := t.ExecuteTemplate(file, tmpl.Name(), a.Project); err != nil {
				log.Fatal(err)
			}
		}

		if err := a.InjectTaskConstantName(); err != nil {
			log.Fatal(err)
		}

		if err := a.InjectTask(); err != nil {
			log.Fatal(err)
		}

		if err := a.GoFmt(); err != nil {
			log.Fatal(err)
		}

		fmt.Printf(InfoColor, "New task scaffold is completed.\n")
	},
}

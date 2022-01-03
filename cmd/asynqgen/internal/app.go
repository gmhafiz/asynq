package internal

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"github.com/rogpeppe/go-internal/modfile"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"text/template"
)

var (
	//go:embed template/*
	templates embed.FS

	ErrFileExists  = errors.New("file exists")
	ErrBlacklisted = errors.New("blacklisted")
)

type App struct {
	Project Project
	Embed   embed.FS

	Templates map[string]*fs.DirEntry
}

type Project struct {
	TaskConstant      string
	TaskConstantValue string
	TaskName          string

	Domain string

	DomainLowerCase string
	ModuleName      string
}

func New(p Project) *App {
	a := &App{
		Project: p,
		Embed:   templates,

		Templates: make(map[string]*fs.DirEntry),
	}
	name, err := getProjectName()
	if err != nil {
		log.Fatalln(err)
	}
	a.Project.ModuleName = name

	return a
}

// CreateDomainDirectory create a new directory for the domain at ./task/{{domain}}
// and ./tasks/internal/domain/{{domain}}
func (a *App) CreateDomainDirectory(domain string) error {
	var directories []string
	directories = append(directories, path.Join("task", domain))
	directories = append(directories, path.Join("internal", "domain", domain, "handler", "http"))
	directories = append(directories, path.Join("internal", "domain", domain, "usecase"))

	for _, directory := range directories {
		if exists(directory) {
			continue
		}
		err := os.MkdirAll(directory, 0750)
		if err != nil {
			return err
		}
	}

	return nil
}

func exists(path string) bool {
	_, err := os.Open(path)
	return !errors.Is(err, os.ErrNotExist)
}

func (a *App) CreateFile(fileName string) (*os.File, error) {
	base := fileNameWithoutExtension(fileName)
	fileName = fmt.Sprintf("%s.go", base)

	switch fileName {
	case "create.go":
		fallthrough
	case "process.go":
		fileName = path.Join("task", a.Project.DomainLowerCase, fileName)
	case "handler.go":
		fallthrough
	case "register.go":
		fileName = path.Join("internal", "domain", a.Project.DomainLowerCase, "handler", "http", fileName)
	case "usecase.go":
		fileName = path.Join("internal", "domain", a.Project.DomainLowerCase, "usecase", fileName)
	case "request.go":
		fileName = path.Join("internal", "domain", a.Project.DomainLowerCase, fileName)
	}

	if exists(fileName) {
		return nil, ErrFileExists
	}

	blacklist := []string{"newDomain.go", "newDomainImport.go"}
	if inArray(blacklist, fileName) {
		return nil, ErrBlacklisted
	}

	file, err := os.Create(fileName)
	if err != nil {
		return nil, fmt.Errorf("error creating file: %s", fileName)
	}

	return file, nil
}

func inArray(blacklist []string, name string) bool {
	for _, i := range blacklist {
		if i == name {
			return true
		}
	}
	return false
}

func (a *App) InjectTaskConstantName() error {
	const taskListFile = "task/tasks.go"
	const injectImport = "//generate:asynq:gen"

	newTaskConstant := fmt.Sprintf("\n\t// %s defines a Unit of Work\n\t%s = %s", a.Project.TaskConstant, a.Project.TaskConstant, a.Project.TaskConstantValue)

	fileContent, err := ioutil.ReadFile(taskListFile)
	if err != nil {
		return err
	}

	if strings.Contains(string(fileContent), newTaskConstant) {
		return nil
	}

	var newFile []string
	temp := strings.Split(string(fileContent), "\n")
	for _, line := range temp {
		newFile = append(newFile, line)
		stripped := strings.Trim(line, "\t")
		stripped = strings.Trim(stripped, "\n")
		if stripped == injectImport {
			newFile = append(newFile, newTaskConstant)
		}
	}

	fCreate, err := os.Create(taskListFile)
	if err != nil {
		return fmt.Errorf("error opening file: %s. %w", taskListFile, err)
	}
	_, err = fCreate.WriteString(strings.Join(newFile, "\n"))
	if err != nil {
		return fmt.Errorf("error writing file: %s, %w", taskListFile, err)
	}
	return nil
}

func (a *App) InjectDomainRegistration() error {
	const initDomainsFile = "internal/server/initDomains.go"
	const importLine = "import ("
	const callFunc = "func (s *Server) initDomains() {"

	checkDomainRegistration := fmt.Sprintf("\"tasks/internal/domain/%s/", a.Project.DomainLowerCase)

	fileContent, err := ioutil.ReadFile(initDomainsFile)
	if err != nil {
		return err
	}

	if strings.Contains(string(fileContent), checkDomainRegistration) {
		return nil
	}

	importTemplate, err := template.ParseFS(a.Embed, "template/newDomainImport.tmpl")
	if err != nil {
		return err
	}
	importLines := bytes.NewBuffer([]byte{})
	if err := importTemplate.ExecuteTemplate(importLines, importTemplate.Name(), a.Project); err != nil {
		log.Fatal(err)
	}

	initTemplate, err := template.ParseFS(a.Embed, "template/newDomain.tmpl")
	if err != nil {
		return err
	}
	initLines := bytes.NewBuffer([]byte{})
	if err := initTemplate.ExecuteTemplate(initLines, initTemplate.Name(), a.Project); err != nil {
		log.Fatal(err)
	}

	var newFile []string
	temp := strings.Split(string(fileContent), "\n")
	for _, line := range temp {
		newFile = append(newFile, line)
		stripped := strings.Trim(line, "\t")
		stripped = strings.Trim(stripped, "\n")
		if stripped == importLine {
			newFile = append(newFile, importLines.String())
		}
		if stripped == callFunc {
			newFile = append(newFile, fmt.Sprintf("\ts.init%s()", a.Project.Domain))
		}
	}
	newFile = append(newFile, initLines.String())

	fCreate, err := os.Create(initDomainsFile)
	if err != nil {
		return fmt.Errorf("error opening file: %s. %w", initDomainsFile, err)
	}
	_, err = fCreate.WriteString(strings.Join(newFile, "\n"))
	if err != nil {
		return fmt.Errorf("error writing file: %s, %w", initDomainsFile, err)
	}
	return nil

}

func (a *App) GoFmt() error {
	command := exec.Command("go", "fmt", "./...")
	_ = command.Run()

	return nil
}

func fileNameWithoutExtension(fileName string) string {
	if pos := strings.LastIndexByte(fileName, '.'); pos != -1 {
		return fileName[:pos]
	}
	return fileName
}

// adapted from https://stackoverflow.com/a/63393712/1033134
func getProjectName() (string, error) {
	goModBytes, err := ioutil.ReadFile("go.mod")
	if err != nil {
		return "", errors.New("error go mod is not initialized")
	}
	return modfile.ModulePath(goModBytes), nil
}

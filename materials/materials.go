package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/materials-commons/materials"
	"github.com/materials-commons/materials/autoupdate"
	"github.com/materials-commons/materials/site"
	"github.com/materials-commons/materials/wsmaterials"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

var mcuser, _ = materials.NewCurrentUser()

type ServerOptions struct {
	AsServer bool   `long:"server" description:"Run as webserver"`
	Port     uint   `long:"port" description:"The port the server listens on"`
	Address  string `long:"address" description:"The address to bind to"`
	Retry    int    `long:"retry" description:"Number of times to retry connecting to address/port"`
}

type ProjectOptions struct {
	Project   string `long:"project" description:"Specify the project"`
	Directory string `long:"directory" description:"The directory path to the project"`
	Add       bool   `long:"add" description:"Add the project to the project config file"`
	Delete    bool   `long:"delete" description:"Delete the project from the project config file"`
	List      bool   `long:"list" description:"List all known projects and their locations"`
	Upload    bool   `long:"upload" description:"Uploads a new project. Cannot be used on existing projects"`
	Convert   bool   `long:"convert" description:"Converts projects to new layout"`
}

type Options struct {
	Server     ServerOptions  `group:"Server Options"`
	Project    ProjectOptions `group:"Project Options"`
	Initialize bool           `long:"init" description:"Create configuration"`
}

func initialize() {
	usr, err := user.Current()
	checkError(err)

	dirPath := filepath.Join(usr.HomeDir, ".materials")
	err = os.MkdirAll(dirPath, 0777)
	checkError(err)

	if downloadedTo, err := site.Download(); err == nil {
		if site.IsNew(downloadedTo) {
			site.Deploy(downloadedTo)
		}
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err.Error())
		os.Exit(1)
	}
}

func listProjects() {
	projects := materials.CurrentUserProjectDB()
	for _, p := range projects.Projects() {
		fmt.Printf("%s, %s\n", p.Name, p.Path)
	}
}

func convertProjects() {
	setupProjectsDir()
	convertProjectsFile()
}

func setupProjectsDir() {
	projectDB := filepath.Join(materials.Config.User.DotMaterialsPath(), "projectdb")
	err := os.MkdirAll(projectDB, os.ModePerm)
	checkError(err)
}

func convertProjectsFile() {
	projectsPath := filepath.Join(materials.Config.User.DotMaterialsPath(), "projects")
	projectsFile, err := os.Open(projectsPath)
	projectdbPath := filepath.Join(materials.Config.User.DotMaterialsPath(), "projectdb")
	checkError(err)
	defer projectsFile.Close()

	scanner := bufio.NewScanner(projectsFile)
	for scanner.Scan() {
		splitLine := strings.Split(scanner.Text(), "|")
		if len(splitLine) == 3 {
			project := materials.Project{
				Name:    strings.TrimSpace(splitLine[0]),
				Path:    strings.TrimSpace(splitLine[1]),
				Status:  strings.TrimSpace(splitLine[2]),
				ModTime: time.Now(),
				Changes: map[string]materials.ProjectFileChange{},
				Ignore:  []string{},
			}
			b, err := json.MarshalIndent(&project, "", "  ")
			if err != nil {
				fmt.Printf("Could not convert '%s' to new project format\n", scanner.Text())
				continue
			}
			path := filepath.Join(projectdbPath, project.Name+".project")
			if err := ioutil.WriteFile(path, b, os.ModePerm); err != nil {
				fmt.Printf("Unable to write project file %s\n", path)
			}
		}
	}
}

func uploadProject(projectName string) {
	projects := materials.CurrentUserProjectDB()
	project, _ := projects.Find(projectName)
	err := project.Upload()
	if err != nil {
		fmt.Println(err)
	} else {
		projects.Update(func() *materials.Project {
			project.Status = "Loaded"
			return project
		})
	}
}

func startServer(serverOpts ServerOptions) {
	autoupdate.StartUpdateMonitor()

	if serverOpts.Address != "" {
		materials.Config.Server.Address = serverOpts.Address
	}

	if serverOpts.Port != 0 {
		materials.Config.Server.Port = serverOpts.Port
	}

	if serverOpts.Retry != 0 {
		wsmaterials.StartRetry(serverOpts.Retry)
	} else {
		wsmaterials.Start()
	}
}

func main() {
	materials.ConfigInitialize(mcuser)
	var opts Options
	flags.Parse(&opts)

	switch {
	case opts.Initialize:
		initialize()
	case opts.Project.List:
		listProjects()
	case opts.Project.Convert:
		convertProjects()
	case opts.Project.Upload:
		uploadProject(opts.Project.Project)
	case opts.Server.AsServer:
		startServer(opts.Server)
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type Project struct {
	Name              string
	Status            string
	LastBuildDate     time.Time
	CmdPath           string
	ErrorPath         string
	LastBuildDuration int
}

type Config struct {
	Port      int
	Projects  []*Project
	LogToFile bool
}

var config Config

func main() {
	setWorkingDirectory()
	readConfig()
	if config.LogToFile {
		logfile, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer logfile.Close()
		wrt := io.MultiWriter(os.Stdout, logfile)
		log.SetOutput(wrt)
	}
	tmpl := template.Must(template.ParseFiles("www/index.gohtml"))

	http.HandleFunc("/build/", build)
	http.HandleFunc("/status/", status)
	http.HandleFunc("/errors/", errors)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, config)
	})

	// fileServer := http.FileServer(http.Dir("./www"))
	// http.Handle("/", fileServer)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Printf("Starting server on port %d\n", config.Port)
	go func() {
		if err := http.ListenAndServe(":"+strconv.Itoa(config.Port), nil); err != nil {
			log.Fatal(err)
		}
	}()
	<-done
	log.Println()
	log.Println("Shutting down...")
	writeConfig()
}

func readConfig() {
	config = Config{
		Port:      2000,
		LogToFile: true,
	}

	config.Projects = []*Project{
		{
			Name:              "build project1",
			Status:            "",
			LastBuildDate:     time.Now(),
			CmdPath:           "c:/folder/build.cmd",
			ErrorPath:         "c:/folder/error.txt",
			LastBuildDuration: 1000,
		},
	}

	if fileExists("config.json") {
		log.Println("reading config")
		b, _ := os.ReadFile("config.json")
		json.Unmarshal(b, &config)
	}
	writeConfig()

	// reset project status if server crashed
	for _, i := range config.Projects {
		if i.Status == "building" {
			i.Status = ""
		}
	}
}

func writeConfig() {
	b, _ := json.MarshalIndent(config, "", "  ")
	os.WriteFile("config.json", b, 0644)
}

func errors(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/errors/")
	prj := config.findProject(name)
	if prj == nil {
		return
	}

	b, err := os.ReadFile(prj.ErrorPath)
	if err != nil {
		log.Println("ERROR", err.Error())
		return
	}
	fmt.Fprintf(w, "%s", string(b))
}

func status(w http.ResponseWriter, r *http.Request) {

	name := strings.TrimPrefix(r.URL.Path, "/status/")

	prj := config.findProject(name)
	if prj == nil {
		return
	}
	fmt.Fprintf(w, "status : %s = %s", name, prj.Status)
}

func build(w http.ResponseWriter, r *http.Request) {

	name := strings.TrimPrefix(r.URL.Path, "/build/")

	prj := config.findProject(name)
	if prj == nil {
		fmt.Fprintln(w, "project not found")
		return
	}
	if prj.Status == "building" {
		log.Println("already building")
		fmt.Fprintln(w, "already building")
		return
	}
	prj.Status = "building"
	prj.LastBuildDate = time.Now()
	fmt.Fprintf(w, "build : %s = %s\n", name, prj.Status)
	log.Println("building :", prj.Name)
	// rr, ww, _ := os.Pipe()

	// go io.Copy(w, rr)
	ex := exec.Command(prj.CmdPath)
	ex.Dir = path.Dir(prj.CmdPath)

	go func() {
		e := ex.Run()
		failed := false
		if e != nil {
			log.Println(e)
			failed = true
		}
		prj.LastBuildDuration = int(time.Since(prj.LastBuildDate).Seconds())
		log.Println("build done :", prj.Name)
		log.Println("duration :", prj.LastBuildDuration)
		st, err := os.Stat(prj.ErrorPath)
		if err != nil {
			log.Println("ERROR" + err.Error())
			return
		}
		if st.Size() > 0 && failed {
			prj.Status = "error"
			return
		}
		prj.Status = "done"
	}()
}

func (c *Config) findProject(name string) *Project {
	for _, i := range c.Projects {
		if i.Name != name {
			continue
		}
		return i
	}
	return nil
}

func fileExists(fn string) bool {
	_, e := os.Stat(fn)
	return e == nil
}

func setWorkingDirectory() {
	f, _ := os.Executable()
	os.Chdir(filepath.Dir(f))
}

package main

import (
	"fmt"
	"templates"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const defaultListenPort = "8000"

type HelloWorldConfig struct {
	Hostname string
	Services map[string]string
	Environment	map[string]string
	Headers  http.Header
	Request  http.Request
	Host     string
}

func (config *HelloWorldConfig) GetManifest() (string, error) {
	return templates.CompileTemplateFromMap(templates.HelloWorldTemplate, config)
}

func (config *HelloWorldConfig) getServices() {
	k8sServices := make(map[string]string)
	environment := make(map[string]string)

	for _, evar := range os.Environ() {
		show := strings.Split(evar, "=")
		regName := regexp.MustCompile("^.*_PORT$")
		regLink := regexp.MustCompile("^(tcp|udp)://.*")
		if regName.MatchString(show[0]) && regLink.MatchString(show[1]) {
			k8sServices[strings.TrimSuffix(show[0], "_PORT")] = show[1]
		}
		environment[show[0]] = show[1]
	}

	config.Services = k8sServices
	config.Environment = environment
}

func (config *HelloWorldConfig) Init(r *http.Request) {
	config.Hostname, _ = os.Hostname()
	config.Host = r.Host
	config.Headers = r.Header
	config.getServices()
	config.Request = *r
}

func handler(w http.ResponseWriter, r *http.Request) {
	config := &HelloWorldConfig{}
	config.Init(r)
	data, err := config.GetManifest()
	if err != nil {
		fmt.Fprintln(w, err)
	}

	fmt.Fprint(w, data)
}

func main() {
	webPort := os.Getenv("HTTP_PORT")
	if webPort == "" {
		webPort = defaultListenPort
	}

	fmt.Println("Running http service at", webPort, "port")
	http.HandleFunc("/", handler)
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir(os.Getenv("PWD")))))
	http.ListenAndServe(":"+webPort, nil)
}

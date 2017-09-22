package main

import (
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/golang/protobuf/jsonpb"
)

type ArgT struct {
	TemplateName string
	FileName     string
}

var (
	ticker      = "AAPL"
	templateDir = "/src/github.com/gooption/gooption-gen/templates"
	projectMap  = map[string][]ArgT{
		"service": []ArgT{
			{TemplateName: "handler", FileName: "handlers.go"},
			{TemplateName: "handler_test", FileName: "handlers_test.go"},
		},
		"gobs": []ArgT{
			{TemplateName: "service", FileName: "gobs-service_test.go"},
		},
		"client": []ArgT{},
	}
)

func main() {
	project := os.Args[1]
	requests := os.Args[2:]

	generateJSON(requests, func(path string) (io.WriteCloser, error) {
		f, err := os.Create(path)
		return f, err
	})

	generateTemplate(projectMap[project], requests, func(path string) (io.WriteCloser, error) {
		f, err := os.Create(path)
		return f, err
	})
}

func generateTemplate(projectConfig []ArgT, requests []string, newWriter func(path string) (io.WriteCloser, error)) {
	funcMap := template.FuncMap{
		"ToLower": strings.ToLower,
	}

	vals := map[string]interface{}{
		"Package":  os.Getenv("GOPACKAGE"),
		"Requests": requests,
	}

	root := os.Getenv("GOPATH")
	templates, _ := template.New("templates").Funcs(funcMap).ParseGlob(root + templateDir + "/*")
	for _, config := range projectConfig {
		if w, err := newWriter(config.FileName); err == nil {
			templates.ExecuteTemplate(w, config.TemplateName, vals)
			w.Close()
		} else {
			panic(err)
		}
	}
}

func generateJSON(requests []string, newWriter func(path string) (io.WriteCloser, error)) {
	marshaler := jsonpb.Marshaler{EnumsAsInts: true}
	for _, requestArg := range requests {
		request := strings.ToLower(requestArg)
		if mockRequest, err := NewRequest(request); err == nil {
			if w, err := newWriter("./testdata/" + requestArg + "Request.json"); err == nil {
				defer w.Close()
				marshaler.Marshal(w, mockRequest)
			} else {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
}

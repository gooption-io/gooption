package main

import (
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/golang/protobuf/jsonpb"
	"github.com/mkideal/cli"
)

type argT struct {
	cli.Helper
	Project  string   `cli:"*p, project" usage:"target project eg. service or gobs"`
	Requests []string `cli:"*r, requests" usage:"requets to generate eg. price greek or impliedvol"`
}

type config struct {
	TemplateName string
	FileName     string
}

var (
	ticker      = "AAPL"
	templateDir = "/src/github.com/gooption/gooption-cli/templates"
	configMap   = map[string][]config{
		"service": []config{
			{TemplateName: "handler", FileName: "handlers.go"},
			{TemplateName: "handler_test", FileName: "handlers_test.go"},
		},
		"gobs": []config{
			{TemplateName: "service", FileName: "gobs-srv_test.go"},
		},
		"client": []config{}, //json only
	}
)

func main() {
	cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)

		newFileWriter := func(path string) (io.WriteCloser, error) {
			f, err := os.Create(path)
			return f, err
		}

		generateJSON(argv.Requests, newFileWriter)
		generateTemplate(configMap[argv.Project], argv.Requests, newFileWriter)
		return nil
	})
}

func generateTemplate(projectConfig []config, requests []string, newWriter func(path string) (io.WriteCloser, error)) {
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

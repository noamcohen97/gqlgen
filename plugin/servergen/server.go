package servergen

import (
	_ "embed"
	"errors"
	"io/fs"
	"log"
	"os"

	"github.com/99designs/gqlgen/codegen"
	"github.com/99designs/gqlgen/codegen/templates"
	"github.com/99designs/gqlgen/plugin"
)

//go:embed server.gotpl
var serverTemplate string

func New(filename string) plugin.Plugin {
	return &Plugin{filename}
}

type Plugin struct {
	filename string
}

var _ plugin.CodeGenerator = &Plugin{}

func (m *Plugin) Name() string {
	return "servergen"
}

func (m *Plugin) GenerateCode(data *codegen.Data) error {
	serverBuild := &ServerBuild{
		ExecPackageName:     data.Config.Exec.ImportPath(),
		ResolverPackageName: data.Config.Resolver.ImportPath(),
	}

	if _, err := os.Stat(m.filename); errors.Is(err, fs.ErrNotExist) {
		return templates.Render(templates.Options{
			PackageName:           "main",
			Filename:              m.filename,
			Data:                  serverBuild,
			Packages:              data.Config.Packages,
			Template:              serverTemplate,
			RemoveFileBeforeWrite: data.Config.RemoveFileBeforeWrite,
		})
	}

	log.Printf("Skipped server: %s already exists\n", m.filename)
	return nil
}

type ServerBuild struct {
	codegen.Data

	ExecPackageName     string
	ResolverPackageName string
}

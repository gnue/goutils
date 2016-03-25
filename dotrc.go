package dotrc

import (
	"bytes"
	"os"
	"os/exec"
	"text/template"
)

const execScript = `
export SHELL={{ .Shell }}

for f in {{range .Files}}{{.}} {{end}}; do
  [ -f "$f" ] && source "$f"
done

exec "$@"`

var execTmpl = template.Must(template.New("exec").Parse(execScript))

type Dotrc struct {
	Shell string
	Files []string
}

func New(files ...string) *Dotrc {
	sh := os.Getenv("SHELL")
	if sh == "" {
		sh = "/bin/sh"
	}

	return &Dotrc{Shell: sh, Files: files}
}

func (rc *Dotrc) Script() string {
	var b bytes.Buffer

	err := execTmpl.Execute(&b, rc)
	if err != nil {
		panic(err)
	}

	return b.String()
}

func (rc *Dotrc) Command(name string, arg ...string) *exec.Cmd {
	args := append([]string{"-c", rc.Script(), "--", name}, arg...)
	return exec.Command(rc.Shell, args...)
}

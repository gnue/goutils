package dotrc

import (
	"bytes"
	"os/exec"
	"text/template"
)

const execScript = `
for f in {{range .Files}}{{.}} {{end}}; do
  [ -f "$f" ] && source "$f"
done
exec "$@"`

var execTmpl = template.Must(template.New("exec").Parse(execScript))

type Dotrc struct {
	Files []string
}

func New(files ...string) *Dotrc {
	return &Dotrc{Files: files}
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
	return exec.Command("/bin/sh", args...)
}

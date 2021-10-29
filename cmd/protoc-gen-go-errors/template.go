package main

import (
	"bytes"
	"text/template"
)

var errorsMap = `
var (
	errorsEnMap map[string]string
	errorsChMap map[string]string
	once        sync.Once
	onceCh      sync.Once
)

func ch(in string) (out string) {
	onceCh.Do(func() {
		errorsChMap = map[string]string{
		{{ range .Errors }}
			{{.Name}}_{{.Value}}.String(): "{{.ReasonCh}}",
		{{- end }}
		}
	})
	out = in
	if r := errorsChMap[in]; len(r) > 0 {
		out = r
	}
	return
}

func en(in string) (out string) {
	once.Do(func() {
		errorsEnMap = map[string]string{
		{{ range .Errors }}
			{{.Name}}_{{.Value}}.String(): "{{.Reason}}",
		{{- end }}
		}
	})
	out = in
	if r := errorsEnMap[in]; len(r) > 0 {
		out = r
	}
	return
}

func convert(ctx context.Context, s string) string {
	if lang, ok := errorctx.FromErrorsContext(ctx); ok {
		if strings.EqualFold(lang, "ch") {
			return ch(s)
		}
	}
	return en(s)
}
`
var errorsTemplate = `
{{ range .Errors }}

func Is{{.CamelValue}}(err error) bool {
	e := errors.FromError(err)
	return e.Reason == {{.Name}}_{{.Value}}.String() && e.Code == {{.HTTPCode}} 
}

func Error{{.CamelValue}}(ctx context.Context, format string, args ...interface{}) *errors.Error {
  	reason := convert(ctx, {{.Name}}_{{.Value}}.String())
	return errors.New({{.HTTPCode}}, reason, fmt.Sprintf(format, args...))
}

{{- end }}
`

type errorInfo struct {
	Name       string
	Value      string
	Reason     string
	ReasonCh   string
	HTTPCode   int
	CamelValue string
}

type errorWrapper struct {
	Errors []*errorInfo
}

func (e *errorWrapper) executeMap() string {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("errors").Parse(errorsMap)
	if err != nil {
		panic(err)
	}
	if err := tmpl.Execute(buf, e); err != nil {
		panic(err)
	}
	return buf.String()
}

func (e *errorWrapper) execute() string {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("errors").Parse(errorsTemplate)
	if err != nil {
		panic(err)
	}
	if err := tmpl.Execute(buf, e); err != nil {
		panic(err)
	}
	return buf.String()
}

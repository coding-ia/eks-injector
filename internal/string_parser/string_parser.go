package string_parser

import (
	"bytes"
	"text/template"
)

func ParseString(value string, variables map[string]string) (string, error) {
	tmpl := template.New("variables")

	_, err := tmpl.Parse(value)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, variables)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

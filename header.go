package main

import (
	"html/template"
	"io"
)

type Header struct {
	Title string
}

const (
	headerTemplateStr string = `<!DOCTYPE html>
<head>
<link rel="stylesheet" href="_/static/css/bootstrap.min.css" >
<script src="_/static/js/jquery-3.4.1.min.js"></script>
<script src="_/static/js/popper.min.js"></script>
<title>{{.Title}}</title>
</head>
`
)

var headerTemplate *template.Template = nil

//WriteHeader write html header.
func WriteHeader(writer io.Writer, header Header) error {
	if headerTemplate == nil {
		var err error
		headerTemplate, err = template.New("header").Parse(headerTemplateStr)
		if err != nil {
			return err
		}
	}

	return headerTemplate.Execute(writer, header)
}

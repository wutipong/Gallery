package main

import (
	"html/template"
	"io"
	"log"
	"os"
)

type Header struct {
	Title string
}

func init() {
	var err error
	headerTemplate, err = template.New("header.gohtml").ParseFiles("template/header.gohtml")
	if err != nil {
		log.Panic(err)
		os.Exit(-1)
	}
	err = headerTemplate.Execute(os.Stdout, Header{Title: "Test"})
}

var headerTemplate *template.Template = nil

//WriteHeader write html header.
func WriteHeader(writer io.Writer, header Header) error {
	return headerTemplate.Execute(writer, header)
}

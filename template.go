package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"text/template"

	readability "github.com/mauidude/go-readability"
	"github.com/motemen/go-pocket/api"
)

const templ = `
<!DOCTYPE html>
<meta charset="UTF-8">
<h1>{{.Title}}</h1>
{{.Body}}
`

type Article struct {
	api.Item
}

func (a Article) String() string {
	r, err := a.getReadability()
	if err != nil {
		panic(err)
	}

	t := template.Must(template.New("article").Parse(templ))

	var buf bytes.Buffer
	t.Execute(&buf, map[string]string{
		"Body":  r.Content(),
		"Title": a.Title(),
	})

	return buf.String()
}

func (a Article) getReadability() (readability.Document, error) {
	resp, err := http.Get(a.URL())
	if err != nil {
		return readability.Document{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return readability.Document{}, err
	}

	var doc *readability.Document
	doc, err = readability.NewDocument(string(body))
	if err != nil {
		return readability.Document{}, err
	}

	return *doc, nil
}

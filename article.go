package main

import (
	"html/template"
	"io"
	"net/http"
	"time"

	readability "github.com/go-shiori/go-readability"
	"github.com/motemen/go-pocket/api"
)

const timeOut = 30 * time.Second

const templ = `
<!DOCTYPE html>
<html>
<head>
  <title>{{.Title}}</title>
  <meta charset="UTF-8">
  <link rel="canonical" href="{{.URL}}" />
</head>
<body>
  <h1><a href="{{.URL}}">{{.Title}}</a></h1>
  {{ if ne .Image "" }}
    <img src="{{.Image}}" />
  {{ end }}
  {{.Body}}
</body>
`

var t *template.Template

func init() {
	t = template.Must(template.New("article").Parse(templ))
}

type Article struct {
	api.Item
}

func (a Article) Download() (io.Reader, error) {
	r, w := io.Pipe()

	resp, err := http.Get(a.URL())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	parser := readability.NewParser()
	article, err := parser.Parse(resp.Body, a.URL())
	if err != nil {
		return nil, err
	}

	go func() {
		t.Execute(w, map[string]interface{}{
			"Body":  template.HTML(article.Content),
			"Url":   a.URL(),
			"Title": a.Title(),
			"Image": article.Image,
		})
		w.Close()
	}()

	return r, err
}

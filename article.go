package main

import (
	"bytes"
	"html/template"
	"time"

	readability "github.com/RadhiFadlillah/go-readability"
	"github.com/motemen/go-pocket/api"
)

const timeOut = 30 * time.Second

const templ = `
<!DOCTYPE html>
<html>
<head>
  <title>{{.Title}}</title>
  <meta charset="UTF-8">
  <meta name="author" content="{{.Author}}">
  <link rel="canonical" href="{{.URL}}" />
</head>
<body>
  <h1><a href="{{.URL}}">{{.Title}}</a></h1>
  {{ if ne .Image "" }}
    <img src="{{.Image}}" />
  {{ end }}
  {{.Body }}
</body>
`

type Article struct {
	api.Item
}

func (a Article) Download() ([]byte, error) {
	article, err := readability.Parse(a.URL(), timeOut)
	if err != nil {
		return []byte{}, err
	}

	t := template.Must(template.New("article").Parse(templ))

	var buf bytes.Buffer
	err = t.Execute(&buf, map[string]interface{}{
		"Body":   template.HTML(article.RawContent),
		"Url":    a.URL(),
		"Title":  a.Title(),
		"Image":  article.Meta.Image,
		"Author": article.Meta.Author,
	})

	return buf.Bytes(), err
}

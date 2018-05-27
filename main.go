package main

import (
	"fmt"
	"io/ioutil"

	"github.com/motemen/go-pocket/api"
)

const configDir = "/tmp"

const outputDir = "./articles"

const clearScreen = "\033[H\033[2J"

func main() {
	client, err := GetClient()
	if err != nil {
		panic(err)
	}

	options := &api.RetrieveOption{}
	res, err := client.Retrieve(options)
	if err != nil {
		panic(err)
	}

	i := 0
	for _, item := range res.List {
		i++
		print(clearScreen)
		fmt.Printf("[%d\t/\t%d] Downloading article: %s", i, len(res.List), item.Title())

		fileName := outputDir + "/" + item.Title() + ".html"
		article := Article{Item: item}
		err := ioutil.WriteFile(fileName, []byte(article.String()), 0644)
		if err != nil {
			panic(err)
		}
	}
}

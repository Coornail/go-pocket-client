package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

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

		fileName := outputDir + "/" + cleanFileName(item.Title()) + ".html"
		if _, err := os.Stat(fileName); os.IsNotExist(err) {
			article := Article{Item: item}
			res, err := article.Download()
			if err != nil {
				fmt.Errorf("Error downloading article: %s\n", err.Error())
				continue
			}

			if err := ioutil.WriteFile(fileName, res, 0644); err != nil {
				fmt.Errorf("Error writing to file: %s\n", err.Error())
				continue
			}
		}
	}
}

var fileNameCharacters = regexp.MustCompile(`(?m)[^\w]`)
var multipleDashes = regexp.MustCompile(`[-]+`)
var trailingDash = regexp.MustCompile(`-$`)

func cleanFileName(in string) string {
	// All non-filename characters to dashes.
	res := fileNameCharacters.ReplaceAll([]byte(in), []byte("-"))

	// Multiple dashes to single dash.
	res = multipleDashes.ReplaceAll(res, []byte("-"))

	// Remove trailing dash.
	res = trailingDash.ReplaceAll(res, []byte(""))
	return string(res)
}

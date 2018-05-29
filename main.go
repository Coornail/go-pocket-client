package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/motemen/go-pocket/api"
)

/*
TODO:
- [ ] Set creation time to the added time
*/

const clearScreen = "\033[H\033[2J"

// Flags
var (
	state     string
	outputDir string
	domain    string
	tag       string
	search    string
	force     bool
)

// @TODO set consumer key

func init() {
	flag.StringVar(&state, "state", string(api.StateUnread), "Type of article to download")
	flag.StringVar(&outputDir, "outputDir", "./articles", "Directory to download the articles to")
	flag.StringVar(&domain, "domain", "", "Domain to limit the archiving to")
	flag.StringVar(&search, "search", "", "Search to limit the archiving to")
	flag.StringVar(&tag, "tag", "", "Tag to limit the archiving to")
	flag.BoolVar(&force, "force", false, "Download already downloaded articles")
	flag.Parse()

	if state != string(api.StateUnread) && state != string(api.StateAll) && state != string(api.StateArchive) {
		fmt.Printf("State should be: %s, %s or %s\n", string(api.StateUnread), string(api.StateAll), string(api.StateUnread))
		os.Exit(1)
	}
}

func main() {
	client, err := GetClient()
	if err != nil {
		panic(err)
	}

	options := &api.RetrieveOption{
		State:  api.State(state),
		Domain: domain,
		Search: search,
		Tag:    tag,
	}
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
		if _, err := os.Stat(fileName); os.IsNotExist(err) || force {
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

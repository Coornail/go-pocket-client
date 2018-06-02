package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"regexp"

	"github.com/gosuri/uiprogress"
	"github.com/motemen/go-pocket/api"
)

/*
TODO:
- [ ] Set creation time to the added time
*/

// Flags
var (
	state       string
	outputDir   string
	domain      string
	tag         string
	search      string
	force       bool
	parallelism int

	consumerKey string
)

func init() {
	flag.StringVar(&state, "state", string(api.StateUnread), "Type of article to download")
	flag.StringVar(&outputDir, "outputDir", "./articles", "Directory to download the articles to")
	flag.StringVar(&domain, "domain", "", "Domain to limit the archiving to")
	flag.StringVar(&search, "search", "", "Search to limit the archiving to")
	flag.StringVar(&tag, "tag", "", "Tag to limit the archiving to")
	flag.BoolVar(&force, "force", false, "Redownload already downloaded articles")
	flag.IntVar(&parallelism, "parallelism", 8, "Number of threads to download the articles")
	flag.Parse()

	if state != string(api.StateUnread) && state != string(api.StateAll) && state != string(api.StateArchive) {
		fmt.Printf("State should be: %s, %s or %s\n", string(api.StateUnread), string(api.StateAll), string(api.StateArchive))
		os.Exit(1)
	}

	consumerKey = os.Getenv("POCKET_CONSUMER_KEY")
}

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

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
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	// Create download directory if it doesn't exist.
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err = os.MkdirAll(outputDir, 0744); err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(1)
		}
	}

	uiprogress.Start()
	bar := uiprogress.AddBar(len(res.List))
	bar.AppendCompleted()

	// Append the currently processed title
	var currentTitle string
	bar.AppendFunc(func(b *uiprogress.Bar) string {
		renderedTitle := currentTitle

		// Do not render longer lines than the width of the terminal.
		w, _ := getWidth()
		maxTitleLength := int(w) - b.Width - 32
		if len(renderedTitle) > maxTitleLength {
			renderedTitle = renderedTitle[:maxTitleLength] + "..."
		}

		res := fmt.Sprintf("[%d/%d] %s", bar.Current(), len(res.List), renderedTitle)
		for i := 0; i < maxTitleLength-len(currentTitle)+len("..."); i++ {
			res += " "
		}
		return res
	})

	// Prepare the job queue
	jobs := make(chan work, len(res.List))
	results := make(chan work, len(res.List))

	// Worker threads
	for w := 0; w < parallelism; w++ {
		go worker(jobs, results)
	}

	// Enqueue jobs.
	for _, item := range res.List {
		jobs <- work{input: item}
	}
	close(jobs)

	// Write results to file.
	for range res.List {
		res := <-results

		if res.alreadyDownloaded {
			fmt.Printf("Skipping %s\n\n", res.input.Title())
			continue
		}

		currentTitle = res.input.Title()
		if res.err != nil {
			fmt.Printf("Error downloading article: %s\n\n", res.err.Error())
			continue
		}

		if err := ioutil.WriteFile(res.outputFileName, res.output, 0644); err != nil {
			fmt.Printf("Error writing to file: %s\n\n", err.Error())
			continue
		}
		bar.Incr()
	}

}

// work stores a downloadable article, and the readability result.
type work struct {
	input             api.Item
	err               error
	output            []byte
	alreadyDownloaded bool
	outputFileName    string
}

// worker handles downloading articles from a job queue.
func worker(jobs <-chan work, results chan<- work) {
	for j := range jobs {
		title := j.input.Title()
		outputFileName := outputDir + "/" + cleanFileName(title) + ".html"
		if _, err := os.Stat(outputFileName); !os.IsNotExist(err) && !force {
			results <- work{alreadyDownloaded: true, input: j.input}
			continue
		}

		article := Article{Item: j.input}
		res, err := article.Download()
		results <- work{
			input:          j.input,
			output:         res,
			err:            err,
			outputFileName: outputFileName,
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

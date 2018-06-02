# go-pocket-client

[Pocket](https://getpocket.com/) client to download articles as readable HTML.

Great for using it on an ebook reader or to keep an archive for quickly grepping through it.

## Set up
* Create a [new app in pocket](https://getpocket.com/developer/apps/new) with at least _Retrieve_ permissions
* Export the resulting consumer key in your terminal: `export POCKET_CONSUMER_KEY=[YOUR-CONSUMER-KEY]`
* On first run, you will be asked to authorize the app by clicking the url

## Usage

```
Usage of ./go-pocket-client:
  -domain string
        Domain to limit the archiving to
  -force
        Download already downloaded articles
  -httptest.serve string
        if non-empty, httptest.NewServer serves on this address and blocks
  -outputDir string
        Directory to download the articles to (default "./articles")
  -parallelism int
        Number of threads to download the articles (default 8)
  -search string
        Search to limit the archiving to
  -state string
        Type of article to download (default "unread")
  -tag string
        Tag to limit the archiving to
```

# License
[MIT](https://raw.githubusercontent.com/Coornail/vim-go-conceal/master/LICENSE)

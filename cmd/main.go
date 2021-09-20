package main

import (
	"flag"
	"github.com/A1esandr/linkcheck"
	"os"
)

var urlFlag = flag.String("url", "", "URL of the site, for example, https://golang.org")
var htmlOnlyFlag = flag.Bool("htmlonly", true, "Scan only html pages, true or false")

func main() {
	flag.Parse()
	url := os.Getenv("URL")
	if len(url) == 0 {
		url = *urlFlag
	}
	htmlOnlyParam := true
	htmlOnly := os.Getenv("HTML_ONLY")
	if htmlOnly == "0" {
		htmlOnlyParam = false
	}
	if len(htmlOnly) == 0 {
		htmlOnlyParam = *htmlOnlyFlag
	}
	linkcheck.New(htmlOnlyParam).Start(url)
}

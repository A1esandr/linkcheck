package main

import (
	"flag"
	"os"

	"github.com/A1esandr/linkcheck"
)

var urlFlag = flag.String("url", "", "URL of the site, for example, https://golang.org")

func main() {
	flag.Parse()
	url := os.Getenv("URL")
	if len(url) == 0 {
		url = *urlFlag
	}
	linkcheck.New().Check(url)
}

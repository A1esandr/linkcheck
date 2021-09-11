package main

import (
	"flag"
	"github.com/A1esandr/linkcheck"
	"os"
)

var urlFlag = flag.String("url", "", "URL of the site, for example, https://golang.org")

func main() {
	flag.Parse()
	url := os.Getenv("URL")
	if len(url) == 0 {
		url = *urlFlag
	}
	linkcheck.New().Start(url)
}

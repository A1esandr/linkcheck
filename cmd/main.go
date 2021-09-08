package main

import (
	"flag"
	"fmt"
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
	results, err := linkcheck.New().Check(url)
	if err != nil {
		fmt.Println(err)
	}
	for from, state := range results {
		fmt.Printf("%s : %s \n", from, state)
	}
}

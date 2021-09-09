package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/A1esandr/linkcheck"
)

var urlFlag = flag.String("url", "", "URL of the site, for example, https://golang.org")

func main() {
	flag.Parse()
	url := os.Getenv("URL")
	if len(url) == 0 {
		url = *urlFlag
	}
	loaded := make(map[string]struct{})
	tocheck := make(map[string]struct{})
	errors := make(map[string]string)
	tocheck[url] = struct{}{}

	execute(tocheck, loaded, errors, url)

	for from, state := range errors {
		fmt.Printf("%s : %s \n", state, from)
	}
}

func execute(tocheck map[string]struct{}, loaded map[string]struct{}, errors map[string]string, url string) {
	for {
		collect := make(map[string]struct{})
		for key := range tocheck {
			results, err := linkcheck.New().Check(key)
			if err != nil {
				fmt.Println(err)
			}
			loaded[key] = struct{}{}
			newcheck := parseResults(results, loaded, errors, url)
			for k := range newcheck {
				collect[k] = struct{}{}
			}
		}
		if len(collect) == 0 {
			break
		}
		tocheck = collect
	}
}

func parseResults(results map[string]string, loaded map[string]struct{}, errors map[string]string, url string) map[string]struct{} {
	tocheck := make(map[string]struct{})
	for from, state := range results {
		fmt.Printf("%s : %s \n", state, from)
		if _, ok := loaded[from]; !ok && strings.HasPrefix(from, url) && strings.HasSuffix(from, ".html") {
			tocheck[from] = struct{}{}
		}
		if state != "OK" {
			errors[from] = state
		}
	}
	return tocheck
}

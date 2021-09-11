package linkcheck

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/A1esandr/crawler"
)

type (
	checker struct {
	}

	Checker interface {
		Start(url string)
		Check(url string) (map[string]string, error)
	}
)

func New() Checker {
	return &checker{}
}

func (c *checker) Start(url string) {
	loaded := make(map[string]struct{})
	tocheck := make(map[string]struct{})
	errors := make(map[string]string)
	tocheck[url] = struct{}{}

	execute(tocheck, loaded, errors, url)

	for from, state := range errors {
		fmt.Printf("%s : %s \n", state, from)
	}
}

func (c *checker) Check(url string) (map[string]string, error) {
	links, err := crawler.New().Run(url)
	if err != nil {
		return nil, err
	}
	result := make(map[string]string)
	for _, link := range links {
		err = c.check(link, 0)
		if err != nil {
			result[link] = err.Error()
			continue
		}
		result[link] = "OK"
	}
	return result, nil
}

func (c *checker) check(url string, count int) error {
	client := http.Client{
		Timeout: 2 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	if resp == nil {
		return errors.New("nil response")
	}
	if resp.StatusCode != http.StatusOK && count < 3 {
		if count == 2 {
			return fmt.Errorf("not downloaded, status %d", resp.StatusCode)
		}
		log.Println("Error loading", url)
		time.Sleep(time.Duration(300+rand.Intn(1000)) * time.Millisecond)
		return c.check(url, count+1)
	}
	return nil
}

func execute(tocheck map[string]struct{}, loaded map[string]struct{}, errors map[string]string, url string) {
	for {
		collect := make(map[string]struct{})
		for key := range tocheck {
			results, err := New().Check(key)
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

package linkcheck

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/A1esandr/crawler"
)

type (
	checker struct {
		checked *syncSet
		url     string
		htmlOnly bool
	}

	Checker interface {
		Start(url string)
		Check(url string) (map[string]string, error)
	}

	syncSet struct {
		items map[string]struct{}
		mu    sync.Mutex
	}

	syncMap struct {
		items map[string]string
		mu    sync.Mutex
	}
)

func New(htmlOnly bool) Checker {
	return &checker{checked: &syncSet{items: make(map[string]struct{})}, htmlOnly: htmlOnly}
}

func (c *checker) Start(url string) {
	loaded := &syncSet{items: make(map[string]struct{})}
	tocheck := make(map[string]struct{})
	errs := &syncMap{items: make(map[string]string)}
	tocheck[url] = struct{}{}

	c.execute(tocheck, loaded, errs, url)

	for from, state := range errs.items {
		fmt.Printf("Not OK: %s : %s \n", state, from)
	}
	if len(errs.items) == 0 {
		fmt.Println("Finished without errors")
	}
}

func (c *checker) Check(url string) (map[string]string, error) {
	links, err := crawler.New().Run(url)
	if err != nil {
		return nil, err
	}
	result := make(map[string]string)
	for _, link := range links {
		c.checked.mu.Lock()
		_, ok := c.checked.items[link]
		c.checked.mu.Unlock()
		if ok {
			continue
		}
		err = c.check(link, 0)
		c.checked.mu.Lock()
		c.checked.items[link] = struct{}{}
		c.checked.mu.Unlock()
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

func (c *checker) valid(from, url string) bool {
	if strings.HasPrefix(from, url) {
		if c.htmlOnly {
			if strings.HasSuffix(from, ".html") {
				return true
			}
		} else {
			return true
		}
	}
	return false
}

func (c *checker) execute(tocheck map[string]struct{}, loaded *syncSet, errs *syncMap, url string) {
	var wg sync.WaitGroup
	for {
		collect := &syncSet{items: make(map[string]struct{})}
		for key := range tocheck {
			wg.Add(1)
			go func(key string) {
				results, err := c.Check(key)
				if err != nil {
					fmt.Println(err)
				}
				loaded.mu.Lock()
				errs.mu.Lock()
				loaded.items[key] = struct{}{}
				newcheck := c.parseResults(results, loaded.items, errs.items, url)
				loaded.mu.Unlock()
				errs.mu.Unlock()

				collect.mu.Lock()
				for k := range newcheck {
					collect.items[k] = struct{}{}
				}
				collect.mu.Unlock()
				wg.Done()
			}(key)
		}
		wg.Wait()
		if len(collect.items) == 0 {
			break
		}
		tocheck = collect.items
	}
}

func (c *checker) parseResults(results map[string]string, loaded map[string]struct{}, errs map[string]string, url string) map[string]struct{} {
	tocheck := make(map[string]struct{})
	for from, state := range results {
		fmt.Printf("%s : %s \n", state, from)
		if _, ok := loaded[from]; !ok && c.valid(from, url) {
			tocheck[from] = struct{}{}
		}
		if state != "OK" {
			errs[from] = state
		}
	}
	return tocheck
}

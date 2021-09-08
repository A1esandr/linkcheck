package linkcheck

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/A1esandr/crawler"
)

type (
	checker struct {
	}

	Checker interface {
		Check(url string) (map[string]string, error)
	}
)

func New() Checker {
	return &checker{}
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

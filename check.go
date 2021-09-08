package linkcheck

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
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
	return nil, nil
}

func (c *checker) check(url string, count int) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	if resp == nil {
		return fmt.Errorf("nil response from %s", url)
	}
	if resp.StatusCode != http.StatusOK && count < 3 {
		if count == 2 {
			return fmt.Errorf("not downloaded %s", url)
		}
		log.Println("Error loading", url)
		time.Sleep(time.Duration(300+rand.Intn(1000)) * time.Millisecond)
		return c.check(url, count+1)
	}
	return nil
}

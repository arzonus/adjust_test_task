package walker

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"
)

//go:generate mockgen -destination=walker_mock.go -package=walker "github.com/arzonus/adjust_test_task" HTTPClient
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

const defaultParallel = 10

func NewWalker(parallel uint, httpClient HTTPClient) *Walker {
	if parallel == 0 {
		parallel = defaultParallel
	}
	return &Walker{
		parallel:   parallel,
		httpClient: httpClient,
	}
}

type Walker struct {
	parallel   uint
	httpClient HTTPClient
}

func (w Walker) Walk(addresses ...string) error {
	for i := range addresses {
		if _, err := url.Parse(addresses[i]); err != nil {
			return fmt.Errorf("couldn't parse %s as url, %w", addresses[i], err)
		}
	}

	return w.walk(addresses...)
}

func (w Walker) walk(addresses ...string) error {
	var (
		wg       sync.WaitGroup
		parallel = make(chan struct{}, w.parallel)
	)

	wg.Add(len(addresses))

	for i := range addresses {
		go func(i int) {
			parallel <- struct{}{}
			defer func() {
				wg.Done()
				<-parallel
			}()

			hash, err := GetMD5Response(w.httpClient, addresses[i])
			if err != nil {
				fmt.Printf("%s %s\n", addresses[i], err)
			} else {
				fmt.Printf("%s %s\n", addresses[i], hex.EncodeToString(hash[:]))
			}
		}(i)
	}

	wg.Wait()
	return nil
}

func GetMD5Response(httpClient HTTPClient, address string) ([md5.Size]byte, error) {
	u, err := url.Parse(address)
	if err != nil {
		return [md5.Size]byte{}, fmt.Errorf("address %s is invalid, %w", address, err)
	}

	if u.Scheme == "" {
		u.Scheme = "http"
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return [md5.Size]byte{}, fmt.Errorf("couldn't create http request to %s, %w", address, err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return [md5.Size]byte{}, fmt.Errorf("couldn't request to %s, %w", address, err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("couldn't close request's body to %s, %s\n", address, err)
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return [md5.Size]byte{}, fmt.Errorf("couldn't read request's body to %s, %w", address, err)
	}

	return md5.Sum(body), nil
}

package http

import (
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	url2 "net/url"
	"time"
)

func Get(proxy *url2.URL, header map[string]string, url string) (*http.Response, error) {

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	if proxy != nil {
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxy),
		}
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if header != nil {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}
	if err != nil {
		log.Errorf(err.Error())
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf(err.Error())
		return nil, err
	}

	return resp, nil
}

func Post(url string, bodyReader io.Reader) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodPost, url, bodyReader)
	if err != nil {
		log.Errorf(err.Error())
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Errorf(err.Error())
		return nil, err
	}

	return resp, nil

}

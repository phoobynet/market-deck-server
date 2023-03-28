package sechttp

import (
	"github.com/gojek/heimdall/v7/httpclient"
	"io"
	"net/http"
	"time"
)

// Get - Makes a request with the SEC required headers
func Get(url string) ([]byte, error) {
	timeout := 10 * time.Second

	client := httpclient.NewClient(httpclient.WithHTTPTimeout(timeout))

	headers := make(http.Header)
	headers.Set("User-Agent", "demo project sec@g1nger.com")

	res, err := client.Get(url, headers)

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)

	return body, nil
}

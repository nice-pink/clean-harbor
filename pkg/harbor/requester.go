package harbor

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// HttpClient is a Requester which returns an error for any non-2xx response.
// This is essential: silently treating 401/404/HTML responses as empty
// results makes the cleaner a no-op without any indication of failure.

type HttpClient struct {
	client   *http.Client
	user     string
	password string
	logLevel int
}

func NewHttpClient(user string, password string, timeoutSeconds int, logLevel int) *HttpClient {
	return &HttpClient{
		client:   &http.Client{Timeout: time.Duration(timeoutSeconds) * time.Second},
		user:     user,
		password: password,
		logLevel: logLevel,
	}
}

func (c *HttpClient) Get(url string, printBody bool) ([]byte, error) {
	if c.logLevel > 0 {
		fmt.Println("Get:", url)
	}
	body, status, err := c.do(http.MethodGet, url)
	if err != nil {
		return nil, err
	}
	if status < 200 || status > 299 {
		return nil, fmt.Errorf("HTTP %d for GET %s: %s", status, url, truncateBody(body))
	}
	if printBody {
		fmt.Println(string(body))
	}
	return body, nil
}

func (c *HttpClient) Delete(url string) (bool, error) {
	if c.logLevel > 0 {
		fmt.Println("Delete:", url)
	}
	body, status, err := c.do(http.MethodDelete, url)
	if err != nil {
		return false, err
	}
	if status != http.StatusOK && status != http.StatusAccepted {
		return false, fmt.Errorf("HTTP %d for DELETE %s: %s", status, url, truncateBody(body))
	}
	fmt.Println("Success! Deleted:", url)
	return true, nil
}

func (c *HttpClient) do(method string, url string) ([]byte, int, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, 0, err
	}
	if c.user != "" {
		req.SetBasicAuth(c.user, c.password)
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	return body, resp.StatusCode, nil
}

func truncateBody(body []byte) string {
	const max = 256
	if len(body) > max {
		return string(body[:max]) + "..."
	}
	return string(body)
}

package request

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/nice-pink/clean-harbor/pkg/network"
)

type Requester struct {
}

// request

func (r *Requester) Get(url string, auth network.Auth, printBody bool) ([]byte, error) {
	// build request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// add basic auth
	if auth.BasicUser != "" && auth.BasicPassword != "" {
		req.SetBasicAuth(auth.BasicUser, auth.BasicPassword)
	}

	// request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	// read and return
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	if printBody {
		fmt.Println(string(body))
	}

	return body, err
}

func (r *Requester) Delete(url string, auth network.Auth) (bool, error) {
	// build request
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	// add basic auth
	if auth.BasicUser != "" && auth.BasicPassword != "" {
		req.SetBasicAuth(auth.BasicUser, auth.BasicPassword)
	}

	// request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	defer resp.Body.Close()

	// read and return
	if resp.StatusCode != 200 && resp.StatusCode != 202 {
		fmt.Println("Could not delete. Status code:", strconv.Itoa(resp.StatusCode))
		return false, nil
	}
	return true, nil
}

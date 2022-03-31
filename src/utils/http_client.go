package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

// Get 发起一个GET请求
func Get(url string, headers map[string]string) (response []byte, err error) {
	client := http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	if len(headers) > 0 {
		for key, val := range headers {
			req.Header.Set(key, val)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var i int
TRY:
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		i++
		if i < 3 {
			time.Sleep(time.Second)
			goto TRY
		}
		return nil, err
	}

	return body, nil
}

// Post 发起一个Post请求
func Post(url string, data interface{}, contentType string) (content []byte, err error) {
	jsonStr, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}
	req.Header.Add("content-type", contentType)
	defer req.Body.Close()

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return result, nil
}

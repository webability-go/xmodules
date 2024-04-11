package tools

import (
	"crypto/sha256"
	"encoding/hex"

	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func CallAPI(apiservice string, postparams map[string]string, forwarded string) (map[string]interface{}, error) {

	formdata := url.Values{}
	for key, val := range postparams {
		formdata[key] = []string{val}
	}

	req, err := http.NewRequest(http.MethodPost, apiservice, strings.NewReader(formdata.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-Forwarded-For", forwarded)

	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 120 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 120 * time.Second,
	}
	var netClient = &http.Client{
		Timeout:   time.Second * 30,
		Transport: netTransport,
	}
	data, err := netClient.Do(req)
	if err != nil {
		// Log error (timeout? no server?)
		fmt.Println(apiservice, postparams, err)
		return nil, err
	}
	defer data.Body.Close()

	var result map[string]interface{}
	if data.StatusCode != 200 {
		return nil, errors.New("Error: wrong API response")
	}

	err = json.NewDecoder(data.Body).Decode(&result)
	if err != nil {
		fmt.Println(apiservice, postparams, err)
		return nil, err
	}
	return result, nil
}

func GetTimeDigest(digestkey string) string {
	currentTime := time.Now()
	sum := sha256.Sum256([]byte(digestkey + currentTime.Format("2006-01-02")))
	return hex.EncodeToString(sum[:])
}

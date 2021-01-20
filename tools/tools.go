package tools

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func UUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}

/* Creates a string of random chars and digits of length
   chartype = -1: lowers, uppers and digits
   chartype = 0: only digits
   chartype = 1 only uppers
   chartype = 2 only lowers
*/
func CreateKey(length int, chartype int) string {
	rand.Seed(time.Now().UnixNano())
	result, rd := "", 0
	for i := 0; i < length; i++ {
		if chartype == -1 {
			rd = rand.Intn(3)
		} else {
			rd = chartype
		}
		key := 0
		switch rd {
		case 0:
			key = 48 + rand.Intn(10)
		case 1:
			key = 65 + rand.Intn(26)
		case 2:
			key = 97 + rand.Intn(26)
		}
		result += string(key)
	}
	return result
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func calljson(jsonservice string, jsonparams map[string]string, forwarded string) (map[string]interface{}, error) {

	formdata := url.Values{}
	for key, val := range jsonparams {
		formdata[key] = []string{val}
	}

	req, err := http.NewRequest(http.MethodPost, jsonservice, strings.NewReader(formdata.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-Forwarded-For", forwarded)
	data, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer data.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(data.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

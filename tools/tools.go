package tools

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/webability-go/xdominion"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
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

/*
Creates a string of random chars and digits of length

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

func CallJSON(jsonservice string, jsonparams map[string]string, forwarded string) (map[string]interface{}, error) {

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

// StripTags ...
func StripTags(data string) string {
	reg, _ := regexp.Compile("(<([^>]+)>)")
	safe := reg.ReplaceAllString(data, "")
	return safe
}

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}

// RemoveAccents ...
func RemoveAccents(data string) string {
	b := make([]byte, len(data))
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	_, _, e := t.Transform(b, []byte(data), true)
	if e != nil {
		return ""
	}
	return string(b)
}

// RemoveNonLetters ...
func RemoveNonLetters(data string) string {
	reg, _ := regexp.Compile("([^一-龠ぁ-ゔァ-ヴーa-zA-Z0-9_ ａ-ｚＡ-Ｚ０-９々〆〤ヶ])")
	safe := reg.ReplaceAllString(data, " ")
	return safe
}

// ReplaceSpaces ...
func ReplaceSpaces(data string) string {
	reg, _ := regexp.Compile("(\\s{1,})")
	safe := reg.ReplaceAllString(data, "-")
	return safe
}

// StripWord ...
func StripWord(data string) string {
	// transform escaped things
	safe := html.UnescapeString(data)
	// remove tags
	safe = StripTags(safe)
	// to lowers
	safe = strings.ToLower(safe)
	// replace all accents
	safe = RemoveAccents(safe)
	// strip anything that is not a-z, 0-9, _
	safe = RemoveNonLetters(safe)
	safe = strings.Trim(safe, " \r\n\t")
	return safe
}

// BuildLink ...
func BuildLink(data string) string {
	safe := StripWord(data)
	safe = ReplaceSpaces(safe)
	return safe
}

func CompareRecords(r1, r2 *xdominion.XRecord, ignore []string) string {
	if r1 == nil && r2 != nil {
		return CompareRecords(r2, nil, ignore)
	}
	if r1 == nil {
		return "Nothing to compare"
	}
	if r2 == nil {
		// prints the full record
		return r1.String()
	}
	message := ""
	some := false
	for k := range *r1 {
		// ignore some fields
		if ignore != nil {
			found := false
			for _, i := range ignore {
				if i == k {
					found = true
					break
				}
			}
			if found {
				continue
			}
		}
		v1, _ := r1.GetString(k)
		v2, _ := r2.GetString(k)
		if v1 != v2 {
			message += " - @" + k + "[" + v1 + "]>>[" + v2 + "]"
			some = true
		}
	}
	if !some {
		message = " - No changes"
	}
	return message
}

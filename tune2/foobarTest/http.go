package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
)

type test struct {
	test,
	test1 string
}

func (test *test) Read(b []byte) (int, error) {
	n := bytes.NewReader(b)

	nResult := make([]byte, 2046)

	n1, err := n.Read(nResult)
	if err != nil {
		return 0, err
	}
	return n1, err

}

func main4() {
	// get, err := http.Get("http://www.google.com/robots.txt")
	// if err != nil {
	// 	log.Panic(err)
	// }
	// defer get.Body.Close()
	// fmt.Printf("get %+v\n", get.Body)
	// head, err := http.Head("http://www.google.com/robots.txt")
	// if err != nil {
	// 	log.Panic(err)
	// }
	// defer head.Body.Close()
	// fmt.Printf("header %+v\n", head.Body)
	// // form := test{"test","Test"}
	// form := get.Request.Form
	// form.Add("username", "test")
	// form.Add("password", "test")
	// post, err := http.PostForm("http://www.google.com/robots.txt",form )
	// if err != nil {
	// 	log.Panic(err)
	// }
	// defer post.Body.Close()
	// fmt.Printf("post %+v\n", post.Body)

	g, err := http.Get("http://www.google.com/robots.txt")

	if err != nil {
		log.Panic(err)
	}

	fmt.Println("status code ", g.Status)
	defer g.Body.Close()

	all, err := io.ReadAll(g.Body)

	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("%+s\n", all)
}

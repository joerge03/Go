package main

import "fmt"

type servOpts (func(detail *Detail))

type Detail struct {
	Name  string
	Url   string
	Info  string
	isTCP bool
}

type Server struct {
	Detail Detail
}

func DefaultDetail() Detail {
	return Detail{
		Name: "hihi",
		Url:  "hihi.com",
		Info: "meow",
	}
}

func isTCP(d *Detail) {

	d.isTCP = true
}

func updateInfo(s string) servOpts {
	return func(d *Detail) {
		d.Info = s
	}

}

func NewServer(opts ...servOpts) *Server {
	defaultOptions := DefaultDetail()
	for _, fn := range opts {
		fn(&defaultOptions)
	}
	return &Server{
		Detail: defaultOptions,
	}
}

func main() {

	server := NewServer(updateInfo("asdfasdf"), isTCP)

	fmt.Printf("%+v\n", server)
}

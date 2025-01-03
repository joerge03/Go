package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)



func NewClient(key string) *Client{
	return &Client{
		API_KEY: key,
	}
}

func (c *Client) APIInfo() {
	fmt.Printf("api key %v",c.API_KEY)
	infoLink := fmt.Sprintf("%s/api-info?key=%v",SHODAN_URL,c.API_KEY)
	
	r, err := http.Get(infoLink)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(r.Status)
	defer r.Body.Close()
	data:= new(any)
	
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil{
		log.Panic(err)
	}

	marshalledData,err := json.MarshalIndent(data,"", "\t")

	if err != nil {
		log.Panic(err)
	}
	// test := json.Unmarshal(data,&data)

	fmt.Printf("resutl: %+v\n", string(marshalledData))	
}

func (c *Client) HostSearch(){
	hostLink := fmt.Sprintf("%v/shodan/host/search?key=%v&query=product:nginx&facets=country", SHODAN_URL,c.API_KEY)

	res, err := http.Get(hostLink)
	if err != nil {
		log.Panic(err)
	}
	defer res.Body.Close()	
	data := new(any)
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		log.Panic(err)
	}
	formattedJsonData, err  := json.MarshalIndent(data,"","\t")
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("%+v\n", string(formattedJsonData))
}